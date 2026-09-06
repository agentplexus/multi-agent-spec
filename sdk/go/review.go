package multiagentspec

import (
	"fmt"
	"time"
)

// ReviewStatus represents the state of a human review checkpoint.
type ReviewStatus string

const (
	// ReviewPending means the item has been proposed but not yet reviewed.
	ReviewPending ReviewStatus = "pending"

	// ReviewConfirmed means a human reviewed the proposed value and accepted it as-is.
	ReviewConfirmed ReviewStatus = "confirmed"

	// ReviewEdited means a human reviewed the proposed value and replaced it with their own.
	ReviewEdited ReviewStatus = "edited"

	// ReviewRejected means a human reviewed the item and rejected it outright, with no final value.
	ReviewRejected ReviewStatus = "rejected"
)

// IsResolved reports whether the status reflects a completed human decision.
func (s ReviewStatus) IsResolved() bool {
	switch s {
	case ReviewConfirmed, ReviewEdited, ReviewRejected:
		return true
	default:
		return false
	}
}

// ReviewItem is a single agent-produced output flagged for human confirmation
// before it is trusted as ground truth or allowed to proceed.
//
// This is distinct from Loop's escalation:human, which pauses a retry loop once
// max attempts are exhausted. A ReviewItem is a per-item checkpoint usable
// anywhere an agent's output needs a human decision before acceptance — for
// example, confirming a gold-label candidate, approving a flagged finding, or
// gating an action before it takes effect.
type ReviewItem struct {
	// ID identifies the subject being reviewed (e.g. an incident id).
	ID string `json:"id" yaml:"id"`

	// ProposedValue is what the agent proposed, as a JSON-serializable value.
	ProposedValue interface{} `json:"proposed_value" yaml:"proposed_value"`

	// Reason explains why this item was produced or flagged (e.g. low
	// confidence, an ambiguous case, a conflicting signal).
	Reason string `json:"reason,omitempty" yaml:"reason,omitempty"`

	// Confidence is the agent's confidence in ProposedValue, in [0,1].
	Confidence float64 `json:"confidence,omitempty" yaml:"confidence,omitempty"`

	// Status is the current review state. Defaults to "pending".
	Status ReviewStatus `json:"status,omitempty" yaml:"status,omitempty"`

	// FinalValue is the human-confirmed or human-edited value: a copy of
	// ProposedValue once Status is "confirmed", or the human's replacement once
	// Status is "edited". Unset while Status is "pending" or "rejected".
	FinalValue interface{} `json:"final_value,omitempty" yaml:"final_value,omitempty"`

	// Reviewer identifies who made the decision (e.g. a name or handle).
	Reviewer string `json:"reviewer,omitempty" yaml:"reviewer,omitempty"`

	// Notes are optional reviewer remarks, particularly useful on rejection.
	Notes string `json:"notes,omitempty" yaml:"notes,omitempty"`

	// ReviewedAt is when the decision was made.
	ReviewedAt *time.Time `json:"reviewed_at,omitempty" yaml:"reviewed_at,omitempty"`
}

// EffectiveStatus returns Status, defaulting to "pending" when unset.
func (r *ReviewItem) EffectiveStatus() ReviewStatus {
	if r.Status == "" {
		return ReviewPending
	}
	return r.Status
}

// Confirm accepts ProposedValue as FinalValue.
func (r *ReviewItem) Confirm(reviewer string, at time.Time) {
	r.Status = ReviewConfirmed
	r.FinalValue = r.ProposedValue
	r.Reviewer = reviewer
	r.ReviewedAt = &at
}

// Edit replaces ProposedValue with a human-supplied value.
func (r *ReviewItem) Edit(value interface{}, reviewer string, at time.Time) {
	r.Status = ReviewEdited
	r.FinalValue = value
	r.Reviewer = reviewer
	r.ReviewedAt = &at
}

// Reject marks the item as rejected, with no final value.
func (r *ReviewItem) Reject(reviewer, notes string, at time.Time) {
	r.Status = ReviewRejected
	r.Reviewer = reviewer
	r.Notes = notes
	r.ReviewedAt = &at
}

// ReviewBatch is a set of ReviewItems produced by one agent run, awaiting or
// having completed human review.
type ReviewBatch struct {
	// ID identifies this batch (e.g. a UUID or a caller-assigned slug).
	ID string `json:"id" yaml:"id"`

	// Source identifies what produced the batch (e.g. an agent name or rubric id).
	Source string `json:"source,omitempty" yaml:"source,omitempty"`

	// CreatedAt is when the batch was produced.
	CreatedAt time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`

	// Items are the individual review checkpoints in this batch.
	Items []ReviewItem `json:"items" yaml:"items"`
}

// NewReviewBatch creates an empty ReviewBatch.
func NewReviewBatch(id, source string, createdAt time.Time) *ReviewBatch {
	return &ReviewBatch{ID: id, Source: source, CreatedAt: createdAt}
}

// AddItem appends a review item and returns the batch for chaining.
func (b *ReviewBatch) AddItem(item ReviewItem) *ReviewBatch {
	b.Items = append(b.Items, item)
	return b
}

// Pending returns items that have not yet received a human decision.
func (b *ReviewBatch) Pending() []ReviewItem {
	var out []ReviewItem
	for _, it := range b.Items {
		if !it.EffectiveStatus().IsResolved() {
			out = append(out, it)
		}
	}
	return out
}

// Resolved returns items that have received a human decision.
func (b *ReviewBatch) Resolved() []ReviewItem {
	var out []ReviewItem
	for _, it := range b.Items {
		if it.EffectiveStatus().IsResolved() {
			out = append(out, it)
		}
	}
	return out
}

// AllResolved reports whether every item in the batch has been reviewed.
func (b *ReviewBatch) AllResolved() bool {
	return len(b.Pending()) == 0
}

// Validate checks the batch for structural consistency: a non-empty id, and
// items with non-empty, unique ids and in-range confidence.
func (b *ReviewBatch) Validate() error {
	if b.ID == "" {
		return fmt.Errorf("review batch id is required")
	}
	seen := make(map[string]bool, len(b.Items))
	for i, it := range b.Items {
		if it.ID == "" {
			return fmt.Errorf("item %d: id is required", i)
		}
		if seen[it.ID] {
			return fmt.Errorf("item %d: duplicate id %q", i, it.ID)
		}
		seen[it.ID] = true
		if it.Confidence < 0 || it.Confidence > 1 {
			return fmt.Errorf("item %q: confidence %v out of [0,1]", it.ID, it.Confidence)
		}
	}
	return nil
}

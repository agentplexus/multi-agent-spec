package multiagentspec

import (
	"encoding/json"
	"testing"
	"time"
)

func TestReviewStatusIsResolved(t *testing.T) {
	tests := []struct {
		status ReviewStatus
		want   bool
	}{
		{ReviewPending, false},
		{"", false},
		{ReviewConfirmed, true},
		{ReviewEdited, true},
		{ReviewRejected, true},
	}
	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsResolved(); got != tt.want {
				t.Errorf("IsResolved() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReviewItem_EffectiveStatus(t *testing.T) {
	it := ReviewItem{}
	if got := it.EffectiveStatus(); got != ReviewPending {
		t.Errorf("EffectiveStatus() = %v, want %v", got, ReviewPending)
	}
	it.Status = ReviewConfirmed
	if got := it.EffectiveStatus(); got != ReviewConfirmed {
		t.Errorf("EffectiveStatus() = %v, want %v", got, ReviewConfirmed)
	}
}

func TestReviewItem_Confirm(t *testing.T) {
	it := ReviewItem{ID: "x", ProposedValue: "LLM", Confidence: 0.4}
	at := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	it.Confirm("alice", at)

	if it.Status != ReviewConfirmed {
		t.Fatalf("Status = %v, want confirmed", it.Status)
	}
	if it.FinalValue != "LLM" {
		t.Fatalf("FinalValue = %v, want LLM", it.FinalValue)
	}
	if it.Reviewer != "alice" || it.ReviewedAt == nil || !it.ReviewedAt.Equal(at) {
		t.Fatalf("reviewer/time not set: %+v", it)
	}
}

func TestReviewItem_Edit(t *testing.T) {
	it := ReviewItem{ID: "x", ProposedValue: "LLM"}
	at := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	it.Edit("Identity", "alice", at)

	if it.Status != ReviewEdited {
		t.Fatalf("Status = %v, want edited", it.Status)
	}
	if it.FinalValue != "Identity" {
		t.Fatalf("FinalValue = %v, want Identity", it.FinalValue)
	}
	if it.ProposedValue != "LLM" {
		t.Fatalf("ProposedValue should be unchanged, got %v", it.ProposedValue)
	}
}

func TestReviewItem_Reject(t *testing.T) {
	it := ReviewItem{ID: "x", ProposedValue: "LLM"}
	at := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	it.Reject("alice", "not a real incident", at)

	if it.Status != ReviewRejected {
		t.Fatalf("Status = %v, want rejected", it.Status)
	}
	if it.FinalValue != nil {
		t.Fatalf("FinalValue should be nil on rejection, got %v", it.FinalValue)
	}
	if it.Notes != "not a real incident" {
		t.Fatalf("Notes = %q", it.Notes)
	}
}

func TestReviewBatch_PendingResolvedAllResolved(t *testing.T) {
	b := NewReviewBatch("batch-1", "domain-classifier", time.Now())
	b.AddItem(ReviewItem{ID: "a", ProposedValue: "LLM"})
	b.AddItem(ReviewItem{ID: "b", ProposedValue: "Agent"})

	if len(b.Pending()) != 2 || len(b.Resolved()) != 0 {
		t.Fatalf("expected 2 pending, 0 resolved; got pending=%d resolved=%d", len(b.Pending()), len(b.Resolved()))
	}
	if b.AllResolved() {
		t.Fatalf("AllResolved() = true, want false")
	}

	b.Items[0].Confirm("alice", time.Now())
	if len(b.Pending()) != 1 || len(b.Resolved()) != 1 {
		t.Fatalf("expected 1 pending, 1 resolved; got pending=%d resolved=%d", len(b.Pending()), len(b.Resolved()))
	}

	b.Items[1].Reject("alice", "n/a", time.Now())
	if !b.AllResolved() {
		t.Fatalf("AllResolved() = false, want true")
	}
}

func TestReviewBatch_Validate(t *testing.T) {
	tests := []struct {
		name    string
		batch   ReviewBatch
		wantErr bool
	}{
		{
			name:  "valid",
			batch: ReviewBatch{ID: "b1", Items: []ReviewItem{{ID: "a", Confidence: 0.5}, {ID: "b", Confidence: 1}}},
		},
		{
			name:    "missing batch id",
			batch:   ReviewBatch{Items: []ReviewItem{{ID: "a"}}},
			wantErr: true,
		},
		{
			name:    "missing item id",
			batch:   ReviewBatch{ID: "b1", Items: []ReviewItem{{ID: ""}}},
			wantErr: true,
		},
		{
			name:    "duplicate item id",
			batch:   ReviewBatch{ID: "b1", Items: []ReviewItem{{ID: "a"}, {ID: "a"}}},
			wantErr: true,
		},
		{
			name:    "confidence out of range",
			batch:   ReviewBatch{ID: "b1", Items: []ReviewItem{{ID: "a", Confidence: 1.5}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.batch.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReviewBatch_JSONRoundTrip(t *testing.T) {
	at := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	b := NewReviewBatch("batch-1", "domain-classifier", at)
	b.AddItem(ReviewItem{ID: "aiid-1", ProposedValue: "LLM", Confidence: 0.4, Reason: "ambiguous"})

	data, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out ReviewBatch
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.ID != "batch-1" || len(out.Items) != 1 || out.Items[0].ID != "aiid-1" {
		t.Fatalf("round-trip mismatch: %+v", out)
	}
	if err := out.Validate(); err != nil {
		t.Fatalf("round-tripped batch invalid: %v", err)
	}
}

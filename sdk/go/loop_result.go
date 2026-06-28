package multiagentspec

import "time"

// LoopStatus represents the outcome of a loop execution.
type LoopStatus string

const (
	// LoopStatusSuccess indicates the loop completed successfully.
	LoopStatusSuccess LoopStatus = "success"

	// LoopStatusMaxAttempts indicates max attempts were reached.
	LoopStatusMaxAttempts LoopStatus = "max_attempts"

	// LoopStatusEscalated indicates the loop was escalated per policy.
	LoopStatusEscalated LoopStatus = "escalated"

	// LoopStatusAborted indicates the loop was aborted due to error.
	LoopStatusAborted LoopStatus = "aborted"

	// LoopStatusInProgress indicates the loop is still running.
	LoopStatusInProgress LoopStatus = "in_progress"
)

// IsTerminal returns true if this status represents a terminal state.
func (s LoopStatus) IsTerminal() bool {
	switch s {
	case LoopStatusSuccess, LoopStatusMaxAttempts, LoopStatusEscalated, LoopStatusAborted:
		return true
	default:
		return false
	}
}

// IsSuccess returns true if this status represents successful completion.
func (s LoopStatus) IsSuccess() bool {
	return s == LoopStatusSuccess
}

// CheckStatus represents the outcome of a single check.
type CheckStatus string

const (
	// CheckStatusGO indicates the check passed.
	CheckStatusGO CheckStatus = "GO"

	// CheckStatusWARN indicates the check passed with warnings.
	CheckStatusWARN CheckStatus = "WARN"

	// CheckStatusNOGO indicates the check failed.
	CheckStatusNOGO CheckStatus = "NO-GO"

	// CheckStatusSKIP indicates the check was skipped.
	CheckStatusSKIP CheckStatus = "SKIP"
)

// IsPass returns true if this status represents a passing check.
func (s CheckStatus) IsPass() bool {
	return s == CheckStatusGO || s == CheckStatusWARN
}

// CheckResult represents the outcome of a single check execution.
type CheckResult struct {
	// ID is the check identifier.
	ID string `json:"id"`

	// Status is the check outcome.
	Status CheckStatus `json:"status"`

	// Detail provides additional information about the result.
	Detail string `json:"detail,omitempty"`

	// Metadata contains structured data about the check.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Duration is how long the check took.
	Duration time.Duration `json:"duration,omitempty"`
}

// IterationResult represents the outcome of a single loop iteration.
type IterationResult struct {
	// Attempt is the iteration number (1-indexed).
	Attempt int `json:"attempt"`

	// ValidationStatus is the overall validation result.
	ValidationStatus CheckStatus `json:"validation_status"`

	// CheckResults are the individual check outcomes.
	CheckResults []CheckResult `json:"check_results,omitempty"`

	// ActionsTaken describes what the actor did.
	ActionsTaken []string `json:"actions_taken,omitempty"`

	// StartedAt is when this iteration started.
	StartedAt time.Time `json:"started_at"`

	// CompletedAt is when this iteration completed.
	CompletedAt time.Time `json:"completed_at"`
}

// Duration returns the iteration duration.
func (r *IterationResult) Duration() time.Duration {
	return r.CompletedAt.Sub(r.StartedAt)
}

// FailedChecks returns checks that did not pass.
func (r *IterationResult) FailedChecks() []CheckResult {
	var failed []CheckResult
	for _, check := range r.CheckResults {
		if !check.Status.IsPass() {
			failed = append(failed, check)
		}
	}
	return failed
}

// LoopResult represents the complete outcome of a loop execution.
type LoopResult struct {
	// LoopName is the name of the loop that was executed.
	LoopName string `json:"loop_name"`

	// LoopType is REAL or VEAL.
	LoopType LoopType `json:"loop_type"`

	// Status is the overall loop outcome.
	Status LoopStatus `json:"status"`

	// Iterations are the results from each loop iteration.
	Iterations []IterationResult `json:"iterations"`

	// TotalAttempts is how many iterations were executed.
	TotalAttempts int `json:"total_attempts"`

	// StartedAt is when the loop started.
	StartedAt time.Time `json:"started_at"`

	// CompletedAt is when the loop completed.
	CompletedAt time.Time `json:"completed_at"`

	// EscalationReason explains why escalation occurred (if applicable).
	EscalationReason string `json:"escalation_reason,omitempty"`

	// Error contains any error message.
	Error string `json:"error,omitempty"`

	// Outputs contains data produced by the loop.
	Outputs map[string]interface{} `json:"outputs,omitempty"`
}

// Duration returns the total loop duration.
func (r *LoopResult) Duration() time.Duration {
	return r.CompletedAt.Sub(r.StartedAt)
}

// LastIteration returns the most recent iteration result, or nil if none.
func (r *LoopResult) LastIteration() *IterationResult {
	if len(r.Iterations) == 0 {
		return nil
	}
	return &r.Iterations[len(r.Iterations)-1]
}

// FinalValidationStatus returns the validation status from the last iteration.
func (r *LoopResult) FinalValidationStatus() CheckStatus {
	if last := r.LastIteration(); last != nil {
		return last.ValidationStatus
	}
	return ""
}

// AllActionsTaken returns all actions taken across all iterations.
func (r *LoopResult) AllActionsTaken() []string {
	var actions []string
	for _, iter := range r.Iterations {
		actions = append(actions, iter.ActionsTaken...)
	}
	return actions
}

// NewLoopResult creates a new LoopResult for the given loop.
func NewLoopResult(loop *Loop) *LoopResult {
	return &LoopResult{
		LoopName:   loop.Name,
		LoopType:   loop.Type,
		Status:     LoopStatusInProgress,
		Iterations: []IterationResult{},
		StartedAt:  time.Now(),
	}
}

// AddIteration adds an iteration result and updates the total attempts.
func (r *LoopResult) AddIteration(iter IterationResult) {
	r.Iterations = append(r.Iterations, iter)
	r.TotalAttempts = len(r.Iterations)
}

// Complete marks the loop as complete with the given status.
func (r *LoopResult) Complete(status LoopStatus) {
	r.Status = status
	r.CompletedAt = time.Now()
}

// CompleteWithError marks the loop as aborted with an error.
func (r *LoopResult) CompleteWithError(err error) {
	r.Status = LoopStatusAborted
	r.Error = err.Error()
	r.CompletedAt = time.Now()
}

// Escalate marks the loop as escalated with a reason.
func (r *LoopResult) Escalate(reason string) {
	r.Status = LoopStatusEscalated
	r.EscalationReason = reason
	r.CompletedAt = time.Now()
}

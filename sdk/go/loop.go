package multiagentspec

import "fmt"

// LoopCategory represents the two loop paradigms.
type LoopCategory string

const (
	// CategoryMissionDriven loops work toward a goal described in a mission statement.
	CategoryMissionDriven LoopCategory = "mission-driven"

	// CategoryStateDriven loops converge toward a valid state through validation.
	CategoryStateDriven LoopCategory = "state-driven"
)

// String returns the string representation of the category.
func (c LoopCategory) String() string {
	return string(c)
}

// LoopType represents the loop execution pattern.
type LoopType string

const (
	// LoopREAL is a mission-driven loop: Read requirements, Eval situation, Act toward goal, Loop.
	// Use for open-ended tasks like "implement feature X" or "migrate to new API".
	LoopREAL LoopType = "REAL"

	// LoopVEAL is a state-driven validation loop: Validate state, Eval findings, Act to correct, Loop.
	// Use for convergent tasks like "code must pass lint" or "docs must be complete".
	LoopVEAL LoopType = "VEAL"
)

// Category returns the loop category for this type.
func (t LoopType) Category() LoopCategory {
	switch t {
	case LoopREAL:
		return CategoryMissionDriven
	case LoopVEAL:
		return CategoryStateDriven
	default:
		return ""
	}
}

// IsMissionDriven returns true if this is a mission-driven loop type.
func (t LoopType) IsMissionDriven() bool {
	return t.Category() == CategoryMissionDriven
}

// IsStateDriven returns true if this is a state-driven loop type.
func (t LoopType) IsStateDriven() bool {
	return t.Category() == CategoryStateDriven
}

// EscalationPolicy defines what to do when max attempts are reached.
type EscalationPolicy string

const (
	// EscalationHuman stops the loop and requests human intervention.
	EscalationHuman EscalationPolicy = "human"

	// EscalationAbort stops the loop and fails the workflow.
	EscalationAbort EscalationPolicy = "abort"

	// EscalationContinue proceeds despite unresolved issues.
	EscalationContinue EscalationPolicy = "continue"

	// EscalationFallback invokes a fallback agent for manual resolution.
	EscalationFallback EscalationPolicy = "fallback"
)

// CheckType represents how a loop check is executed.
type CheckType string

const (
	// CheckTypeCommand executes a shell command and checks exit code.
	CheckTypeCommand CheckType = "command"

	// CheckTypePattern searches for a regex pattern in files.
	CheckTypePattern CheckType = "pattern"

	// CheckTypeFile checks if a file exists and optionally validates content.
	CheckTypeFile CheckType = "file"

	// CheckTypeManual requires human or agent judgment.
	CheckTypeManual CheckType = "manual"
)

// LoopCheck represents a validation check within a loop.
type LoopCheck struct {
	// ID is the unique check identifier within the loop.
	ID string `json:"id" yaml:"id"`

	// Description is a human-readable description of what this check validates.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Type is how the check is executed.
	Type CheckType `json:"type,omitempty" yaml:"type,omitempty"`

	// Command is the shell command to execute (for type: command).
	Command string `json:"command,omitempty" yaml:"command,omitempty"`

	// Pattern is the regex pattern to search for (for type: pattern).
	Pattern string `json:"pattern,omitempty" yaml:"pattern,omitempty"`

	// Files is a glob pattern for files to check (for type: pattern).
	Files string `json:"files,omitempty" yaml:"files,omitempty"`

	// File is the file path to check (for type: file).
	File string `json:"file,omitempty" yaml:"file,omitempty"`

	// Required indicates if check failure causes NO-GO status.
	Required *bool `json:"required,omitempty" yaml:"required,omitempty"`

	// Expected describes the expected outcome for success.
	Expected string `json:"expected,omitempty" yaml:"expected,omitempty"`
}

// IsRequired returns true if this check is required (defaults to true).
func (c *LoopCheck) IsRequired() bool {
	if c.Required == nil {
		return true
	}
	return *c.Required
}

// Loop represents an autonomous loop definition.
type Loop struct {
	// Name is the unique identifier for the loop.
	Name string `json:"name" yaml:"name"`

	// Description is a human-readable description of the loop's purpose.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Type is the loop execution pattern (REAL or VEAL).
	Type LoopType `json:"type" yaml:"type"`

	// Validator is the agent that validates state (read-only).
	// Required for VEAL loops, optional for REAL loops.
	Validator string `json:"validator,omitempty" yaml:"validator,omitempty"`

	// Actor is the agent that acts to correct state or achieve the mission.
	Actor string `json:"actor" yaml:"actor"`

	// Mission is the goal statement for REAL loops.
	// Describes what the loop should achieve.
	Mission string `json:"mission,omitempty" yaml:"mission,omitempty"`

	// Checks are validation checks to run each iteration.
	// Primarily used for VEAL loops to define the expected state.
	Checks []LoopCheck `json:"checks,omitempty" yaml:"checks,omitempty"`

	// MaxAttempts is the maximum loop iterations before escalation.
	// Defaults to 3.
	MaxAttempts int `json:"max_attempts,omitempty" yaml:"max_attempts,omitempty"`

	// Escalation is the policy when max attempts are reached.
	// Defaults to "human".
	Escalation EscalationPolicy `json:"escalation,omitempty" yaml:"escalation,omitempty"`

	// FallbackAgent is the agent to invoke if escalation is "fallback".
	FallbackAgent string `json:"fallback_agent,omitempty" yaml:"fallback_agent,omitempty"`

	// SuccessCriteria describes what constitutes loop success.
	SuccessCriteria string `json:"success_criteria,omitempty" yaml:"success_criteria,omitempty"`

	// Inputs are data inputs required to start the loop.
	Inputs []Port `json:"inputs,omitempty" yaml:"inputs,omitempty"`

	// Outputs are data outputs produced when the loop completes.
	Outputs []Port `json:"outputs,omitempty" yaml:"outputs,omitempty"`
}

// NewLoop creates a new Loop with the given name and type.
func NewLoop(name string, loopType LoopType) *Loop {
	return &Loop{
		Name:        name,
		Type:        loopType,
		MaxAttempts: 3,
		Escalation:  EscalationHuman,
	}
}

// NewVEALLoop creates a new VEAL loop with validator and actor agents.
func NewVEALLoop(name, validator, actor string) *Loop {
	return &Loop{
		Name:        name,
		Type:        LoopVEAL,
		Validator:   validator,
		Actor:       actor,
		MaxAttempts: 3,
		Escalation:  EscalationHuman,
	}
}

// NewREALLoop creates a new REAL loop with actor agent and mission.
func NewREALLoop(name, actor, mission string) *Loop {
	return &Loop{
		Name:        name,
		Type:        LoopREAL,
		Actor:       actor,
		Mission:     mission,
		MaxAttempts: 3,
		Escalation:  EscalationHuman,
	}
}

// WithDescription sets the loop's description and returns the loop for chaining.
func (l *Loop) WithDescription(description string) *Loop {
	l.Description = description
	return l
}

// WithValidator sets the validator agent and returns the loop for chaining.
func (l *Loop) WithValidator(validator string) *Loop {
	l.Validator = validator
	return l
}

// WithActor sets the actor agent and returns the loop for chaining.
func (l *Loop) WithActor(actor string) *Loop {
	l.Actor = actor
	return l
}

// WithMission sets the mission statement and returns the loop for chaining.
func (l *Loop) WithMission(mission string) *Loop {
	l.Mission = mission
	return l
}

// WithChecks sets the validation checks and returns the loop for chaining.
func (l *Loop) WithChecks(checks ...LoopCheck) *Loop {
	l.Checks = checks
	return l
}

// AddCheck adds a validation check and returns the loop for chaining.
func (l *Loop) AddCheck(check LoopCheck) *Loop {
	l.Checks = append(l.Checks, check)
	return l
}

// WithMaxAttempts sets the max attempts and returns the loop for chaining.
func (l *Loop) WithMaxAttempts(maxAttempts int) *Loop {
	l.MaxAttempts = maxAttempts
	return l
}

// WithEscalation sets the escalation policy and returns the loop for chaining.
func (l *Loop) WithEscalation(escalation EscalationPolicy) *Loop {
	l.Escalation = escalation
	return l
}

// WithFallbackAgent sets the fallback agent and returns the loop for chaining.
func (l *Loop) WithFallbackAgent(agent string) *Loop {
	l.FallbackAgent = agent
	return l
}

// WithSuccessCriteria sets the success criteria and returns the loop for chaining.
func (l *Loop) WithSuccessCriteria(criteria string) *Loop {
	l.SuccessCriteria = criteria
	return l
}

// EffectiveMaxAttempts returns the max attempts, defaulting to 3.
func (l *Loop) EffectiveMaxAttempts() int {
	if l.MaxAttempts <= 0 {
		return 3
	}
	return l.MaxAttempts
}

// EffectiveEscalation returns the escalation policy, defaulting to "human".
func (l *Loop) EffectiveEscalation() EscalationPolicy {
	if l.Escalation == "" {
		return EscalationHuman
	}
	return l.Escalation
}

// Category returns the loop category.
func (l *Loop) Category() LoopCategory {
	return l.Type.Category()
}

// IsMissionDriven returns true if this is a mission-driven (REAL) loop.
func (l *Loop) IsMissionDriven() bool {
	return l.Type.IsMissionDriven()
}

// IsStateDriven returns true if this is a state-driven (VEAL) loop.
func (l *Loop) IsStateDriven() bool {
	return l.Type.IsStateDriven()
}

// RequiredChecks returns only the checks marked as required.
func (l *Loop) RequiredChecks() []LoopCheck {
	var required []LoopCheck
	for _, check := range l.Checks {
		if check.IsRequired() {
			required = append(required, check)
		}
	}
	return required
}

// Validate checks loop configuration consistency.
// Returns an error if the configuration is invalid for the loop type.
func (l *Loop) Validate() error {
	if l.Name == "" {
		return fmt.Errorf("loop name is required")
	}

	if l.Actor == "" {
		return fmt.Errorf("loop actor is required")
	}

	switch l.Type {
	case LoopVEAL:
		// VEAL loops require a validator
		if l.Validator == "" {
			return fmt.Errorf("VEAL loop requires a validator agent")
		}
		// VEAL loops should have checks
		if len(l.Checks) == 0 {
			return fmt.Errorf("VEAL loop should have at least one check")
		}
	case LoopREAL:
		// REAL loops should have a mission
		if l.Mission == "" {
			return fmt.Errorf("REAL loop should have a mission statement")
		}
	default:
		return fmt.Errorf("invalid loop type: %s", l.Type)
	}

	// Fallback escalation requires fallback agent
	if l.Escalation == EscalationFallback && l.FallbackAgent == "" {
		return fmt.Errorf("fallback escalation requires fallback_agent")
	}

	return nil
}

// Agents returns the list of agent names involved in this loop.
func (l *Loop) Agents() []string {
	agents := []string{l.Actor}
	if l.Validator != "" && l.Validator != l.Actor {
		agents = append([]string{l.Validator}, agents...)
	}
	if l.FallbackAgent != "" {
		agents = append(agents, l.FallbackAgent)
	}
	return agents
}

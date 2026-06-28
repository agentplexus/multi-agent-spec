package multiagentspec

import (
	"encoding/json"
	"testing"
)

func TestLoopTypeConstants(t *testing.T) {
	tests := []struct {
		lt   LoopType
		want string
	}{
		{LoopREAL, "REAL"},
		{LoopVEAL, "VEAL"},
	}

	for _, tt := range tests {
		if string(tt.lt) != tt.want {
			t.Errorf("LoopType %v = %q, want %q", tt.lt, string(tt.lt), tt.want)
		}
	}
}

func TestLoopType_Category(t *testing.T) {
	tests := []struct {
		lt       LoopType
		expected LoopCategory
	}{
		{LoopREAL, CategoryMissionDriven},
		{LoopVEAL, CategoryStateDriven},
	}
	for _, tt := range tests {
		t.Run(string(tt.lt), func(t *testing.T) {
			if got := tt.lt.Category(); got != tt.expected {
				t.Errorf("Category() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLoopType_IsMissionDriven(t *testing.T) {
	if !LoopREAL.IsMissionDriven() {
		t.Error("REAL should be mission-driven")
	}
	if LoopREAL.IsStateDriven() {
		t.Error("REAL should not be state-driven")
	}

	if LoopVEAL.IsMissionDriven() {
		t.Error("VEAL should not be mission-driven")
	}
	if !LoopVEAL.IsStateDriven() {
		t.Error("VEAL should be state-driven")
	}
}

func TestLoopCategory_String(t *testing.T) {
	if CategoryMissionDriven.String() != "mission-driven" {
		t.Errorf("CategoryMissionDriven.String() = %q, want %q", CategoryMissionDriven.String(), "mission-driven")
	}
	if CategoryStateDriven.String() != "state-driven" {
		t.Errorf("CategoryStateDriven.String() = %q, want %q", CategoryStateDriven.String(), "state-driven")
	}
}

func TestNewLoop(t *testing.T) {
	loop := NewLoop("test-loop", LoopVEAL)

	if loop.Name != "test-loop" {
		t.Errorf("Name = %q, want %q", loop.Name, "test-loop")
	}
	if loop.Type != LoopVEAL {
		t.Errorf("Type = %q, want %q", loop.Type, LoopVEAL)
	}
	if loop.MaxAttempts != 3 {
		t.Errorf("MaxAttempts = %d, want 3", loop.MaxAttempts)
	}
	if loop.Escalation != EscalationHuman {
		t.Errorf("Escalation = %q, want %q", loop.Escalation, EscalationHuman)
	}
}

func TestNewVEALLoop(t *testing.T) {
	loop := NewVEALLoop("qa-fix", "qa", "code-fixer")

	if loop.Name != "qa-fix" {
		t.Errorf("Name = %q, want %q", loop.Name, "qa-fix")
	}
	if loop.Type != LoopVEAL {
		t.Errorf("Type = %q, want %q", loop.Type, LoopVEAL)
	}
	if loop.Validator != "qa" {
		t.Errorf("Validator = %q, want %q", loop.Validator, "qa")
	}
	if loop.Actor != "code-fixer" {
		t.Errorf("Actor = %q, want %q", loop.Actor, "code-fixer")
	}
}

func TestNewREALLoop(t *testing.T) {
	loop := NewREALLoop("feature-impl", "developer", "Implement user authentication")

	if loop.Name != "feature-impl" {
		t.Errorf("Name = %q, want %q", loop.Name, "feature-impl")
	}
	if loop.Type != LoopREAL {
		t.Errorf("Type = %q, want %q", loop.Type, LoopREAL)
	}
	if loop.Actor != "developer" {
		t.Errorf("Actor = %q, want %q", loop.Actor, "developer")
	}
	if loop.Mission != "Implement user authentication" {
		t.Errorf("Mission = %q, want %q", loop.Mission, "Implement user authentication")
	}
}

func TestLoop_WithMethods(t *testing.T) {
	loop := NewLoop("test", LoopVEAL).
		WithDescription("Test loop").
		WithValidator("validator").
		WithActor("actor").
		WithMaxAttempts(5).
		WithEscalation(EscalationAbort).
		WithSuccessCriteria("All checks pass")

	if loop.Description != "Test loop" {
		t.Errorf("Description = %q, want %q", loop.Description, "Test loop")
	}
	if loop.Validator != "validator" {
		t.Errorf("Validator = %q, want %q", loop.Validator, "validator")
	}
	if loop.Actor != "actor" {
		t.Errorf("Actor = %q, want %q", loop.Actor, "actor")
	}
	if loop.MaxAttempts != 5 {
		t.Errorf("MaxAttempts = %d, want 5", loop.MaxAttempts)
	}
	if loop.Escalation != EscalationAbort {
		t.Errorf("Escalation = %q, want %q", loop.Escalation, EscalationAbort)
	}
	if loop.SuccessCriteria != "All checks pass" {
		t.Errorf("SuccessCriteria = %q, want %q", loop.SuccessCriteria, "All checks pass")
	}
}

func TestLoop_AddCheck(t *testing.T) {
	loop := NewVEALLoop("test", "validator", "actor").
		AddCheck(LoopCheck{ID: "build", Type: CheckTypeCommand, Command: "go build ./..."}).
		AddCheck(LoopCheck{ID: "lint", Type: CheckTypeCommand, Command: "golangci-lint run"})

	if len(loop.Checks) != 2 {
		t.Errorf("len(Checks) = %d, want 2", len(loop.Checks))
	}
	if loop.Checks[0].ID != "build" {
		t.Errorf("Checks[0].ID = %q, want %q", loop.Checks[0].ID, "build")
	}
}

func TestLoop_EffectiveMaxAttempts(t *testing.T) {
	tests := []struct {
		name        string
		maxAttempts int
		expected    int
	}{
		{"zero defaults to 3", 0, 3},
		{"negative defaults to 3", -1, 3},
		{"positive returns value", 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loop := &Loop{MaxAttempts: tt.maxAttempts}
			if got := loop.EffectiveMaxAttempts(); got != tt.expected {
				t.Errorf("EffectiveMaxAttempts() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestLoop_EffectiveEscalation(t *testing.T) {
	tests := []struct {
		name       string
		escalation EscalationPolicy
		expected   EscalationPolicy
	}{
		{"empty defaults to human", "", EscalationHuman},
		{"abort returns abort", EscalationAbort, EscalationAbort},
		{"continue returns continue", EscalationContinue, EscalationContinue},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loop := &Loop{Escalation: tt.escalation}
			if got := loop.EffectiveEscalation(); got != tt.expected {
				t.Errorf("EffectiveEscalation() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestLoop_Validate_VEAL(t *testing.T) {
	tests := []struct {
		name    string
		loop    *Loop
		wantErr bool
	}{
		{
			name:    "valid VEAL loop",
			loop:    NewVEALLoop("test", "validator", "actor").AddCheck(LoopCheck{ID: "check1"}),
			wantErr: false,
		},
		{
			name:    "VEAL missing validator",
			loop:    &Loop{Name: "test", Type: LoopVEAL, Actor: "actor", Checks: []LoopCheck{{ID: "check1"}}},
			wantErr: true,
		},
		{
			name:    "VEAL missing checks",
			loop:    &Loop{Name: "test", Type: LoopVEAL, Validator: "validator", Actor: "actor"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.loop.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoop_Validate_REAL(t *testing.T) {
	tests := []struct {
		name    string
		loop    *Loop
		wantErr bool
	}{
		{
			name:    "valid REAL loop",
			loop:    NewREALLoop("test", "actor", "Do something"),
			wantErr: false,
		},
		{
			name:    "REAL missing mission",
			loop:    &Loop{Name: "test", Type: LoopREAL, Actor: "actor"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.loop.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoop_Validate_Fallback(t *testing.T) {
	loop := &Loop{
		Name:       "test",
		Type:       LoopVEAL,
		Validator:  "validator",
		Actor:      "actor",
		Checks:     []LoopCheck{{ID: "check1"}},
		Escalation: EscalationFallback,
		// Missing FallbackAgent
	}

	if err := loop.Validate(); err == nil {
		t.Error("expected error for fallback without fallback_agent")
	}

	loop.FallbackAgent = "fallback"
	if err := loop.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLoop_Agents(t *testing.T) {
	loop := &Loop{
		Name:          "test",
		Validator:     "validator",
		Actor:         "actor",
		FallbackAgent: "fallback",
	}

	agents := loop.Agents()
	if len(agents) != 3 {
		t.Errorf("len(Agents()) = %d, want 3", len(agents))
	}
	if agents[0] != "validator" {
		t.Errorf("Agents()[0] = %q, want %q", agents[0], "validator")
	}
	if agents[1] != "actor" {
		t.Errorf("Agents()[1] = %q, want %q", agents[1], "actor")
	}
	if agents[2] != "fallback" {
		t.Errorf("Agents()[2] = %q, want %q", agents[2], "fallback")
	}
}

func TestLoop_RequiredChecks(t *testing.T) {
	required := true
	notRequired := false

	loop := &Loop{
		Checks: []LoopCheck{
			{ID: "c1", Required: &required},
			{ID: "c2", Required: &notRequired},
			{ID: "c3"}, // nil defaults to true
		},
	}

	reqChecks := loop.RequiredChecks()
	if len(reqChecks) != 2 {
		t.Errorf("len(RequiredChecks()) = %d, want 2", len(reqChecks))
	}
}

func TestCheck_IsRequired(t *testing.T) {
	required := true
	notRequired := false

	tests := []struct {
		name     string
		check    LoopCheck
		expected bool
	}{
		{"nil defaults to true", LoopCheck{ID: "c1"}, true},
		{"explicit true", LoopCheck{ID: "c2", Required: &required}, true},
		{"explicit false", LoopCheck{ID: "c3", Required: &notRequired}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.check.IsRequired(); got != tt.expected {
				t.Errorf("IsRequired() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLoopJSONSerialization(t *testing.T) {
	loop := NewVEALLoop("qa-fix", "qa", "code-fixer").
		WithDescription("QA fix loop").
		WithMaxAttempts(3).
		AddCheck(LoopCheck{ID: "build", Type: CheckTypeCommand, Command: "go build ./..."}).
		AddCheck(LoopCheck{ID: "lint", Type: CheckTypeCommand, Command: "golangci-lint run"})

	data, err := json.Marshal(loop)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded Loop
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Name != loop.Name {
		t.Errorf("Name = %q, want %q", decoded.Name, loop.Name)
	}
	if decoded.Type != LoopVEAL {
		t.Errorf("Type = %q, want %q", decoded.Type, LoopVEAL)
	}
	if decoded.Validator != "qa" {
		t.Errorf("Validator = %q, want %q", decoded.Validator, "qa")
	}
	if len(decoded.Checks) != 2 {
		t.Errorf("len(Checks) = %d, want 2", len(decoded.Checks))
	}
}

func TestStepIsLoopStep(t *testing.T) {
	agentStep := Step{Name: "agent-step", Agent: "my-agent"}
	loopStep := Step{Name: "loop-step", Loop: "my-loop"}

	if agentStep.IsLoopStep() {
		t.Error("agent step should not be a loop step")
	}
	if !agentStep.IsAgentStep() {
		t.Error("agent step should be an agent step")
	}

	if !loopStep.IsLoopStep() {
		t.Error("loop step should be a loop step")
	}
	if loopStep.IsAgentStep() {
		t.Error("loop step should not be an agent step")
	}
}

package multiagentspec

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestWorkflowTypeConstants(t *testing.T) {
	tests := []struct {
		wt   WorkflowType
		want string
	}{
		{WorkflowChain, "chain"},
		{WorkflowScatter, "scatter"},
		{WorkflowGraph, "graph"},
		{WorkflowCrew, "crew"},
		{WorkflowSwarm, "swarm"},
		{WorkflowCouncil, "council"},
	}

	for _, tt := range tests {
		if string(tt.wt) != tt.want {
			t.Errorf("WorkflowType %v = %q, want %q", tt.wt, string(tt.wt), tt.want)
		}
	}
}

func TestWorkflowType_Category(t *testing.T) {
	tests := []struct {
		wt       WorkflowType
		expected WorkflowCategory
	}{
		{WorkflowChain, CategoryDeterministic},
		{WorkflowScatter, CategoryDeterministic},
		{WorkflowGraph, CategoryDeterministic},
		{WorkflowCrew, CategorySelfDirected},
		{WorkflowSwarm, CategorySelfDirected},
		{WorkflowCouncil, CategorySelfDirected},
	}
	for _, tt := range tests {
		t.Run(string(tt.wt), func(t *testing.T) {
			if got := tt.wt.Category(); got != tt.expected {
				t.Errorf("Category() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestWorkflowType_IsDeterministic(t *testing.T) {
	deterministic := []WorkflowType{
		WorkflowChain, WorkflowScatter, WorkflowGraph,
	}
	selfDirected := []WorkflowType{
		WorkflowCrew, WorkflowSwarm, WorkflowCouncil,
	}

	for _, wt := range deterministic {
		if !wt.IsDeterministic() {
			t.Errorf("%s should be deterministic", wt)
		}
		if wt.IsSelfDirected() {
			t.Errorf("%s should not be self-directed", wt)
		}
	}
	for _, wt := range selfDirected {
		if wt.IsDeterministic() {
			t.Errorf("%s should not be deterministic", wt)
		}
		if !wt.IsSelfDirected() {
			t.Errorf("%s should be self-directed", wt)
		}
	}
}

func TestWorkflowCategory_String(t *testing.T) {
	if CategoryDeterministic.String() != "deterministic" {
		t.Errorf("CategoryDeterministic.String() = %q, want %q", CategoryDeterministic.String(), "deterministic")
	}
	if CategorySelfDirected.String() != "self-directed" {
		t.Errorf("CategorySelfDirected.String() = %q, want %q", CategorySelfDirected.String(), "self-directed")
	}
}

func TestTeam_WorkflowCategory(t *testing.T) {
	tests := []struct {
		name     string
		team     *Team
		expected WorkflowCategory
	}{
		{
			name:     "nil workflow",
			team:     &Team{Name: "test", Version: "1.0.0"},
			expected: "",
		},
		{
			name:     "chain workflow",
			team:     &Team{Name: "test", Version: "1.0.0", Workflow: &Workflow{Type: WorkflowChain}},
			expected: CategoryDeterministic,
		},
		{
			name:     "crew workflow",
			team:     &Team{Name: "test", Version: "1.0.0", Workflow: &Workflow{Type: WorkflowCrew}},
			expected: CategorySelfDirected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.team.WorkflowCategory(); got != tt.expected {
				t.Errorf("WorkflowCategory() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTeam_IsDeterministic(t *testing.T) {
	deterministicTeam := &Team{
		Name:     "test",
		Version:  "1.0.0",
		Workflow: &Workflow{Type: WorkflowGraph},
	}
	selfDirectedTeam := &Team{
		Name:     "test",
		Version:  "1.0.0",
		Workflow: &Workflow{Type: WorkflowSwarm},
	}
	noWorkflowTeam := &Team{
		Name:    "test",
		Version: "1.0.0",
	}

	if !deterministicTeam.IsDeterministic() {
		t.Error("graph workflow team should be deterministic")
	}
	if deterministicTeam.IsSelfDirected() {
		t.Error("graph workflow team should not be self-directed")
	}

	if selfDirectedTeam.IsDeterministic() {
		t.Error("swarm workflow team should not be deterministic")
	}
	if !selfDirectedTeam.IsSelfDirected() {
		t.Error("swarm workflow team should be self-directed")
	}

	if noWorkflowTeam.IsDeterministic() {
		t.Error("team without workflow should not be deterministic")
	}
	if noWorkflowTeam.IsSelfDirected() {
		t.Error("team without workflow should not be self-directed")
	}
}

func TestStepSerialization(t *testing.T) {
	step := Step{
		Name:      "research",
		Agent:     "researcher",
		DependsOn: []string{"init"},
		Inputs: []Port{
			{Name: "topic", Type: PortTypeString, From: "init.topic"},
		},
		Outputs: []Port{
			{Name: "results", Type: PortTypeObject},
		},
	}

	data, err := json.Marshal(step)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded Step
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Name != step.Name {
		t.Errorf("Name = %q, want %q", decoded.Name, step.Name)
	}
	if decoded.Agent != step.Agent {
		t.Errorf("Agent = %q, want %q", decoded.Agent, step.Agent)
	}
	if len(decoded.DependsOn) != 1 {
		t.Errorf("len(DependsOn) = %d, want 1", len(decoded.DependsOn))
	}
	if len(decoded.Inputs) != 1 || decoded.Inputs[0].Name != "topic" {
		t.Errorf("Inputs[0].Name = %q, want %q", decoded.Inputs[0].Name, "topic")
	}
	if decoded.Inputs[0].From != "init.topic" {
		t.Errorf("Inputs[0].From = %q, want %q", decoded.Inputs[0].From, "init.topic")
	}
}

func TestStepOmitEmpty(t *testing.T) {
	step := Step{
		Name:  "minimal",
		Agent: "agent1",
	}

	data, err := json.Marshal(step)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if _, ok := m["depends_on"]; ok {
		t.Error("depends_on should be omitted when nil")
	}
	if _, ok := m["inputs"]; ok {
		t.Error("inputs should be omitted when nil")
	}
	if _, ok := m["outputs"]; ok {
		t.Error("outputs should be omitted when nil")
	}
}

func TestWorkflowSerialization(t *testing.T) {
	workflow := Workflow{
		Type: WorkflowGraph,
		Steps: []Step{
			{Name: "s1", Agent: "a1"},
			{Name: "s2", Agent: "a2", DependsOn: []string{"s1"}},
		},
	}

	data, err := json.Marshal(workflow)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded Workflow
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Type != WorkflowGraph {
		t.Errorf("Type = %q, want %q", decoded.Type, WorkflowGraph)
	}
	if len(decoded.Steps) != 2 {
		t.Errorf("len(Steps) = %d, want 2", len(decoded.Steps))
	}
}

func TestNewTeam(t *testing.T) {
	team := NewTeam("test-team", "1.0.0")

	if team.Name != "test-team" {
		t.Errorf("Name = %q, want %q", team.Name, "test-team")
	}
	if team.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", team.Version, "1.0.0")
	}
	if len(team.Agents) != 0 {
		t.Errorf("len(Agents) = %d, want 0", len(team.Agents))
	}
}

func TestTeamWithAgents(t *testing.T) {
	team := NewTeam("test", "1.0.0").WithAgents("agent1", "agent2", "agent3")

	if len(team.Agents) != 3 {
		t.Errorf("len(Agents) = %d, want 3", len(team.Agents))
	}
	if team.Agents[0] != "agent1" {
		t.Errorf("Agents[0] = %q, want %q", team.Agents[0], "agent1")
	}
}

func TestTeamWithOrchestrator(t *testing.T) {
	team := NewTeam("test", "1.0.0").
		WithAgents("agent1", "agent2").
		WithOrchestrator("agent1")

	if team.Orchestrator != "agent1" {
		t.Errorf("Orchestrator = %q, want %q", team.Orchestrator, "agent1")
	}
}

func TestTeamWithWorkflow(t *testing.T) {
	workflow := &Workflow{Type: WorkflowCrew}
	team := NewTeam("test", "1.0.0").WithWorkflow(workflow)

	if team.Workflow == nil {
		t.Error("Workflow should not be nil")
	}
	if team.Workflow.Type != WorkflowCrew {
		t.Errorf("Workflow.Type = %q, want %q", team.Workflow.Type, WorkflowCrew)
	}
}

func TestTeamChaining(t *testing.T) {
	team := NewTeam("chained", "2.0.0").
		WithAgents("a1", "a2").
		WithOrchestrator("a1").
		WithWorkflow(&Workflow{Type: WorkflowChain})

	if team.Name != "chained" {
		t.Errorf("Name = %q, want %q", team.Name, "chained")
	}
	if team.Version != "2.0.0" {
		t.Errorf("Version = %q, want %q", team.Version, "2.0.0")
	}
	if len(team.Agents) != 2 {
		t.Errorf("len(Agents) = %d, want 2", len(team.Agents))
	}
	if team.Orchestrator != "a1" {
		t.Errorf("Orchestrator = %q, want %q", team.Orchestrator, "a1")
	}
}

func TestTeamJSONSerialization(t *testing.T) {
	team := &Team{
		Name:         "json-team",
		Version:      "1.0.0",
		Description:  "JSON test team",
		Agents:       []string{"a1", "a2"},
		Orchestrator: "a1",
		Context:      "Shared context",
	}

	data, err := json.Marshal(team)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded Team
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Name != team.Name {
		t.Errorf("Name = %q, want %q", decoded.Name, team.Name)
	}
	if decoded.Version != team.Version {
		t.Errorf("Version = %q, want %q", decoded.Version, team.Version)
	}
	if len(decoded.Agents) != 2 {
		t.Errorf("len(Agents) = %d, want 2", len(decoded.Agents))
	}
}

func TestTeam_ValidateAgentReferences_AllResolve(t *testing.T) {
	team := NewTeam("test", "1.0.0").
		WithAgents("qa", "release").
		WithOrchestrator("release").
		WithWorkflow(&Workflow{
			Type: WorkflowChain,
			Steps: []Step{
				{Name: "qa-validation", Agent: "qa"},
				{Name: "ship", Agent: "release"},
			},
		})

	if err := team.ValidateAgentReferences([]string{"qa", "release"}); err != nil {
		t.Errorf("ValidateAgentReferences() = %v, want nil", err)
	}
}

func TestTeam_ValidateAgentReferences_MissingRosterAgent(t *testing.T) {
	team := NewTeam("test", "1.0.0").WithAgents("qa", "documentation")

	err := team.ValidateAgentReferences([]string{"qa"})
	if err == nil {
		t.Fatal("ValidateAgentReferences() = nil, want error for missing roster agent")
	}
	if !strings.Contains(err.Error(), `roster references agent "documentation"`) {
		t.Errorf("error = %q, want it to mention the missing roster agent", err.Error())
	}
}

func TestTeam_ValidateAgentReferences_MissingOrchestrator(t *testing.T) {
	team := NewTeam("test", "1.0.0").
		WithAgents("qa").
		WithOrchestrator("release-coordinator")

	err := team.ValidateAgentReferences([]string{"qa"})
	if err == nil {
		t.Fatal("ValidateAgentReferences() = nil, want error for missing orchestrator")
	}
	if !strings.Contains(err.Error(), `orchestrator references agent "release-coordinator"`) {
		t.Errorf("error = %q, want it to mention the missing orchestrator", err.Error())
	}
}

func TestTeam_ValidateAgentReferences_MissingWorkflowStepAgent(t *testing.T) {
	team := NewTeam("test", "1.0.0").
		WithAgents("qa").
		WithWorkflow(&Workflow{
			Type: WorkflowChain,
			Steps: []Step{
				{Name: "docs-validation", Agent: "documentation"},
			},
		})

	err := team.ValidateAgentReferences([]string{"qa"})
	if err == nil {
		t.Fatal("ValidateAgentReferences() = nil, want error for missing workflow step agent")
	}
	if !strings.Contains(err.Error(), `workflow step "docs-validation" references agent "documentation"`) {
		t.Errorf("error = %q, want it to mention the missing step agent", err.Error())
	}
}

func TestTeam_ValidateAgentReferences_LoopStepsAreSkipped(t *testing.T) {
	// A step assigned to a Loop (not an Agent) must not be flagged, even
	// though its Agent field is empty.
	team := NewTeam("test", "1.0.0").
		WithAgents("qa").
		WithWorkflow(&Workflow{
			Type: WorkflowChain,
			Steps: []Step{
				{Name: "docs-fix", Loop: "docs-fix"},
			},
		})

	if err := team.ValidateAgentReferences([]string{"qa"}); err != nil {
		t.Errorf("ValidateAgentReferences() = %v, want nil (loop steps have no agent to validate)", err)
	}
}

func TestTeam_ValidateAgentReferences_DuplicateReferencesReportedOnce(t *testing.T) {
	// "documentation" appears in both the roster and a workflow step; it
	// should be reported once per distinct reference kind, not silently
	// deduplicated across kinds (roster vs. step are different findings).
	team := NewTeam("test", "1.0.0").
		WithAgents("documentation").
		WithWorkflow(&Workflow{
			Type: WorkflowChain,
			Steps: []Step{
				{Name: "docs-validation", Agent: "documentation"},
				{Name: "docs-validation-again", Agent: "documentation"},
			},
		})

	err := team.ValidateAgentReferences(nil)
	if err == nil {
		t.Fatal("ValidateAgentReferences() = nil, want error")
	}
	joined, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("error does not support errors.Join Unwrap() []error: %T", err)
	}
	// roster + step "docs-validation" + step "docs-validation-again" = 3;
	// the second step reusing "documentation" is still a distinct step
	// name, so it is not deduplicated with the first.
	if got := len(joined.Unwrap()); got != 3 {
		t.Errorf("len(joined errors) = %d, want 3, got: %v", got, err)
	}
}

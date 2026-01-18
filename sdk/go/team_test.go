package multiagentspec

import (
	"encoding/json"
	"testing"
)

func TestWorkflowTypeConstants(t *testing.T) {
	tests := []struct {
		wt   WorkflowType
		want string
	}{
		{WorkflowSequential, "sequential"},
		{WorkflowParallel, "parallel"},
		{WorkflowDAG, "dag"},
		{WorkflowOrchestrated, "orchestrated"},
	}

	for _, tt := range tests {
		if string(tt.wt) != tt.want {
			t.Errorf("WorkflowType %v = %q, want %q", tt.wt, string(tt.wt), tt.want)
		}
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
		Type: WorkflowDAG,
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

	if decoded.Type != WorkflowDAG {
		t.Errorf("Type = %q, want %q", decoded.Type, WorkflowDAG)
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
	workflow := &Workflow{Type: WorkflowOrchestrated}
	team := NewTeam("test", "1.0.0").WithWorkflow(workflow)

	if team.Workflow == nil {
		t.Error("Workflow should not be nil")
	}
	if team.Workflow.Type != WorkflowOrchestrated {
		t.Errorf("Workflow.Type = %q, want %q", team.Workflow.Type, WorkflowOrchestrated)
	}
}

func TestTeamChaining(t *testing.T) {
	team := NewTeam("chained", "2.0.0").
		WithAgents("a1", "a2").
		WithOrchestrator("a1").
		WithWorkflow(&Workflow{Type: WorkflowSequential})

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

package multiagentspec

import (
	"encoding/json"
	"testing"
)

func TestModelConstants(t *testing.T) {
	tests := []struct {
		model Model
		want  string
	}{
		{ModelHaiku, "haiku"},
		{ModelSonnet, "sonnet"},
		{ModelOpus, "opus"},
	}

	for _, tt := range tests {
		if string(tt.model) != tt.want {
			t.Errorf("Model %v = %q, want %q", tt.model, string(tt.model), tt.want)
		}
	}
}

func TestToolConstants(t *testing.T) {
	tests := []struct {
		tool Tool
		want string
	}{
		{ToolWebSearch, "WebSearch"},
		{ToolWebFetch, "WebFetch"},
		{ToolRead, "Read"},
		{ToolWrite, "Write"},
		{ToolGlob, "Glob"},
		{ToolGrep, "Grep"},
		{ToolBash, "Bash"},
		{ToolEdit, "Edit"},
		{ToolTask, "Task"},
	}

	for _, tt := range tests {
		if string(tt.tool) != tt.want {
			t.Errorf("Tool %v = %q, want %q", tt.tool, string(tt.tool), tt.want)
		}
	}
}

func TestNewAgent(t *testing.T) {
	agent := NewAgent("test-agent", "A test agent")

	if agent.Name != "test-agent" {
		t.Errorf("Name = %q, want %q", agent.Name, "test-agent")
	}
	if agent.Description != "A test agent" {
		t.Errorf("Description = %q, want %q", agent.Description, "A test agent")
	}
	if agent.Model != ModelSonnet {
		t.Errorf("Model = %q, want %q", agent.Model, ModelSonnet)
	}
}

func TestAgentWithModel(t *testing.T) {
	agent := NewAgent("test", "Test").WithModel(ModelHaiku)

	if agent.Model != ModelHaiku {
		t.Errorf("Model = %q, want %q", agent.Model, ModelHaiku)
	}
}

func TestAgentWithTools(t *testing.T) {
	agent := NewAgent("test", "Test").WithTools("Read", "Write", "Bash")

	if len(agent.Tools) != 3 {
		t.Errorf("len(Tools) = %d, want 3", len(agent.Tools))
	}
	if agent.Tools[0] != "Read" {
		t.Errorf("Tools[0] = %q, want %q", agent.Tools[0], "Read")
	}
}

func TestAgentWithInstructions(t *testing.T) {
	instructions := "You are a helpful agent."
	agent := NewAgent("test", "Test").WithInstructions(instructions)

	if agent.Instructions != instructions {
		t.Errorf("Instructions = %q, want %q", agent.Instructions, instructions)
	}
}

func TestAgentChaining(t *testing.T) {
	agent := NewAgent("chained", "Chained agent").
		WithModel(ModelOpus).
		WithTools("Read", "Write").
		WithInstructions("Be helpful.")

	if agent.Name != "chained" {
		t.Errorf("Name = %q, want %q", agent.Name, "chained")
	}
	if agent.Model != ModelOpus {
		t.Errorf("Model = %q, want %q", agent.Model, ModelOpus)
	}
	if len(agent.Tools) != 2 {
		t.Errorf("len(Tools) = %d, want 2", len(agent.Tools))
	}
	if agent.Instructions != "Be helpful." {
		t.Errorf("Instructions = %q, want %q", agent.Instructions, "Be helpful.")
	}
}

func TestAgentJSONSerialization(t *testing.T) {
	agent := &Agent{
		Name:         "json-test",
		Description:  "JSON test agent",
		Model:        ModelSonnet,
		Tools:        []string{"Read", "Write"},
		Skills:       []string{"skill1"},
		Dependencies: []string{"other-agent"},
		Instructions: "Test instructions.",
	}

	data, err := json.Marshal(agent)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded Agent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Name != agent.Name {
		t.Errorf("Name = %q, want %q", decoded.Name, agent.Name)
	}
	if decoded.Model != agent.Model {
		t.Errorf("Model = %q, want %q", decoded.Model, agent.Model)
	}
	if len(decoded.Tools) != len(agent.Tools) {
		t.Errorf("len(Tools) = %d, want %d", len(decoded.Tools), len(agent.Tools))
	}
}

func TestAgentJSONOmitEmpty(t *testing.T) {
	agent := &Agent{
		Name:        "minimal",
		Description: "Minimal agent",
	}

	data, err := json.Marshal(agent)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	// Check that omitempty fields are not present
	if _, ok := m["instructions"]; ok {
		t.Error("instructions should be omitted when empty")
	}
	if _, ok := m["model"]; ok {
		t.Error("model should be omitted when empty")
	}
}

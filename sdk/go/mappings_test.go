package multiagentspec

import (
	"strings"
	"testing"
)

func TestClaudeCodeModels(t *testing.T) {
	tests := []struct {
		model Model
		want  string
	}{
		{ModelHaiku, "haiku"},
		{ModelSonnet, "sonnet"},
		{ModelOpus, "opus"},
	}

	for _, tt := range tests {
		got, ok := ClaudeCodeModels[tt.model]
		if !ok {
			t.Errorf("ClaudeCodeModels[%q] not found", tt.model)
			continue
		}
		if got != tt.want {
			t.Errorf("ClaudeCodeModels[%q] = %q, want %q", tt.model, got, tt.want)
		}
	}
}

func TestKiroCLIModels(t *testing.T) {
	tests := []struct {
		model Model
		want  string
	}{
		{ModelHaiku, "claude-haiku-35"},
		{ModelSonnet, "claude-sonnet-4"},
		{ModelOpus, "claude-opus-4"},
	}

	for _, tt := range tests {
		got, ok := KiroCLIModels[tt.model]
		if !ok {
			t.Errorf("KiroCLIModels[%q] not found", tt.model)
			continue
		}
		if got != tt.want {
			t.Errorf("KiroCLIModels[%q] = %q, want %q", tt.model, got, tt.want)
		}
	}
}

func TestBedrockModels(t *testing.T) {
	tests := []struct {
		model    Model
		contains string
	}{
		{ModelHaiku, "haiku"},
		{ModelSonnet, "sonnet"},
		{ModelOpus, "opus"},
	}

	for _, tt := range tests {
		got, ok := BedrockModels[tt.model]
		if !ok {
			t.Errorf("BedrockModels[%q] not found", tt.model)
			continue
		}
		if !strings.Contains(got, tt.contains) {
			t.Errorf("BedrockModels[%q] = %q, should contain %q", tt.model, got, tt.contains)
		}
	}
}

func TestKiroCLITools(t *testing.T) {
	tests := []struct {
		tool Tool
		want string
	}{
		{ToolWebSearch, "web_search"},
		{ToolWebFetch, "web_fetch"},
		{ToolRead, "read"},
		{ToolWrite, "write"},
		{ToolGlob, "glob"},
		{ToolGrep, "grep"},
		{ToolBash, "bash"},
		{ToolEdit, "edit"},
		{ToolTask, "task"},
	}

	for _, tt := range tests {
		got, ok := KiroCLITools[tt.tool]
		if !ok {
			t.Errorf("KiroCLITools[%q] not found", tt.tool)
			continue
		}
		if got != tt.want {
			t.Errorf("KiroCLITools[%q] = %q, want %q", tt.tool, got, tt.want)
		}
	}
}

func TestAgentKitTools(t *testing.T) {
	tests := []struct {
		tool Tool
		want string
	}{
		{ToolWebSearch, "shell"},
		{ToolWebFetch, "shell"},
		{ToolRead, "read"},
		{ToolWrite, "write"},
		{ToolGlob, "glob"},
		{ToolGrep, "grep"},
		{ToolBash, "shell"},
		{ToolEdit, "write"},
		{ToolTask, "shell"},
	}

	for _, tt := range tests {
		got, ok := AgentKitTools[tt.tool]
		if !ok {
			t.Errorf("AgentKitTools[%q] not found", tt.tool)
			continue
		}
		if got != tt.want {
			t.Errorf("AgentKitTools[%q] = %q, want %q", tt.tool, got, tt.want)
		}
	}
}

func TestMapModelToClaudeCode(t *testing.T) {
	tests := []struct {
		model Model
		want  string
	}{
		{ModelHaiku, "haiku"},
		{ModelSonnet, "sonnet"},
		{ModelOpus, "opus"},
		{Model("unknown"), "unknown"}, // Fallback case
	}

	for _, tt := range tests {
		got := MapModelToClaudeCode(tt.model)
		if got != tt.want {
			t.Errorf("MapModelToClaudeCode(%q) = %q, want %q", tt.model, got, tt.want)
		}
	}
}

func TestMapModelToKiroCLI(t *testing.T) {
	tests := []struct {
		model Model
		want  string
	}{
		{ModelHaiku, "claude-haiku-35"},
		{ModelSonnet, "claude-sonnet-4"},
		{ModelOpus, "claude-opus-4"},
		{Model("unknown"), "unknown"}, // Fallback case
	}

	for _, tt := range tests {
		got := MapModelToKiroCLI(tt.model)
		if got != tt.want {
			t.Errorf("MapModelToKiroCLI(%q) = %q, want %q", tt.model, got, tt.want)
		}
	}
}

func TestMapModelToBedrock(t *testing.T) {
	tests := []struct {
		model    Model
		contains string
	}{
		{ModelHaiku, "haiku"},
		{ModelSonnet, "sonnet"},
		{ModelOpus, "opus"},
	}

	for _, tt := range tests {
		got := MapModelToBedrock(tt.model)
		if !strings.Contains(got, tt.contains) {
			t.Errorf("MapModelToBedrock(%q) = %q, should contain %q", tt.model, got, tt.contains)
		}
	}

	// Test fallback
	unknown := MapModelToBedrock(Model("unknown"))
	if unknown != "unknown" {
		t.Errorf("MapModelToBedrock(unknown) = %q, want %q", unknown, "unknown")
	}
}

func TestMapToolToKiroCLI(t *testing.T) {
	tests := []struct {
		tool Tool
		want string
	}{
		{ToolWebSearch, "web_search"},
		{ToolRead, "read"},
		{ToolBash, "bash"},
		{Tool("unknown"), "unknown"}, // Fallback case
	}

	for _, tt := range tests {
		got := MapToolToKiroCLI(tt.tool)
		if got != tt.want {
			t.Errorf("MapToolToKiroCLI(%q) = %q, want %q", tt.tool, got, tt.want)
		}
	}
}

func TestMapToolToAgentKit(t *testing.T) {
	tests := []struct {
		tool Tool
		want string
	}{
		{ToolWebSearch, "shell"},
		{ToolRead, "read"},
		{ToolBash, "shell"},
		{Tool("unknown"), "unknown"}, // Fallback case
	}

	for _, tt := range tests {
		got := MapToolToAgentKit(tt.tool)
		if got != tt.want {
			t.Errorf("MapToolToAgentKit(%q) = %q, want %q", tt.tool, got, tt.want)
		}
	}
}

// Test completeness of mappings
func TestMappingCompleteness(t *testing.T) {
	models := []Model{ModelHaiku, ModelSonnet, ModelOpus}
	tools := []Tool{ToolWebSearch, ToolWebFetch, ToolRead, ToolWrite, ToolGlob, ToolGrep, ToolBash, ToolEdit, ToolTask}

	// Check ClaudeCodeModels
	for _, m := range models {
		if _, ok := ClaudeCodeModels[m]; !ok {
			t.Errorf("ClaudeCodeModels missing %q", m)
		}
	}

	// Check KiroCLIModels
	for _, m := range models {
		if _, ok := KiroCLIModels[m]; !ok {
			t.Errorf("KiroCLIModels missing %q", m)
		}
	}

	// Check BedrockModels
	for _, m := range models {
		if _, ok := BedrockModels[m]; !ok {
			t.Errorf("BedrockModels missing %q", m)
		}
	}

	// Check KiroCLITools
	for _, tool := range tools {
		if _, ok := KiroCLITools[tool]; !ok {
			t.Errorf("KiroCLITools missing %q", tool)
		}
	}

	// Check AgentKitTools
	for _, tool := range tools {
		if _, ok := AgentKitTools[tool]; !ok {
			t.Errorf("AgentKitTools missing %q", tool)
		}
	}
}

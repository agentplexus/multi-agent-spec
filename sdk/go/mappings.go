package multiagentspec

// ClaudeCodeModels maps canonical model names to Claude Code identifiers.
var ClaudeCodeModels = map[Model]string{
	ModelHaiku:  "haiku",
	ModelSonnet: "sonnet",
	ModelOpus:   "opus",
}

// KiroCLIModels maps canonical model names to Kiro CLI identifiers.
var KiroCLIModels = map[Model]string{
	ModelHaiku:  "claude-haiku-35",
	ModelSonnet: "claude-sonnet-4",
	ModelOpus:   "claude-opus-4",
}

// BedrockModels maps canonical model names to AWS Bedrock identifiers.
var BedrockModels = map[Model]string{
	ModelHaiku:  "anthropic.claude-3-haiku-20240307-v1:0",
	ModelSonnet: "anthropic.claude-3-5-sonnet-20241022-v2:0",
	ModelOpus:   "anthropic.claude-3-opus-20240229-v1:0",
}

// KiroCLITools maps canonical tool names to Kiro CLI identifiers.
var KiroCLITools = map[Tool]string{
	ToolWebSearch: "web_search",
	ToolWebFetch:  "web_fetch",
	ToolRead:      "read",
	ToolWrite:     "write",
	ToolGlob:      "glob",
	ToolGrep:      "grep",
	ToolBash:      "bash",
	ToolEdit:      "edit",
	ToolTask:      "task",
}

// AgentKitTools maps canonical tool names to AgentKit local identifiers.
var AgentKitTools = map[Tool]string{
	ToolWebSearch: "shell",
	ToolWebFetch:  "shell",
	ToolRead:      "read",
	ToolWrite:     "write",
	ToolGlob:      "glob",
	ToolGrep:      "grep",
	ToolBash:      "shell",
	ToolEdit:      "write",
	ToolTask:      "shell",
}

// MapModelToClaudeCode converts a canonical model to Claude Code format.
func MapModelToClaudeCode(model Model) string {
	if mapped, ok := ClaudeCodeModels[model]; ok {
		return mapped
	}
	return string(model)
}

// MapModelToKiroCLI converts a canonical model to Kiro CLI format.
func MapModelToKiroCLI(model Model) string {
	if mapped, ok := KiroCLIModels[model]; ok {
		return mapped
	}
	return string(model)
}

// MapModelToBedrock converts a canonical model to AWS Bedrock format.
func MapModelToBedrock(model Model) string {
	if mapped, ok := BedrockModels[model]; ok {
		return mapped
	}
	return string(model)
}

// MapToolToKiroCLI converts a canonical tool to Kiro CLI format.
func MapToolToKiroCLI(tool Tool) string {
	if mapped, ok := KiroCLITools[tool]; ok {
		return mapped
	}
	return string(tool)
}

// MapToolToAgentKit converts a canonical tool to AgentKit local format.
func MapToolToAgentKit(tool Tool) string {
	if mapped, ok := AgentKitTools[tool]; ok {
		return mapped
	}
	return string(tool)
}

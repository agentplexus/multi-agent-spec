// Package multiagentspec provides Go types for Multi-Agent Spec definitions.
//
// This package provides structs and utilities for defining multi-agent systems
// with full JSON serialization support.
//
// Example:
//
//	agent := multiagentspec.Agent{
//	    Name:        "my-agent",
//	    Description: "A helpful agent",
//	    Model:       multiagentspec.ModelSonnet,
//	    Tools:       []string{"Read", "Write"},
//	}
//	data, _ := json.MarshalIndent(agent, "", "  ")
package multiagentspec

// Model represents the model capability tier.
type Model string

const (
	ModelHaiku  Model = "haiku"
	ModelSonnet Model = "sonnet"
	ModelOpus   Model = "opus"
)

// Tool represents canonical tool names available to agents.
type Tool string

const (
	ToolWebSearch Tool = "WebSearch"
	ToolWebFetch  Tool = "WebFetch"
	ToolRead      Tool = "Read"
	ToolWrite     Tool = "Write"
	ToolGlob      Tool = "Glob"
	ToolGrep      Tool = "Grep"
	ToolBash      Tool = "Bash"
	ToolEdit      Tool = "Edit"
	ToolTask      Tool = "Task"
)

// Agent represents an agent definition.
type Agent struct {
	// Name is the unique identifier for the agent (lowercase, hyphenated).
	Name string `json:"name"`

	// Description is a brief summary of what the agent does.
	Description string `json:"description,omitempty"`

	// Model is the capability tier (haiku, sonnet, opus).
	Model Model `json:"model,omitempty"`

	// Tools are the tools available to this agent.
	Tools []string `json:"tools,omitempty"`

	// Skills are capabilities the agent can invoke.
	Skills []string `json:"skills,omitempty"`

	// Dependencies are other agents this agent depends on.
	Dependencies []string `json:"dependencies,omitempty"`

	// Instructions is the system prompt for the agent.
	Instructions string `json:"instructions,omitempty"`
}

// NewAgent creates a new Agent with the given name and description.
func NewAgent(name, description string) *Agent {
	return &Agent{
		Name:        name,
		Description: description,
		Model:       ModelSonnet,
	}
}

// WithModel sets the agent's model and returns the agent for chaining.
func (a *Agent) WithModel(model Model) *Agent {
	a.Model = model
	return a
}

// WithTools sets the agent's tools and returns the agent for chaining.
func (a *Agent) WithTools(tools ...string) *Agent {
	a.Tools = tools
	return a
}

// WithInstructions sets the agent's instructions and returns the agent for chaining.
func (a *Agent) WithInstructions(instructions string) *Agent {
	a.Instructions = instructions
	return a
}

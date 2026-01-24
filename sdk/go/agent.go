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

// TaskType represents how a task is executed.
type TaskType string

const (
	TaskTypeCommand TaskType = "command"
	TaskTypePattern TaskType = "pattern"
	TaskTypeFile    TaskType = "file"
	TaskTypeManual  TaskType = "manual"
)

// Task represents a task that an agent can perform.
type Task struct {
	// ID is the unique task identifier within this agent.
	ID string `json:"id"`

	// Description describes what this task validates or accomplishes.
	Description string `json:"description,omitempty"`

	// Type is how the task is executed (command, pattern, file, manual).
	Type TaskType `json:"type,omitempty"`

	// Command is the shell command to execute (for type: command).
	Command string `json:"command,omitempty"`

	// Pattern is the regex pattern to search for (for type: pattern).
	Pattern string `json:"pattern,omitempty"`

	// File is the file path to check (for type: file).
	File string `json:"file,omitempty"`

	// Files is a glob pattern for files to check (for type: pattern).
	Files string `json:"files,omitempty"`

	// Required indicates if task failure causes agent to report NO-GO.
	Required *bool `json:"required,omitempty"`

	// ExpectedOutput describes what constitutes success.
	ExpectedOutput string `json:"expected_output,omitempty"`

	// HumanInLoop describes when to prompt for human intervention.
	HumanInLoop string `json:"human_in_loop,omitempty"`
}

// Agent represents an agent definition.
type Agent struct {
	// Name is the unique identifier for the agent (lowercase, hyphenated).
	Name string `json:"name" yaml:"name"`

	// Namespace is the optional namespace for organizing agents.
	// Derived from subdirectory path if not explicitly set in frontmatter.
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`

	// Description is a brief summary of what the agent does.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Model is the capability tier (haiku, sonnet, opus).
	Model Model `json:"model,omitempty" yaml:"model,omitempty"`

	// Tools are the tools available to this agent.
	Tools []string `json:"tools,omitempty" yaml:"tools,omitempty"`

	// AllowedTools are tools that can execute without user confirmation.
	AllowedTools []string `json:"allowedTools,omitempty" yaml:"allowedTools,omitempty"`

	// Skills are capabilities the agent can invoke.
	Skills []string `json:"skills,omitempty" yaml:"skills,omitempty"`

	// Dependencies are other agents this agent depends on.
	Dependencies []string `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`

	// Requires lists external tools or binaries required (e.g., go, git).
	Requires []string `json:"requires,omitempty" yaml:"requires,omitempty"`

	// Instructions is the system prompt for the agent.
	Instructions string `json:"instructions,omitempty" yaml:"instructions,omitempty"`

	// Tasks are the tasks this agent can perform.
	Tasks []Task `json:"tasks,omitempty" yaml:"tasks,omitempty"`
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

// WithNamespace sets the agent's namespace and returns the agent for chaining.
func (a *Agent) WithNamespace(namespace string) *Agent {
	a.Namespace = namespace
	return a
}

// QualifiedName returns the fully qualified agent name.
// Returns "namespace/name" if namespace is set, otherwise just "name".
func (a *Agent) QualifiedName() string {
	if a.Namespace == "" {
		return a.Name
	}
	return a.Namespace + "/" + a.Name
}

// ParseQualifiedName splits a qualified agent name into namespace and name parts.
// Returns empty namespace if no "/" is present.
//
// Examples:
//
//	ParseQualifiedName("agent-name")        → ("", "agent-name")
//	ParseQualifiedName("prd/lead")          → ("prd", "lead")
//	ParseQualifiedName("shared/review")     → ("shared", "review")
func ParseQualifiedName(qualifiedName string) (namespace, name string) {
	for i := 0; i < len(qualifiedName); i++ {
		if qualifiedName[i] == '/' {
			return qualifiedName[:i], qualifiedName[i+1:]
		}
	}
	return "", qualifiedName
}

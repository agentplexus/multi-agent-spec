package multiagentspec

// WorkflowType represents the workflow execution pattern.
type WorkflowType string

const (
	WorkflowSequential  WorkflowType = "sequential"
	WorkflowParallel    WorkflowType = "parallel"
	WorkflowDAG         WorkflowType = "dag"
	WorkflowOrchestrated WorkflowType = "orchestrated"
)

// Step represents a workflow step definition.
type Step struct {
	// Name is the step identifier.
	Name string `json:"name"`

	// Agent is the agent to execute this step.
	Agent string `json:"agent"`

	// DependsOn lists steps that must complete before this step.
	DependsOn []string `json:"depends_on,omitempty"`

	// Inputs are mappings from previous step outputs.
	Inputs map[string]string `json:"inputs,omitempty"`

	// Outputs are named outputs from this step.
	Outputs []string `json:"outputs,omitempty"`
}

// Workflow represents a workflow definition.
type Workflow struct {
	// Type is the workflow execution pattern.
	Type WorkflowType `json:"type,omitempty"`

	// Steps are the ordered steps in the workflow.
	Steps []Step `json:"steps,omitempty"`
}

// Team represents a team definition.
type Team struct {
	// Name is the team identifier (e.g., stats-agent-team).
	Name string `json:"name"`

	// Version is the semantic version of the team definition.
	Version string `json:"version"`

	// Description is a brief summary of the team's purpose.
	Description string `json:"description,omitempty"`

	// Agents is the list of agent names in the team.
	Agents []string `json:"agents"`

	// Orchestrator is the name of the orchestrator agent.
	Orchestrator string `json:"orchestrator,omitempty"`

	// Workflow is the workflow definition for agent coordination.
	Workflow *Workflow `json:"workflow,omitempty"`

	// Context is shared background information for all agents.
	Context string `json:"context,omitempty"`
}

// NewTeam creates a new Team with the given name and version.
func NewTeam(name, version string) *Team {
	return &Team{
		Name:    name,
		Version: version,
		Agents:  []string{},
	}
}

// WithAgents sets the team's agents and returns the team for chaining.
func (t *Team) WithAgents(agents ...string) *Team {
	t.Agents = agents
	return t
}

// WithOrchestrator sets the orchestrator and returns the team for chaining.
func (t *Team) WithOrchestrator(orchestrator string) *Team {
	t.Orchestrator = orchestrator
	return t
}

// WithWorkflow sets the workflow and returns the team for chaining.
func (t *Team) WithWorkflow(workflow *Workflow) *Team {
	t.Workflow = workflow
	return t
}

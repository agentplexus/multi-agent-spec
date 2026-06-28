package multiagentspec

import (
	"encoding/json"
	"fmt"
)

// WorkflowCategory represents the two workflow paradigms.
type WorkflowCategory string

const (
	// CategoryDeterministic workflows have execution paths defined in the schema.
	CategoryDeterministic WorkflowCategory = "deterministic"

	// CategorySelfDirected workflows let agents decide execution paths.
	CategorySelfDirected WorkflowCategory = "self-directed"
)

// String returns the string representation of the category.
func (c WorkflowCategory) String() string {
	return string(c)
}

// WorkflowType represents the workflow execution pattern.
type WorkflowType string

const (
	// Deterministic workflows (schema controls execution)

	// WorkflowChain executes steps sequentially: A → B → C
	WorkflowChain WorkflowType = "chain"

	// WorkflowScatter executes steps in parallel with fan-out/fan-in: A → [B,C,D] → E
	WorkflowScatter WorkflowType = "scatter"

	// WorkflowGraph executes steps as a DAG with conditional paths.
	WorkflowGraph WorkflowType = "graph"

	// Self-directed workflows (agents control execution)
	// These will be fully implemented in v0.8.0 - see PRD_SELFDIRECTED.md

	// WorkflowCrew has a lead agent delegating to specialists.
	WorkflowCrew WorkflowType = "crew"

	// WorkflowSwarm is self-organizing with shared task queue.
	WorkflowSwarm WorkflowType = "swarm"

	// WorkflowCouncil uses peer debate and consensus.
	WorkflowCouncil WorkflowType = "council"
)

// Category returns the workflow category for this type.
func (w WorkflowType) Category() WorkflowCategory {
	switch w {
	case WorkflowChain, WorkflowScatter, WorkflowGraph:
		return CategoryDeterministic
	case WorkflowCrew, WorkflowSwarm, WorkflowCouncil:
		return CategorySelfDirected
	default:
		return ""
	}
}

// IsDeterministic returns true if this is a deterministic workflow type.
func (w WorkflowType) IsDeterministic() bool {
	return w.Category() == CategoryDeterministic
}

// IsSelfDirected returns true if this is a self-directed workflow type.
func (w WorkflowType) IsSelfDirected() bool {
	return w.Category() == CategorySelfDirected
}

// PortType represents the data type of a port.
type PortType string

const (
	PortTypeString  PortType = "string"
	PortTypeNumber  PortType = "number"
	PortTypeBoolean PortType = "boolean"
	PortTypeObject  PortType = "object"
	PortTypeArray   PortType = "array"
	PortTypeFile    PortType = "file"
)

// Port represents a typed input or output for a workflow step.
type Port struct {
	// Name is the port identifier (e.g., version_recommendation, test_results).
	Name string `json:"name"`

	// Type is the data type of this port.
	Type PortType `json:"type,omitempty"`

	// Description is a human-readable description of this data.
	Description string `json:"description,omitempty"`

	// Required indicates whether this input is required (inputs only).
	Required *bool `json:"required,omitempty"`

	// From is the source reference as 'step_name.output_name' (inputs only).
	From string `json:"from,omitempty"`

	// Schema is a JSON Schema for validating this port's data.
	Schema json.RawMessage `json:"schema,omitempty"`

	// Default is the default value if not provided (inputs only).
	Default interface{} `json:"default,omitempty"`
}

// Step represents a workflow step definition.
type Step struct {
	// Name is the step identifier.
	Name string `json:"name"`

	// Agent is the agent to execute this step (mutually exclusive with Loop).
	Agent string `json:"agent,omitempty"`

	// Loop is the loop to execute this step (mutually exclusive with Agent).
	// Use this for REAL or VEAL loop steps.
	Loop string `json:"loop,omitempty"`

	// DependsOn lists steps that must complete before this step.
	DependsOn []string `json:"depends_on,omitempty"`

	// Inputs are typed data inputs consumed by this step.
	Inputs []Port `json:"inputs,omitempty"`

	// Outputs are typed data outputs produced by this step.
	Outputs []Port `json:"outputs,omitempty"`
}

// IsLoopStep returns true if this step executes a loop.
func (s *Step) IsLoopStep() bool {
	return s.Loop != ""
}

// IsAgentStep returns true if this step executes an agent.
func (s *Step) IsAgentStep() bool {
	return s.Agent != ""
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

	// Self-directed workflow fields

	// Collaboration defines how agents interact in self-directed workflows.
	Collaboration *CollaborationConfig `json:"collaboration,omitempty"`

	// SelfClaim allows agents to self-claim tasks from a shared queue (swarm).
	SelfClaim bool `json:"self_claim,omitempty"`

	// PlanApproval requires plan approval before implementation (crew).
	PlanApproval bool `json:"plan_approval,omitempty"`
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

// WorkflowCategory returns the category of this team's workflow.
// Returns empty string if no workflow is defined.
func (t *Team) WorkflowCategory() WorkflowCategory {
	if t.Workflow == nil {
		return ""
	}
	return t.Workflow.Type.Category()
}

// IsDeterministic returns true if this team uses a deterministic workflow.
func (t *Team) IsDeterministic() bool {
	return t.WorkflowCategory() == CategoryDeterministic
}

// IsSelfDirected returns true if this team uses a self-directed workflow.
func (t *Team) IsSelfDirected() bool {
	return t.WorkflowCategory() == CategorySelfDirected
}

// WithCollaboration sets the collaboration config and returns the team for chaining.
func (t *Team) WithCollaboration(collab *CollaborationConfig) *Team {
	t.Collaboration = collab
	return t
}

// WithSelfClaim enables self-claim and returns the team for chaining.
func (t *Team) WithSelfClaim(selfClaim bool) *Team {
	t.SelfClaim = selfClaim
	return t
}

// WithPlanApproval enables plan approval and returns the team for chaining.
func (t *Team) WithPlanApproval(planApproval bool) *Team {
	t.PlanApproval = planApproval
	return t
}

// Validate checks team configuration consistency.
// Returns an error if the configuration is invalid for the workflow type.
func (t *Team) Validate() error {
	if t.Workflow == nil {
		return nil
	}

	wt := t.Workflow.Type

	switch wt {
	case WorkflowCrew:
		// Crew workflow requires a lead agent
		if t.Collaboration == nil || t.Collaboration.Lead == "" {
			// Fall back to orchestrator if set
			if t.Orchestrator == "" {
				return fmt.Errorf("crew workflow requires collaboration.lead or orchestrator")
			}
		}
	case WorkflowSwarm:
		// Swarm workflow requires task_queue or self_claim
		hasTaskQueue := t.Collaboration != nil && t.Collaboration.TaskQueue
		if !hasTaskQueue && !t.SelfClaim {
			return fmt.Errorf("swarm workflow requires collaboration.task_queue or self_claim")
		}
	case WorkflowCouncil:
		// Council workflow requires consensus rules
		if t.Collaboration == nil || t.Collaboration.Consensus == nil {
			return fmt.Errorf("council workflow requires collaboration.consensus")
		}
	}

	return nil
}

// EffectiveLead returns the lead agent name for self-directed workflows.
// Checks collaboration.lead first, then falls back to orchestrator.
func (t *Team) EffectiveLead() string {
	if t.Collaboration != nil && t.Collaboration.Lead != "" {
		return t.Collaboration.Lead
	}
	return t.Orchestrator
}

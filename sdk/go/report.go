package multiagentspec

import (
	"encoding/json"
	"time"
)

// Status represents the validation status following NASA Go/No-Go terminology.
type Status string

const (
	StatusGo   Status = "GO"
	StatusWarn Status = "WARN"
	StatusNoGo Status = "NO-GO"
	StatusSkip Status = "SKIP"
)

// Icon returns the UTF-8 icon for the status.
func (s Status) Icon() string {
	switch s {
	case StatusGo:
		return "\U0001F7E2" // 🟢
	case StatusWarn:
		return "\U0001F7E1" // 🟡
	case StatusNoGo:
		return "\U0001F534" // 🔴
	case StatusSkip:
		return "\u26AA" // ⚪
	default:
		return "?"
	}
}

// TaskResult represents the result of executing a single task.
// Each task corresponds to a task defined in the agent's task list.
type TaskResult struct {
	// ID is the task identifier (matches task id in agent definition)
	ID string `json:"id"`

	// Status is GO, WARN, NO-GO, or SKIP
	Status Status `json:"status"`

	// Detail is optional additional information about the result
	Detail string `json:"detail,omitempty"`

	// DurationMs is the task execution time in milliseconds
	DurationMs int64 `json:"duration_ms,omitempty"`

	// Metadata allows tasks to include structured data
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// AgentResult is the JSON-serializable output from each validation agent.
// This is the intermediate representation that agents produce and the
// coordinator consumes to build the final TeamReport.
type AgentResult struct {
	// Schema is the JSON schema URL for validation
	Schema string `json:"$schema,omitempty"`

	// AgentID identifies the agent (e.g., "pm", "qa", "documentation")
	AgentID string `json:"agent_id"`

	// StepID is the workflow step name (e.g., "pm-validation", "qa-validation")
	StepID string `json:"step_id"`

	// Inputs are values received from upstream agents in the DAG
	Inputs map[string]interface{} `json:"inputs,omitempty"`

	// Outputs are values produced by this agent for downstream agents
	Outputs map[string]interface{} `json:"outputs,omitempty"`

	// Tasks are the individual task results (one per line in reports)
	Tasks []TaskResult `json:"tasks"`

	// Status is the overall status for this agent (computed from tasks)
	Status Status `json:"status"`

	// ExecutedAt is when the agent completed execution
	ExecutedAt time.Time `json:"executed_at"`

	// AgentModel is the LLM model used (e.g., "sonnet", "haiku", "opus")
	AgentModel string `json:"agent_model,omitempty"`

	// Duration is how long the agent took to execute
	Duration string `json:"duration,omitempty"`

	// Error is set if the agent failed to execute
	Error string `json:"error,omitempty"`
}

// TeamSection represents a team/agent section in the report.
type TeamSection struct {
	// ID is the workflow step ID (e.g., "pm-validation")
	ID string `json:"id"`

	// Name is the agent name (e.g., "pm")
	Name string `json:"name"`

	// AgentID matches the agent definition in team.json
	AgentID string `json:"agent_id,omitempty"`

	// Model is the LLM model used
	Model string `json:"model,omitempty"`

	// DependsOn lists the IDs of upstream teams in the DAG
	DependsOn []string `json:"depends_on,omitempty"`

	// Tasks are the task results for this team (one per line in reports)
	Tasks []TaskResult `json:"tasks"`

	// Status is the overall status (computed from tasks)
	Status Status `json:"status"`
}

// TeamReport is the complete JSON-serializable report.
// This is what the coordinator produces by aggregating AgentResults.
type TeamReport struct {
	// Schema is the JSON schema URL for validation
	Schema string `json:"$schema,omitempty"`

	// Project is the repository identifier
	Project string `json:"project"`

	// Version is the target release version
	Version string `json:"version"`

	// Target is a human-readable target description
	Target string `json:"target,omitempty"`

	// Phase is the workflow phase (e.g., "PHASE 1: REVIEW")
	Phase string `json:"phase"`

	// Teams are the validation teams/agents
	Teams []TeamSection `json:"teams"`

	// Status is the overall status (computed from teams)
	Status Status `json:"status"`

	// GeneratedAt is when the report was generated
	GeneratedAt time.Time `json:"generated_at"`

	// GeneratedBy identifies the coordinator
	GeneratedBy string `json:"generated_by,omitempty"`
}

// ComputeStatus computes the overall status from tasks.
func (a *AgentResult) ComputeStatus() Status {
	return computeStatusFromTasks(a.Tasks)
}

// ToTeamSection converts an AgentResult to a TeamSection for the report.
func (a *AgentResult) ToTeamSection() TeamSection {
	return TeamSection{
		ID:      a.StepID,
		Name:    a.AgentID,
		AgentID: a.AgentID,
		Model:   a.AgentModel,
		Tasks:   a.Tasks,
		Status:  a.ComputeStatus(),
	}
}

// OverallStatus computes the overall status for a team section.
func (t *TeamSection) OverallStatus() Status {
	return computeStatusFromTasks(t.Tasks)
}

// ComputeOverallStatus computes the overall status from all teams.
func (r *TeamReport) ComputeOverallStatus() Status {
	hasNoGo := false
	hasWarn := false

	for _, t := range r.Teams {
		switch t.Status {
		case StatusNoGo:
			hasNoGo = true
		case StatusWarn:
			hasWarn = true
		}
	}

	if hasNoGo {
		return StatusNoGo
	}
	if hasWarn {
		return StatusWarn
	}
	return StatusGo
}

// IsGo returns true if all teams pass validation.
func (r *TeamReport) IsGo() bool {
	for _, t := range r.Teams {
		if t.Status == StatusNoGo {
			return false
		}
	}
	return true
}

// SortByDAG sorts teams in topological order based on DependsOn relationships.
// Teams with no dependencies appear first, followed by teams whose dependencies
// have been satisfied. This uses Kahn's algorithm for topological sorting.
// Teams at the same level are sorted alphabetically by ID for deterministic output.
// If the DAG has cycles, teams in cycles appear at the end in their original order.
func (r *TeamReport) SortByDAG() {
	if len(r.Teams) <= 1 {
		return
	}

	// Build ID -> index mapping and in-degree count
	idToTeam := make(map[string]*TeamSection)
	inDegree := make(map[string]int)

	for i := range r.Teams {
		id := r.Teams[i].ID
		idToTeam[id] = &r.Teams[i]
		inDegree[id] = 0
	}

	// Calculate in-degrees (count of dependencies)
	for i := range r.Teams {
		for _, dep := range r.Teams[i].DependsOn {
			if _, exists := idToTeam[dep]; exists {
				inDegree[r.Teams[i].ID]++
			}
		}
	}

	// Build adjacency list (downstream teams for each team)
	downstream := make(map[string][]string)
	for i := range r.Teams {
		for _, dep := range r.Teams[i].DependsOn {
			if _, exists := idToTeam[dep]; exists {
				downstream[dep] = append(downstream[dep], r.Teams[i].ID)
			}
		}
	}

	// Kahn's algorithm: start with teams that have no dependencies
	// Collect all teams with zero in-degree, sort alphabetically
	queue := make([]string, 0)
	for i := range r.Teams {
		if inDegree[r.Teams[i].ID] == 0 {
			queue = append(queue, r.Teams[i].ID)
		}
	}
	sortStrings(queue)

	sorted := make([]TeamSection, 0, len(r.Teams))
	processed := make(map[string]bool)

	for len(queue) > 0 {
		// Dequeue
		id := queue[0]
		queue = queue[1:]

		if processed[id] {
			continue
		}
		processed[id] = true

		sorted = append(sorted, *idToTeam[id])

		// Collect newly ready teams
		newlyReady := make([]string, 0)
		for _, downID := range downstream[id] {
			inDegree[downID]--
			if inDegree[downID] == 0 {
				newlyReady = append(newlyReady, downID)
			}
		}
		// Sort and append to maintain alphabetical order at each level
		sortStrings(newlyReady)
		queue = append(queue, newlyReady...)
	}

	// Add any remaining teams (cycles or missing dependencies) at the end
	for i := range r.Teams {
		if !processed[r.Teams[i].ID] {
			sorted = append(sorted, r.Teams[i])
		}
	}

	r.Teams = sorted
}

// sortStrings sorts a slice of strings in place.
func sortStrings(s []string) {
	for i := 0; i < len(s)-1; i++ {
		for j := i + 1; j < len(s); j++ {
			if s[j] < s[i] {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}

// FinalMessage returns the final status message for display.
func (r *TeamReport) FinalMessage() string {
	if r.IsGo() {
		return "\U0001F680 TEAM: GO for " + r.Version + " \U0001F680" // 🚀 TEAM: GO for vX.Y.Z 🚀
	}
	return "\U0001F6D1 TEAM: NO-GO for " + r.Version + " \U0001F6D1" // 🛑 TEAM: NO-GO for vX.Y.Z 🛑
}

// ToJSON serializes the report to JSON.
func (r *TeamReport) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// AggregateResults combines multiple AgentResults into a TeamReport.
func AggregateResults(results []AgentResult, project, version, phase string) *TeamReport {
	teams := make([]TeamSection, 0, len(results))
	for _, r := range results {
		teams = append(teams, r.ToTeamSection())
	}

	report := &TeamReport{
		Schema:      "https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/report/team-report.schema.json",
		Project:     project,
		Version:     version,
		Target:      version,
		Phase:       phase,
		Teams:       teams,
		GeneratedAt: time.Now().UTC(),
		GeneratedBy: "release-coordinator",
	}

	report.Status = report.ComputeOverallStatus()

	return report
}

// ParseAgentResult parses JSON into an AgentResult.
func ParseAgentResult(data []byte) (*AgentResult, error) {
	var result AgentResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ParseTeamReport parses JSON into a TeamReport.
func ParseTeamReport(data []byte) (*TeamReport, error) {
	var report TeamReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, err
	}
	return &report, nil
}

// computeStatusFromTasks is a helper to compute status from a slice of task results.
func computeStatusFromTasks(tasks []TaskResult) Status {
	hasNoGo := false
	hasWarn := false
	allSkipped := true

	for _, t := range tasks {
		if t.Status != StatusSkip {
			allSkipped = false
		}
		switch t.Status {
		case StatusNoGo:
			hasNoGo = true
		case StatusWarn:
			hasWarn = true
		}
	}

	if allSkipped {
		return StatusSkip
	}
	if hasNoGo {
		return StatusNoGo
	}
	if hasWarn {
		return StatusWarn
	}
	return StatusGo
}

// Backward compatibility aliases
// Deprecated: Use TaskResult instead
type Check = TaskResult

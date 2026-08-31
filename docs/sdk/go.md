# Go SDK

The Go SDK provides typed structs and utilities for working with multi-agent-spec.

## Installation

```bash
go get github.com/plexusone/multi-agent-spec@latest
```

## Import

```go
import mas "github.com/plexusone/multi-agent-spec/sdk/go"
```

## Core Types

### Agent

```go
type Agent struct {
    Name         string            `json:"name"`
    Namespace    string            `json:"namespace,omitempty"`
    Description  string            `json:"description,omitempty"`
    Icon         string            `json:"icon,omitempty"`
    Model        Model             `json:"model,omitempty"`
    Tools        []string          `json:"tools,omitempty"`
    AllowedTools []string          `json:"allowedTools,omitempty"`
    Skills       []string          `json:"skills,omitempty"`
    Dependencies []string          `json:"dependencies,omitempty"`
    Requires     []string          `json:"requires,omitempty"`
    Instructions string            `json:"instructions,omitempty"`
    Tasks        []Task            `json:"tasks,omitempty"`

    // Self-directed workflow fields
    Role         string            `json:"role,omitempty"`
    Goal         string            `json:"goal,omitempty"`
    Backstory    string            `json:"backstory,omitempty"`
    Delegation   *DelegationConfig `json:"delegation,omitempty"`
}

// Builder methods
agent := mas.NewAgent("security-reviewer", "Reviews code for security").
    WithModel(mas.ModelSonnet).
    WithTools("Read", "Grep", "Glob").
    WithRole("Security Analyst").
    WithGoal("Identify vulnerabilities").
    WithBackstory("10 years security experience").
    WithDelegation(&mas.DelegationConfig{
        AllowDelegation: false,
        CanReceiveFrom:  []string{"architect"},
    })

// Delegation helpers
agent.CanDelegate()              // true if can delegate work
agent.CanDelegateTo("frontend")  // true if can delegate to frontend
agent.CanReceiveFrom("architect") // true if can receive from architect
agent.QualifiedName()            // "namespace/name" or "name"
```

### DelegationConfig

```go
type DelegationConfig struct {
    AllowDelegation bool     `json:"allow_delegation,omitempty"`
    CanDelegateTo   []string `json:"can_delegate_to,omitempty"`
    CanReceiveFrom  []string `json:"can_receive_from,omitempty"`
}
```

### Team

```go
type Team struct {
    Name          string              `json:"name"`
    Version       string              `json:"version"`
    Description   string              `json:"description,omitempty"`
    Agents        []string            `json:"agents"`
    Orchestrator  string              `json:"orchestrator,omitempty"`
    Workflow      *Workflow           `json:"workflow,omitempty"`
    Context       string              `json:"context,omitempty"`

    // Self-directed workflow fields
    Collaboration *CollaborationConfig `json:"collaboration,omitempty"`
    SelfClaim     bool                 `json:"self_claim,omitempty"`
    PlanApproval  bool                 `json:"plan_approval,omitempty"`
}

// Builder methods
team := mas.NewTeam("dev-team", "1.0.0").
    WithAgents("architect", "frontend", "backend").
    WithWorkflow(&mas.Workflow{Type: mas.WorkflowCrew}).
    WithCollaboration(&mas.CollaborationConfig{
        Lead:        "architect",
        Specialists: []string{"frontend", "backend"},
    }).
    WithPlanApproval(true)

// Workflow helpers
team.WorkflowCategory()  // CategoryDeterministic or CategorySelfDirected
team.IsDeterministic()   // true for chain, scatter, graph
team.IsSelfDirected()    // true for crew, swarm, council
team.EffectiveLead()     // returns lead agent name
team.Validate()          // validates workflow-specific requirements
team.ValidateAgentReferences(names) // returns joined error for unresolved agent refs
```

### Workflow Types

```go
// Workflow categories
const (
    CategoryDeterministic WorkflowCategory = "deterministic"
    CategorySelfDirected  WorkflowCategory = "self-directed"
)

// Workflow types
const (
    // Deterministic (schema controls execution)
    WorkflowChain   WorkflowType = "chain"   // A → B → C
    WorkflowScatter WorkflowType = "scatter" // A → [B,C,D] → E
    WorkflowGraph   WorkflowType = "graph"   // DAG with dependencies

    // Self-directed (agents control execution)
    WorkflowCrew    WorkflowType = "crew"    // Lead delegates to specialists
    WorkflowSwarm   WorkflowType = "swarm"   // Self-claiming from queue
    WorkflowCouncil WorkflowType = "council" // Peer debate + consensus
)

// Helper methods
wt := mas.WorkflowCrew
wt.Category()        // CategorySelfDirected
wt.IsDeterministic() // false
wt.IsSelfDirected()  // true
```

### Workflow

```go
type Workflow struct {
    Type  WorkflowType `json:"type,omitempty"`
    Steps []Step       `json:"steps,omitempty"`
}

type Step struct {
    Name      string   `json:"name"`
    Agent     string   `json:"agent"`
    DependsOn []string `json:"depends_on,omitempty"`
    Inputs    []Port   `json:"inputs,omitempty"`
    Outputs   []Port   `json:"outputs,omitempty"`
}

type Port struct {
    Name        string `json:"name"`
    Type        string `json:"type,omitempty"`
    Description string `json:"description,omitempty"`
    From        string `json:"from,omitempty"`
}
```

### CollaborationConfig

Configuration for self-directed workflows.

```go
type CollaborationConfig struct {
    Lead        string          `json:"lead,omitempty"`
    Specialists []string        `json:"specialists,omitempty"`
    TaskQueue   bool            `json:"task_queue,omitempty"`
    Consensus   *ConsensusRules `json:"consensus,omitempty"`
    Channels    []Channel       `json:"channels,omitempty"`
}

type ConsensusRules struct {
    RequiredAgreement float64 `json:"required_agreement,omitempty"` // 0.0-1.0
    MaxRounds         int     `json:"max_rounds,omitempty"`
    TieBreaker        string  `json:"tie_breaker,omitempty"`
}

type Channel struct {
    Name         string      `json:"name"`
    Type         ChannelType `json:"type"`
    Participants []string    `json:"participants,omitempty"`
}

const (
    ChannelDirect    ChannelType = "direct"
    ChannelBroadcast ChannelType = "broadcast"
    ChannelPubSub    ChannelType = "pub-sub"
)
```

### Message

Inter-agent messages for self-directed workflows.

```go
type Message struct {
    ID          string                 `json:"id"`
    Type        MessageType            `json:"type"`
    From        string                 `json:"from"`
    To          string                 `json:"to,omitempty"`
    Subject     string                 `json:"subject,omitempty"`
    Content     string                 `json:"content"`
    Attachments []Attachment           `json:"attachments,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
    Timestamp   time.Time              `json:"timestamp"`
}

// Message types
const (
    MsgDelegateWork     MessageType = "delegate_work"
    MsgAskQuestion      MessageType = "ask_question"
    MsgShareFinding     MessageType = "share_finding"
    MsgRequestApproval  MessageType = "request_approval"
    MsgApproval         MessageType = "approval"
    MsgRejection        MessageType = "rejection"
    MsgChallenge        MessageType = "challenge"
    MsgVote             MessageType = "vote"
    MsgTaskClaimed      MessageType = "task_claimed"
    MsgTaskCompleted    MessageType = "task_completed"
    MsgShutdownRequest  MessageType = "shutdown_request"
    MsgShutdownApproved MessageType = "shutdown_approved"
)

// Create a new message
msg := mas.NewMessage(mas.MsgDelegateWork, "architect", "frontend", "Implement login form")
msg.IsBroadcast() // true if To is "*" or empty
```

### Deployment

```go
type Deployment struct {
    Team    string   `json:"team"`
    Targets []Target `json:"targets"`
}

type Target struct {
    Name         string              `json:"name"`
    Platform     Platform            `json:"platform"`
    Output       string              `json:"output,omitempty"`
    ClaudeCode   *ClaudeCodeConfig   `json:"claudeCode,omitempty"`
    KiroCLI      *KiroCLIConfig      `json:"kiroCli,omitempty"`
    CrewAI       *CrewAIConfig       `json:"crewai,omitempty"`
    // ... other platform configs
}

type ClaudeCodeConfig struct {
    AgentDir     string `json:"agentDir"`
    Format       string `json:"format"`
    TeamMode     string `json:"team_mode,omitempty"`     // subagent or team
    TeammateMode string `json:"teammate_mode,omitempty"` // in-process, tmux, auto
    EnableTeams  bool   `json:"enable_teams,omitempty"`
}

type CrewAIConfig struct {
    Model           string `json:"model,omitempty"`
    ProcessType     string `json:"processType,omitempty"`
    AllowDelegation bool   `json:"allowDelegation,omitempty"`
    ManagerLLM      string `json:"managerLlm,omitempty"`
}
```

### TeamReport

```go
type TeamReport struct {
    Title         string            `json:"title,omitempty"`
    Project       string            `json:"project"`
    Version       string            `json:"version"`
    Phase         string            `json:"phase"`
    Tags          map[string]string `json:"tags,omitempty"`
    Teams         []TeamSection     `json:"teams"`
    Status        Status            `json:"status"`
    Summary       string            `json:"summary,omitempty"`
    Conclusion    string            `json:"conclusion,omitempty"`
    SummaryBlocks []ContentBlock    `json:"summary_blocks,omitempty"`
    FooterBlocks  []ContentBlock    `json:"footer_blocks,omitempty"`
    GeneratedAt   time.Time         `json:"generated_at"`
}

type TeamSection struct {
    ID            string         `json:"id"`
    Name          string         `json:"name"`
    Status        Status         `json:"status"`
    Verdict       string         `json:"verdict,omitempty"`
    Tasks         []TaskResult   `json:"tasks,omitempty"`
    ContentBlocks []ContentBlock `json:"content_blocks,omitempty"`
    Narrative     string         `json:"narrative,omitempty"`
}
```

## Status Constants

```go
const (
    StatusGo   Status = "GO"
    StatusWarn Status = "WARN"
    StatusNoGo Status = "NO-GO"
    StatusSkip Status = "SKIP"
)

status := mas.StatusGo
icon := status.Icon() // "🟢"
```

## Platform Constants

```go
const (
    PlatformClaudeCode   Platform = "claude-code"
    PlatformKiroCLI      Platform = "kiro-cli"
    PlatformGeminiCLI    Platform = "gemini-cli"
    PlatformAWSAgentCore Platform = "aws-agentcore"
    PlatformCrewAI       Platform = "crewai"
    PlatformAutoGen      Platform = "autogen"
)
```

## Content Blocks

Rich content for reports.

```go
// Create content blocks
textBlock := mas.NewTextBlock("Analysis complete")
listBlock := mas.NewListBlock([]string{"Item 1", "Item 2", "Item 3"})
tableBlock := mas.NewTableBlock(
    []string{"Name", "Status"},
    [][]string{{"Test A", "Pass"}, {"Test B", "Fail"}},
)
kvBlock := mas.NewKVPairsBlock(map[string]string{
    "Coverage": "85%",
    "Tests":    "42 passed",
})
metricBlock := mas.NewMetricBlock("Coverage", 85.5, "%")
```

## Creating Self-Directed Teams

### Crew Workflow

```go
team := mas.NewTeam("dev-team", "1.0.0").
    WithAgents("architect", "frontend", "backend", "qa").
    WithWorkflow(&mas.Workflow{Type: mas.WorkflowCrew}).
    WithCollaboration(&mas.CollaborationConfig{
        Lead:        "architect",
        Specialists: []string{"frontend", "backend", "qa"},
    }).
    WithPlanApproval(true)

if err := team.Validate(); err != nil {
    log.Fatal(err) // "crew workflow requires collaboration.lead"
}

// Check the team's roster, orchestrator, and workflow steps against the
// agents actually available, e.g. after scanning an agent spec directory.
availableAgents := []string{"architect", "frontend", "backend", "qa"}
if err := team.ValidateAgentReferences(availableAgents); err != nil {
    log.Fatal(err) // e.g. `team "dev-team" roster references agent "backend-old" with no matching spec`
}
```

### Swarm Workflow

```go
team := mas.NewTeam("triage-team", "1.0.0").
    WithAgents("triager-1", "triager-2", "triager-3").
    WithWorkflow(&mas.Workflow{Type: mas.WorkflowSwarm}).
    WithCollaboration(&mas.CollaborationConfig{
        TaskQueue: true,
    })
```

### Council Workflow

```go
team := mas.NewTeam("review-council", "1.0.0").
    WithAgents("reviewer-1", "reviewer-2", "reviewer-3").
    WithWorkflow(&mas.Workflow{Type: mas.WorkflowCouncil}).
    WithCollaboration(&mas.CollaborationConfig{
        Consensus: &mas.ConsensusRules{
            RequiredAgreement: 0.66,
            MaxRounds:         3,
            TieBreaker:        "reviewer-1",
        },
        Channels: []mas.Channel{
            {Name: "findings", Type: mas.ChannelBroadcast, Participants: []string{"*"}},
        },
    })
```

## Rendering Reports

### Box Format

```go
renderer := mas.NewBoxRenderer()
output := renderer.Render(report)
fmt.Println(output)
```

### Narrative Format

```go
renderer := mas.NewNarrativeRenderer()
markdown := renderer.Render(report)
os.WriteFile("report.md", []byte(markdown), 0644)
```

## Loading Definitions

```go
// Load agent from file
agent, err := mas.LoadAgentFromFile("specs/agents/security.md")

// Load team from file
team, err := mas.LoadTeamFromFile("specs/teams/release.json")

// Load deployment from file
deployment, err := mas.LoadDeploymentFromFile("specs/deployments/local.json")

// Load all agents from directory (recursive, with namespaces)
agents, err := mas.LoadAgentsFromDir("specs/agents")

// Load agents flat (non-recursive)
agents, err := mas.LoadAgentsFromDirFlat("specs/agents")
```

## See Also

- [Agent Schema](../schemas/agent.md) - Agent fields and role-based config
- [Team Schema](../schemas/team.md) - Workflow types and collaboration
- [Deployment Schema](../schemas/deployment.md) - Platform configuration
- [mas CLI](../cli/mas.md) - Command-line interface

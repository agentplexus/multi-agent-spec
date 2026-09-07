# Schema Overview

Multi-Agent Spec uses JSON Schema (draft 2020-12) for all definitions. This enables validation, IDE support, and code generation.

## Core Schemas

| Schema | Purpose | File |
|--------|---------|------|
| [Agent](agent.md) | Individual agent definitions | `agent/agent.schema.json` |
| [Team](team.md) | Multi-agent workflows | `orchestration/team.schema.json` |
| [Deployment](deployment.md) | Platform-specific configs | `deployment/deployment.schema.json` |
| [Report](report.md) | Execution results | `report/team-report.schema.json` |
| [Skill](skill.md) | Reusable agent capabilities | `skill/skill.schema.json` |
| [Loop](loop.md) | Autonomous bounded execution | `orchestration/loop.schema.json` |
| [Review](review.md) | Human-in-the-loop confirmation checkpoints | `review/review.schema.json` |
| Message | Inter-agent messaging | `message/message.schema.json` |

## Workflow Categories

Multi-Agent Spec supports two workflow paradigms:

| Category | Description | Types |
|----------|-------------|-------|
| **Deterministic** | Schema controls execution paths | `chain`, `scatter`, `graph` |
| **Self-directed** | Agents control execution paths | `crew`, `swarm`, `council` |

See [Team Schema](team.md) for details on each workflow type.

## Schema URLs

All schemas are hosted on GitHub:

```
https://raw.githubusercontent.com/plexusone/multi-agent-spec/main/schema/{category}/{name}.schema.json
```

| Schema | URL |
|--------|-----|
| Agent | `.../schema/agent/agent.schema.json` |
| Team | `.../schema/orchestration/team.schema.json` |
| Deployment | `.../schema/deployment/deployment.schema.json` |
| Skill | `.../schema/skill/skill.schema.json` |
| Loop | `.../schema/orchestration/loop.schema.json` |
| Review | `.../schema/review/review.schema.json` |
| Message | `.../schema/message/message.schema.json` |

## Using Schemas

### In JSON Files

Add `$schema` to enable validation:

```json
{
  "$schema": "https://raw.githubusercontent.com/plexusone/multi-agent-spec/main/schema/agent/agent.schema.json",
  "name": "my-agent",
  "description": "My agent description"
}
```

### In VS Code

Add to `.vscode/settings.json`:

```json
{
  "json.schemas": [
    {
      "fileMatch": ["**/agents/*.json"],
      "url": "https://raw.githubusercontent.com/plexusone/multi-agent-spec/main/schema/agent/agent.schema.json"
    },
    {
      "fileMatch": ["**/teams/*.json"],
      "url": "https://raw.githubusercontent.com/plexusone/multi-agent-spec/main/schema/orchestration/team.schema.json"
    },
    {
      "fileMatch": ["**/deployments/*.json"],
      "url": "https://raw.githubusercontent.com/plexusone/multi-agent-spec/main/schema/deployment/deployment.schema.json"
    }
  ]
}
```

## Key Types

### Workflow Types

| Type | Category | Pattern | Use Case |
|------|----------|---------|----------|
| `chain` | Deterministic | A → B → C | Sequential pipeline |
| `scatter` | Deterministic | A → [B,C,D] → E | Parallel fan-out |
| `graph` | Deterministic | DAG | Complex dependencies |
| `crew` | Self-directed | Lead → Specialists | Delegation hierarchy |
| `swarm` | Self-directed | Shared queue | Self-organizing |
| `council` | Self-directed | Peer debate | Consensus voting |

### Agent Role Fields

For self-directed workflows, agents use role-based fields:

| Field | Purpose |
|-------|---------|
| `role` | Agent's job title |
| `goal` | What the agent aims to achieve |
| `backstory` | Context for autonomous decisions |
| `delegation` | Delegation permissions |

### Message Types

For inter-agent communication:

| Type | Description |
|------|-------------|
| `delegate_work` | Lead assigns task |
| `ask_question` | Request information |
| `share_finding` | Broadcast discovery |
| `request_approval` | Seek plan approval |
| `challenge` | Dispute a finding |
| `vote` | Cast consensus vote |

## Go SDK

The Go SDK provides typed structs for all schemas:

```go
import mas "github.com/plexusone/multi-agent-spec/sdk/go"

// Basic agent
agent := mas.NewAgent("analyzer", "Analyzes code").
    WithModel(mas.ModelSonnet).
    WithTools("Read", "Grep")

// Self-directed agent with role
reviewer := mas.NewAgent("security-reviewer", "Reviews security").
    WithRole("Security Analyst").
    WithGoal("Find vulnerabilities").
    WithDelegation(&mas.DelegationConfig{
        AllowDelegation: false,
        CanReceiveFrom:  []string{"architect"},
    })

// Team with self-directed workflow
team := mas.NewTeam("dev-team", "1.0.0").
    WithAgents("architect", "frontend", "backend").
    WithWorkflow(&mas.Workflow{Type: mas.WorkflowCrew}).
    WithCollaboration(&mas.CollaborationConfig{
        Lead:        "architect",
        Specialists: []string{"frontend", "backend"},
    })

// Check workflow category
if team.IsSelfDirected() {
    // Agents control execution
}
```

## Design Principles

1. **Go-first** - Schemas are generated from Go types for perfect alignment
2. **Portable** - Platform-agnostic core with platform-specific extensions
3. **Validated** - All fields have types and constraints
4. **Extensible** - `metadata` fields allow custom data without schema changes
5. **Two paradigms** - Supports both deterministic and self-directed workflows

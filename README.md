# Multi-Agent Spec

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/multi-agent-spec/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/multi-agent-spec/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/multi-agent-spec/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/multi-agent-spec/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/multi-agent-spec/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/multi-agent-spec/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/multi-agent-spec
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/multi-agent-spec
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/multi-agent-spec
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/multi-agent-spec
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://plexusone.dev/multi-agent-spec
 [viz-svg]: https://img.shields.io/badge/Go-visualizaton-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fmulti-agent-spec
 [loc-svg]: https://tokei.rs/b1/github/plexusone/multi-agent-spec
 [repo-url]: https://github.com/plexusone/multi-agent-spec
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/multi-agent-spec/blob/main/LICENSE

A specification for defining multi-agent AI systems with platform-agnostic agent definitions and deployment configurations.

## Overview

Multi-Agent Spec provides a standardized way to define:

1. **Agents** - Individual AI agents with capabilities, tools, and instructions
2. **Teams** - Groups of agents with orchestration patterns
3. **Deployments** - Target platforms and configurations

## Architecture

```
┌───────────────────────────────────────────────────────────────────────────┐
│                     Definition Layer                                      │
├───────────────────────────────────────────────────────────────────────────┤
│  specs/agents/*.md    │  specs/teams/*.json  │  specs/deployments/*.json  │
│  (Markdown + YAML)    │  (Orchestration)     │  (Targets)                 │
└───────────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│                     Deployment Layer                             │
├────────────┬────────────┬────────────┬────────────┬──────────────┤
│ Claude     │ Kiro       │ AWS        │ AWS        │ Kubernetes   │
│ Code       │ CLI        │ AgentCore  │ EKS        │ / Helm       │
├────────────┼────────────┼────────────┼────────────┼──────────────┤
│ .claude/   │ plugins/   │ cdk/       │ eks/       │ helm/        │
│ agents/    │ kiro/      │            │            │              │
└────────────┴────────────┴────────────┴────────────┴──────────────┘
```

## Directory Structure

Multi-agent-spec definitions should be organized in a dedicated directory (default: `specs/`) to keep them separate from other project files:

```
specs/                          # or custom directory name
├── agents/                     # Agent definitions
│   ├── orchestrator.md
│   ├── researcher.md
│   └── writer.md
├── teams/                      # Team definitions
│   └── my-team.json
└── deployments/                # Deployment configurations
    └── my-team.json
```

**Benefits:**

- Clean separation from other project files
- Tooling can process the entire directory: `genagents --spec-dir=specs/`
- Portable across repositories

**Directory conventions:**

| Directory | Contents | Format |
|-----------|----------|--------|
| `specs/agents/` | Agent definitions | Markdown with YAML frontmatter |
| `specs/teams/` | Team/workflow definitions | JSON (per `team.schema.json`) |
| `specs/deployments/` | Deployment targets | JSON (per `deployment.schema.json`) |

For simple projects, files can also be placed at the repository root:

```
agents/
team.json
deployment.json
```

## Nested Agent Directories (Optional)

For projects with many agents, you can organize them into subdirectories by namespace:

```
specs/
├── agents/
│   ├── shared/                 # Shared agents (namespace: "shared")
│   │   ├── review-board.md     # → shared/review-board
│   │   └── scoring.md          # → shared/scoring
│   ├── mrd/                    # MRD agents (namespace: "mrd")
│   │   └── market-sizing.md    # → mrd/market-sizing
│   ├── prd/                    # PRD agents (namespace: "prd")
│   │   ├── lead.md             # → prd/lead
│   │   └── requirements.md     # → prd/requirements
│   └── orchestrator.md         # Root-level agent (no namespace)
├── teams/
│   └── requirements-team.json
└── deployments/
    └── requirements-team.json
```

**Key points:**

- Subdirectory name becomes the agent's namespace automatically
- Root-level agents have no namespace (fully backward compatible)
- Team definitions reference agents by qualified name: `"namespace/agent-name"`
- Explicit `namespace` in frontmatter overrides the directory-derived namespace

### Referencing Namespaced Agents in Teams

```json
{
  "name": "requirements-team",
  "version": "1.0.0",
  "agents": [
    "orchestrator",
    "prd/lead",
    "prd/requirements",
    "shared/review-board"
  ],
  "orchestrator": "prd/lead",
  "workflow": {
    "type": "graph",
    "steps": [
      {
        "name": "discovery",
        "agent": "prd/lead"
      },
      {
        "name": "review",
        "agent": "shared/review-board",
        "depends_on": ["discovery"]
      }
    ]
  }
}
```

### Explicit Namespace Override

You can override the directory-derived namespace in the agent's frontmatter:

```markdown
---
name: special-agent
namespace: custom
description: Agent with explicit namespace
---
```

This agent will be referenced as `custom/special-agent` regardless of which subdirectory it resides in.

## Schemas

### Agent Definition Schema

Defines individual agents with capabilities and instructions.

- **Schema**: [`schema/agent/agent.schema.json`](schema/agent/agent.schema.json)
- **Format**: Hugo-compatible Markdown with YAML front matter

```markdown
---
name: my-agent
description: A helpful agent
model: sonnet
tools: [WebSearch, Read, Write]
---

You are a helpful agent...
```

### Team Schema

Defines agent teams with orchestration patterns.

- **Schema**: [`schema/orchestration/team.schema.json`](schema/orchestration/team.schema.json)

```json
{
  "name": "my-team",
  "version": "1.0.0",
  "agents": ["orchestrator", "researcher", "writer"],
  "orchestrator": "orchestrator",
  "workflow": {
    "type": "crew"
  }
}
```

### Deployment Schema

Defines target platforms and configurations.

- **Schema**: [`schema/deployment/deployment.schema.json`](schema/deployment/deployment.schema.json)

### Skill Schema (v0.8.0+)

Defines reusable capabilities that agents can invoke.

- **Schema**: [`schema/skill/skill.schema.json`](schema/skill/skill.schema.json)
- **Format**: Markdown with YAML front matter

```markdown
---
name: version-analysis
description: Analyze git history for semantic versioning
model: haiku
triggers: [version, semver]
dependencies: [git]
tools: [Bash, Read]
---

# Version Analysis

Analyze the git commit history to determine the next semantic version.
```

### Message Schema (v0.6.0+)

Defines inter-agent messages for self-directed workflows.

- **Schema**: [`schema/message/message.schema.json`](schema/message/message.schema.json)

```json
{
  "team": "my-team",
  "targets": [
    {
      "name": "local-claude",
      "platform": "claude-code",
      "priority": "p1",
      "output": ".claude/agents"
    },
    {
      "name": "aws-production",
      "platform": "aws-agentcore",
      "priority": "p1",
      "output": "cdk/",
      "config": {
        "region": "us-east-1",
        "iac": "cdk"
      }
    }
  ]
}
```

## Workflow Types

Multi-agent-spec supports two categories of workflow patterns:

### Workflow Categories

| Category | Description | Control |
|----------|-------------|---------|
| **Deterministic** | Execution paths defined in the schema | Schema controls agent execution |
| **Self-Directed** | Agents decide execution paths at runtime | Agents control their own coordination |

### Deterministic Workflows

Schema-controlled execution patterns where the workflow definition determines the order and dependencies:

| Type | Pattern | Description |
|------|---------|-------------|
| `chain` | A → B → C | Sequential execution, each step waits for the previous |
| `scatter` | A → [B,C,D] → E | Parallel fan-out with fan-in aggregation |
| `graph` | DAG | Directed acyclic graph with explicit dependencies |

### Self-Directed Workflows

Agent-controlled execution patterns where agents coordinate autonomously (requires v0.8.0+):

| Type | Pattern | Description |
|------|---------|-------------|
| `crew` | Lead + Specialists | Lead agent delegates tasks to specialist agents |
| `swarm` | Shared Queue | Self-organizing agents with shared task queue |
| `council` | Peer Debate | Peer agents debate and reach consensus |

### Usage Example

```json
{
  "name": "release-team",
  "version": "1.0.0",
  "agents": ["version-analyzer", "changelog-writer", "release-publisher"],
  "workflow": {
    "type": "graph",
    "steps": [
      {
        "name": "analyze",
        "agent": "version-analyzer"
      },
      {
        "name": "changelog",
        "agent": "changelog-writer",
        "depends_on": ["analyze"]
      },
      {
        "name": "publish",
        "agent": "release-publisher",
        "depends_on": ["changelog"]
      }
    ]
  }
}
```

### Choosing a Workflow Type

| Use Case | Recommended Type |
|----------|-----------------|
| Simple pipeline with ordered steps | `chain` |
| Independent tasks that can run concurrently | `scatter` |
| Complex dependencies between steps | `graph` |
| Dynamic task delegation by a lead agent | `crew` |
| Autonomous agent collaboration | `swarm` |
| Multi-perspective decision making | `council` |

## Self-Directed Workflows (v0.8.0+)

Self-directed workflows enable agents to coordinate autonomously at runtime. This requires additional configuration for agent roles, collaboration patterns, and inter-agent communication.

### Agent Role Configuration

For self-directed workflows, agents should define role-based fields:

```markdown
---
name: security-reviewer
role: Security Analyst
goal: Identify security vulnerabilities and recommend mitigations
backstory: |
  Senior security engineer with 10 years of experience in
  application security, penetration testing, and secure code review.
delegation:
  allow_delegation: false
model: sonnet
tools: [Read, Grep, Glob, Bash]
---

# Security Review Instructions

When reviewing code:
1. Check for injection vulnerabilities (SQL, XSS, command injection)
2. Validate authentication and authorization logic
3. Look for hardcoded secrets or credentials

Share findings via the 'findings' channel.
Challenge other reviewers if you see security implications they missed.
```

### Delegation Configuration

Control which agents can delegate work to others:

```yaml
delegation:
  allow_delegation: true           # Can this agent delegate?
  can_delegate_to:                 # Restrict delegation targets (empty = any)
    - researcher
    - writer
  can_receive_from:                # Restrict delegation sources (empty = any)
    - lead
```

### Crew Workflow Example

A lead agent coordinates specialist agents:

```json
{
  "name": "code-review-crew",
  "version": "1.0.0",
  "agents": ["lead-reviewer", "security-reviewer", "performance-reviewer"],
  "workflow": {
    "type": "crew"
  },
  "collaboration": {
    "lead": "lead-reviewer",
    "specialists": ["security-reviewer", "performance-reviewer"]
  },
  "plan_approval": true
}
```

### Swarm Workflow Example

Self-organizing agents with a shared task queue:

```json
{
  "name": "research-swarm",
  "version": "1.0.0",
  "agents": ["researcher-1", "researcher-2", "researcher-3", "synthesizer"],
  "workflow": {
    "type": "swarm"
  },
  "collaboration": {
    "task_queue": true,
    "channels": [
      {
        "name": "discoveries",
        "type": "broadcast",
        "participants": ["*"]
      }
    ]
  },
  "self_claim": true
}
```

### Council Workflow Example

Peer agents debate and reach consensus:

```json
{
  "name": "code-review-council",
  "version": "1.0.0",
  "agents": ["security-reviewer", "performance-reviewer", "test-reviewer"],
  "workflow": {
    "type": "council"
  },
  "collaboration": {
    "consensus": {
      "required_agreement": 0.66,
      "max_rounds": 3,
      "tie_breaker": "security-reviewer"
    },
    "channels": [
      {
        "name": "findings",
        "type": "broadcast",
        "participants": ["*"]
      },
      {
        "name": "challenges",
        "type": "direct",
        "participants": ["*"]
      }
    ]
  }
}
```

### Collaboration Configuration

| Field | Description | Used By |
|-------|-------------|---------|
| `lead` | Lead agent name | `crew` |
| `specialists` | Non-delegating specialist agent names | `crew` |
| `task_queue` | Enable shared task queue | `swarm` |
| `consensus.required_agreement` | Fraction of agents that must agree (0.0-1.0) | `council` |
| `consensus.max_rounds` | Maximum debate rounds | `council` |
| `consensus.tie_breaker` | Agent to break ties | `council` |
| `channels` | Communication channels between agents | All |

### Channel Types

| Type | Description | Use Case |
|------|-------------|----------|
| `direct` | Point-to-point messaging | Private communication between two agents |
| `broadcast` | One-to-all messaging | Sharing findings with all agents |
| `pub-sub` | Topic-based messaging | Agents subscribe to relevant topics |

### Message Types

Inter-agent messages support the following types:

| Type | Description | Workflow |
|------|-------------|----------|
| `delegate_work` | Assign task to another agent | `crew` |
| `ask_question` | Request information | All |
| `share_finding` | Broadcast a discovery | All |
| `request_approval` | Request approval for action | `crew` |
| `approval` / `rejection` | Respond to approval request | `crew` |
| `challenge` | Challenge another agent's finding | `council` |
| `vote` | Cast a vote on a proposal | `council` |
| `task_claimed` / `task_completed` | Task queue updates | `swarm` |
| `shutdown_request` / `shutdown_approved` | Graceful shutdown | All |

### Platform Mappings

#### Claude Code Agent Teams

| Team Config | Claude Code Config |
|-------------|-------------------|
| `workflow.type: crew` | `team_mode: team`, lead spawns teammates |
| `workflow.type: swarm` | `team_mode: team`, `self_claim: true` |
| `workflow.type: council` | `team_mode: team`, broadcast + challenge |
| `collaboration.lead` | Lead session |
| `plan_approval: true` | `--require-plan-approval` |

**Deployment example:**

```json
{
  "name": "claude-code-teams",
  "platform": "claude-code",
  "mode": "multi-process",
  "claudeCode": {
    "agentDir": ".claude/agents",
    "format": "markdown",
    "team_mode": "team",
    "teammate_mode": "in-process",
    "enable_teams": true
  }
}
```

#### CrewAI

| Team Config | CrewAI Config |
|-------------|---------------|
| `workflow.type: crew` | `process: hierarchical` |
| `workflow.type: swarm` | `process: consensual` |
| `workflow.type: council` | `process: consensual` + voting |
| `collaboration.lead` | Manager agent |
| `agent.delegation.allow_delegation` | `allow_delegation=True` |

**Deployment example:**

```json
{
  "name": "crewai-local",
  "platform": "crewai",
  "mode": "single-process",
  "crewai": {
    "model": "claude-3-5-sonnet",
    "processType": "hierarchical",
    "memory": true,
    "allowDelegation": true,
    "managerLlm": "claude-3-5-sonnet"
  }
}
```

## Supported Platforms

### P1 - Primary Targets (CLI Assistants)

| Platform | Description | Mode | Output Format |
|----------|-------------|------|---------------|
| `claude-code` | Claude Code CLI sub-agents | `single-process` | Markdown |
| `gemini-cli` | Google Gemini CLI Assistant | `single-process` | Config |
| `kiro-cli` | Kiro CLI sub-agents | `single-process` | JSON |

### P2 - Agent Frameworks

| Platform | Description | Mode | Output Format |
|----------|-------------|------|---------------|
| `adk-go` | Google Agent Development Kit (Go) | `distributed` | Go |
| `crewai` | CrewAI multi-agent framework | `single-process` | Python |
| `autogen` | Microsoft AutoGen framework | `single-process` | Python |
| `aws-agentcore` | AWS Bedrock AgentCore | `serverless` | CDK/Pulumi |

### P3 - Container Orchestration

| Platform | Description | Mode | Output Format |
|----------|-------------|------|---------------|
| `kubernetes` | Generic Kubernetes | `distributed` | Helm |
| `aws-eks` | AWS Elastic Kubernetes Service | `distributed` | Helm |
| `azure-aks` | Azure Kubernetes Service | `distributed` | Helm |
| `gcp-gke` | Google Kubernetes Engine | `distributed` | Helm |
| `docker-compose` | Local Docker deployment | `multi-process` | YAML |

## Deployment Modes

Multi-agent-spec supports different deployment modes to match the runtime characteristics of each platform:

| Mode | Description | Use Case |
|------|-------------|----------|
| `single-process` | All agents run in one process | CLI assistants (Claude Code, Gemini CLI) |
| `multi-process` | Agents run as separate local processes | Local development with isolation |
| `distributed` | Agents run on separate servers/containers | Production microservices (K8s, ADK) |
| `serverless` | Agents run as serverless functions | AWS Lambda, Cloud Functions |

### Mode Selection by Platform

| Platform | Recommended Mode | Runtime Config |
|----------|------------------|----------------|
| `claude-code` | `single-process` | Not needed |
| `gemini-cli` | `single-process` | Not needed |
| `kiro-cli` | `single-process` | Not needed |
| `crewai` | `single-process` | Optional (memory, iterations) |
| `autogen` | `single-process` | Optional (human input mode) |
| `adk-go` | `distributed` | Recommended (retry, timeout, observability) |
| `aws-agentcore` | `serverless` | Recommended (timeout, resources) |
| `kubernetes` | `distributed` | Required (resources, retry, observability) |

## Runtime Configuration

Runtime configuration is specified in the **deployment schema**, not the team schema. This separation allows the same team definition to work across different platforms with platform-appropriate runtime settings.

### Schema Separation

```
team.schema.json (Logical)          deployment.schema.json (Runtime)
├── agents                          ├── platform
├── workflow                        ├── mode
│   ├── steps                       └── runtime
│   │   ├── depends_on (DAG)            ├── defaults
│   │   ├── inputs (data flow)          │   ├── timeout
│   │   └── outputs (data flow)         │   ├── retry
└── context                             │   └── resources
                                        ├── steps (per-step overrides)
                                        └── observability
```

### Runtime Settings

| Setting | Description | Example |
|---------|-------------|---------|
| `timeout` | Step execution limit | `"5m"`, `"1h"` |
| `retry.max_attempts` | Max retry attempts | `3` |
| `retry.backoff` | Backoff strategy | `"exponential"` |
| `condition` | Conditional execution | `"inputs.ready == true"` |
| `concurrency` | Max parallel executions | `2` |
| `resources.cpu` | CPU limit | `"500m"`, `"2"` |
| `resources.memory` | Memory limit | `"512Mi"`, `"2Gi"` |

### Example: Single-Process Deployment (No Runtime Config)

```json
{
  "name": "local-claude",
  "platform": "claude-code",
  "mode": "single-process",
  "output": ".claude/agents"
}
```

### Example: Distributed Deployment (Full Runtime Config)

```json
{
  "name": "k8s-production",
  "platform": "kubernetes",
  "mode": "distributed",
  "runtime": {
    "defaults": {
      "timeout": "5m",
      "retry": {
        "max_attempts": 3,
        "backoff": "exponential"
      }
    },
    "steps": {
      "qa-validation": {
        "timeout": "15m",
        "resources": { "cpu": "2", "memory": "2Gi" }
      }
    },
    "observability": {
      "tracing": { "enabled": true, "exporter": "otlp" },
      "metrics": { "enabled": true, "exporter": "prometheus" }
    }
  }
}
```

## JSON Schema Guidelines for Go Compatibility

The multi-agent-spec schemas are designed to be compatible with Go code generation. When modifying or extending the schemas, follow these guidelines:

### Avoid Problematic Patterns

| Pattern | Problem | Alternative |
|---------|---------|-------------|
| `anyOf` without discriminator | Generates `interface{}` | Add `const` discriminator field |
| `oneOf` without discriminator | Generates `interface{}` | Add `const` discriminator field |
| Nested unions >2 levels | Complex unmarshaling | Flatten hierarchy |
| Large unions (>10 variants) | Unwieldy switch statements | Split into smaller unions |

### Union Types with Discriminators

If using `anyOf` or `oneOf`, add a discriminator field with unique `const` values:

```json
{
  "oneOf": [
    {
      "type": "object",
      "properties": {
        "platform": { "const": "claude-code" },
        "agentDir": { "type": "string" }
      }
    },
    {
      "type": "object",
      "properties": {
        "platform": { "const": "kubernetes" },
        "namespace": { "type": "string" }
      }
    }
  ]
}
```

### Validating Schemas

Use [`schemago`](https://github.com/grokify/schemago) to check schemas for Go compatibility:

```bash
schemago lint schema/agent/agent.schema.json
schemago lint schema/orchestration/team.schema.json
schemago lint schema/deployment/deployment.schema.json
```

### Preferred Patterns

- Use `$ref` to reference definitions (schemago handles these correctly)
- Keep union variants to `$ref` references when possible
- Use nullable pattern `anyOf [T, null]` for optional types
- Prefer simple string enums over complex union types

## Model Mappings

Canonical model names map to platform-specific identifiers:

| Canonical | Claude Code | Kiro CLI | AWS Bedrock |
|-----------|-------------|----------|-------------|
| `haiku` | `haiku` | `claude-haiku-35` | `anthropic.claude-3-haiku-*` |
| `sonnet` | `sonnet` | `claude-sonnet-4` | `anthropic.claude-3-sonnet-*` |
| `opus` | `opus` | `claude-opus-4` | `anthropic.claude-3-opus-*` |

## Tool Mappings

Canonical tool names map to platform-specific identifiers:

| Canonical | Claude Code | Kiro CLI | Description |
|-----------|-------------|----------|-------------|
| `WebSearch` | `WebSearch` | `web_search` | Search the web |
| `WebFetch` | `WebFetch` | `web_fetch` | Fetch web pages |
| `Read` | `Read` | `read` | Read files |
| `Write` | `Write` | `write` | Write files |
| `Glob` | `Glob` | `glob` | Find files by pattern |
| `Grep` | `Grep` | `grep` | Search file contents |
| `Bash` | `Bash` | `bash` | Execute commands |
| `Edit` | `Edit` | `edit` | Edit files |
| `Task` | `Task` | `task` | Spawn sub-agents |

## Installation

### Go SDK

```bash
go get github.com/plexusone/multi-agent-spec@v0.8.0
```

Import the SDK package:

```go
import multiagentspec "github.com/plexusone/multi-agent-spec/sdk/go"
```

### Python SDK

```bash
pip install multi-agent-spec
```

### TypeScript SDK

```bash
npm install @plexusone/multi-agent-spec
```

## Usage with aiassistkit

Generate platform-specific agents using `genagents`:

```bash
# Generate from specs/ directory (recommended)
genagents --spec-dir=specs/ --output=.claude/agents --format=claude

# Generate for multiple targets
genagents --spec-dir=specs/ \
  --targets="claude:.claude/agents,kiro:plugins/kiro/agents"

# Process a custom directory
genagents --spec-dir=my-agents/ --output=.claude/agents --format=claude

# Verbose output
genagents --spec-dir=specs/ --targets="..." --verbose
```

## Examples

See the [`examples/`](examples/) directory for complete examples:

- [`stats-agent-team/`](examples/stats-agent-team/) - Statistics research and verification team

## Related Projects

- [aiassistkit](https://github.com/plexusone/aiassistkit) - Agent generation and deployment tooling
- [agentkit](https://github.com/plexusone/agentkit) - Multi-platform agent runtime
- [stats-agent-team](https://github.com/plexusone/stats-agent-team) - Reference implementation

## License

MIT

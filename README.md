# Multi-Agent Spec

A specification for defining multi-agent AI systems with platform-agnostic agent definitions and deployment configurations.

## Overview

Multi-Agent Spec provides a standardized way to define:

1. **Agents** - Individual AI agents with capabilities, tools, and instructions
2. **Teams** - Groups of agents with orchestration patterns
3. **Deployments** - Target platforms and configurations

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     Definition Layer                             │
├─────────────────────────────────────────────────────────────────┤
│  agents/*.md          │  team.json           │  deployment.json  │
│  (Markdown + YAML)    │  (Orchestration)     │  (Targets)        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Deployment Layer                             │
├────────────┬────────────┬────────────┬────────────┬─────────────┤
│ Claude     │ Kiro       │ AWS        │ AWS        │ Kubernetes  │
│ Code       │ CLI        │ AgentCore  │ EKS        │ / Helm      │
├────────────┼────────────┼────────────┼────────────┼─────────────┤
│ .claude/   │ plugins/   │ cdk/       │ eks/       │ helm/       │
│ agents/    │ kiro/      │            │            │             │
└────────────┴────────────┴────────────┴────────────┴─────────────┘
```

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
    "type": "orchestrated"
  }
}
```

### Deployment Schema

Defines target platforms and configurations.

- **Schema**: [`schema/deployment/deployment.schema.json`](schema/deployment/deployment.schema.json)

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

## Supported Platforms

### P1 - Primary Targets

| Platform | Description | Output Format |
|----------|-------------|---------------|
| `claude-code` | Claude Code CLI sub-agents | Markdown |
| `kiro-cli` | Kiro CLI sub-agents | JSON |
| `aws-agentcore` | AWS Bedrock AgentCore | CDK/Pulumi |

### P2 - Secondary Targets

| Platform | Description | Output Format |
|----------|-------------|---------------|
| `aws-eks` | AWS Elastic Kubernetes Service | Helm |
| `azure-aks` | Azure Kubernetes Service | Helm |
| `gcp-gke` | Google Kubernetes Engine | Helm |
| `kubernetes` | Generic Kubernetes | Helm |
| `docker-compose` | Local Docker deployment | YAML |

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

## Usage with aiassistkit

Generate platform-specific agents using `genagents`:

```bash
# Generate for Claude Code
genagents -spec=plugins/spec/agents -output=.claude/agents -format=claude

# Generate for multiple targets
genagents -spec=plugins/spec/agents \
  -targets="claude:.claude/agents,kiro:plugins/kiro/agents"

# Verbose output
genagents -spec=plugins/spec/agents -targets="..." -verbose
```

## Examples

See the [`examples/`](examples/) directory for complete examples:

- [`stats-agent-team/`](examples/stats-agent-team/) - Statistics research and verification team

## Related Projects

- [aiassistkit](https://github.com/agentplexus/aiassistkit) - Agent generation and deployment tooling
- [agentkit](https://github.com/agentplexus/agentkit) - Multi-platform agent runtime
- [stats-agent-team](https://github.com/agentplexus/stats-agent-team) - Reference implementation

## License

MIT

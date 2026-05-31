# Multi-Agent Spec

A portable specification format for multi-agent AI systems. Define agents, teams, and deployments once — deploy to Claude Code, Kiro CLI, AWS AgentCore, CrewAI, and more.

## Why Multi-Agent Spec?

AI coding assistants and agent frameworks each have their own configuration formats. Multi-Agent Spec provides a unified way to define:

- **Agents** — Individual AI agents with roles, tools, and tasks
- **Teams** — Multi-agent workflows with DAG-based orchestration
- **Deployments** — Platform-specific deployment configurations
- **Reports** — Structured output from agent executions

## Key Features

- **Write once, deploy anywhere** — Define agents in a portable format, generate platform-native configs
- **DAG-based orchestration** — Express complex workflows with dependencies
- **Go/No-Go semantics** — Built-in status tracking inspired by NASA mission control
- **Rich reporting** — Structured reports with severity, verdicts, and aggregation tags
- **Schema-validated** — JSON Schema for all formats with Go and TypeScript SDKs

## Supported Platforms

| Platform | Type | Status |
|----------|------|--------|
| Claude Code | AI Coding Assistant | Supported |
| Kiro CLI | AI Coding Assistant | Supported |
| Gemini CLI | AI Coding Assistant | Planned |
| AWS AgentCore | Cloud Agent Framework | Supported |
| CrewAI | Python Agent Framework | Supported |
| AutoGen | Python Agent Framework | Supported |
| Kubernetes | Container Orchestration | Supported |

## Quick Example

Define an agent:

```yaml
# agents/reviewer.md
---
name: reviewer
description: Reviews code changes for quality and security
model: sonnet
tools:
  - Read
  - Grep
  - Glob
tasks:
  - id: security-check
    description: Check for security vulnerabilities
    type: pattern
    pattern: "(password|secret|api_key)\\s*="
---

You are a code reviewer specializing in security analysis.
```

Define a team workflow:

```json
{
  "name": "review-team",
  "workflow": {
    "steps": [
      {
        "name": "static-analysis",
        "agent": "reviewer",
        "depends_on": []
      },
      {
        "name": "security-audit",
        "agent": "security",
        "depends_on": ["static-analysis"]
      }
    ]
  }
}
```

Deploy to multiple platforms:

```json
{
  "team": "review-team",
  "targets": [
    {
      "name": "local-claude",
      "platform": "claude-code",
      "output": ".claude/agents"
    },
    {
      "name": "local-kiro",
      "platform": "kiro-cli",
      "output": "~/.kiro/agents",
      "kiroCli": {
        "prefix": "rev_"
      }
    }
  ]
}
```

## Getting Started

- [Installation](getting-started/installation.md) — Install the SDK and CLI
- [Quick Start](getting-started/quickstart.md) — Create your first multi-agent team

## Documentation

- [Schemas](schemas/overview.md) — JSON Schema reference
- [CLI](cli/mas.md) — Command-line interface
- [SDK](sdk/go.md) — Go SDK reference

## Guides

- [GitHub PR Review](guides/github-pr-review.md) — Review PRs with Claude Code and post comments

# Loop Schema

The loop schema defines autonomous loop patterns for multi-agent systems.

**Schema**: [`schema/orchestration/loop.schema.json`](https://github.com/plexusone/multi-agent-spec/blob/main/schema/orchestration/loop.schema.json)

## Overview

Loops are bounded execution patterns that run autonomously until completion or escalation. They are inspired by the REPL (Read-Eval-Print-Loop) pattern from interactive programming.

## Loop Types

| Type | Category | Description |
|------|----------|-------------|
| `REAL` | Mission-driven | Read requirements, Eval situation, Act toward goal, Loop |
| `VEAL` | State-driven | Validate state, Eval findings, Act to correct, Loop |

## Schema Definition

```yaml
name: qa-fix
type: VEAL
description: QA validation and fix loop
validator: qa
actor: code-fixer
max_attempts: 3
escalation: human
success_criteria: All checks pass with GO status
checks:
  - id: build
    type: command
    command: go build ./...
    required: true
    expected: Exit code 0
  - id: lint
    type: command
    command: golangci-lint run
    required: true
```

## Fields

### Common Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Loop identifier (kebab-case) |
| `type` | string | Yes | `REAL` or `VEAL` |
| `description` | string | No | Human-readable description |
| `actor` | string | Yes | Agent that performs actions |
| `max_attempts` | integer | No | Maximum iterations (default: 3) |
| `escalation` | string | No | Policy when max reached (default: `human`) |
| `success_criteria` | string | No | Description of success conditions |
| `inputs` | array | No | Data inputs required to start the loop |
| `outputs` | array | No | Data outputs produced when loop completes |

### VEAL-Specific Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `validator` | string | Yes | Agent that validates state (read-only) |
| `checks` | array | Yes | Validation checks to run |

### REAL-Specific Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `mission` | string | Yes | Goal statement for the loop |

## Check Definition

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Check identifier (kebab-case) |
| `type` | string | No | `command`, `pattern`, `file`, or `manual` |
| `command` | string | No | Shell command (for type: command) |
| `pattern` | string | No | Regex pattern (for type: pattern) |
| `files` | string | No | Glob pattern for files to check |
| `file` | string | No | File path (for type: file) |
| `required` | boolean | No | Causes NO-GO if failed (default: true) |
| `expected` | string | No | Description of success |

## Escalation Policies

| Policy | Description |
|--------|-------------|
| `human` | Stop and request human intervention |
| `abort` | Stop and fail the workflow |
| `continue` | Proceed despite unresolved issues |
| `fallback` | Invoke a fallback agent |

## Examples

### VEAL Loop (State-Driven)

```yaml
name: qa-fix
type: VEAL
description: QA validation and fix loop
validator: qa
actor: code-fixer
max_attempts: 3
escalation: human
checks:
  - id: build
    type: command
    command: go build ./...
    required: true
  - id: tests
    type: command
    command: go test -v ./...
    required: true
  - id: lint
    type: command
    command: golangci-lint run
    required: true
```

### REAL Loop (Mission-Driven)

```yaml
name: feature-impl
type: REAL
description: Feature implementation loop
actor: developer
mission: |
  Implement user authentication with OAuth2 support.
  Must include login, logout, and session management.
max_attempts: 5
escalation: human
success_criteria: |
  All authentication flows work correctly.
  Tests pass. Documentation updated.
```

## Using Loops in Workflows

Loops can be referenced as steps in team workflows:

```json
{
  "workflow": {
    "type": "chain",
    "steps": [
      { "name": "validate-scope", "agent": "pm" },
      { "name": "qa-fix", "loop": "qa-fix" },
      { "name": "docs-fix", "loop": "docs-fix" },
      { "name": "tag-release", "agent": "release" }
    ]
  }
}
```

## Go SDK Usage

```go
import mas "github.com/plexusone/multi-agent-spec/sdk/go"

// Create a VEAL loop
qaLoop := mas.NewVEALLoop("qa-fix", "qa", "code-fixer").
    WithDescription("QA validation and fix loop").
    WithMaxAttempts(3).
    AddCheck(mas.LoopCheck{
        ID:      "build",
        Type:    mas.CheckTypeCommand,
        Command: "go build ./...",
    })

// Create a REAL loop
featureLoop := mas.NewREALLoop(
    "feature-impl",
    "developer",
    "Implement user authentication",
)

// Validate
if err := qaLoop.Validate(); err != nil {
    log.Fatal(err)
}

// Load from file
loop, err := mas.LoadLoopFromFile("specs/loops/qa-fix.yaml")

// Load from directory
loops, err := mas.LoadLoopsFromDir("specs/loops")
```

## Related

- [Loop Pattern Guide](../guides/loops.md) - Detailed guide on REAL/VEAL patterns
- [Team Schema](team.md) - Using loops in workflow steps

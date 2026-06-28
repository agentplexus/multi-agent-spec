# REAL/VEAL Loop Patterns

This guide describes autonomous loop patterns for multi-agent systems.

## Overview

Loops are autonomous, bounded execution patterns that run without human intervention until they complete or escalate. They are inspired by the REPL (Read-Eval-Print-Loop) pattern from interactive programming environments.

## Loop Hierarchy

```
REAL (abstract)
  │
  └── VEAL (state-driven)
        ├── qa-fix-loop      (code state → valid code)
        ├── docs-fix-loop    (docs state → complete docs)
        └── security-fix-loop (security state → secure code)
```

## REAL — Read Eval Act Loop

The generic, mission-driven loop pattern. Use for open-ended tasks where the goal is defined by a mission statement rather than specific validation checks.

```
REAL Loop:
  R - Read (mission, requirements, goals)
  E - Eval (assess situation against goals)
  A - Act (take action toward goals)
  L - Loop (until mission complete or escalate)
```

**Characteristics:**

- Open-ended goals
- Requirements-driven
- May involve discovery/exploration
- Single actor agent (no separate validator)

**Example use cases:**

- "Implement feature X"
- "Migrate from API v1 to v2"
- "Refactor authentication system"

**Example definition:**

```yaml
name: feature-impl
type: REAL
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

## VEAL — Validate Eval Act Loop

A specialization of REAL for **state convergence**. Validates current state against expected state and acts to correct deviations.

```
VEAL Loop:
  V - Validate (check current state)
  E - Eval (compare to expected state → GO/NO-GO)
  A - Act (correct state if NO-GO)
  L - Loop (until state valid or escalate)
```

**Characteristics:**

- Convergent (state approaches valid)
- Idempotent checks
- Bounded attempts
- Dual-agent pattern (validator + fixer)

**Example use cases:**

- QA validation (build, test, lint, format)
- Documentation completeness
- Security compliance
- Configuration validation

**Example definition:**

```yaml
name: qa-fix
type: VEAL
description: QA validation and fix loop
validator: qa
actor: code-fixer
max_attempts: 3
escalation: human
success_criteria: All checks pass (GO status)
checks:
  - id: build
    type: command
    command: go build ./...
    required: true
    expected: Exit code 0

  - id: tests
    type: command
    command: go test -v ./...
    required: true
    expected: All tests pass

  - id: lint
    type: command
    command: golangci-lint run
    required: true
    expected: No lint errors

  - id: format
    type: command
    command: gofmt -l .
    required: true
    expected: No files need formatting
```

## Loop Configuration

### Common Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Unique identifier (kebab-case) |
| `type` | string | Yes | `REAL` or `VEAL` |
| `description` | string | No | Human-readable description |
| `actor` | string | Yes | Agent that performs actions |
| `max_attempts` | int | No | Maximum iterations (default: 3) |
| `escalation` | string | No | Policy when max reached (default: `human`) |
| `success_criteria` | string | No | Description of success conditions |

### VEAL-Specific Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `validator` | string | Yes | Agent that validates state (read-only) |
| `checks` | array | Yes | Validation checks to run |

### REAL-Specific Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `mission` | string | Yes | Goal statement for the loop |

### Check Definition

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Check identifier (kebab-case) |
| `type` | string | No | `command`, `pattern`, `file`, or `manual` |
| `command` | string | No | Shell command (for type: command) |
| `pattern` | string | No | Regex pattern (for type: pattern) |
| `files` | string | No | Glob pattern for files to check |
| `file` | string | No | File path (for type: file) |
| `required` | bool | No | Causes NO-GO if failed (default: true) |
| `expected` | string | No | Description of success |

### Escalation Policies

| Policy | Description |
|--------|-------------|
| `human` | Stop and request human intervention |
| `abort` | Stop and fail the workflow |
| `continue` | Proceed despite unresolved issues |
| `fallback` | Invoke a fallback agent |

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

### Creating Loops Programmatically

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
    }).
    AddCheck(mas.LoopCheck{
        ID:      "lint",
        Type:    mas.CheckTypeCommand,
        Command: "golangci-lint run",
    })

// Create a REAL loop
featureLoop := mas.NewREALLoop(
    "feature-impl",
    "developer",
    "Implement user authentication with OAuth2",
).WithMaxAttempts(5)

// Validate
if err := qaLoop.Validate(); err != nil {
    log.Fatal(err)
}
```

### Loading Loops from Files

```go
loader := mas.NewLoader()

// Load single loop
loop, err := loader.LoadLoop("loops/qa-fix.yaml")

// Load all loops from directory
loops, err := mas.LoadLoopsFromDir("loops/")
```

## Best Practices

1. **Keep checks idempotent** — Running the same check multiple times should produce the same result
2. **Use specific checks** — Prefer multiple small checks over one large check
3. **Set reasonable max_attempts** — 3 is usually sufficient; more may indicate a deeper issue
4. **Separate concerns** — Validator agents should be read-only; actor agents should fix issues
5. **Document success criteria** — Make it clear what "done" looks like
6. **Use fallback agents** — For complex issues that may need escalation to a specialist

# Skill Schema

Defines reusable capabilities that agents can invoke. Skills encapsulate instructions, scripts, and resources for specific tasks.

## Schema URL

```
https://raw.githubusercontent.com/plexusone/multi-agent-spec/main/schema/skill/skill.schema.json
```

## Structure

```json
{
  "$schema": "...",
  "name": "string",
  "description": "string",
  "instructions": "string",
  "scripts": ["string"],
  "references": ["string"],
  "assets": ["string"],
  "triggers": ["string"],
  "dependencies": ["string"],
  "model": "haiku|sonnet|opus",
  "tools": ["string"]
}
```

## Fields

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Unique skill identifier (lowercase, hyphenated) |

### Core Fields

| Field | Type | Description |
|-------|------|-------------|
| `description` | string | Brief summary of what the skill does |
| `instructions` | string | Main prompt/instructions executed when skill is invoked |
| `model` | string | Recommended model tier: `haiku`, `sonnet`, `opus` |
| `tools` | string[] | Tools this skill requires access to |

### Resource Fields

| Field | Type | Description |
|-------|------|-------------|
| `scripts` | string[] | Paths to executable scripts relative to skill directory |
| `references` | string[] | Paths to reference documentation |
| `assets` | string[] | Paths to templates, config files, or other resources |

### Invocation Fields

| Field | Type | Description |
|-------|------|-------------|
| `triggers` | string[] | Keywords or patterns that invoke this skill |
| `dependencies` | string[] | External tools/binaries required (e.g., `git`, `npm`) |

## File Format

Skills are defined in Markdown with YAML frontmatter:

```markdown
---
name: version-analysis
description: Analyze git history for semantic versioning
model: haiku
triggers: [version, semver, changelog]
dependencies: [git]
tools: [Bash, Read]
scripts:
  - scripts/analyze.sh
references:
  - docs/semver.md
---

# Version Analysis

Analyze the git commit history to determine the next semantic version.

## Instructions

1. Look at commits since the last tag
2. Classify each commit type (feat, fix, etc.)
3. Determine major/minor/patch bump based on conventional commits
```

## Directory Structure

Skills are organized in a directory structure:

```
skills/
├── version-analysis/
│   ├── SKILL.md           # Skill definition
│   ├── scripts/           # Executable scripts
│   │   └── analyze.sh
│   ├── references/        # Reference documentation
│   │   └── semver.md
│   └── assets/            # Templates, configs
│       └── categories.json
├── commit-classification/
│   └── SKILL.md
└── changelog/
    └── SKILL.md
```

## Example: Commit Classification Skill

```markdown
---
name: commit-classification
description: Classify commits according to conventional commits specification
model: haiku
triggers: [classify, commit, conventional]
dependencies: [git]
tools: [Bash, Read]
---

# Commit Classification

Classify git commits using the conventional commits specification.

## Instructions

Analyze each commit message and classify it into one of:

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **style**: Code style changes (formatting, semicolons)
- **refactor**: Code refactoring
- **perf**: Performance improvements
- **test**: Adding or updating tests
- **build**: Build system changes
- **ci**: CI configuration changes
- **chore**: Other changes (dependencies, configs)

## Output Format

Return a JSON array of classified commits:

```json
[
  {"hash": "abc123", "type": "feat", "scope": "auth", "subject": "add OAuth support"},
  {"hash": "def456", "type": "fix", "scope": null, "subject": "resolve race condition"}
]
```
```

## Using Skills in Agents

Agents reference skills by name in their `skills` field:

```yaml
---
name: release-agent
description: Manages releases
skills:
  - version-analysis
  - commit-classification
  - changelog
---
```

When an agent encounters a trigger keyword (e.g., "version" or "classify"), it should consider invoking the matching skill.

## Go SDK Usage

```go
import multiagentspec "github.com/plexusone/multi-agent-spec/sdk/go"

// Create a skill programmatically
skill := multiagentspec.NewSkill("version-analysis", "Analyze versions").
    WithInstructions("Analyze git history...").
    WithModel(multiagentspec.ModelHaiku).
    WithTools("Bash", "Read").
    AddTrigger("version").
    AddTrigger("semver").
    AddDependency("git").
    AddScript("scripts/analyze.sh")

// Load skills from directory
skills, err := multiagentspec.LoadSkillsFromDir("skills/")

// Load into a SkillSet for querying
skillSet, err := multiagentspec.LoadSkillSetFromDir("skills/")
matches := skillSet.FindByTrigger("version")
```

## Platform Mappings

Skills are mapped to platform-specific formats:

| Platform | Skill Format |
|----------|--------------|
| Claude Code | `skills/<name>/SKILL.md` |
| Kiro CLI | `steering/<name>.md` |
| OpenAI Codex | `skills/<name>/SKILL.md` |

## Best Practices

1. **Single Responsibility**: Each skill should do one thing well
2. **Clear Triggers**: Use specific, unambiguous trigger keywords
3. **Document Dependencies**: List all required external tools
4. **Include Examples**: Show expected input/output in instructions
5. **Reference Documentation**: Link to external docs for complex topics

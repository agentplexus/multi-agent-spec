# CLAUDE.md

Project-specific instructions for Claude Code working on multi-agent-spec.

## Project Overview

Multi-Agent Spec is a specification for defining multi-agent AI systems with platform-agnostic agent definitions and deployment configurations.

## Architecture

```
multi-agent-spec/
├── schema/                    # JSON Schema files (generated from Go types)
│   ├── agent/
│   ├── deployment/
│   ├── message/
│   ├── orchestration/
│   ├── report/
│   └── skill/
├── sdk/
│   ├── go/                    # Go SDK (source of truth for types)
│   ├── typescript/            # TypeScript SDK (generated from JSON Schema)
│   └── python/                # Python SDK
├── tools/
│   └── generate/              # Schema generation tool
├── cmd/
│   └── mas/                   # CLI tool
└── docs/                      # MkDocs documentation
```

## Go-First Schema Generation

**Go types are the source of truth.** JSON Schemas are generated from Go structs.

### Workflow

1. Define or modify Go types in `sdk/go/*.go`
2. Run schema generator:
   ```bash
   cd tools/generate && go run main.go
   ```
3. Verify generated schemas in `schema/`
4. Run TypeScript generator (if needed):
   ```bash
   cd sdk/typescript && npm run generate
   ```
5. Update Python SDK manually if needed

### Adding a New Type

1. Create Go struct in `sdk/go/` with proper JSON/YAML tags
2. Add to `tools/generate/main.go`:
   ```go
   if err := generateSchema(
       &multiagentspec.NewType{},
       filepath.Join(outputDir, "category", "newtype.schema.json"),
       "Multi-Agent Spec - NewType Definition",
       "Description of the type",
       "https://raw.githubusercontent.com/plexusone/multi-agent-spec/main/schema/category/newtype.schema.json",
   ); err != nil {
       return fmt.Errorf("generating newtype schema: %w", err)
   }
   ```
3. Run generator and commit both Go types and generated schemas together

## SDK Generation

### TypeScript SDK

Auto-generated from JSON Schema using `json-schema-to-zod`:

```bash
cd sdk/typescript
npm run generate   # Generates Zod schemas in src/generated/
npm run build      # Compiles TypeScript
```

To add a new type, update `sdk/typescript/scripts/generate-schemas.mjs`:
```javascript
const SCHEMAS = [
  // ... existing schemas
  { file: 'category/newtype.schema.json', name: 'NewType' },
];
```

### Python SDK

Can be auto-generated from JSON Schema using `datamodel-code-generator`:

```bash
cd sdk/python
pip install -e ".[dev]"
python -m multi_agent_spec.scripts.generate
```

To add a new type, update `sdk/python/src/multi_agent_spec/scripts/generate.py`:
```python
SCHEMAS = [
    # ... existing schemas
    ("category/newtype.schema.json", "newtype"),
]
```

**Note:** The main `models.py` contains manually curated models. Generated models go to `src/multi_agent_spec/generated/`.

## Release Process

### Pre-Release Checklist

1. **Run tests**: `go test ./...`
2. **Run linter**: `golangci-lint run`
3. **Regenerate schemas**: `cd tools/generate && go run main.go`
4. **Update CLI version**: `cmd/mas/cmd/root.go`
5. **Update CHANGELOG.json**: Add new version entry
6. **Generate CHANGELOG.md**: `schangelog generate CHANGELOG.json -o CHANGELOG.md`
7. **Create release notes**: `docs/releases/vX.Y.Z.md`
8. **Update mkdocs.yml**: Add release to nav

### Versioning

- Project uses semantic versioning
- Single top-level Go module: `github.com/plexusone/multi-agent-spec`
- SDK import path: `github.com/plexusone/multi-agent-spec/sdk/go`

### Tagging

```bash
git push origin main
git tag vX.Y.Z
git push origin vX.Y.Z
```

## Commit Conventions

Follow [Conventional Commits](https://www.conventionalcommits.org/):

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation |
| `refactor` | Code refactoring |
| `test` | Adding tests |
| `chore` | Maintenance |

### Scopes

- `sdk/go` - Go SDK changes
- `schema` - JSON Schema changes
- `cli` - CLI tool changes
- `docs` - Documentation changes

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test ./sdk/go/...
```

## Documentation

- Site: https://plexusone.github.io/multi-agent-spec/
- Built with MkDocs Material
- Source in `docs/`

### Building Docs Locally

```bash
pip install mkdocs-material
mkdocs serve
```

## Key Types

| Type | Description | File |
|------|-------------|------|
| Agent | Individual agent definition | `sdk/go/agent.go` |
| Team | Multi-agent workflow | `sdk/go/team.go` |
| Deployment | Platform configuration | `sdk/go/deployment.go` |
| Skill | Reusable capability | `sdk/go/skill.go` |
| TeamReport | Execution results | `sdk/go/report.go` |

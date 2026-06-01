// Package main generates JSON Schema files from Go types.
//
// Usage:
//
//	go run ./tools/generate
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/invopop/jsonschema"

	multiagentspec "github.com/plexusone/multi-agent-spec/sdk/go"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Output directory for generated schemas (relative to repo root)
	outputDir := "../../schema"

	// Generate Agent schema
	if err := generateSchema(
		&multiagentspec.Agent{},
		filepath.Join(outputDir, "agent", "agent.schema.json"),
		"Multi-Agent Spec - Agent Definition",
		"Schema for defining an AI agent in a multi-agent system",
		"https://raw.githubusercontent.com/plexusone/multi-agent-spec/main/schema/agent/agent.schema.json",
	); err != nil {
		return fmt.Errorf("generating agent schema: %w", err)
	}

	// Generate Team schema
	if err := generateSchema(
		&multiagentspec.Team{},
		filepath.Join(outputDir, "orchestration", "team.schema.json"),
		"Multi-Agent Spec - Team Definition",
		"Schema for defining a team of AI agents with orchestration",
		"https://raw.githubusercontent.com/plexusone/multi-agent-spec/main/schema/orchestration/team.schema.json",
	); err != nil {
		return fmt.Errorf("generating team schema: %w", err)
	}

	// Generate Deployment schema
	if err := generateSchema(
		&multiagentspec.Deployment{},
		filepath.Join(outputDir, "deployment", "deployment.schema.json"),
		"Multi-Agent Spec - Deployment Definition",
		"Schema for defining deployment targets for multi-agent systems",
		"https://raw.githubusercontent.com/plexusone/multi-agent-spec/main/schema/deployment/deployment.schema.json",
	); err != nil {
		return fmt.Errorf("generating deployment schema: %w", err)
	}

	// Generate TeamReport schema
	if err := generateSchema(
		&multiagentspec.TeamReport{},
		filepath.Join(outputDir, "report", "team-report.schema.json"),
		"Multi-Agent Spec - Team Report",
		"Schema for team validation reports",
		"https://raw.githubusercontent.com/plexusone/multi-agent-spec/main/schema/report/team-report.schema.json",
	); err != nil {
		return fmt.Errorf("generating team-report schema: %w", err)
	}

	// Generate Skill schema
	if err := generateSchema(
		&multiagentspec.Skill{},
		filepath.Join(outputDir, "skill", "skill.schema.json"),
		"Multi-Agent Spec - Skill Definition",
		"Schema for defining a reusable skill that agents can invoke",
		"https://raw.githubusercontent.com/plexusone/multi-agent-spec/main/schema/skill/skill.schema.json",
	); err != nil {
		return fmt.Errorf("generating skill schema: %w", err)
	}

	fmt.Println("Schema generation complete!")
	return nil
}

func generateSchema(v interface{}, outputPath, title, description, id string) error {
	// Create reflector with options
	r := &jsonschema.Reflector{
		DoNotReference:            false, // Use $ref for named types
		ExpandedStruct:            false,
		AllowAdditionalProperties: false,
	}

	// Generate schema
	schema := r.Reflect(v)
	schema.Title = title
	schema.Description = description
	schema.ID = jsonschema.ID(id)

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling schema: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	// Write file
	if err := os.WriteFile(outputPath, data, 0600); err != nil {
		return fmt.Errorf("writing file %s: %w", outputPath, err)
	}

	fmt.Printf("Generated: %s\n", outputPath)
	return nil
}

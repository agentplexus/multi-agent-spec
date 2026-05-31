package multiagentspec

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/invopop/jsonschema"
)

// TestGenerateSkillSchema generates the skill JSON schema.
// Run with: go test -run TestGenerateSkillSchema -v
func TestGenerateSkillSchema(t *testing.T) {
	if os.Getenv("GENERATE_SCHEMA") != "1" {
		t.Skip("Set GENERATE_SCHEMA=1 to generate schema")
	}

	r := jsonschema.Reflector{
		DoNotReference:             false,
		ExpandedStruct:             false,
		AllowAdditionalProperties:  false,
		RequiredFromJSONSchemaTags: true,
	}

	schema := r.Reflect(&Skill{})
	schema.ID = jsonschema.ID("https://raw.githubusercontent.com/plexusone/multi-agent-spec/main/schema/skill/skill.schema.json")
	schema.Title = "Multi-Agent Spec - Skill Definition"
	schema.Description = "Schema for defining a reusable skill that agents can invoke"

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		t.Fatalf("marshaling schema: %v", err)
	}

	// Output path relative to repo root (sdk/go is 2 levels deep)
	outputPath := filepath.Join("..", "..", "schema", "skill", "skill.schema.json")

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("creating directory: %v", err)
	}

	if err := os.WriteFile(outputPath, data, 0600); err != nil {
		t.Fatalf("writing schema: %v", err)
	}

	t.Logf("Generated: %s", outputPath)
}

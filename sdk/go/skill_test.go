package multiagentspec

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewSkill(t *testing.T) {
	skill := NewSkill("version-analysis", "Analyze git history for semantic versioning")

	if skill.Name != "version-analysis" {
		t.Errorf("expected name 'version-analysis', got '%s'", skill.Name)
	}
	if skill.Description != "Analyze git history for semantic versioning" {
		t.Errorf("expected description to be set, got '%s'", skill.Description)
	}
}

func TestSkillBuilder(t *testing.T) {
	skill := NewSkill("commit-classification", "Classify commits").
		WithInstructions("Analyze commit messages and classify them...").
		WithModel(ModelHaiku).
		WithTools("Bash", "Read").
		AddTrigger("classify").
		AddTrigger("commit").
		AddDependency("git").
		AddScript("scripts/analyze.sh").
		AddReference("docs/conventional-commits.md").
		AddAsset("templates/categories.json")

	if skill.Instructions == "" {
		t.Error("expected instructions to be set")
	}
	if skill.Model != ModelHaiku {
		t.Errorf("expected model haiku, got %s", skill.Model)
	}
	if len(skill.Tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(skill.Tools))
	}
	if len(skill.Triggers) != 2 {
		t.Errorf("expected 2 triggers, got %d", len(skill.Triggers))
	}
	if len(skill.Dependencies) != 1 {
		t.Errorf("expected 1 dependency, got %d", len(skill.Dependencies))
	}
	if len(skill.Scripts) != 1 {
		t.Errorf("expected 1 script, got %d", len(skill.Scripts))
	}
	if len(skill.References) != 1 {
		t.Errorf("expected 1 reference, got %d", len(skill.References))
	}
	if len(skill.Assets) != 1 {
		t.Errorf("expected 1 asset, got %d", len(skill.Assets))
	}
}

func TestSkillHasTrigger(t *testing.T) {
	skill := NewSkill("test", "Test skill").
		AddTrigger("foo").
		AddTrigger("bar")

	if !skill.HasTrigger("foo") {
		t.Error("expected HasTrigger('foo') to return true")
	}
	if !skill.HasTrigger("bar") {
		t.Error("expected HasTrigger('bar') to return true")
	}
	if skill.HasTrigger("baz") {
		t.Error("expected HasTrigger('baz') to return false")
	}
}

func TestSkillHasDependency(t *testing.T) {
	skill := NewSkill("test", "Test skill").
		AddDependency("git").
		AddDependency("npm")

	if !skill.HasDependency("git") {
		t.Error("expected HasDependency('git') to return true")
	}
	if skill.HasDependency("docker") {
		t.Error("expected HasDependency('docker') to return false")
	}
}

func TestSkillJSONSerialization(t *testing.T) {
	skill := NewSkill("version-analysis", "Analyze versions").
		WithInstructions("Instructions here").
		WithModel(ModelSonnet).
		AddTrigger("version").
		AddDependency("git")

	data, err := json.Marshal(skill)
	if err != nil {
		t.Fatalf("failed to marshal skill: %v", err)
	}

	var decoded Skill
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal skill: %v", err)
	}

	if decoded.Name != skill.Name {
		t.Errorf("name mismatch: got %s, want %s", decoded.Name, skill.Name)
	}
	if decoded.Description != skill.Description {
		t.Errorf("description mismatch: got %s, want %s", decoded.Description, skill.Description)
	}
	if decoded.Model != skill.Model {
		t.Errorf("model mismatch: got %s, want %s", decoded.Model, skill.Model)
	}
	if len(decoded.Triggers) != len(skill.Triggers) {
		t.Errorf("triggers count mismatch: got %d, want %d", len(decoded.Triggers), len(skill.Triggers))
	}
}

func TestSkillYAMLSerialization(t *testing.T) {
	skill := NewSkill("test-skill", "A test skill").
		WithInstructions("Do the thing").
		AddTrigger("test")

	data, err := yaml.Marshal(skill)
	if err != nil {
		t.Fatalf("failed to marshal skill to YAML: %v", err)
	}

	var decoded Skill
	if err := yaml.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal skill from YAML: %v", err)
	}

	if decoded.Name != skill.Name {
		t.Errorf("name mismatch: got %s, want %s", decoded.Name, skill.Name)
	}
}

func TestSkillSet(t *testing.T) {
	ss := NewSkillSet()

	skill1 := NewSkill("skill-a", "Skill A").AddTrigger("alpha")
	skill2 := NewSkill("skill-b", "Skill B").AddTrigger("beta").AddTrigger("alpha")

	ss.Add(skill1)
	ss.Add(skill2)

	// Test Get
	if ss.Get("skill-a") == nil {
		t.Error("expected to find skill-a")
	}
	if ss.Get("skill-c") != nil {
		t.Error("expected skill-c to be nil")
	}

	// Test Names
	names := ss.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}

	// Test FindByTrigger
	alphaSkills := ss.FindByTrigger("alpha")
	if len(alphaSkills) != 2 {
		t.Errorf("expected 2 skills with 'alpha' trigger, got %d", len(alphaSkills))
	}

	betaSkills := ss.FindByTrigger("beta")
	if len(betaSkills) != 1 {
		t.Errorf("expected 1 skill with 'beta' trigger, got %d", len(betaSkills))
	}

	unknownSkills := ss.FindByTrigger("unknown")
	if len(unknownSkills) != 0 {
		t.Errorf("expected 0 skills with 'unknown' trigger, got %d", len(unknownSkills))
	}
}

func TestSkillJSONOmitEmpty(t *testing.T) {
	skill := NewSkill("minimal", "Minimal skill")

	data, err := json.Marshal(skill)
	if err != nil {
		t.Fatalf("failed to marshal skill: %v", err)
	}

	// Check that empty fields are omitted
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	// Name and description should be present
	if _, ok := m["name"]; !ok {
		t.Error("expected 'name' field to be present")
	}
	if _, ok := m["description"]; !ok {
		t.Error("expected 'description' field to be present")
	}

	// Empty fields should be omitted
	if _, ok := m["instructions"]; ok {
		t.Error("expected 'instructions' field to be omitted when empty")
	}
	if _, ok := m["triggers"]; ok {
		t.Error("expected 'triggers' field to be omitted when empty")
	}
	if _, ok := m["dependencies"]; ok {
		t.Error("expected 'dependencies' field to be omitted when empty")
	}
}

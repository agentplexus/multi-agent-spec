package multiagentspec

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Loader loads multi-agent-spec definitions from files.
type Loader struct{}

// LoaderOption configures the loader.
type LoaderOption func(*Loader)

// NewLoader creates a new loader with the given options.
func NewLoader(opts ...LoaderOption) *Loader {
	l := &Loader{}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// LoadTeam loads a Team from a JSON file.
func (l *Loader) LoadTeam(path string) (*Team, error) {
	return LoadTeamFromFile(path)
}

// LoadAgent loads an Agent from a markdown file.
func (l *Loader) LoadAgent(path string) (*Agent, error) {
	return LoadAgentFromFile(path)
}

// LoadDeployment loads a Deployment from a JSON file.
func (l *Loader) LoadDeployment(path string) (*Deployment, error) {
	return LoadDeploymentFromFile(path)
}

// LoadSkill loads a Skill from a markdown file.
func (l *Loader) LoadSkill(path string) (*Skill, error) {
	return LoadSkillFromFile(path)
}

// LoadLoop loads a Loop from a JSON or YAML file.
func (l *Loader) LoadLoop(path string) (*Loop, error) {
	return LoadLoopFromFile(path)
}

// LoadAgentFromFile loads an Agent from a markdown file with YAML frontmatter.
//
// The file format is:
//
//	---
//	name: agent-name
//	description: Agent description
//	model: sonnet
//	tools: [Read, Write, Bash]
//	tasks:
//	  - id: task-id
//	    description: Task description
//	---
//
//	# Agent Name
//
//	Instructions in markdown...
func LoadAgentFromFile(path string) (*Agent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", path, err)
	}

	return ParseAgentMarkdown(data)
}

// ParseAgentMarkdown parses an Agent from markdown bytes with YAML frontmatter.
func ParseAgentMarkdown(data []byte) (*Agent, error) {
	frontmatter, body, err := splitFrontmatter(data)
	if err != nil {
		return nil, fmt.Errorf("parse frontmatter: %w", err)
	}

	var agent Agent
	if err := yaml.Unmarshal(frontmatter, &agent); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	// Set instructions from markdown body
	agent.Instructions = strings.TrimSpace(string(body))

	return &agent, nil
}

// LoadAgentsFromDir loads all Agent definitions from a directory.
// It recursively scans subdirectories. Agents in subdirectories have their
// namespace set to the subdirectory name (relative to the root dir), unless
// an explicit namespace is specified in the agent's frontmatter.
//
// Example structure:
//
//	agents/
//	├── shared/
//	│   └── review-board.md    → namespace: "shared", name: "review-board"
//	├── prd/
//	│   └── lead.md            → namespace: "prd", name: "lead"
//	└── orchestrator.md        → namespace: "", name: "orchestrator"
func LoadAgentsFromDir(dir string) ([]*Agent, error) {
	var agents []*Agent

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Skip non-markdown files
		if filepath.Ext(d.Name()) != ".md" {
			return nil
		}

		agent, err := LoadAgentFromFile(path)
		if err != nil {
			return fmt.Errorf("load %s: %w", path, err)
		}

		// Derive namespace from subdirectory if not explicitly set
		if agent.Namespace == "" {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return fmt.Errorf("relative path %s: %w", path, err)
			}

			relDir := filepath.Dir(relPath)
			if relDir != "." {
				// Convert path separators to forward slash for consistency
				agent.Namespace = filepath.ToSlash(relDir)
			}
		}

		agents = append(agents, agent)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walk dir %s: %w", dir, err)
	}

	return agents, nil
}

// LoadAgentsFromDirFlat loads agents from a single directory without recursion.
// This preserves the original non-recursive behavior for cases where
// subdirectories should be ignored.
func LoadAgentsFromDirFlat(dir string) ([]*Agent, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir %s: %w", dir, err)
	}

	var agents []*Agent
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		agent, err := LoadAgentFromFile(path)
		if err != nil {
			return nil, fmt.Errorf("load %s: %w", entry.Name(), err)
		}
		agents = append(agents, agent)
	}

	return agents, nil
}

// LoadSkillFromFile loads a Skill from a markdown file with YAML frontmatter.
//
// The file format is:
//
//	---
//	name: skill-name
//	description: Skill description
//	triggers: [keyword1, keyword2]
//	dependencies: [git, npm]
//	---
//
//	# Skill Name
//
//	Instructions in markdown...
func LoadSkillFromFile(path string) (*Skill, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", path, err)
	}

	return ParseSkillMarkdown(data)
}

// ParseSkillMarkdown parses a Skill from markdown bytes with YAML frontmatter.
func ParseSkillMarkdown(data []byte) (*Skill, error) {
	frontmatter, body, err := splitFrontmatter(data)
	if err != nil {
		return nil, fmt.Errorf("parse frontmatter: %w", err)
	}

	var skill Skill
	if err := yaml.Unmarshal(frontmatter, &skill); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	// Set instructions from markdown body
	skill.Instructions = strings.TrimSpace(string(body))

	return &skill, nil
}

// LoadSkillsFromDir loads all Skill definitions from a directory.
// It recursively scans subdirectories looking for SKILL.md files or
// markdown files with skill frontmatter.
//
// Example structure:
//
//	skills/
//	├── version-analysis/
//	│   ├── SKILL.md
//	│   └── scripts/
//	├── commit-classification/
//	│   └── SKILL.md
//	└── changelog/
//	    └── SKILL.md
func LoadSkillsFromDir(dir string) ([]*Skill, error) {
	var skills []*Skill

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Only load SKILL.md files or .md files at skill directory level
		name := d.Name()
		if name != "SKILL.md" && filepath.Ext(name) != ".md" {
			return nil
		}

		// Skip files in nested directories (scripts/, references/, assets/)
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return fmt.Errorf("relative path %s: %w", path, err)
		}

		// Only process direct children or SKILL.md in immediate subdirs
		parts := strings.Split(filepath.ToSlash(relPath), "/")
		if len(parts) > 2 {
			return nil // Skip nested files like skills/foo/scripts/bar.md
		}
		if len(parts) == 2 && name != "SKILL.md" {
			return nil // In subdir but not SKILL.md
		}

		skill, err := LoadSkillFromFile(path)
		if err != nil {
			return fmt.Errorf("load %s: %w", path, err)
		}

		// Infer name from directory if not set
		if skill.Name == "" && len(parts) == 2 {
			skill.Name = parts[0]
		}

		skills = append(skills, skill)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walk dir %s: %w", dir, err)
	}

	return skills, nil
}

// LoadSkillSetFromDir loads skills into a SkillSet from a directory.
func LoadSkillSetFromDir(dir string) (*SkillSet, error) {
	skills, err := LoadSkillsFromDir(dir)
	if err != nil {
		return nil, err
	}

	ss := NewSkillSet()
	for _, skill := range skills {
		ss.Add(skill)
	}

	return ss, nil
}

// LoadTeamFromFile loads a Team from a JSON file.
func LoadTeamFromFile(path string) (*Team, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", path, err)
	}

	var team Team
	if err := json.Unmarshal(data, &team); err != nil {
		return nil, fmt.Errorf("parse json: %w", err)
	}

	return &team, nil
}

// LoadDeploymentFromFile loads a Deployment from a JSON file.
func LoadDeploymentFromFile(path string) (*Deployment, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", path, err)
	}

	var deployment Deployment
	if err := json.Unmarshal(data, &deployment); err != nil {
		return nil, fmt.Errorf("parse json: %w", err)
	}

	return &deployment, nil
}

// LoadLoopFromFile loads a Loop from a JSON or YAML file.
//
// Supported formats:
//   - JSON: Standard JSON object
//   - YAML: Standard YAML document
//
// Example YAML:
//
//	name: qa-fix
//	type: VEAL
//	validator: qa
//	actor: code-fixer
//	max_attempts: 3
//	escalation: human
//	checks:
//	  - id: build
//	    type: command
//	    command: go build ./...
func LoadLoopFromFile(path string) (*Loop, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", path, err)
	}

	ext := filepath.Ext(path)
	var loop Loop

	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &loop); err != nil {
			return nil, fmt.Errorf("parse json: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &loop); err != nil {
			return nil, fmt.Errorf("parse yaml: %w", err)
		}
	default:
		// Try JSON first, then YAML
		if err := json.Unmarshal(data, &loop); err != nil {
			if err := yaml.Unmarshal(data, &loop); err != nil {
				return nil, fmt.Errorf("parse file (tried json and yaml): %w", err)
			}
		}
	}

	return &loop, nil
}

// LoadLoopsFromDir loads all Loop definitions from a directory.
// It looks for .json, .yaml, and .yml files.
//
// Example structure:
//
//	loops/
//	├── qa-fix.yaml
//	├── docs-fix.yaml
//	└── security-fix.json
func LoadLoopsFromDir(dir string) ([]*Loop, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir %s: %w", dir, err)
	}

	var loops []*Loop
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".json" && ext != ".yaml" && ext != ".yml" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		loop, err := LoadLoopFromFile(path)
		if err != nil {
			return nil, fmt.Errorf("load %s: %w", entry.Name(), err)
		}
		loops = append(loops, loop)
	}

	return loops, nil
}

// splitFrontmatter splits YAML frontmatter from markdown body.
// Frontmatter is delimited by --- at the start and end.
func splitFrontmatter(data []byte) (frontmatter, body []byte, err error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))

	// Check for opening delimiter
	if !scanner.Scan() {
		return nil, nil, fmt.Errorf("empty file")
	}
	if strings.TrimSpace(scanner.Text()) != "---" {
		return nil, nil, fmt.Errorf("missing frontmatter delimiter")
	}

	// Read frontmatter until closing delimiter
	var fm bytes.Buffer
	foundEnd := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			foundEnd = true
			break
		}
		fm.WriteString(line)
		fm.WriteString("\n")
	}

	if !foundEnd {
		return nil, nil, fmt.Errorf("missing closing frontmatter delimiter")
	}

	// Rest is body
	var bd bytes.Buffer
	for scanner.Scan() {
		bd.WriteString(scanner.Text())
		bd.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("scan error: %w", err)
	}

	return fm.Bytes(), bd.Bytes(), nil
}

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

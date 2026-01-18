package multiagentspec

import (
	"bytes"
	"testing"
)

func TestSortByDAG(t *testing.T) {
	tests := []struct {
		name     string
		teams    []TeamSection
		expected []string // expected order of IDs
	}{
		{
			name:     "empty",
			teams:    []TeamSection{},
			expected: []string{},
		},
		{
			name: "single team",
			teams: []TeamSection{
				{ID: "pm-validation", Name: "pm"},
			},
			expected: []string{"pm-validation"},
		},
		{
			name: "linear chain",
			teams: []TeamSection{
				{ID: "release-validation", Name: "release", DependsOn: []string{"qa-validation"}},
				{ID: "qa-validation", Name: "qa", DependsOn: []string{"pm-validation"}},
				{ID: "pm-validation", Name: "pm"},
			},
			expected: []string{"pm-validation", "qa-validation", "release-validation"},
		},
		{
			name: "diamond DAG",
			teams: []TeamSection{
				{ID: "release-validation", Name: "release", DependsOn: []string{"qa-validation", "docs-validation", "security-validation"}},
				{ID: "security-validation", Name: "security", DependsOn: []string{"pm-validation"}},
				{ID: "docs-validation", Name: "documentation", DependsOn: []string{"pm-validation"}},
				{ID: "qa-validation", Name: "qa", DependsOn: []string{"pm-validation"}},
				{ID: "pm-validation", Name: "pm"},
			},
			// pm first, then docs/qa/security (alphabetically), then release last
			expected: []string{"pm-validation", "docs-validation", "qa-validation", "security-validation", "release-validation"},
		},
		{
			name: "no dependencies - preserve order",
			teams: []TeamSection{
				{ID: "a", Name: "a"},
				{ID: "b", Name: "b"},
				{ID: "c", Name: "c"},
			},
			expected: []string{"a", "b", "c"},
		},
		{
			name: "missing dependency - graceful handling",
			teams: []TeamSection{
				{ID: "b", Name: "b", DependsOn: []string{"missing"}},
				{ID: "a", Name: "a"},
			},
			// "b" depends on "missing" which doesn't exist, so both have 0 in-degree
			// alphabetical ordering: a, b
			expected: []string{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &TeamReport{Teams: tt.teams}
			report.SortByDAG()

			if len(report.Teams) != len(tt.expected) {
				t.Fatalf("expected %d teams, got %d", len(tt.expected), len(report.Teams))
			}

			for i, expectedID := range tt.expected {
				if report.Teams[i].ID != expectedID {
					t.Errorf("position %d: expected %q, got %q", i, expectedID, report.Teams[i].ID)
				}
			}
		})
	}
}

func TestComputeStatus(t *testing.T) {
	tests := []struct {
		name     string
		tasks    []TaskResult
		expected Status
	}{
		{
			name:     "all GO",
			tasks:    []TaskResult{{Status: StatusGo}, {Status: StatusGo}},
			expected: StatusGo,
		},
		{
			name:     "has WARN",
			tasks:    []TaskResult{{Status: StatusGo}, {Status: StatusWarn}},
			expected: StatusWarn,
		},
		{
			name:     "has NO-GO",
			tasks:    []TaskResult{{Status: StatusGo}, {Status: StatusNoGo}},
			expected: StatusNoGo,
		},
		{
			name:     "all SKIP",
			tasks:    []TaskResult{{Status: StatusSkip}, {Status: StatusSkip}},
			expected: StatusSkip,
		},
		{
			name:     "NO-GO takes precedence over WARN",
			tasks:    []TaskResult{{Status: StatusWarn}, {Status: StatusNoGo}},
			expected: StatusNoGo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &AgentResult{Tasks: tt.tasks}
			if got := result.ComputeStatus(); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestRendererSortsTeams(t *testing.T) {
	// Create a report with teams in wrong order
	report := &TeamReport{
		Project: "test-project",
		Version: "v1.0.0",
		Target:  "v1.0.0",
		Phase:   "PHASE 1: REVIEW",
		Teams: []TeamSection{
			{
				ID:        "release-validation",
				Name:      "release",
				DependsOn: []string{"qa-validation"},
				Status:    StatusGo,
				Tasks:     []TaskResult{{ID: "task", Status: StatusGo}},
			},
			{
				ID:        "qa-validation",
				Name:      "qa",
				DependsOn: []string{"pm-validation"},
				Status:    StatusGo,
				Tasks:     []TaskResult{{ID: "task", Status: StatusGo}},
			},
			{
				ID:     "pm-validation",
				Name:   "pm",
				Status: StatusGo,
				Tasks:  []TaskResult{{ID: "task", Status: StatusGo}},
			},
		},
		Status: StatusGo,
	}

	var buf bytes.Buffer
	renderer := NewRenderer(&buf)
	if err := renderer.Render(report); err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// After rendering, teams should be sorted
	if report.Teams[0].ID != "pm-validation" {
		t.Errorf("expected first team to be pm-validation, got %s", report.Teams[0].ID)
	}
	if report.Teams[1].ID != "qa-validation" {
		t.Errorf("expected second team to be qa-validation, got %s", report.Teams[1].ID)
	}
	if report.Teams[2].ID != "release-validation" {
		t.Errorf("expected third team to be release-validation, got %s", report.Teams[2].ID)
	}
}

func TestStatusIcon(t *testing.T) {
	tests := []struct {
		status Status
		icon   string
	}{
		{StatusGo, "\U0001F7E2"},   // 🟢
		{StatusWarn, "\U0001F7E1"}, // 🟡
		{StatusNoGo, "\U0001F534"}, // 🔴
		{StatusSkip, "\u26AA"},     // ⚪
		{Status("UNKNOWN"), "?"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.Icon(); got != tt.icon {
				t.Errorf("expected %q, got %q", tt.icon, got)
			}
		})
	}
}

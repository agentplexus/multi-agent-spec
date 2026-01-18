package multiagentspec

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

const (
	// boxWidth is the inner width of the box (between the border characters).
	boxWidth = 76
)

// Renderer renders TeamReport to various formats.
type Renderer struct {
	w io.Writer
}

// NewRenderer creates a new Renderer writing to w.
func NewRenderer(w io.Writer) *Renderer {
	return &Renderer{w: w}
}

// Render renders the report using the box template.
// It automatically sorts teams by DAG order before rendering.
func (r *Renderer) Render(report *TeamReport) error {
	report.SortByDAG()
	return r.renderBox(report)
}

// renderBox renders the report in the box format.
func (r *Renderer) renderBox(report *TeamReport) error {
	tmpl, err := template.New("report").Funcs(templateFuncs()).Parse(BoxTemplate)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}
	return tmpl.Execute(r.w, report)
}

// templateFuncs returns the template function map.
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"header":       header,
		"separator":    separator,
		"footer":       footer,
		"teamHeader":   teamHeader,
		"taskLine":     taskLine,
		"centerLine":   centerLine,
		"paddedLine":   paddedLine,
		"finalMessage": finalMessage,
	}
}

// header returns the top border of the box.
func header() string {
	return "╔" + strings.Repeat("═", boxWidth) + "╗"
}

// separator returns a separator line.
func separator() string {
	return "╠" + strings.Repeat("═", boxWidth) + "╣"
}

// footer returns the bottom border of the box.
func footer() string {
	return "╚" + strings.Repeat("═", boxWidth) + "╝"
}

// centerLine centers text within the box.
func centerLine(text string) string {
	visualLen := visualLength(text)
	padding := max(0, boxWidth-visualLen)
	left := padding / 2
	right := padding - left
	return "║" + strings.Repeat(" ", left) + text + strings.Repeat(" ", right) + "║"
}

// paddedLine left-aligns text with padding.
func paddedLine(text string) string {
	visualLen := visualLength(text)
	padding := max(0, boxWidth-visualLen-1)
	return "║ " + text + strings.Repeat(" ", padding) + "║"
}

// teamHeader formats a team header line.
func teamHeader(team TeamSection) string {
	text := fmt.Sprintf("%s (%s)", team.ID, team.Name)
	return paddedLine(text)
}

// taskLine formats a single task result line.
func taskLine(task TaskResult) string {
	id := task.ID
	if len(id) > 24 {
		id = id[:21] + "..."
	}

	icon := task.Status.Icon()
	statusText := string(task.Status)

	detail := task.Detail
	maxDetail := boxWidth - 40
	if len(detail) > maxDetail {
		detail = detail[:maxDetail-3] + "..."
	}

	line := fmt.Sprintf("  %-24s %s %-5s %s", id, icon, statusText, detail)
	return paddedLine(line)
}

// finalMessage formats the final status message line.
func finalMessage(report *TeamReport) string {
	return centerLine(report.FinalMessage())
}

// visualLength calculates the visual length of a string,
// accounting for emoji characters that take 2 columns.
func visualLength(s string) int {
	length := 0
	for _, r := range s {
		if r >= 0x1F300 && r <= 0x1FAFF {
			length += 2
		} else if r >= 0x2600 && r <= 0x27BF {
			length += 2
		} else {
			length++
		}
	}
	return length
}

// BoxTemplate is the text/template for the box format report.
// This is the reference implementation for rendering TeamReport to text.
// Each task is rendered on its own line with status indicator.
const BoxTemplate = `{{ header }}
{{ centerLine "TEAM STATUS REPORT" }}
{{ separator }}
{{ paddedLine (printf "Project: %s" .Project) }}
{{ paddedLine (printf "Target:  %s" .Target) }}
{{ separator }}
{{ paddedLine .Phase }}
{{- range .Teams }}
{{ separator }}
{{ teamHeader . }}
{{- range .Tasks }}
{{ taskLine . }}
{{- end }}
{{- end }}
{{ separator }}
{{ finalMessage . }}
{{ footer }}
`

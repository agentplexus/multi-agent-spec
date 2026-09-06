package multiagentspec

import (
	"github.com/invopop/jsonschema"
)

// JSONSchema implements jsonschema.Schema for Model type.
func (Model) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Enum:        []interface{}{"haiku", "sonnet", "opus"},
		Default:     "sonnet",
		Description: "Model capability tier (mapped to platform-specific models)",
	}
}

// JSONSchema implements jsonschema.Schema for Tool type.
func (Tool) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []interface{}{
			"WebSearch", "WebFetch", "Read", "Write",
			"Glob", "Grep", "Bash", "Edit", "Task",
		},
		Description: "Canonical tool name (mapped to platform-specific names)",
	}
}

// JSONSchema implements jsonschema.Schema for TaskType type.
func (TaskType) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Enum:        []interface{}{"command", "pattern", "file", "manual"},
		Default:     "manual",
		Description: "How the task is executed",
	}
}

// JSONSchema implements jsonschema.Schema for WorkflowType type.
func (WorkflowType) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Enum:        []interface{}{"sequential", "parallel", "dag", "orchestrated"},
		Default:     "orchestrated",
		Description: "Workflow execution pattern",
	}
}

// JSONSchema implements jsonschema.Schema for Platform type.
func (Platform) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []interface{}{
			"claude-code", "gemini-cli", "kiro-cli", "adk-go",
			"crewai", "autogen", "aws-agentcore", "aws-eks",
			"azure-aks", "gcp-gke", "kubernetes", "docker-compose",
			"agentkit-local",
		},
		Description: "Supported deployment platform",
	}
}

// JSONSchema implements jsonschema.Schema for Priority type.
func (Priority) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Enum:        []interface{}{"p1", "p2", "p3"},
		Default:     "p2",
		Description: "Deployment priority level",
	}
}

// JSONSchema implements jsonschema.Schema for DeploymentMode type.
func (DeploymentMode) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Enum:        []interface{}{"single-process", "multi-process", "distributed", "serverless"},
		Description: "Deployment execution mode",
	}
}

// JSONSchema implements jsonschema.Schema for PortType type.
func (PortType) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Enum:        []interface{}{"string", "number", "boolean", "object", "array", "file"},
		Description: "Data type of a port",
	}
}

// JSONSchema implements jsonschema.Schema for Status type.
func (Status) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Enum:        []interface{}{"GO", "NO-GO", "WARN", "SKIP"},
		Description: "Validation status",
	}
}

// JSONSchema implements jsonschema.Schema for ContentBlockType type.
func (ContentBlockType) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Enum:        []interface{}{"kv_pairs", "list", "table", "text", "metric"},
		Description: "Content block type discriminator",
	}
}

// JSONSchema implements jsonschema.Schema for ReviewStatus type.
func (ReviewStatus) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Enum:        []interface{}{"pending", "confirmed", "edited", "rejected"},
		Default:     "pending",
		Description: "State of a human review checkpoint",
	}
}

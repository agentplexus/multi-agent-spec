package multiagentspec

import "encoding/json"

// Platform represents supported deployment platforms.
type Platform string

const (
	PlatformClaudeCode    Platform = "claude-code"
	PlatformKiroCLI       Platform = "kiro-cli"
	PlatformAWSAgentCore  Platform = "aws-agentcore"
	PlatformAWSEKS        Platform = "aws-eks"
	PlatformAzureAKS      Platform = "azure-aks"
	PlatformGCPGKE        Platform = "gcp-gke"
	PlatformKubernetes    Platform = "kubernetes"
	PlatformDockerCompose Platform = "docker-compose"
	PlatformAgentKitLocal Platform = "agentkit-local"
)

// Priority represents deployment priority levels.
type Priority string

const (
	PriorityP1 Priority = "p1"
	PriorityP2 Priority = "p2"
	PriorityP3 Priority = "p3"
)

// Target represents a deployment target definition.
type Target struct {
	// Name is the unique name for this deployment target.
	Name string `json:"name"`

	// Platform is the target platform for deployment.
	Platform Platform `json:"platform"`

	// Priority is the deployment priority.
	Priority Priority `json:"priority,omitempty"`

	// Output is the directory for generated deployment artifacts.
	Output string `json:"output"`

	// Config is platform-specific configuration.
	Config json.RawMessage `json:"config,omitempty"`
}

// Deployment represents a deployment definition.
type Deployment struct {
	// Schema is the JSON Schema reference.
	Schema string `json:"$schema,omitempty"`

	// Team is the reference to the team definition (team name).
	Team string `json:"team"`

	// Targets is the list of deployment targets.
	Targets []Target `json:"targets"`
}

// NewDeployment creates a new Deployment for the given team.
func NewDeployment(team string) *Deployment {
	return &Deployment{
		Team:    team,
		Targets: []Target{},
	}
}

// AddTarget adds a deployment target and returns the deployment for chaining.
func (d *Deployment) AddTarget(target Target) *Deployment {
	d.Targets = append(d.Targets, target)
	return d
}

// ClaudeCodeConfig is the configuration for Claude Code platform.
type ClaudeCodeConfig struct {
	AgentDir string `json:"agentDir"`
	Format   string `json:"format"`
}

// KiroCLIConfig is the configuration for Kiro CLI platform.
type KiroCLIConfig struct {
	PluginDir string `json:"pluginDir"`
	Format    string `json:"format"`
}

// AWSAgentCoreConfig is the configuration for AWS AgentCore platform.
type AWSAgentCoreConfig struct {
	Region          string `json:"region"`
	FoundationModel string `json:"foundationModel"`
	IAC             string `json:"iac"`
	LambdaRuntime   string `json:"lambdaRuntime"`
}

// KubernetesConfig is the configuration for Kubernetes platforms.
type KubernetesConfig struct {
	Namespace      string          `json:"namespace"`
	HelmChart      bool            `json:"helmChart"`
	ImageRegistry  string          `json:"imageRegistry,omitempty"`
	ResourceLimits *ResourceLimits `json:"resourceLimits,omitempty"`
}

// ResourceLimits defines Kubernetes resource limits.
type ResourceLimits struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

// AgentKitLocalConfig is the configuration for AgentKit local platform.
type AgentKitLocalConfig struct {
	Transport string `json:"transport"`
	Port      int    `json:"port,omitempty"`
}

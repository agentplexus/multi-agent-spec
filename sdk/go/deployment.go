package multiagentspec

// Platform represents supported deployment platforms.
type Platform string

const (
	PlatformClaudeCode    Platform = "claude-code"
	PlatformGeminiCLI     Platform = "gemini-cli"
	PlatformKiroCLI       Platform = "kiro-cli"
	PlatformADKGo         Platform = "adk-go"
	PlatformCrewAI        Platform = "crewai"
	PlatformAutoGen       Platform = "autogen"
	PlatformAWSAgentCore  Platform = "aws-agentcore"
	PlatformAWSEKS        Platform = "aws-eks"
	PlatformAzureAKS      Platform = "azure-aks"
	PlatformGCPGKE        Platform = "gcp-gke"
	PlatformKubernetes    Platform = "kubernetes"
	PlatformDockerCompose Platform = "docker-compose"
	PlatformAgentKitLocal Platform = "agentkit-local"
)

// DeploymentMode represents the deployment execution mode.
type DeploymentMode string

const (
	ModeSingleProcess DeploymentMode = "single-process"
	ModeMultiProcess  DeploymentMode = "multi-process"
	ModeDistributed   DeploymentMode = "distributed"
	ModeServerless    DeploymentMode = "serverless"
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

	// Mode is the deployment mode affecting runtime behavior.
	Mode DeploymentMode `json:"mode,omitempty"`

	// Priority is the deployment priority.
	Priority Priority `json:"priority,omitempty"`

	// Output is the directory for generated deployment artifacts.
	Output string `json:"output,omitempty"`

	// Runtime is the runtime configuration for workflow execution.
	Runtime *RuntimeConfig `json:"runtime,omitempty"`

	// Platform-specific configurations (use the one matching Platform field)
	ClaudeCode    *ClaudeCodeConfig    `json:"claudeCode,omitempty"`
	GeminiCLI     *GeminiCLIConfig     `json:"geminiCli,omitempty"`
	KiroCLI       *KiroCLIConfig       `json:"kiroCli,omitempty"`
	ADKGo         *ADKGoConfig         `json:"adkGo,omitempty"`
	CrewAI        *CrewAIConfig        `json:"crewai,omitempty"`
	AutoGen       *AutoGenConfig       `json:"autogen,omitempty"`
	AWSAgentCore  *AWSAgentCoreConfig  `json:"awsAgentCore,omitempty"`
	Kubernetes    *KubernetesConfig    `json:"kubernetes,omitempty"`
	DockerCompose *DockerComposeConfig `json:"dockerCompose,omitempty"`
	AgentKitLocal *AgentKitLocalConfig `json:"agentKitLocal,omitempty"`
}

// RuntimeConfig holds runtime configuration for workflow execution.
type RuntimeConfig struct {
	// Defaults are the default runtime settings for all steps.
	Defaults *StepRuntime `json:"defaults,omitempty"`

	// Steps contains per-step runtime overrides keyed by step name.
	Steps map[string]*StepRuntime `json:"steps,omitempty"`

	// Observability contains monitoring and tracing settings.
	Observability *ObservabilityConfig `json:"observability,omitempty"`
}

// StepRuntime holds runtime settings for a workflow step.
type StepRuntime struct {
	// Timeout is the step timeout (e.g., 30s, 5m, 1h).
	Timeout string `json:"timeout,omitempty"`

	// Retry is the retry policy for step failures.
	Retry *RetryPolicy `json:"retry,omitempty"`

	// Condition is a condition expression for conditional execution.
	Condition string `json:"condition,omitempty"`

	// Concurrency is the max concurrent executions of this step.
	Concurrency int `json:"concurrency,omitempty"`

	// Resources are resource limits for this step.
	Resources *ResourceLimits `json:"resources,omitempty"`
}

// RetryPolicy defines the retry behavior for step failures.
type RetryPolicy struct {
	// MaxAttempts is the maximum number of retry attempts.
	MaxAttempts int `json:"max_attempts,omitempty"`

	// Backoff is the backoff strategy (fixed, exponential, linear).
	Backoff string `json:"backoff,omitempty"`

	// InitialDelay is the initial delay before first retry.
	InitialDelay string `json:"initial_delay,omitempty"`

	// MaxDelay is the maximum delay between retries.
	MaxDelay string `json:"max_delay,omitempty"`

	// RetryableErrors are error types that should trigger retry.
	RetryableErrors []string `json:"retryable_errors,omitempty"`
}

// ObservabilityConfig holds observability and monitoring configuration.
type ObservabilityConfig struct {
	// Tracing contains distributed tracing settings.
	Tracing *TracingConfig `json:"tracing,omitempty"`

	// Metrics contains metrics collection settings.
	Metrics *MetricsConfig `json:"metrics,omitempty"`

	// Logging contains logging configuration.
	Logging *LoggingConfig `json:"logging,omitempty"`
}

// TracingConfig holds distributed tracing configuration.
type TracingConfig struct {
	Enabled    bool    `json:"enabled,omitempty"`
	Exporter   string  `json:"exporter,omitempty"`
	Endpoint   string  `json:"endpoint,omitempty"`
	SampleRate float64 `json:"sample_rate,omitempty"`
}

// MetricsConfig holds metrics collection configuration.
type MetricsConfig struct {
	Enabled  bool   `json:"enabled,omitempty"`
	Exporter string `json:"exporter,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
}

// LoggingConfig holds logging configuration.
type LoggingConfig struct {
	Level  string `json:"level,omitempty"`
	Format string `json:"format,omitempty"`
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

// ResourceLimits defines resource limits for step execution.
type ResourceLimits struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
	GPU    int    `json:"gpu,omitempty"`
}

// AgentKitLocalConfig is the configuration for AgentKit local platform.
type AgentKitLocalConfig struct {
	Transport string `json:"transport"`
	Port      int    `json:"port,omitempty"`
}

// GeminiCLIConfig is the configuration for Google Gemini CLI Assistant.
type GeminiCLIConfig struct {
	Model     string `json:"model,omitempty"`
	ConfigDir string `json:"configDir,omitempty"`
}

// ADKGoConfig is the configuration for Google Agent Development Kit (Go).
type ADKGoConfig struct {
	Model        string `json:"model,omitempty"`
	ServerPort   int    `json:"serverPort,omitempty"`
	SessionStore string `json:"sessionStore,omitempty"`
	ToolRegistry string `json:"toolRegistry,omitempty"`
}

// CrewAIConfig is the configuration for CrewAI deployment.
type CrewAIConfig struct {
	Model         string `json:"model,omitempty"`
	Verbose       bool   `json:"verbose,omitempty"`
	Memory        bool   `json:"memory,omitempty"`
	ProcessType   string `json:"processType,omitempty"`
	MaxIterations int    `json:"maxIterations,omitempty"`
}

// AutoGenConfig is the configuration for Microsoft AutoGen deployment.
type AutoGenConfig struct {
	Model                   string               `json:"model,omitempty"`
	HumanInputMode          string               `json:"humanInputMode,omitempty"`
	MaxConsecutiveAutoReply int                  `json:"maxConsecutiveAutoReply,omitempty"`
	CodeExecutionConfig     *CodeExecutionConfig `json:"codeExecutionConfig,omitempty"`
}

// CodeExecutionConfig holds AutoGen code execution settings.
type CodeExecutionConfig struct {
	WorkDir   string `json:"workDir,omitempty"`
	UseDocker bool   `json:"useDocker,omitempty"`
}

// DockerComposeConfig is the configuration for Docker Compose deployment.
type DockerComposeConfig struct {
	NetworkMode string `json:"networkMode,omitempty"`
}

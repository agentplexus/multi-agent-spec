package multiagentspec

import (
	"encoding/json"
	"testing"
)

func TestPlatformConstants(t *testing.T) {
	tests := []struct {
		platform Platform
		want     string
	}{
		{PlatformClaudeCode, "claude-code"},
		{PlatformKiroCLI, "kiro-cli"},
		{PlatformAWSAgentCore, "aws-agentcore"},
		{PlatformAWSEKS, "aws-eks"},
		{PlatformAzureAKS, "azure-aks"},
		{PlatformGCPGKE, "gcp-gke"},
		{PlatformKubernetes, "kubernetes"},
		{PlatformDockerCompose, "docker-compose"},
		{PlatformAgentKitLocal, "agentkit-local"},
	}

	for _, tt := range tests {
		if string(tt.platform) != tt.want {
			t.Errorf("Platform %v = %q, want %q", tt.platform, string(tt.platform), tt.want)
		}
	}
}

func TestPriorityConstants(t *testing.T) {
	tests := []struct {
		priority Priority
		want     string
	}{
		{PriorityP1, "p1"},
		{PriorityP2, "p2"},
		{PriorityP3, "p3"},
	}

	for _, tt := range tests {
		if string(tt.priority) != tt.want {
			t.Errorf("Priority %v = %q, want %q", tt.priority, string(tt.priority), tt.want)
		}
	}
}

func TestTargetSerialization(t *testing.T) {
	target := Target{
		Name:     "local-claude",
		Platform: PlatformClaudeCode,
		Priority: PriorityP1,
		Output:   ".claude/agents",
		ClaudeCode: &ClaudeCodeConfig{
			AgentDir: ".claude/agents",
			Format:   "markdown",
		},
	}

	data, err := json.Marshal(target)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded Target
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Name != target.Name {
		t.Errorf("Name = %q, want %q", decoded.Name, target.Name)
	}
	if decoded.Platform != PlatformClaudeCode {
		t.Errorf("Platform = %q, want %q", decoded.Platform, PlatformClaudeCode)
	}
	if decoded.Priority != PriorityP1 {
		t.Errorf("Priority = %q, want %q", decoded.Priority, PriorityP1)
	}
}

func TestNewDeployment(t *testing.T) {
	deployment := NewDeployment("test-team")

	if deployment.Team != "test-team" {
		t.Errorf("Team = %q, want %q", deployment.Team, "test-team")
	}
	if len(deployment.Targets) != 0 {
		t.Errorf("len(Targets) = %d, want 0", len(deployment.Targets))
	}
}

func TestDeploymentAddTarget(t *testing.T) {
	deployment := NewDeployment("test").
		AddTarget(Target{
			Name:     "t1",
			Platform: PlatformClaudeCode,
			Output:   "out1",
		}).
		AddTarget(Target{
			Name:     "t2",
			Platform: PlatformKiroCLI,
			Output:   "out2",
		})

	if len(deployment.Targets) != 2 {
		t.Errorf("len(Targets) = %d, want 2", len(deployment.Targets))
	}
	if deployment.Targets[0].Name != "t1" {
		t.Errorf("Targets[0].Name = %q, want %q", deployment.Targets[0].Name, "t1")
	}
	if deployment.Targets[1].Platform != PlatformKiroCLI {
		t.Errorf("Targets[1].Platform = %q, want %q", deployment.Targets[1].Platform, PlatformKiroCLI)
	}
}

func TestDeploymentJSONSerialization(t *testing.T) {
	deployment := &Deployment{
		Schema: "../schema/deployment.schema.json",
		Team:   "json-team",
		Targets: []Target{
			{
				Name:     "target1",
				Platform: PlatformAWSAgentCore,
				Priority: PriorityP1,
				Output:   "cdk",
			},
		},
	}

	data, err := json.Marshal(deployment)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded Deployment
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Schema != deployment.Schema {
		t.Errorf("Schema = %q, want %q", decoded.Schema, deployment.Schema)
	}
	if decoded.Team != deployment.Team {
		t.Errorf("Team = %q, want %q", decoded.Team, deployment.Team)
	}
	if len(decoded.Targets) != 1 {
		t.Errorf("len(Targets) = %d, want 1", len(decoded.Targets))
	}
}

func TestClaudeCodeConfig(t *testing.T) {
	config := ClaudeCodeConfig{
		AgentDir: ".claude/agents",
		Format:   "markdown",
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded ClaudeCodeConfig
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.AgentDir != config.AgentDir {
		t.Errorf("AgentDir = %q, want %q", decoded.AgentDir, config.AgentDir)
	}
	if decoded.Format != config.Format {
		t.Errorf("Format = %q, want %q", decoded.Format, config.Format)
	}
}

func TestKiroCLIConfig(t *testing.T) {
	config := KiroCLIConfig{
		PluginDir: "plugins/kiro",
		Format:    "json",
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded KiroCLIConfig
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.PluginDir != config.PluginDir {
		t.Errorf("PluginDir = %q, want %q", decoded.PluginDir, config.PluginDir)
	}
}

func TestAWSAgentCoreConfig(t *testing.T) {
	config := AWSAgentCoreConfig{
		Region:          "us-west-2",
		FoundationModel: "anthropic.claude-3-sonnet",
		IAC:             "cdk",
		LambdaRuntime:   "python3.11",
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded AWSAgentCoreConfig
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Region != config.Region {
		t.Errorf("Region = %q, want %q", decoded.Region, config.Region)
	}
	if decoded.IAC != config.IAC {
		t.Errorf("IAC = %q, want %q", decoded.IAC, config.IAC)
	}
}

func TestKubernetesConfig(t *testing.T) {
	config := KubernetesConfig{
		Namespace:     "multi-agent",
		HelmChart:     true,
		ImageRegistry: "ghcr.io/agentplexus",
		ResourceLimits: &ResourceLimits{
			CPU:    "1000m",
			Memory: "1Gi",
		},
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded KubernetesConfig
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Namespace != config.Namespace {
		t.Errorf("Namespace = %q, want %q", decoded.Namespace, config.Namespace)
	}
	if decoded.ResourceLimits == nil {
		t.Error("ResourceLimits should not be nil")
	} else {
		if decoded.ResourceLimits.CPU != "1000m" {
			t.Errorf("ResourceLimits.CPU = %q, want %q", decoded.ResourceLimits.CPU, "1000m")
		}
	}
}

func TestResourceLimits(t *testing.T) {
	limits := ResourceLimits{
		CPU:    "500m",
		Memory: "512Mi",
	}

	data, err := json.Marshal(limits)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded ResourceLimits
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.CPU != limits.CPU {
		t.Errorf("CPU = %q, want %q", decoded.CPU, limits.CPU)
	}
	if decoded.Memory != limits.Memory {
		t.Errorf("Memory = %q, want %q", decoded.Memory, limits.Memory)
	}
}

func TestAgentKitLocalConfig(t *testing.T) {
	config := AgentKitLocalConfig{
		Transport: "http",
		Port:      8080,
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded AgentKitLocalConfig
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Transport != config.Transport {
		t.Errorf("Transport = %q, want %q", decoded.Transport, config.Transport)
	}
	if decoded.Port != config.Port {
		t.Errorf("Port = %d, want %d", decoded.Port, config.Port)
	}
}

func TestAgentKitLocalConfigOmitPort(t *testing.T) {
	config := AgentKitLocalConfig{
		Transport: "stdio",
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if _, ok := m["port"]; ok {
		t.Error("port should be omitted when zero")
	}
}

func TestNewPlatformConstants(t *testing.T) {
	tests := []struct {
		platform Platform
		want     string
	}{
		{PlatformGeminiCLI, "gemini-cli"},
		{PlatformADKGo, "adk-go"},
		{PlatformCrewAI, "crewai"},
		{PlatformAutoGen, "autogen"},
	}

	for _, tt := range tests {
		if string(tt.platform) != tt.want {
			t.Errorf("Platform %v = %q, want %q", tt.platform, string(tt.platform), tt.want)
		}
	}
}

func TestGeminiCLIConfig(t *testing.T) {
	config := GeminiCLIConfig{
		Model:     "gemini-2.0-flash",
		ConfigDir: ".gemini",
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded GeminiCLIConfig
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Model != config.Model {
		t.Errorf("Model = %q, want %q", decoded.Model, config.Model)
	}
	if decoded.ConfigDir != config.ConfigDir {
		t.Errorf("ConfigDir = %q, want %q", decoded.ConfigDir, config.ConfigDir)
	}
}

func TestADKGoConfig(t *testing.T) {
	config := ADKGoConfig{
		Model:        "gemini-2.0-flash",
		ServerPort:   8080,
		SessionStore: "memory",
		ToolRegistry: "local",
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded ADKGoConfig
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Model != config.Model {
		t.Errorf("Model = %q, want %q", decoded.Model, config.Model)
	}
	if decoded.ServerPort != config.ServerPort {
		t.Errorf("ServerPort = %d, want %d", decoded.ServerPort, config.ServerPort)
	}
}

func TestCrewAIConfig(t *testing.T) {
	config := CrewAIConfig{
		Model:         "gpt-4",
		Verbose:       true,
		Memory:        true,
		ProcessType:   "sequential",
		MaxIterations: 10,
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded CrewAIConfig
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Model != config.Model {
		t.Errorf("Model = %q, want %q", decoded.Model, config.Model)
	}
	if decoded.Verbose != config.Verbose {
		t.Errorf("Verbose = %v, want %v", decoded.Verbose, config.Verbose)
	}
}

func TestAutoGenConfig(t *testing.T) {
	config := AutoGenConfig{
		Model:                   "gpt-4",
		HumanInputMode:          "NEVER",
		MaxConsecutiveAutoReply: 5,
		CodeExecutionConfig: &CodeExecutionConfig{
			WorkDir:   "/tmp/autogen",
			UseDocker: true,
		},
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded AutoGenConfig
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Model != config.Model {
		t.Errorf("Model = %q, want %q", decoded.Model, config.Model)
	}
	if decoded.CodeExecutionConfig == nil {
		t.Error("CodeExecutionConfig should not be nil")
	} else {
		if decoded.CodeExecutionConfig.UseDocker != true {
			t.Errorf("UseDocker = %v, want true", decoded.CodeExecutionConfig.UseDocker)
		}
	}
}

func TestTargetWithTypedConfigs(t *testing.T) {
	target := Target{
		Name:     "gemini-target",
		Platform: PlatformGeminiCLI,
		GeminiCLI: &GeminiCLIConfig{
			Model: "gemini-2.0-flash",
		},
	}

	data, err := json.Marshal(target)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded Target
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.GeminiCLI == nil {
		t.Error("GeminiCLI config should not be nil")
	} else if decoded.GeminiCLI.Model != "gemini-2.0-flash" {
		t.Errorf("GeminiCLI.Model = %q, want %q", decoded.GeminiCLI.Model, "gemini-2.0-flash")
	}
}

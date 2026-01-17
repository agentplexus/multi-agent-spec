"""Tests for multi_agent_spec models."""

import json

import pytest
from pydantic import ValidationError

from multi_agent_spec import (
    # Enums
    Model,
    Platform,
    Priority,
    Tool,
    WorkflowType,
    # Agent models
    Agent,
    # Team/Orchestration models
    Step,
    Team,
    Workflow,
    # Deployment models
    AgentKitLocalConfig,
    AwsAgentCoreConfig,
    ClaudeCodeConfig,
    Deployment,
    KiroCliConfig,
    KubernetesConfig,
    ResourceLimits,
    Target,
    # Mappings
    BEDROCK_MODELS,
    CLAUDE_CODE_MODELS,
    KIRO_CLI_MODELS,
    KIRO_CLI_TOOLS,
)


# =============================================================================
# Enum Tests
# =============================================================================


class TestTool:
    """Tests for Tool enum."""

    def test_all_tools_exist(self) -> None:
        """Test that all expected tools are defined."""
        expected = [
            "WebSearch",
            "WebFetch",
            "Read",
            "Write",
            "Glob",
            "Grep",
            "Bash",
            "Edit",
            "Task",
        ]
        for tool_name in expected:
            assert hasattr(Tool, tool_name.upper().replace("SEARCH", "_SEARCH").replace("FETCH", "_FETCH"))

    def test_tool_values(self) -> None:
        """Test tool enum values."""
        assert Tool.WEB_SEARCH.value == "WebSearch"
        assert Tool.READ.value == "Read"
        assert Tool.WRITE.value == "Write"


class TestModel:
    """Tests for Model enum."""

    def test_all_models_exist(self) -> None:
        """Test that all expected models are defined."""
        assert Model.HAIKU.value == "haiku"
        assert Model.SONNET.value == "sonnet"
        assert Model.OPUS.value == "opus"


class TestWorkflowType:
    """Tests for WorkflowType enum."""

    def test_all_workflow_types(self) -> None:
        """Test all workflow types."""
        assert WorkflowType.SEQUENTIAL.value == "sequential"
        assert WorkflowType.PARALLEL.value == "parallel"
        assert WorkflowType.DAG.value == "dag"
        assert WorkflowType.ORCHESTRATED.value == "orchestrated"


class TestPlatform:
    """Tests for Platform enum."""

    def test_all_platforms(self) -> None:
        """Test all platforms."""
        expected = [
            "claude-code",
            "kiro-cli",
            "aws-agentcore",
            "aws-eks",
            "azure-aks",
            "gcp-gke",
            "kubernetes",
            "docker-compose",
            "agentkit-local",
        ]
        platform_values = [p.value for p in Platform]
        for exp in expected:
            assert exp in platform_values


class TestPriority:
    """Tests for Priority enum."""

    def test_all_priorities(self) -> None:
        """Test all priorities."""
        assert Priority.P1.value == "p1"
        assert Priority.P2.value == "p2"
        assert Priority.P3.value == "p3"


# =============================================================================
# Agent Tests
# =============================================================================


class TestAgent:
    """Tests for Agent model."""

    def test_valid_agent_all_fields(self) -> None:
        """Test creating agent with all fields."""
        agent = Agent(
            name="test-agent",
            description="A test agent",
            model=Model.SONNET,
            tools=["Read", "Write"],
            skills=["skill1"],
            dependencies=["other-agent"],
            instructions="You are a test agent.",
        )

        assert agent.name == "test-agent"
        assert agent.description == "A test agent"
        assert agent.model == "sonnet"  # enum converted to value
        assert agent.tools == ["Read", "Write"]
        assert agent.skills == ["skill1"]
        assert agent.dependencies == ["other-agent"]
        assert agent.instructions == "You are a test agent."

    def test_minimal_agent(self) -> None:
        """Test creating agent with minimal fields."""
        agent = Agent(name="minimal", description="Minimal agent")

        assert agent.name == "minimal"
        assert agent.model == "sonnet"  # default
        assert agent.tools == []
        assert agent.skills == []
        assert agent.dependencies == []
        assert agent.instructions is None

    def test_invalid_agent_name_uppercase(self) -> None:
        """Test that uppercase names are rejected."""
        with pytest.raises(ValidationError):
            Agent(name="Invalid-Agent", description="Test")

    def test_invalid_agent_name_starts_with_number(self) -> None:
        """Test that names starting with numbers are rejected."""
        with pytest.raises(ValidationError):
            Agent(name="123-agent", description="Test")

    def test_invalid_agent_name_underscore(self) -> None:
        """Test that names with underscores are rejected."""
        with pytest.raises(ValidationError):
            Agent(name="agent_name", description="Test")

    def test_agent_json_serialization(self) -> None:
        """Test JSON serialization."""
        agent = Agent(
            name="test",
            description="Test",
            model=Model.HAIKU,
            tools=["Read"],
        )
        json_str = agent.model_dump_json()
        data = json.loads(json_str)

        assert data["name"] == "test"
        assert data["model"] == "haiku"
        assert data["tools"] == ["Read"]

    def test_agent_from_dict(self) -> None:
        """Test creating agent from dict."""
        agent = Agent.model_validate({
            "name": "from-dict",
            "description": "From dict",
            "model": "opus",
        })

        assert agent.name == "from-dict"
        assert agent.model == "opus"


# =============================================================================
# Step Tests
# =============================================================================


class TestStep:
    """Tests for Step model."""

    def test_valid_step_all_fields(self) -> None:
        """Test creating step with all fields."""
        step = Step(
            name="research",
            agent="researcher",
            depends_on=["init"],
            inputs={"topic": "init.topic"},
            outputs=["results"],
        )

        assert step.name == "research"
        assert step.agent == "researcher"
        assert step.depends_on == ["init"]
        assert step.inputs == {"topic": "init.topic"}
        assert step.outputs == ["results"]

    def test_minimal_step(self) -> None:
        """Test creating step with minimal fields."""
        step = Step(name="step1", agent="agent1")

        assert step.name == "step1"
        assert step.agent == "agent1"
        assert step.depends_on is None
        assert step.inputs is None
        assert step.outputs is None


# =============================================================================
# Workflow Tests
# =============================================================================


class TestWorkflow:
    """Tests for Workflow model."""

    def test_workflow_default_type(self) -> None:
        """Test default workflow type."""
        workflow = Workflow()
        assert workflow.type == "orchestrated"

    def test_workflow_with_type(self) -> None:
        """Test workflow with explicit type."""
        workflow = Workflow(type=WorkflowType.DAG)
        assert workflow.type == "dag"

    def test_workflow_with_steps(self) -> None:
        """Test workflow with steps."""
        workflow = Workflow(
            type=WorkflowType.SEQUENTIAL,
            steps=[
                Step(name="s1", agent="a1"),
                Step(name="s2", agent="a2", depends_on=["s1"]),
            ],
        )

        assert len(workflow.steps) == 2
        assert workflow.steps[1].depends_on == ["s1"]


# =============================================================================
# Team Tests
# =============================================================================


class TestTeam:
    """Tests for Team model."""

    def test_valid_team_all_fields(self) -> None:
        """Test creating team with all fields."""
        team = Team(
            name="test-team",
            version="1.0.0",
            description="A test team",
            agents=["agent1", "agent2"],
            orchestrator="agent1",
            workflow=Workflow(type=WorkflowType.ORCHESTRATED),
            context="Shared context",
        )

        assert team.name == "test-team"
        assert team.version == "1.0.0"
        assert team.agents == ["agent1", "agent2"]
        assert team.orchestrator == "agent1"

    def test_minimal_team(self) -> None:
        """Test creating team with minimal fields."""
        team = Team(name="minimal", version="1.0.0", agents=["agent1"])

        assert team.name == "minimal"
        assert team.orchestrator is None
        assert team.workflow is None

    def test_invalid_version_format(self) -> None:
        """Test that invalid version format is rejected."""
        with pytest.raises(ValidationError):
            Team(name="team", version="1.0", agents=["a1"])

        with pytest.raises(ValidationError):
            Team(name="team", version="v1.0.0", agents=["a1"])

    def test_empty_agents_rejected(self) -> None:
        """Test that empty agents list is rejected."""
        with pytest.raises(ValidationError):
            Team(name="team", version="1.0.0", agents=[])


# =============================================================================
# Config Tests
# =============================================================================


class TestClaudeCodeConfig:
    """Tests for ClaudeCodeConfig."""

    def test_defaults(self) -> None:
        """Test default values."""
        config = ClaudeCodeConfig()
        assert config.agent_dir == ".claude/agents"
        assert config.format == "markdown"

    def test_custom_values(self) -> None:
        """Test custom values with alias."""
        config = ClaudeCodeConfig(agentDir="custom/path")
        assert config.agent_dir == "custom/path"


class TestKiroCliConfig:
    """Tests for KiroCliConfig."""

    def test_defaults(self) -> None:
        """Test default values."""
        config = KiroCliConfig()
        assert config.plugin_dir == "plugins/kiro/agents"
        assert config.format == "json"


class TestAwsAgentCoreConfig:
    """Tests for AwsAgentCoreConfig."""

    def test_defaults(self) -> None:
        """Test default values."""
        config = AwsAgentCoreConfig()
        assert config.region == "us-east-1"
        assert config.iac == "cdk"
        assert config.lambda_runtime == "python3.11"

    def test_custom_values(self) -> None:
        """Test custom values."""
        config = AwsAgentCoreConfig(
            region="eu-west-1",
            iac="pulumi",
            foundationModel="custom-model",
        )
        assert config.region == "eu-west-1"
        assert config.iac == "pulumi"


class TestResourceLimits:
    """Tests for ResourceLimits."""

    def test_defaults(self) -> None:
        """Test default values."""
        limits = ResourceLimits()
        assert limits.cpu == "500m"
        assert limits.memory == "512Mi"

    def test_custom_values(self) -> None:
        """Test custom values."""
        limits = ResourceLimits(cpu="1000m", memory="1Gi")
        assert limits.cpu == "1000m"
        assert limits.memory == "1Gi"


class TestKubernetesConfig:
    """Tests for KubernetesConfig."""

    def test_defaults(self) -> None:
        """Test default values."""
        config = KubernetesConfig()
        assert config.namespace == "multi-agent"
        assert config.helm_chart is True

    def test_with_resource_limits(self) -> None:
        """Test with resource limits."""
        config = KubernetesConfig(
            resourceLimits=ResourceLimits(cpu="2000m", memory="2Gi")
        )
        assert config.resource_limits is not None
        assert config.resource_limits.cpu == "2000m"


class TestAgentKitLocalConfig:
    """Tests for AgentKitLocalConfig."""

    def test_defaults(self) -> None:
        """Test default values."""
        config = AgentKitLocalConfig()
        assert config.transport == "stdio"
        assert config.port is None

    def test_http_transport(self) -> None:
        """Test HTTP transport with port."""
        config = AgentKitLocalConfig(transport="http", port=8080)
        assert config.transport == "http"
        assert config.port == 8080


# =============================================================================
# Target Tests
# =============================================================================


class TestTarget:
    """Tests for Target model."""

    def test_valid_target(self) -> None:
        """Test valid target."""
        target = Target(
            name="local-claude",
            platform=Platform.CLAUDE_CODE,
            priority=Priority.P1,
            output=".claude/agents",
        )

        assert target.name == "local-claude"
        assert target.platform == "claude-code"
        assert target.priority == "p1"

    def test_default_priority(self) -> None:
        """Test default priority."""
        target = Target(
            name="test",
            platform=Platform.KIRO_CLI,
            output="output",
        )
        assert target.priority == "p2"

    def test_with_config(self) -> None:
        """Test target with config."""
        target = Target(
            name="test",
            platform=Platform.AWS_AGENTCORE,
            output="cdk",
            config={"region": "us-west-2"},
        )
        assert target.config == {"region": "us-west-2"}


# =============================================================================
# Deployment Tests
# =============================================================================


class TestDeployment:
    """Tests for Deployment model."""

    def test_valid_deployment(self) -> None:
        """Test valid deployment."""
        deployment = Deployment(
            team="test-team",
            targets=[
                Target(
                    name="local",
                    platform=Platform.CLAUDE_CODE,
                    output=".claude/agents",
                )
            ],
        )

        assert deployment.team == "test-team"
        assert len(deployment.targets) == 1

    def test_empty_targets_rejected(self) -> None:
        """Test that empty targets is rejected."""
        with pytest.raises(ValidationError):
            Deployment(team="test", targets=[])

    def test_with_schema(self) -> None:
        """Test deployment with $schema field."""
        deployment = Deployment.model_validate({
            "$schema": "../schema/deployment.schema.json",
            "team": "test",
            "targets": [
                {"name": "t", "platform": "claude-code", "output": "out"}
            ],
        })
        assert deployment.schema_ == "../schema/deployment.schema.json"


# =============================================================================
# Mapping Tests
# =============================================================================


class TestMappings:
    """Tests for model and tool mappings."""

    def test_claude_code_models(self) -> None:
        """Test Claude Code model mappings."""
        assert CLAUDE_CODE_MODELS["haiku"] == "haiku"
        assert CLAUDE_CODE_MODELS["sonnet"] == "sonnet"
        assert CLAUDE_CODE_MODELS["opus"] == "opus"

    def test_kiro_cli_models(self) -> None:
        """Test Kiro CLI model mappings."""
        assert KIRO_CLI_MODELS["haiku"] == "claude-haiku-35"
        assert KIRO_CLI_MODELS["sonnet"] == "claude-sonnet-4"
        assert KIRO_CLI_MODELS["opus"] == "claude-opus-4"

    def test_bedrock_models(self) -> None:
        """Test Bedrock model mappings."""
        assert "haiku" in BEDROCK_MODELS["haiku"]
        assert "sonnet" in BEDROCK_MODELS["sonnet"]
        assert "opus" in BEDROCK_MODELS["opus"]

    def test_kiro_cli_tools(self) -> None:
        """Test Kiro CLI tool mappings."""
        assert KIRO_CLI_TOOLS["WebSearch"] == "web_search"
        assert KIRO_CLI_TOOLS["WebFetch"] == "web_fetch"
        assert KIRO_CLI_TOOLS["Read"] == "read"
        assert KIRO_CLI_TOOLS["Write"] == "write"
        assert KIRO_CLI_TOOLS["Glob"] == "glob"
        assert KIRO_CLI_TOOLS["Grep"] == "grep"
        assert KIRO_CLI_TOOLS["Bash"] == "bash"
        assert KIRO_CLI_TOOLS["Edit"] == "edit"
        assert KIRO_CLI_TOOLS["Task"] == "task"


# =============================================================================
# Version Tests
# =============================================================================


class TestVersion:
    """Tests for package version."""

    def test_version_exists(self) -> None:
        """Test that version is defined."""
        from multi_agent_spec import __version__
        assert __version__ == "1.0.0"

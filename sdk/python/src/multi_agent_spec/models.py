"""
Multi-Agent Spec - Pydantic Models

Runtime validation and Python type hints for multi-agent system definitions.

Example:
    >>> from multi_agent_spec import Agent, Team, Deployment
    >>>
    >>> agent = Agent(
    ...     name="my-agent",
    ...     description="A helpful agent",
    ...     model="sonnet",
    ...     tools=["Read", "Write"],
    ... )
    >>> print(agent.model_dump_json(indent=2))
"""

from __future__ import annotations

from enum import Enum
from typing import Literal

from pydantic import BaseModel, Field


# =============================================================================
# Enums
# =============================================================================


class Tool(str, Enum):
    """Canonical tool names available to agents."""

    WEB_SEARCH = "WebSearch"
    WEB_FETCH = "WebFetch"
    READ = "Read"
    WRITE = "Write"
    GLOB = "Glob"
    GREP = "Grep"
    BASH = "Bash"
    EDIT = "Edit"
    TASK = "Task"


class Model(str, Enum):
    """Model capability tiers."""

    HAIKU = "haiku"
    SONNET = "sonnet"
    OPUS = "opus"


class WorkflowType(str, Enum):
    """Workflow execution patterns."""

    SEQUENTIAL = "sequential"
    PARALLEL = "parallel"
    DAG = "dag"
    ORCHESTRATED = "orchestrated"


class Platform(str, Enum):
    """Supported deployment platforms."""

    CLAUDE_CODE = "claude-code"
    KIRO_CLI = "kiro-cli"
    AWS_AGENTCORE = "aws-agentcore"
    AWS_EKS = "aws-eks"
    AZURE_AKS = "azure-aks"
    GCP_GKE = "gcp-gke"
    KUBERNETES = "kubernetes"
    DOCKER_COMPOSE = "docker-compose"
    AGENTKIT_LOCAL = "agentkit-local"


class Priority(str, Enum):
    """Deployment priority levels."""

    P1 = "p1"
    P2 = "p2"
    P3 = "p3"


# =============================================================================
# Agent Definition Models
# =============================================================================


class Agent(BaseModel):
    """Agent definition model."""

    name: str = Field(
        ...,
        pattern=r"^[a-z][a-z0-9-]*$",
        description="Unique identifier for the agent (lowercase, hyphenated)",
    )
    description: str = Field(
        ..., description="Brief description of the agent's purpose and capabilities"
    )
    model: Model = Field(
        default=Model.SONNET,
        description="Model capability tier (mapped to platform-specific models)",
    )
    tools: list[str] = Field(
        default_factory=list, description="List of tools the agent can use"
    )
    skills: list[str] = Field(
        default_factory=list, description="List of skills the agent can invoke"
    )
    dependencies: list[str] = Field(
        default_factory=list,
        description="Other agents this agent depends on or can spawn",
    )
    instructions: str | None = Field(
        default=None, description="System prompt / instructions for the agent"
    )

    model_config = {"use_enum_values": True}


# =============================================================================
# Team / Orchestration Models
# =============================================================================


class Step(BaseModel):
    """Workflow step definition."""

    name: str = Field(..., description="Step identifier")
    agent: str = Field(..., description="Agent to execute this step")
    depends_on: list[str] | None = Field(
        default=None, description="Steps that must complete before this step"
    )
    inputs: dict[str, str] | None = Field(
        default=None, description="Input mappings from previous step outputs"
    )
    outputs: list[str] | None = Field(
        default=None, description="Named outputs from this step"
    )


class Workflow(BaseModel):
    """Workflow definition."""

    type: WorkflowType = Field(
        default=WorkflowType.ORCHESTRATED, description="Workflow execution pattern"
    )
    steps: list[Step] | None = Field(
        default=None, description="Ordered steps in the workflow"
    )

    model_config = {"use_enum_values": True}


class Team(BaseModel):
    """Team definition model."""

    name: str = Field(
        ...,
        pattern=r"^[a-z][a-z0-9-]*$",
        description="Team identifier (e.g., stats-agent-team)",
    )
    version: str = Field(
        ...,
        pattern=r"^\d+\.\d+\.\d+$",
        description="Semantic version of the team definition",
    )
    description: str | None = Field(
        default=None, description="Brief description of the team's purpose"
    )
    agents: list[str] = Field(
        ..., min_length=1, description="List of agent names in the team"
    )
    orchestrator: str | None = Field(
        default=None, description="Name of the orchestrator agent"
    )
    workflow: Workflow | None = Field(
        default=None, description="Workflow definition for agent coordination"
    )
    context: str | None = Field(
        default=None,
        description="Shared context or background information for all agents",
    )


# =============================================================================
# Deployment Models
# =============================================================================


class ClaudeCodeConfig(BaseModel):
    """Claude Code platform configuration."""

    agent_dir: str = Field(default=".claude/agents", alias="agentDir")
    format: Literal["markdown"] = "markdown"


class KiroCliConfig(BaseModel):
    """Kiro CLI platform configuration."""

    plugin_dir: str = Field(default="plugins/kiro/agents", alias="pluginDir")
    format: Literal["json"] = "json"


class AwsAgentCoreConfig(BaseModel):
    """AWS AgentCore platform configuration."""

    region: str = "us-east-1"
    foundation_model: str = Field(
        default="anthropic.claude-3-sonnet-20240229-v1:0", alias="foundationModel"
    )
    iac: Literal["cdk", "pulumi", "terraform"] = "cdk"
    lambda_runtime: str = Field(default="python3.11", alias="lambdaRuntime")


class ResourceLimits(BaseModel):
    """Kubernetes resource limits."""

    cpu: str = "500m"
    memory: str = "512Mi"


class KubernetesConfig(BaseModel):
    """Kubernetes platform configuration."""

    namespace: str = "multi-agent"
    helm_chart: bool = Field(default=True, alias="helmChart")
    image_registry: str | None = Field(default=None, alias="imageRegistry")
    resource_limits: ResourceLimits | None = Field(
        default=None, alias="resourceLimits"
    )


class AgentKitLocalConfig(BaseModel):
    """AgentKit local platform configuration."""

    transport: Literal["stdio", "http"] = "stdio"
    port: int | None = None


class Target(BaseModel):
    """Deployment target definition."""

    name: str = Field(..., description="Unique name for this deployment target")
    platform: Platform = Field(..., description="Target platform for deployment")
    priority: Priority = Field(default=Priority.P2, description="Deployment priority")
    output: str = Field(
        ..., description="Output directory for generated deployment artifacts"
    )
    config: dict | None = Field(
        default=None, description="Platform-specific configuration"
    )

    model_config = {"use_enum_values": True}


class Deployment(BaseModel):
    """Deployment definition model."""

    schema_: str | None = Field(default=None, alias="$schema")
    team: str = Field(..., description="Reference to the team definition (team name)")
    targets: list[Target] = Field(
        ..., min_length=1, description="List of deployment targets"
    )


# =============================================================================
# Model Mappings
# =============================================================================

CLAUDE_CODE_MODELS: dict[str, str] = {
    "haiku": "haiku",
    "sonnet": "sonnet",
    "opus": "opus",
}

KIRO_CLI_MODELS: dict[str, str] = {
    "haiku": "claude-haiku-35",
    "sonnet": "claude-sonnet-4",
    "opus": "claude-opus-4",
}

BEDROCK_MODELS: dict[str, str] = {
    "haiku": "anthropic.claude-3-haiku-20240307-v1:0",
    "sonnet": "anthropic.claude-3-5-sonnet-20241022-v2:0",
    "opus": "anthropic.claude-3-opus-20240229-v1:0",
}

KIRO_CLI_TOOLS: dict[str, str] = {
    "WebSearch": "web_search",
    "WebFetch": "web_fetch",
    "Read": "read",
    "Write": "write",
    "Glob": "glob",
    "Grep": "grep",
    "Bash": "bash",
    "Edit": "edit",
    "Task": "task",
}

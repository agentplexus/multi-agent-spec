"""
Multi-Agent Spec SDK for Python

Provides Pydantic models and Python types for defining multi-agent systems.

Example:
    >>> from multi_agent_spec import Agent, Team, Deployment
    >>>
    >>> # Create an agent
    >>> agent = Agent(
    ...     name="my-agent",
    ...     description="A helpful agent",
    ...     model="sonnet",
    ...     tools=["Read", "Write"],
    ... )
    >>>
    >>> # Serialize to JSON
    >>> print(agent.model_dump_json(indent=2))
    >>>
    >>> # Parse from dict
    >>> agent = Agent.model_validate({"name": "test", "description": "Test agent"})
"""

__version__ = "1.0.0"

from multi_agent_spec.models import (
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

__all__ = [
    # Version
    "__version__",
    # Enums
    "Tool",
    "Model",
    "WorkflowType",
    "Platform",
    "Priority",
    # Agent models
    "Agent",
    # Team/Orchestration models
    "Step",
    "Workflow",
    "Team",
    # Deployment models
    "ClaudeCodeConfig",
    "KiroCliConfig",
    "AwsAgentCoreConfig",
    "ResourceLimits",
    "KubernetesConfig",
    "AgentKitLocalConfig",
    "Target",
    "Deployment",
    # Mappings
    "CLAUDE_CODE_MODELS",
    "KIRO_CLI_MODELS",
    "BEDROCK_MODELS",
    "KIRO_CLI_TOOLS",
]

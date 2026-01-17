/**
 * Multi-Agent Spec - Zod Schemas
 *
 * Runtime validation and TypeScript type inference for multi-agent system definitions.
 */

import { z } from "zod";

// =============================================================================
// Agent Definition Schema
// =============================================================================

/** Canonical tool names available to agents */
export const ToolSchema = z.enum([
  "WebSearch",
  "WebFetch",
  "Read",
  "Write",
  "Glob",
  "Grep",
  "Bash",
  "Edit",
  "Task",
]);

/** Model capability tiers */
export const ModelSchema = z.enum(["haiku", "sonnet", "opus"]);

/** Agent definition schema */
export const AgentSchema = z.object({
  /** Unique identifier for the agent (lowercase, hyphenated) */
  name: z
    .string()
    .regex(/^[a-z][a-z0-9-]*$/, "Name must be lowercase with hyphens"),

  /** Brief description of the agent's purpose and capabilities */
  description: z.string(),

  /** Model capability tier (mapped to platform-specific models) */
  model: ModelSchema.default("sonnet"),

  /** List of tools the agent can use */
  tools: z.array(ToolSchema).default([]),

  /** List of skills the agent can invoke */
  skills: z.array(z.string()).default([]),

  /** Other agents this agent depends on or can spawn */
  dependencies: z.array(z.string()).default([]),

  /** System prompt / instructions for the agent */
  instructions: z.string().optional(),
});

// =============================================================================
// Team / Orchestration Schema
// =============================================================================

/** Workflow step definition */
export const StepSchema = z.object({
  /** Step identifier */
  name: z.string(),

  /** Agent to execute this step */
  agent: z.string(),

  /** Steps that must complete before this step */
  depends_on: z.array(z.string()).optional(),

  /** Input mappings from previous step outputs */
  inputs: z.record(z.string()).optional(),

  /** Named outputs from this step */
  outputs: z.array(z.string()).optional(),
});

/** Workflow execution pattern */
export const WorkflowTypeSchema = z.enum([
  "sequential",
  "parallel",
  "dag",
  "orchestrated",
]);

/** Workflow definition */
export const WorkflowSchema = z.object({
  /** Workflow execution pattern */
  type: WorkflowTypeSchema.default("orchestrated"),

  /** Ordered steps in the workflow */
  steps: z.array(StepSchema).optional(),
});

/** Team definition schema */
export const TeamSchema = z.object({
  /** Team identifier (e.g., stats-agent-team) */
  name: z
    .string()
    .regex(/^[a-z][a-z0-9-]*$/, "Name must be lowercase with hyphens"),

  /** Semantic version of the team definition */
  version: z.string().regex(/^\d+\.\d+\.\d+$/, "Must be semver format"),

  /** Brief description of the team's purpose */
  description: z.string().optional(),

  /** List of agent names in the team */
  agents: z.array(z.string()).min(1),

  /** Name of the orchestrator agent */
  orchestrator: z.string().optional(),

  /** Workflow definition for agent coordination */
  workflow: WorkflowSchema.optional(),

  /** Shared context or background information for all agents */
  context: z.string().optional(),
});

// =============================================================================
// Deployment Schema
// =============================================================================

/** Supported deployment platforms */
export const PlatformSchema = z.enum([
  "claude-code",
  "kiro-cli",
  "aws-agentcore",
  "aws-eks",
  "azure-aks",
  "gcp-gke",
  "kubernetes",
  "docker-compose",
  "agentkit-local",
]);

/** Deployment priority levels */
export const PrioritySchema = z.enum(["p1", "p2", "p3"]);

/** Claude Code platform config */
export const ClaudeCodeConfigSchema = z.object({
  agentDir: z.string().default(".claude/agents"),
  format: z.literal("markdown").default("markdown"),
});

/** Kiro CLI platform config */
export const KiroCliConfigSchema = z.object({
  pluginDir: z.string().default("plugins/kiro/agents"),
  format: z.literal("json").default("json"),
});

/** AWS AgentCore platform config */
export const AwsAgentCoreConfigSchema = z.object({
  region: z.string().default("us-east-1"),
  foundationModel: z
    .string()
    .default("anthropic.claude-3-sonnet-20240229-v1:0"),
  iac: z.enum(["cdk", "pulumi", "terraform"]).default("cdk"),
  lambdaRuntime: z.string().default("python3.11"),
});

/** Kubernetes platform config */
export const KubernetesConfigSchema = z.object({
  namespace: z.string().default("multi-agent"),
  helmChart: z.boolean().default(true),
  imageRegistry: z.string().optional(),
  resourceLimits: z
    .object({
      cpu: z.string().default("500m"),
      memory: z.string().default("512Mi"),
    })
    .optional(),
});

/** AgentKit local platform config */
export const AgentKitLocalConfigSchema = z.object({
  transport: z.enum(["stdio", "http"]).default("stdio"),
  port: z.number().optional(),
});

/** Platform-specific config (union type) */
export const PlatformConfigSchema = z.union([
  ClaudeCodeConfigSchema,
  KiroCliConfigSchema,
  AwsAgentCoreConfigSchema,
  KubernetesConfigSchema,
  AgentKitLocalConfigSchema,
  z.record(z.unknown()), // Allow unknown configs for extensibility
]);

/** Deployment target definition */
export const TargetSchema = z.object({
  /** Unique name for this deployment target */
  name: z.string(),

  /** Target platform for deployment */
  platform: PlatformSchema,

  /** Deployment priority */
  priority: PrioritySchema.default("p2"),

  /** Output directory for generated deployment artifacts */
  output: z.string(),

  /** Platform-specific configuration */
  config: PlatformConfigSchema.optional(),
});

/** Deployment definition schema */
export const DeploymentSchema = z.object({
  /** JSON Schema reference */
  $schema: z.string().optional(),

  /** Reference to the team definition (team name) */
  team: z.string(),

  /** List of deployment targets */
  targets: z.array(TargetSchema).min(1),
});

// =============================================================================
// Type Exports (inferred from schemas)
// =============================================================================

export type Tool = z.infer<typeof ToolSchema>;
export type Model = z.infer<typeof ModelSchema>;
export type Agent = z.infer<typeof AgentSchema>;

export type Step = z.infer<typeof StepSchema>;
export type WorkflowType = z.infer<typeof WorkflowTypeSchema>;
export type Workflow = z.infer<typeof WorkflowSchema>;
export type Team = z.infer<typeof TeamSchema>;

export type Platform = z.infer<typeof PlatformSchema>;
export type Priority = z.infer<typeof PrioritySchema>;
export type Target = z.infer<typeof TargetSchema>;
export type Deployment = z.infer<typeof DeploymentSchema>;

export type ClaudeCodeConfig = z.infer<typeof ClaudeCodeConfigSchema>;
export type KiroCliConfig = z.infer<typeof KiroCliConfigSchema>;
export type AwsAgentCoreConfig = z.infer<typeof AwsAgentCoreConfigSchema>;
export type KubernetesConfig = z.infer<typeof KubernetesConfigSchema>;
export type AgentKitLocalConfig = z.infer<typeof AgentKitLocalConfigSchema>;

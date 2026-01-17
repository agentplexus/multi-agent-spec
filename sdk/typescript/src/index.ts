/**
 * Multi-Agent Spec SDK for TypeScript
 *
 * Provides Zod schemas and TypeScript types for defining multi-agent systems.
 *
 * @example
 * ```typescript
 * import { AgentSchema, TeamSchema, DeploymentSchema } from '@agentplexus/multi-agent-spec';
 *
 * // Validate an agent definition
 * const agent = AgentSchema.parse({
 *   name: 'my-agent',
 *   description: 'A helpful agent',
 *   model: 'sonnet',
 *   tools: ['Read', 'Write'],
 * });
 *
 * // TypeScript types are inferred
 * import type { Agent, Team, Deployment } from '@agentplexus/multi-agent-spec';
 * ```
 */

// Re-export all schemas
export {
  // Agent schemas
  ToolSchema,
  ModelSchema,
  AgentSchema,
  // Team/Orchestration schemas
  StepSchema,
  WorkflowTypeSchema,
  WorkflowSchema,
  TeamSchema,
  // Deployment schemas
  PlatformSchema,
  PrioritySchema,
  TargetSchema,
  DeploymentSchema,
  // Platform config schemas
  ClaudeCodeConfigSchema,
  KiroCliConfigSchema,
  AwsAgentCoreConfigSchema,
  KubernetesConfigSchema,
  AgentKitLocalConfigSchema,
  PlatformConfigSchema,
} from "./schemas";

// Re-export all types
export type {
  // Agent types
  Tool,
  Model,
  Agent,
  // Team/Orchestration types
  Step,
  WorkflowType,
  Workflow,
  Team,
  // Deployment types
  Platform,
  Priority,
  Target,
  Deployment,
  // Platform config types
  ClaudeCodeConfig,
  KiroCliConfig,
  AwsAgentCoreConfig,
  KubernetesConfig,
  AgentKitLocalConfig,
} from "./schemas";

// =============================================================================
// Utility Functions
// =============================================================================

import {
  AgentSchema,
  TeamSchema,
  DeploymentSchema,
  type Agent,
  type Team,
  type Deployment,
} from "./schemas";

/**
 * Parse and validate an agent definition.
 * @throws ZodError if validation fails
 */
export function parseAgent(data: unknown): Agent {
  return AgentSchema.parse(data);
}

/**
 * Safely parse an agent definition, returning success/error result.
 */
export function safeParseAgent(data: unknown) {
  return AgentSchema.safeParse(data);
}

/**
 * Parse and validate a team definition.
 * @throws ZodError if validation fails
 */
export function parseTeam(data: unknown): Team {
  return TeamSchema.parse(data);
}

/**
 * Safely parse a team definition, returning success/error result.
 */
export function safeParseTeam(data: unknown) {
  return TeamSchema.safeParse(data);
}

/**
 * Parse and validate a deployment definition.
 * @throws ZodError if validation fails
 */
export function parseDeployment(data: unknown): Deployment {
  return DeploymentSchema.parse(data);
}

/**
 * Safely parse a deployment definition, returning success/error result.
 */
export function safeParseDeployment(data: unknown) {
  return DeploymentSchema.safeParse(data);
}

// =============================================================================
// Model Mappings
// =============================================================================

/** Map canonical model names to Claude Code model identifiers */
export const claudeCodeModels: Record<string, string> = {
  haiku: "haiku",
  sonnet: "sonnet",
  opus: "opus",
};

/** Map canonical model names to Kiro CLI model identifiers */
export const kiroCliModels: Record<string, string> = {
  haiku: "claude-haiku-35",
  sonnet: "claude-sonnet-4",
  opus: "claude-opus-4",
};

/** Map canonical model names to AWS Bedrock model identifiers */
export const bedrockModels: Record<string, string> = {
  haiku: "anthropic.claude-3-haiku-20240307-v1:0",
  sonnet: "anthropic.claude-3-5-sonnet-20241022-v2:0",
  opus: "anthropic.claude-3-opus-20240229-v1:0",
};

// =============================================================================
// Tool Mappings
// =============================================================================

/** Map canonical tool names to Kiro CLI tool identifiers */
export const kiroCliTools: Record<string, string> = {
  WebSearch: "web_search",
  WebFetch: "web_fetch",
  Read: "read",
  Write: "write",
  Glob: "glob",
  Grep: "grep",
  Bash: "bash",
  Edit: "edit",
  Task: "task",
};

/**
 * Auto-generated Zod schemas from JSON Schema.
 * DO NOT EDIT - regenerate with: npm run generate
 * Source: orchestration/team.schema.json
 */

import { z } from 'zod';

export const ChannelSchema = z.object({ "name": z.string(), "type": z.string(), "participants": z.array(z.string()).optional() }).strict();
export type Channel = z.infer<typeof ChannelSchema>;

export const ConsensusRulesSchema = z.object({ "required_agreement": z.number().optional(), "max_rounds": z.number().int().optional(), "tie_breaker": z.string().optional() }).strict();
export type ConsensusRules = z.infer<typeof ConsensusRulesSchema>;

export const CollaborationConfigSchema = z.object({ "lead": z.string().optional(), "specialists": z.array(z.string()).optional(), "task_queue": z.boolean().optional(), "consensus": z.object({ "required_agreement": z.number().optional(), "max_rounds": z.number().int().optional(), "tie_breaker": z.string().optional() }).strict().optional(), "channels": z.array(z.object({ "name": z.string(), "type": z.string(), "participants": z.array(z.string()).optional() }).strict()).optional() }).strict();
export type CollaborationConfig = z.infer<typeof CollaborationConfigSchema>;

export const PortTypeSchema = z.enum(["string","number","boolean","object","array","file"]).describe("Data type of a port");
export type PortType = z.infer<typeof PortTypeSchema>;

export const PortSchema = z.object({ "name": z.string(), "type": z.enum(["string","number","boolean","object","array","file"]).describe("Data type of a port").optional(), "description": z.string().optional(), "required": z.boolean().optional(), "from": z.string().optional(), "schema": z.any().optional(), "default": z.any().optional() }).strict();
export type Port = z.infer<typeof PortSchema>;

export const StepSchema = z.object({ "name": z.string(), "agent": z.string(), "depends_on": z.array(z.string()).optional(), "inputs": z.array(z.object({ "name": z.string(), "type": z.enum(["string","number","boolean","object","array","file"]).describe("Data type of a port").optional(), "description": z.string().optional(), "required": z.boolean().optional(), "from": z.string().optional(), "schema": z.any().optional(), "default": z.any().optional() }).strict()).optional(), "outputs": z.array(z.object({ "name": z.string(), "type": z.enum(["string","number","boolean","object","array","file"]).describe("Data type of a port").optional(), "description": z.string().optional(), "required": z.boolean().optional(), "from": z.string().optional(), "schema": z.any().optional(), "default": z.any().optional() }).strict()).optional() }).strict();
export type Step = z.infer<typeof StepSchema>;

export const WorkflowTypeSchema = z.enum(["sequential","parallel","dag","orchestrated"]).describe("Workflow execution pattern").default("orchestrated");
export type WorkflowType = z.infer<typeof WorkflowTypeSchema>;

export const WorkflowSchema = z.object({ "type": z.enum(["sequential","parallel","dag","orchestrated"]).describe("Workflow execution pattern").default("orchestrated"), "steps": z.array(z.object({ "name": z.string(), "agent": z.string(), "depends_on": z.array(z.string()).optional(), "inputs": z.array(z.object({ "name": z.string(), "type": z.enum(["string","number","boolean","object","array","file"]).describe("Data type of a port").optional(), "description": z.string().optional(), "required": z.boolean().optional(), "from": z.string().optional(), "schema": z.any().optional(), "default": z.any().optional() }).strict()).optional(), "outputs": z.array(z.object({ "name": z.string(), "type": z.enum(["string","number","boolean","object","array","file"]).describe("Data type of a port").optional(), "description": z.string().optional(), "required": z.boolean().optional(), "from": z.string().optional(), "schema": z.any().optional(), "default": z.any().optional() }).strict()).optional() }).strict()).optional() }).strict();
export type Workflow = z.infer<typeof WorkflowSchema>;

export const TeamSchema = z.object({ "name": z.string(), "version": z.string(), "description": z.string().optional(), "agents": z.array(z.string()), "orchestrator": z.string().optional(), "workflow": z.object({ "type": z.enum(["sequential","parallel","dag","orchestrated"]).describe("Workflow execution pattern").default("orchestrated"), "steps": z.array(z.object({ "name": z.string(), "agent": z.string(), "depends_on": z.array(z.string()).optional(), "inputs": z.array(z.object({ "name": z.string(), "type": z.enum(["string","number","boolean","object","array","file"]).describe("Data type of a port").optional(), "description": z.string().optional(), "required": z.boolean().optional(), "from": z.string().optional(), "schema": z.any().optional(), "default": z.any().optional() }).strict()).optional(), "outputs": z.array(z.object({ "name": z.string(), "type": z.enum(["string","number","boolean","object","array","file"]).describe("Data type of a port").optional(), "description": z.string().optional(), "required": z.boolean().optional(), "from": z.string().optional(), "schema": z.any().optional(), "default": z.any().optional() }).strict()).optional() }).strict()).optional() }).strict().optional(), "context": z.string().optional(), "collaboration": z.object({ "lead": z.string().optional(), "specialists": z.array(z.string()).optional(), "task_queue": z.boolean().optional(), "consensus": z.object({ "required_agreement": z.number().optional(), "max_rounds": z.number().int().optional(), "tie_breaker": z.string().optional() }).strict().optional(), "channels": z.array(z.object({ "name": z.string(), "type": z.string(), "participants": z.array(z.string()).optional() }).strict()).optional() }).strict().optional(), "self_claim": z.boolean().optional(), "plan_approval": z.boolean().optional() }).strict();
export type Team = z.infer<typeof TeamSchema>;


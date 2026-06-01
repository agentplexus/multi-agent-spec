/**
 * Auto-generated Zod schemas from JSON Schema.
 * DO NOT EDIT - regenerate with: npm run generate
 * Source: agent/agent.schema.json
 */

import { z } from 'zod';

export const ModelSchema = z.enum(["haiku","sonnet","opus"]).describe("Model capability tier (mapped to platform-specific models)").default("sonnet");
export type Model = z.infer<typeof ModelSchema>;

export const TaskTypeSchema = z.enum(["command","pattern","file","manual"]).describe("How the task is executed").default("manual");
export type TaskType = z.infer<typeof TaskTypeSchema>;

export const TaskSchema = z.object({ "id": z.string(), "description": z.string().optional(), "type": z.enum(["command","pattern","file","manual"]).describe("How the task is executed").default("manual"), "command": z.string().optional(), "pattern": z.string().optional(), "file": z.string().optional(), "files": z.string().optional(), "required": z.boolean().optional(), "expected_output": z.string().optional(), "human_in_loop": z.string().optional() }).strict();
export type Task = z.infer<typeof TaskSchema>;

export const DelegationConfigSchema = z.object({ "allow_delegation": z.boolean().optional(), "can_delegate_to": z.array(z.string()).optional(), "can_receive_from": z.array(z.string()).optional() }).strict();
export type DelegationConfig = z.infer<typeof DelegationConfigSchema>;

export const AgentSchema = z.object({ "name": z.string(), "namespace": z.string().optional(), "description": z.string().optional(), "icon": z.string().optional(), "model": z.enum(["haiku","sonnet","opus"]).describe("Model capability tier (mapped to platform-specific models)").default("sonnet"), "tools": z.array(z.string()).optional(), "allowedTools": z.array(z.string()).optional(), "skills": z.array(z.string()).optional(), "dependencies": z.array(z.string()).optional(), "requires": z.array(z.string()).optional(), "instructions": z.string().optional(), "tasks": z.array(z.object({ "id": z.string(), "description": z.string().optional(), "type": z.enum(["command","pattern","file","manual"]).describe("How the task is executed").default("manual"), "command": z.string().optional(), "pattern": z.string().optional(), "file": z.string().optional(), "files": z.string().optional(), "required": z.boolean().optional(), "expected_output": z.string().optional(), "human_in_loop": z.string().optional() }).strict()).optional(), "role": z.string().optional(), "goal": z.string().optional(), "backstory": z.string().optional(), "delegation": z.object({ "allow_delegation": z.boolean().optional(), "can_delegate_to": z.array(z.string()).optional(), "can_receive_from": z.array(z.string()).optional() }).strict().optional() }).strict();
export type Agent = z.infer<typeof AgentSchema>;


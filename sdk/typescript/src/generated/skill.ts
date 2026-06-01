/**
 * Auto-generated Zod schemas from JSON Schema.
 * DO NOT EDIT - regenerate with: npm run generate
 * Source: skill/skill.schema.json
 */

import { z } from 'zod';

export const ModelSchema = z.enum(["haiku","sonnet","opus"]).describe("Model capability tier (mapped to platform-specific models)").default("sonnet");
export type Model = z.infer<typeof ModelSchema>;

export const SkillSchema = z.object({ "name": z.string(), "description": z.string().optional(), "instructions": z.string().optional(), "scripts": z.array(z.string()).optional(), "references": z.array(z.string()).optional(), "assets": z.array(z.string()).optional(), "triggers": z.array(z.string()).optional(), "dependencies": z.array(z.string()).optional(), "model": z.enum(["haiku","sonnet","opus"]).describe("Model capability tier (mapped to platform-specific models)").default("sonnet"), "tools": z.array(z.string()).optional() }).strict();
export type Skill = z.infer<typeof SkillSchema>;


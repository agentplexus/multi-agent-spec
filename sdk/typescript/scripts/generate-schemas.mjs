#!/usr/bin/env node
/**
 * Generate Zod schemas from JSON Schema files.
 *
 * Usage: node scripts/generate-schemas.mjs
 */

import { readFileSync, writeFileSync, mkdirSync } from 'fs';
import { dirname, join, resolve } from 'path';
import { fileURLToPath } from 'url';
import { jsonSchemaToZod } from 'json-schema-to-zod';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const SCHEMA_DIR = resolve(__dirname, '../../../schema');
const OUTPUT_DIR = resolve(__dirname, '../src/generated');

// Schema files to process
const SCHEMAS = [
  { file: 'agent/agent.schema.json', name: 'Agent' },
  { file: 'orchestration/team.schema.json', name: 'Team' },
  { file: 'deployment/deployment.schema.json', name: 'Deployment' },
  { file: 'skill/skill.schema.json', name: 'Skill' },
];

/**
 * Resolve $ref references in a JSON Schema by inlining $defs.
 */
function resolveRefs(schema, defs) {
  if (typeof schema !== 'object' || schema === null) {
    return schema;
  }

  // Handle boolean schemas
  if (typeof schema === 'boolean') {
    return schema ? {} : { not: {} };
  }

  // Handle $ref
  if (schema.$ref) {
    const refPath = schema.$ref;
    if (refPath.startsWith('#/$defs/')) {
      const defName = refPath.slice(8);
      if (defs[defName]) {
        return resolveRefs(defs[defName], defs);
      }
    }
    // Return as-is if ref cannot be resolved
    return schema;
  }

  // Recursively resolve refs in nested objects and arrays
  const resolved = {};
  for (const [key, value] of Object.entries(schema)) {
    if (key === '$defs' || key === 'definitions') {
      continue; // Skip $defs, they're already being used for resolution
    }
    if (Array.isArray(value)) {
      resolved[key] = value.map(item => resolveRefs(item, defs));
    } else if (typeof value === 'object' && value !== null) {
      resolved[key] = resolveRefs(value, defs);
    } else {
      resolved[key] = value;
    }
  }
  return resolved;
}

/**
 * Generate Zod schema string from a JSON Schema.
 */
function generateZodFromSchema(jsonSchema) {
  const defs = jsonSchema.$defs || jsonSchema.definitions || {};

  // Resolve the root schema (usually a $ref)
  const resolvedSchema = resolveRefs(jsonSchema, defs);

  // Generate Zod code using the library
  const zodCode = jsonSchemaToZod(resolvedSchema, {
    module: 'none',
    name: undefined,
    withJsdoc: true,
  });

  return zodCode;
}

/**
 * Generate Zod definitions for each $def in the schema.
 */
function generateAllDefs(jsonSchema, mainName) {
  const defs = jsonSchema.$defs || jsonSchema.definitions || {};
  const output = [];
  const generatedNames = new Set();

  // Collect all definition names and their dependencies
  const defOrder = topologicalSort(defs);

  for (const defName of defOrder) {
    const def = defs[defName];
    if (!def) continue;

    try {
      // Resolve refs within this definition
      const resolvedDef = resolveRefs(def, defs);

      const zodCode = jsonSchemaToZod(resolvedDef, {
        module: 'none',
        name: undefined,
        withJsdoc: true,
      });

      output.push(`export const ${defName}Schema = ${zodCode};`);
      output.push(`export type ${defName} = z.infer<typeof ${defName}Schema>;`);
      output.push('');
      generatedNames.add(defName);
    } catch (err) {
      console.error(`Error generating ${defName}:`, err.message);
    }
  }

  return output.join('\n');
}

/**
 * Topological sort of definitions based on their $ref dependencies.
 */
function topologicalSort(defs) {
  const visited = new Set();
  const result = [];

  function getDeps(schema) {
    const deps = [];
    const traverse = (obj) => {
      if (typeof obj !== 'object' || obj === null) return;
      if (obj.$ref && obj.$ref.startsWith('#/$defs/')) {
        deps.push(obj.$ref.slice(8));
      }
      for (const value of Object.values(obj)) {
        if (Array.isArray(value)) {
          value.forEach(traverse);
        } else if (typeof value === 'object') {
          traverse(value);
        }
      }
    };
    traverse(schema);
    return deps;
  }

  function visit(name) {
    if (visited.has(name)) return;
    visited.add(name);

    const deps = getDeps(defs[name]);
    for (const dep of deps) {
      if (defs[dep]) {
        visit(dep);
      }
    }
    result.push(name);
  }

  for (const name of Object.keys(defs)) {
    visit(name);
  }

  return result;
}

function main() {
  // Ensure output directory exists
  mkdirSync(OUTPUT_DIR, { recursive: true });

  for (const { file, name } of SCHEMAS) {
    const schemaPath = join(SCHEMA_DIR, file);
    console.log(`Processing ${file}...`);

    try {
      const jsonSchema = JSON.parse(readFileSync(schemaPath, 'utf-8'));
      const zodCode = generateAllDefs(jsonSchema, name);

      const outputPath = join(OUTPUT_DIR, `${name.toLowerCase()}.ts`);
      const fileContent = `/**
 * Auto-generated Zod schemas from JSON Schema.
 * DO NOT EDIT - regenerate with: npm run generate
 * Source: ${file}
 */

import { z } from 'zod';

${zodCode}
`;

      writeFileSync(outputPath, fileContent);
      console.log(`  -> Generated ${outputPath}`);
    } catch (err) {
      console.error(`Error processing ${file}:`, err.message);
      process.exit(1);
    }
  }

  // Generate index file that re-exports everything
  const indexContent = SCHEMAS.map(({ name }) =>
    `export * from './${name.toLowerCase()}.js';`
  ).join('\n');

  writeFileSync(join(OUTPUT_DIR, 'index.ts'), `/**
 * Auto-generated Zod schemas index.
 * DO NOT EDIT - regenerate with: npm run generate
 */

${indexContent}
`);

  console.log('Done!');
}

main();

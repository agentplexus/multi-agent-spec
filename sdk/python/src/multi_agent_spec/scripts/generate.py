#!/usr/bin/env python3
"""
Generate Pydantic models from JSON Schema files.

Usage:
    pip install datamodel-code-generator
    python -m multi_agent_spec.scripts.generate

Or via entry point:
    generate-schemas
"""

import subprocess
import sys
from pathlib import Path


# Schema files to process
SCHEMAS = [
    ("agent/agent.schema.json", "agent"),
    ("orchestration/team.schema.json", "team"),
    ("deployment/deployment.schema.json", "deployment"),
    ("skill/skill.schema.json", "skill"),
]


def main() -> int:
    """Generate Pydantic models from JSON Schema."""
    # Find paths relative to this script
    script_dir = Path(__file__).parent
    sdk_dir = script_dir.parent.parent.parent  # sdk/python/
    schema_dir = sdk_dir.parent.parent / "schema"  # ../../schema/
    output_dir = sdk_dir / "src" / "multi_agent_spec" / "generated"

    # Ensure output directory exists
    output_dir.mkdir(parents=True, exist_ok=True)

    # Check if datamodel-codegen is available
    try:
        subprocess.run(
            ["datamodel-codegen", "--version"],
            capture_output=True,
            check=True,
        )
    except FileNotFoundError:
        print("Error: datamodel-codegen not found.", file=sys.stderr)
        print("Install with: pip install datamodel-code-generator", file=sys.stderr)
        return 1

    for schema_file, name in SCHEMAS:
        schema_path = schema_dir / schema_file
        output_path = output_dir / f"{name}.py"

        if not schema_path.exists():
            print(f"Warning: {schema_path} not found, skipping", file=sys.stderr)
            continue

        print(f"Processing {schema_file}...")

        result = subprocess.run(
            [
                "datamodel-codegen",
                "--input", str(schema_path),
                "--output", str(output_path),
                "--input-file-type", "jsonschema",
                "--output-model-type", "pydantic_v2.BaseModel",
                "--use-standard-collections",
                "--use-union-operator",
                "--field-constraints",
                "--target-python-version", "3.10",
            ],
            capture_output=True,
            text=True,
        )

        if result.returncode != 0:
            print(f"Error generating {name}: {result.stderr}", file=sys.stderr)
            return 1

        print(f"  -> Generated {output_path}")

    # Generate __init__.py for generated module
    init_path = output_dir / "__init__.py"
    exports = [f"from .{name} import *" for _, name in SCHEMAS]
    init_content = '"""Auto-generated Pydantic models from JSON Schema."""\n\n'
    init_content += "\n".join(exports) + "\n"
    init_path.write_text(init_content)
    print(f"  -> Generated {init_path}")

    print("Done!")
    return 0


if __name__ == "__main__":
    sys.exit(main())

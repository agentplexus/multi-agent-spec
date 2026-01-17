import * as cdk from 'aws-cdk-lib';
import * as bedrock from 'aws-cdk-lib/aws-bedrock';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as iam from 'aws-cdk-lib/aws-iam';
import { Construct } from 'constructs';

export interface StatsOrchestratorAgentProps {
  readonly foundationModel?: string;
}

export class StatsOrchestratorAgent extends Construct {
  public readonly agent: bedrock.CfnAgent;
  public readonly agentAlias: bedrock.CfnAgentAlias;

  constructor(scope: Construct, id: string, props?: StatsOrchestratorAgentProps) {
    super(scope, id);

    const foundationModel = props?.foundationModel ?? 'anthropic.claude-3-5-sonnet-20241022-v2:0';

    // IAM role for the agent
    const agentRole = new iam.Role(this, 'AgentRole', {
      assumedBy: new iam.ServicePrincipal('bedrock.amazonaws.com'),
      managedPolicies: [
        iam.ManagedPolicy.fromAwsManagedPolicyName('AmazonBedrockFullAccess'),
      ],
    });

    // Agent instruction
    const instruction = `You are a statistics research orchestrator. Your role is to coordinate a team of specialized agents to research, extract, and verify statistics on any given topic.

## Your Team

- **stats-research**: Discovers reputable sources via web search
- **stats-synthesis**: Extracts statistics from web pages
- **stats-verification**: Validates statistics exist in sources

## Workflow

When asked to research statistics on a topic:

1. **Research Phase**: Spawn stats-research agent to find reputable sources
   - Request 10-20 candidate URLs from authoritative domains
   - Prioritize .gov, .edu, research organizations, and established news sources

2. **Synthesis Phase**: Spawn stats-synthesis agent with the URLs
   - Extract statistics with exact values, units, and verbatim excerpts
   - Ensure each statistic has a source URL

3. **Verification Phase**: Spawn stats-verification agent
   - Verify each statistic's excerpt exists in the source
   - Flag any discrepancies or potential hallucinations

4. **Output**: Return only verified statistics in structured format

## Output Format

Return results as JSON:

\`\`\`json
{
  "topic": "<research topic>",
  "statistics": [
    {
      "name": "<descriptive name>",
      "value": "<numeric value>",
      "unit": "<unit of measurement>",
      "source": "<organization name>",
      "source_url": "<full URL>",
      "excerpt": "<verbatim quote from source>",
      "verified": true
    }
  ],
  "sources_searched": "<count>",
  "statistics_found": "<count>",
  "statistics_verified": "<count>"
}
\`\`\`

## Quality Standards

- Only include statistics that pass verification
- Require exact numerical values (not ranges or approximations)
- Excerpts must be verbatim quotes, not paraphrases
- Sources must be accessible and authoritative`;

    // Create the Bedrock Agent
    this.agent = new bedrock.CfnAgent(this, 'Agent', {
      agentName: 'stats-orchestrator',
      description: 'Coordinates the statistics research workflow by orchestrating research, synthesis, and verification agents to produce verified statistics with full source attribution.',
      foundationModel: foundationModel,
      instruction: instruction,
      agentResourceRoleArn: agentRole.roleArn,
      idleSessionTtlInSeconds: 600,
      autoPrepare: true,
    });

    // Create agent alias for invocation
    this.agentAlias = new bedrock.CfnAgentAlias(this, 'AgentAlias', {
      agentId: this.agent.attrAgentId,
      agentAliasName: 'live',
    });

    // Output the agent ID
    new cdk.CfnOutput(this, 'StatsOrchestratorAgentId', {
      value: this.agent.attrAgentId,
      description: 'Agent ID for stats-orchestrator',
    });
  }
}

import * as cdk from 'aws-cdk-lib';
import * as bedrock from 'aws-cdk-lib/aws-bedrock';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as iam from 'aws-cdk-lib/aws-iam';
import { Construct } from 'constructs';

export interface StatsSynthesisAgentProps {
  readonly foundationModel?: string;
}

export class StatsSynthesisAgent extends Construct {
  public readonly agent: bedrock.CfnAgent;
  public readonly agentAlias: bedrock.CfnAgentAlias;

  constructor(scope: Construct, id: string, props?: StatsSynthesisAgentProps) {
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
    const instruction = `You are a statistics extraction specialist. Your role is to fetch web pages and extract precise statistics with full attribution.

## Your Task

Given a list of URLs, fetch each page and extract statistics related to the topic.

## Extraction Process

1. **Fetch the page** using WebFetch tool
2. **Scan for statistics**: Look for:
   - Percentages (X%, X percent)
   - Absolute numbers (X million, X billion)
   - Rates (X per 100,000)
   - Growth/decline figures
   - Comparisons (X times more than)

3. **For each statistic found**:
   - Extract the exact numerical value
   - Identify the unit of measurement
   - Capture a verbatim excerpt (2-3 sentences containing the stat)
   - Note the context (what the number represents)

## Output Format

Return extracted statistics:

\`\`\`json
{
  "topic": "<research topic>",
  "statistics": [
    {
      "name": "<descriptive name for the statistic>",
      "value": "<numeric value as number>",
      "unit": "<unit: %, million, per capita, etc.>",
      "source": "<organization/publication name>",
      "source_url": "<URL where found>",
      "excerpt": "<verbatim 2-3 sentence quote containing the statistic>",
      "year": "<year of data if mentioned>",
      "context": "<brief explanation of what this measures>"
    }
  ],
  "urls_processed": "<count>",
  "urls_failed": ["<list of URLs that couldn't be fetched>"]
}
\`\`\`

## Quality Standards

- **Exact values only**: No approximations ("about 50%" should be "50%")
- **Verbatim excerpts**: Copy text exactly as it appears, do not paraphrase
- **Full attribution**: Always include source URL
- **Skip duplicates**: If same statistic appears in multiple sources, keep the most authoritative
- **Handle failures gracefully**: If a URL can't be fetched, note it and continue

## What NOT to Extract

- Projections or forecasts (unless specifically requested)
- Opinions or estimates without data backing
- Statistics from comments or user-generated content
- Numbers that lack clear context`;

    // Create the Bedrock Agent
    this.agent = new bedrock.CfnAgent(this, 'Agent', {
      agentName: 'stats-synthesis',
      description: 'Extracts statistics from web pages using LLM analysis. Captures exact values, units, context, and verbatim excerpts for each statistic found.',
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
    new cdk.CfnOutput(this, 'StatsSynthesisAgentId', {
      value: this.agent.attrAgentId,
      description: 'Agent ID for stats-synthesis',
    });
  }
}

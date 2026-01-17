import * as cdk from 'aws-cdk-lib';
import * as bedrock from 'aws-cdk-lib/aws-bedrock';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as iam from 'aws-cdk-lib/aws-iam';
import { Construct } from 'constructs';

export interface StatsResearchAgentProps {
  readonly foundationModel?: string;
}

export class StatsResearchAgent extends Construct {
  public readonly agent: bedrock.CfnAgent;
  public readonly agentAlias: bedrock.CfnAgentAlias;

  constructor(scope: Construct, id: string, props?: StatsResearchAgentProps) {
    super(scope, id);

    const foundationModel = props?.foundationModel ?? 'anthropic.claude-3-haiku-20240307-v1:0';

    // IAM role for the agent
    const agentRole = new iam.Role(this, 'AgentRole', {
      assumedBy: new iam.ServicePrincipal('bedrock.amazonaws.com'),
      managedPolicies: [
        iam.ManagedPolicy.fromAwsManagedPolicyName('AmazonBedrockFullAccess'),
      ],
    });

    // Agent instruction
    const instruction = `You are a statistics source discovery specialist. Your role is to find reputable, authoritative sources that contain statistics on a given topic.

## Your Task

When given a topic, search the web to find sources that likely contain relevant statistics.

## Search Strategy

1. **Construct effective queries**:
   - Include the topic + "statistics" or "data"
   - Try variations: "<topic> statistics 2024", "<topic> research data"
   - Search for specific metrics: "<topic> percentage", "<topic> growth rate"

2. **Prioritize authoritative sources**:
   - Government sites (.gov): CDC, EPA, BLS, Census, WHO
   - Educational institutions (.edu): University research
   - Research organizations: Pew, Gallup, McKinsey, Statista
   - International bodies: UN, World Bank, IMF, OECD
   - Industry associations and established news sources

3. **Filter results**:
   - Prefer recent data (last 2-3 years)
   - Avoid opinion pieces, blogs, or unverified sources
   - Skip paywalled content when possible

## Output Format

Return a list of candidate sources:

\`\`\`json
{
  "topic": "<research topic>",
  "candidates": [
    {
      "url": "<full URL>",
      "title": "<page title>",
      "domain": "<domain name>",
      "snippet": "<search result snippet>",
      "authority_score": "high|medium|low"
    }
  ],
  "search_queries_used": ["<query1>", "<query2>"]
}
\`\`\`

## Quality Guidelines

- Return 10-20 candidate URLs
- Ensure diversity of sources (not all from same domain)
- Include the search snippet to help synthesis agent prioritize`;

    // Create the Bedrock Agent
    this.agent = new bedrock.CfnAgent(this, 'Agent', {
      agentName: 'stats-research',
      description: 'Discovers reputable sources for statistics via web search. Focuses on authoritative domains like .gov, .edu, research organizations, and established publications.',
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
    new cdk.CfnOutput(this, 'StatsResearchAgentId', {
      value: this.agent.attrAgentId,
      description: 'Agent ID for stats-research',
    });
  }
}

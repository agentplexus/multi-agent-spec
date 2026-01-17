import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';

import { StatsOrchestratorAgent } from './agents/stats-orchestrator';

import { StatsResearchAgent } from './agents/stats-research';

import { StatsSynthesisAgent } from './agents/stats-synthesis';

import { StatsVerificationAgent } from './agents/stats-verification';


export interface StatsAgentTeamStackProps extends cdk.StackProps {
  readonly foundationModel?: string;
}

export class StatsAgentTeamStack extends cdk.Stack {
  public readonly statsOrchestratorAgent: StatsOrchestratorAgent;
  public readonly statsResearchAgent: StatsResearchAgent;
  public readonly statsSynthesisAgent: StatsSynthesisAgent;
  public readonly statsVerificationAgent: StatsVerificationAgent;


  constructor(scope: Construct, id: string, props?: StatsAgentTeamStackProps) {
    super(scope, id, props);

    const foundationModel = props?.foundationModel ?? 'anthropic.claude-3-sonnet-20240229-v1:0';

    // StatsOrchestrator Agent
    this.statsOrchestratorAgent = new StatsOrchestratorAgent(this, 'StatsOrchestrator', {
      foundationModel,
    });

    // StatsResearch Agent
    this.statsResearchAgent = new StatsResearchAgent(this, 'StatsResearch', {
      foundationModel,
    });

    // StatsSynthesis Agent
    this.statsSynthesisAgent = new StatsSynthesisAgent(this, 'StatsSynthesis', {
      foundationModel,
    });

    // StatsVerification Agent
    this.statsVerificationAgent = new StatsVerificationAgent(this, 'StatsVerification', {
      foundationModel,
    });

  }
}

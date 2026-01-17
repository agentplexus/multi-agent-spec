#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { StatsAgentTeamStack } from '../lib/stats-agent-team-stack';

const app = new cdk.App();

new StatsAgentTeamStack(app, 'StatsAgentTeamStack', {
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION ?? 'us-east-1',
  },
});

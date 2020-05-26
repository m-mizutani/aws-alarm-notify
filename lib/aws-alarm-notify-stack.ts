import * as cdk from "@aws-cdk/core";
import * as lambda from "@aws-cdk/aws-lambda";
import * as sqs from "@aws-cdk/aws-sqs";
import * as iam from "@aws-cdk/aws-iam";
import * as sns from "@aws-cdk/aws-sns";

import { SqsEventSource, SnsDlq } from "@aws-cdk/aws-lambda-event-sources";
import { SqsSubscription } from "@aws-cdk/aws-sns-subscriptions";

interface snsTopic {
  id: string;
  arn: string;
}
interface AlarmNotifyArguments {
  snsTopics: Array<snsTopic>;
  iamRole: string;
  slackURL: string;
  sentryDSN?: string;
}

export class AlarmNotifyStack extends cdk.Stack {
  constructor(
    scope: cdk.Construct,
    id: string,
    args: AlarmNotifyArguments,
    props?: cdk.StackProps
  ) {
    super(scope, id, props);

    const buildPath = lambda.Code.asset("./build");
    const lambdaRole = iam.Role.fromRoleArn(this, "LambdaRole", args.iamRole, {
      mutable: false,
    });

    // SQS
    const notifyQueue = new sqs.Queue(this, "NotifyQueue", {
      visibilityTimeout: cdk.Duration.seconds(10),
    });

    for (let topic of args.snsTopics) {
      const snsTopic = sns.Topic.fromTopicArn(this, topic.id, topic.arn);
      snsTopic.addSubscription(new SqsSubscription(notifyQueue));
    }

    new lambda.Function(this, "Handler", {
      runtime: lambda.Runtime.GO_1_X,
      handler: "Handler",
      code: buildPath,
      role: lambdaRole,
      timeout: cdk.Duration.seconds(10),
      events: [new SqsEventSource(notifyQueue, { batchSize: 1 })],
      environment: {
        SLACK_URL: args.slackURL,
        SENTRY_DSN: args.sentryDSN || "",
      },
    });
  }
}

#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "@aws-cdk/core";

import { AlarmNotifyStack } from "../lib/aws-alarm-notify-stack";

const app = new cdk.App();

new AlarmNotifyStack(app, "AwsAlarmNotifyStack", {
  snsTopics: [
    { id: "no1", arn: "arn:aws:sns:ap-northeast-1:1234567890:topic-1" },
    { id: "no2", arn: "arn:aws:sns:ap-northeast-1:1234567890:topic-2" },
  ],
  iamRole: 'arn:aws:iam::1234567890:role/MyNotifyRole"',
  slackURL: "https://example.com/hogehoge",
});

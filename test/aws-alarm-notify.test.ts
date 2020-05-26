import {
  expect as expectCDK,
  matchTemplate,
  MatchStyle,
} from "@aws-cdk/assert";
import * as cdk from "@aws-cdk/core";
// import * as AwsAlarmNotify from '../lib/aws-alarm-notify-stack';

test("Empty Stack", () => {
  const app = new cdk.App();
  // WHEN
  /*
    const stack = new AwsAlarmNotify.AwsAlarmNotifyStack(app, 'MyTestStack');
    // THEN
    expectCDK(stack).to(matchTemplate({
      "Resources": {}
    }, MatchStyle.EXACT))
    */
});

package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeSampleEvent() events.SQSEvent {
	// Data from https://github.com/blueimp/aws-lambda/blob/master/cloudwatch-alarm-to-slack/test-event.json
	eventData := `{
		"Records": [
		{
			"EventSource": "aws:sns",
			"EventVersion": "1.0",
			"EventSubscriptionArn": "arn:aws:sns:eu-west-1:000000000000:cloudwatch-alarms:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			"Sns": {
			"Type": "Notification",
			"MessageId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			"TopicArn": "arn:aws:sns:eu-west-1:000000000000:cloudwatch-alarms",
			"Subject": "ALARM: \"Example alarm name\" in EU - Ireland",
			"Message": "{\"AlarmName\":\"Example alarm name\",\"AlarmDescription\":\"Example alarm description.\",\"AWSAccountId\":\"000000000000\",\"NewStateValue\":\"ALARM\",\"NewStateReason\":\"Threshold Crossed: 1 datapoint (10.0) was greater than or equal to the threshold (1.0).\",\"StateChangeTime\":\"2017-01-12T16:30:42.236+0000\",\"Region\":\"EU - Ireland\",\"OldStateValue\":\"OK\",\"Trigger\":{\"MetricName\":\"DeliveryErrors\",\"Namespace\":\"ExampleNamespace\",\"Statistic\":\"SUM\",\"Unit\":null,\"Dimensions\":[],\"Period\":300,\"EvaluationPeriods\":1,\"ComparisonOperator\":\"GreaterThanOrEqualToThreshold\",\"Threshold\":1.0}}",
			"Timestamp": "2017-01-12T16:30:42.318Z",
			"SignatureVersion": "1",
			"Signature": "Cg==",
			"SigningCertUrl": "https://sns.eu-west-1.amazonaws.com/SimpleNotificationService-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.pem",
			"UnsubscribeUrl": "https://sns.eu-west-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:eu-west-1:000000000000:cloudwatch-alarms:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			"MessageAttributes": {}
			}
		}
		]
	}`

	var snsEvent events.SNSEvent
	if err := json.Unmarshal([]byte(eventData), &snsEvent); err != nil {
		panic(err)
	}
	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: eventData,
			},
		},
	}

	return sqsEvent
}

func TestHandler(t *testing.T) {
	called := 0
	sqsEvent := makeSampleEvent()
	args := arguments{
		Event:    sqsEvent,
		SlackURL: "https://example.com/hoge",
		post: func(url string, msg *slack.WebhookMessage) error {
			assert.Equal(t, "https://example.com/hoge", url)
			assert.Equal(t, 1, len(msg.Attachments))
			assert.Equal(t, "#e03e2f", msg.Attachments[0].Color)
			assert.Equal(t, "CloudWatch Alarm: Example alarm name", msg.Attachments[0].Title)
			assert.Equal(t, "Example alarm description.", msg.Attachments[0].Text)
			called++
			return nil
		},
	}
	require.NoError(t, handler(args))
	assert.Equal(t, 1, called)
}

func TestHandlerWithSlackURL(t *testing.T) {
	url, ok := os.LookupEnv("TEST_SLACK_URL")
	if !ok {
		t.Skip("TEST_SLACK_URL is not found")
	}

	sqsEvent := makeSampleEvent()
	args := arguments{
		Event:    sqsEvent,
		SlackURL: url,
		post:     slack.PostWebhook,
	}
	require.NoError(t, handler(args))
}

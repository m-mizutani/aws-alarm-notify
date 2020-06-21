package main

import (
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
	"Message": "{\"AlarmName\":\"Example alarm name\",\"AlarmDescription\":\"Example alarm description.\",\"AWSAccountId\":\"000000000000\",\"NewStateValue\":\"ALARM\",\"NewStateReason\":\"Threshold Crossed: 1 datapoint (10.0) was greater than or equal to the threshold (1.0).\",\"StateChangeTime\":\"2017-01-12T16:30:42.236+0000\",\"Region\":\"EU - Ireland\",\"OldStateValue\":\"OK\",\"Trigger\":{\"MetricName\":\"DeliveryErrors\",\"Namespace\":\"ExampleNamespace\",\"Statistic\":\"SUM\",\"Unit\":null,\"Dimensions\":[],\"Period\":300,\"EvaluationPeriods\":1,\"ComparisonOperator\":\"GreaterThanOrEqualToThreshold\",\"Threshold\":1.0}}"
}`

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

package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

// func PostWebhook(url string, msg *WebhookMessage) error
type slackPostWebhook func(url string, msg *slack.WebhookMessage) error

func postAlarm(post slackPostWebhook, url string, alarm cloudWatchAlarm) error {
	colorMap := map[string]string{
		"ALARM": "#e03e2f",
	}
	color := "#2eb886"
	if newColor, ok := colorMap[aws.StringValue(alarm.NewStateValue)]; ok {
		color = newColor
	}

	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{
			{
				Title: "CloudWatch Alarm: " + aws.StringValue(alarm.AlarmName),
				Text:  aws.StringValue(alarm.AlarmDescription),
				Color: color,
				Fields: []slack.AttachmentField{
					{
						Title: "Region",
						Value: aws.StringValue(alarm.Region),
					},
					{
						Title: "StateChangeTime",
						Value: aws.StringValue(alarm.StateChangeTime),
					},
					{
						Title: "NewStateValue",
						Value: aws.StringValue(alarm.NewStateValue),
					},
					{
						Title: "NewStateReason",
						Value: aws.StringValue(alarm.NewStateReason),
					},
				},
			},
		},
	}

	if err := post(url, &msg); err != nil {
		return errors.Wrap(err, "Fail to post CW alarm to slack")
	}

	return nil
}

package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

var logger = logrus.New()

type arguments struct {
	SlackURL string `env:"SLACK_URL"`
	Event    events.SQSEvent

	post slackPostWebhook
}

func main() {
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	lambda.Start(func(ctx context.Context, event events.SQSEvent) error {
		args := arguments{
			Event: event,
			post:  slack.PostWebhook,
		}
		if err := env.Parse(&args); err != nil {
			return errors.Wrap(err, "Fail to parse environment variable")
		}

		if err := handler(args); err != nil {
			logger.WithError(err).WithField("args", args).Error("Failed handler")
			return err
		}

		return nil
	})
}

type snsMessage struct {
	Message          string `json:"Message"`
	MessageID        string `json:"MessageId"`
	Signature        string `json:"Signature"`
	SignatureVersion string `json:"SignatureVersion"`
	SigningCertURL   string `json:"SigningCertURL"`
	Timestamp        string `json:"Timestamp"`
	TopicArn         string `json:"TopicArn"`
	Type             string `json:"Type"`
	UnsubscribeURL   string `json:"UnsubscribeURL"`
}

func handler(args arguments) error {
	logger.WithField("args", args).Info("Start handler")

	for _, sqsRecord := range args.Event.Records {
		var message snsMessage

		if err := json.Unmarshal([]byte(sqsRecord.Body), &message); err != nil {
			logger.WithField("body", sqsRecord.Body).WithError(err).Error("Failed")
			return errors.Wrap(err, "Fail to unmarshal SQS body")
		}

		var alarm cloudWatchAlarm
		if err := json.Unmarshal([]byte(message.Message), &alarm); err != nil {
			logger.WithField("message", message.Message).WithError(err).Error("Failed")
			return errors.Wrap(err, "Fail to unmarshal SNS message in SQS")
		}

		if err := postAlarm(args.post, args.SlackURL, alarm); err != nil {
			return err

		}
	}
	return nil
}

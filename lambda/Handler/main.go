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

func handler(args arguments) error {
	logger.WithField("args", args).Info("Start handler")

	for _, sqsRecord := range args.Event.Records {
		var snsEvent events.SNSEvent

		if err := json.Unmarshal([]byte(sqsRecord.Body), &snsEvent); err != nil {
			logger.WithField("body", sqsRecord.Body).WithError(err).Error("Failed")
			return errors.Wrap(err, "Fail to unmarshal SQS body")
		}

		for _, snsRecord := range snsEvent.Records {
			var alarm cloudWatchAlarm
			if err := json.Unmarshal([]byte(snsRecord.SNS.Message), &alarm); err != nil {
				logger.WithField("message", snsRecord.SNS.Message).WithError(err).Error("Failed")
				return errors.Wrap(err, "Fail to unmarshal SNS message in SQS")
			}

			if err := postAlarm(args.post, args.SlackURL, alarm); err != nil {
				return err
			}
		}

	}
	return nil
}

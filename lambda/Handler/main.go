package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

type arguments struct {
	SlackURL string `env:"SLACK_URL"`
	Event    events.SQSEvent
}

func main() {
	lambda.Start(func(ctx context.Context, event events.SQSEvent) error {
		args := arguments{
			Event: event,
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
	return nil
}

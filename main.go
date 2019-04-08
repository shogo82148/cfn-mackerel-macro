package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shogo82148/cfn-mackerel-macro/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel/apikey/aws"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	s := os.Getenv("FORWARD_LOG_LEVEL")
	if s != "" {
		level, err := logrus.ParseLevel(s)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"input": level,
				"error": err,
			}).Error("fail to parse log level")
		} else {
			logrus.SetLevel(level)
		}
	}
}

func main() {
	provider, err := aws.LoadDefaultProvider()
	if err != nil {
		logrus.WithError(err).Error("fail to load aws config")
	}
	f := cfn.Function{
		APIKeyProvider: provider,
	}
	lambda.Start(f.LambdaWrap())
}

package main

import (
	"net/url"
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
	logrus.Infof("cfn_mackerel_macro v%s", version)

	provider, err := aws.LoadDefaultProvider()
	if err != nil {
		logrus.WithError(err).Error("fail to load aws config")
		os.Exit(1)
	}

	var u *url.URL
	if base := os.Getenv("MACKEREL_APIURL"); base != "" {
		var err error
		u, err = url.Parse(base)
		if err != nil {
			logrus.WithError(err).Error("fail to parse base url")
			os.Exit(1)
		}
	}

	f := cfn.Function{
		APIKeyProvider: provider,
		BaseURL:        u,
	}
	lambda.Start(f.LambdaWrap())
}

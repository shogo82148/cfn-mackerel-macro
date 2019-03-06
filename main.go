package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shogo82148/cfn-mackerel-macro/cfn"
)

func main() {
	apikey := os.Getenv("MACKEREL_APIKEY")
	f := cfn.Function{
		APIKey: apikey,
	}
	lambda.Start(f.LambdaWrap())
}

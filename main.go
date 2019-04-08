package main

import (
	"github.com/shogo82148/cfn-mackerel-macro/mackerel/apikey"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shogo82148/cfn-mackerel-macro/cfn"
)

func main() {
	f := cfn.Function{
		APIKeyProvider: apikey.NewEnvironment("MACKEREL_APIKEY"),
	}
	lambda.Start(f.LambdaWrap())
}

package main

import (
	"context"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, event *cfn.Event) error {
	resp := cfn.NewResponse(event)

	// TODO: implement
	resp.Status = cfn.StatusSuccess
	resp.PhysicalResourceID = "mackerel:foobar"

	return resp.Send()
}

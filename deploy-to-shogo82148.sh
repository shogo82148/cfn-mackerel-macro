#!/bin/bash

# deploy the author's AWS Account for testing

CURRENT=$(cd "$(dirname "$0")" && pwd)

sam package \
    --region ap-northeast-1 \
    --template-file "$CURRENT/template.yaml" \
    --output-template-file "$CURRENT/packaged.yaml" \
    --s3-bucket "shogo82148-test" \
    --s3-prefix cfn-mackerel-macro/resource

sam deploy \
    --region ap-northeast-1 \
    --template-file "$CURRENT/packaged.yaml" \
    --capabilities CAPABILITY_IAM \
    --stack-name cfn-mackerel-macro

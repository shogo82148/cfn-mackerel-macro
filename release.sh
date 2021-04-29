#!/usr/bin/env bash

ROOT=$(cd "$(dirname "$0")" && pwd)

set -uex

VERSION=$(cat VERSION)
MAJOR=$(echo "$VERSION" | cut -d. -f 1)
MINOR=$(echo "$VERSION" | cut -d. -f 2)
PATCH=$(echo "$VERSION" | cut -d. -f 3)

DIST=$ROOT/.build/$VERSION
mkdir -p "$DIST"

make all
cp macro.zip "$DIST"
cp resource.zip "$DIST"
cp template.yaml "$DIST"
cp README.md "$DIST"
cp LICENSE "$DIST"

: publish as a CloudFormation template
cd "$DIST"
while read -r REGION; do
    BUCKET=shogo82148-cloudformation-template-$REGION
    aws cloudformation package \
        --region "$REGION" \
        --template-file "$DIST/template.yaml" \
        --output-template-file "$DIST/$REGION.yaml" \
        --s3-bucket "$BUCKET" \
        --s3-prefix cfn-mackerel-macro/resource
    aws s3 cp --region "$REGION" "$DIST/$REGION.yaml" "s3://$BUCKET/cfn-mackerel-macro/latest.yaml"
    aws s3 cp --region "$REGION" "$DIST/$REGION.yaml" "s3://$BUCKET/cfn-mackerel-macro/v$MAJOR.$MINOR.$PATCH.yaml"
    aws s3 cp --region "$REGION" "$DIST/$REGION.yaml" "s3://$BUCKET/cfn-mackerel-macro/v$MAJOR.$MINOR.yaml"
    aws s3 cp --region "$REGION" "$DIST/$REGION.yaml" "s3://$BUCKET/cfn-mackerel-macro/v$MAJOR.yaml"
done << EOS
af-south-1
ap-east-1
ap-northeast-1
ap-northeast-2
ap-northeast-3
ap-south-1
ap-southeast-1
ap-southeast-2
ca-central-1
eu-central-1
eu-north-1
eu-south-1
eu-west-1
eu-west-2
eu-west-3
me-south-1
sa-east-1
us-east-1
us-east-2
us-west-1
us-west-2
EOS

cd "$ROOT"
( git add . && git commit -m "bump up to v$VERSION" && git push ) || true
ghr -u shogo82148 --draft --replace "v$VERSION" "$DIST"

: publish to AWS Serverless Application Repository
DIST_SAM=$ROOT/.build-sam/$VERSION
mkdir -p "$DIST_SAM"
cp README.md "$DIST_SAM"
cp LICENSE "$DIST_SAM"
sam package \
    --region us-east-1 \
    --template-file "template.yaml" \
    --output-template-file "$DIST_SAM/packaged.yaml" \
    --s3-bucket shogo82148-sam
sam publish \
    --region us-east-1 \
    --template "$DIST_SAM/packaged.yaml"

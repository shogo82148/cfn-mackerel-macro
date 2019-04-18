#!/usr/bin/env bash

ROOT=$(cd "$(dirname "$0")" && pwd)

set -ue

MAJOR=0
MINOR=0
PATCH=2
VERSION="$MAJOR.$MINOR.$PATCH"

cat << EOS > version.go
package main

const version = "$VERSION"
EOS

DIST=$ROOT/.build/$VERSION
mkdir -p "$DIST"
perl -pe 's!__VERSION__!'"$VERSION"'!g' template.yaml > "$DIST/template.yaml"

make all
cp macro.zip "$DIST"
cp resource.zip "$DIST"

cd "$DIST"
while read -r REGION; do
    BUCKET=shogo82148-cloudformation-template-$REGION
    aws cloudformation package \
        --template-file "$DIST/template.yaml" \
        --output-template-file "$DIST/$REGION.yaml" \
        --s3-bucket "$BUCKET" \
        --s3-prefix cfn-mackerel-macro/resource
    aws s3 cp "$DIST/$REGION.yaml" "s3://$BUCKET/cfn-mackerel-macro/latest.yaml"
    aws s3 cp "$DIST/$REGION.yaml" "s3://$BUCKET/cfn-mackerel-macro/v$MAJOR.$MINOR.$PATCH.yaml"
    aws s3 cp "$DIST/$REGION.yaml" "s3://$BUCKET/cfn-mackerel-macro/v$MAJOR.$MINOR.yaml"
    aws s3 cp "$DIST/$REGION.yaml" "s3://$BUCKET/cfn-mackerel-macro/v$MAJOR.yaml"
done << EOS
ap-northeast-1
ap-northeast-2
ap-south-1
ap-southeast-1
ap-southeast-2
ca-central-1
eu-central-1
eu-west-1
eu-west-2
eu-west-3
sa-east-1
us-east-1
us-east-2
us-west-1
us-west-2
EOS

cd "$ROOT"
( git add . && git commit -m "bump up to v$VERSION" && git push ) || true
ghr -u shogo82148 --draft --replace "v$VERSION" "$DIST"

#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

# get the HTTP URL endpoint for the upload function from the CF stack:
uploadHTTPURL=$(aws cloudformation describe-stacks --stack-name imgnstack | jq '.Stacks[].Outputs[] | select(.OutputKey=="UploadImageAPIEndpoint").OutputValue' -r)
# temporary update the JS file with it:
sed -i '.tmp' 's|UPLOAD_HTTP|$uploadHTTPURL|' ui/upload.js
# upload to the S3 bucket:
aws s3 sync ui/ s3://imgn-static --exclude ".DS_Store" --region eu-west-1
# clean up, reinstate original (for next iteration):
mv ui/upload.js.tmp ui/upload.js
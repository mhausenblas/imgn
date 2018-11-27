#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

# get the HTTP API endpoint from the CF stack:
imgnHTTPAPI=$(aws cloudformation describe-stacks --stack-name imgnstack | jq '.Stacks[].Outputs[] | select(.OutputKey=="ImgnAPIEndpoint").OutputValue' -r)
echo The HTTP API endpoint is: $imgnHTTPAPI
# temporary update the JS files with it:
sed -i '.tmp' "s|HTTP_API|$imgnHTTPAPI|" ui/upload.js
sed -i '.tmp' "s|HTTP_API|$imgnHTTPAPI|" ui/gallery.js
# upload to the S3 bucket:
aws s3 sync ui/ s3://imgn-static --exclude ".DS_Store" --region eu-west-1
# clean up, reinstate originals (for next iteration):
mv ui/upload.js.tmp ui/upload.js
mv ui/gallery.js.tmp ui/gallery.js
echo Done. Now go to: http://imgn-static.s3-website-eu-west-1.amazonaws.com/
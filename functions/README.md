# imgn serverless

Clone this repo and work in the `functions/` subdirectory.

Make sure you've got the `aws` CLI and the [SAM CLI](https://github.com/awslabs/aws-sam-cli) installed.

## Preparation

Create S3 buckets, one for the UI (`imgn-static`) and one for storing the uploaded images (`imgn-gallery`):

```bash
$ aws s3api create-bucket \
  --bucket imgn-static \
  --create-bucket-configuration LocationConstraint=eu-west-1 \
  --region eu-west-1
$ aws s3api put-bucket-policy \
  --bucket imgn-static \
  --policy file://s3-ui-bucket-policy.json \
  --region eu-west-1
$ aws s3 website s3://imgn-static/ --index-document index.html

$ aws s3api create-bucket \
  --bucket imgn-gallery \
  --create-bucket-configuration LocationConstraint=eu-west-1 \
  --region eu-west-1
```

## UI

Deploy and update the UI as a static HTML site:

```bash
$ aws s3 sync ui/ s3://imgn-static --exclude ".DS_Store" --region eu-west-1
```

Now the UI is available via http://imgn-static.s3-website-eu-west-1.amazonaws.com/

## Lambda functions and HTTP API

How to build and deploy the Lambda functions and a HTTP API with [SAM](https://github.com/awslabs/serverless-application-model). If you're interested in how to set up and deploy a function and the API Gateway manually, using the `aws` CLI, check out the [notes](low-level/) here.

Note that to get started I did `sam init --name app --runtime go1.x` and in order to do local development I'd do a `sam local start-api`. Also, in order to work, you need to have Docker running, locally.


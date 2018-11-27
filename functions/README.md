# imgn serverless

Clone this repo and work in the `functions/` subdirectory.

Make sure you've got the `aws` CLI and the [SAM CLI](https://github.com/awslabs/aws-sam-cli) installed.

## Preparation

Create S3 buckets, one for the UI (`imgn-static`), one for the Lambda functions (`imgn-app`), and one for storing the uploaded images (`imgn-gallery`):

```bash
# setting up the UI bucket:
$ aws s3api create-bucket \
      --bucket imgn-static \
      --create-bucket-configuration LocationConstraint=eu-west-1 \
      --region eu-west-1
$ aws s3api put-bucket-policy \
      --bucket imgn-static \
      --policy file://s3-ui-bucket-policy.json \
      --region eu-west-1
$ aws s3 website s3://imgn-static/ \
      --index-document index.html

# setting up the app bucket holding the Lambda functions:
$ aws s3api create-bucket \
      --bucket imgn-gallery \
      --create-bucket-configuration LocationConstraint=eu-west-1 \
      --region eu-west-1

# setting up the content bucket for the images to be uploaded:
$ aws s3api create-bucket \
      --bucket imgn-app \
      --create-bucket-configuration LocationConstraint=eu-west-1 \
      --region eu-west-1
```

## UI

Deploy and update the UI as a static HTML site:

```bash
$ aws s3 sync ui/ s3://imgn-static \
      --exclude ".DS_Store" \
      --region eu-west-1
```

Now the UI is available via http://imgn-static.s3-website-eu-west-1.amazonaws.com/

## Lambda functions and HTTP API

How to build and deploy the Lambda functions and a HTTP API with [SAM](https://github.com/awslabs/serverless-application-model). If you're interested in how to set up and deploy a function and the API Gateway manually, using the `aws` CLI, check out the [notes](low-level/) here.

### Local development

To get started I did `sam init --name app --runtime go1.x` initially and developed each of the functions independently. Note that in order to work, you need to have Docker running, locally.

For each code iteration, in `app/` do:

```bash
# 1. run emulation of Lambda and API Gateway locally (via Docker):
$ sam local start-api

# 2. update Go source code (add functionality, fix bugs)

# 3. create a new binary which is automagically synced into SAM runtime:
$ make build
```

Note that if you change anything in the SAM/CF [template file](app/template.yaml) then you need to re-start the local API emulation.

Testing a function (here: [upload image](app/uploadimg)) locally by calling the HTTP API endpoint:

```bash
$ curl -XPOST --header "Content-Type: image/jpeg" \
       --data-binary @test.jpg \
       http://127.0.0.1:3000/upload
```

### Deployment to live environment

```bash
$ sam package \
      --template-file template.yaml --output-template-file imgn-stack.yaml \
      --s3-bucket imgn-app

$ sam deploy \
      --template-file imgn-stack.yaml \
      --stack-name imgnstack \
      --capabilities CAPABILITY_IAM

$ aws cloudformation describe-stacks --stack-name imgnstack
```
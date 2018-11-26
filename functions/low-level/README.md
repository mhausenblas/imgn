# Low-level Lambda setup

Here I show you to build and deploy a Lambda function and set up a trigger via API Gateway manually, using the `aws` CLI.

- [Preparation](#preparation)
- [Managing Lambda functions](#managing-lambda-functions)
- [HTTP API integration](#http-api-integration)

## Preparation

Define a role for Lambda and set permissions:

```bash
$ aws iam create-role --role-name imgn-lambda \
--assume-role-policy-document file://lambda-policy.json

$ aws iam attach-role-policy --role-name imgn-lambda \
--policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
```

Create a HTTP API in the API Gateway:

```bash
$ aws apigateway create-rest-api --name imgn --region eu-west-1
```

From the response, we capture the API ID via `export REST_API_ID=...` and then:

```bash
$ aws apigateway get-resources --rest-api-id $REST_API_ID --region eu-west-1
```

Same here, we capture the root path ID via `export ROOT_PATH_ID=...`.

## Managing Lambda functions

This shows how to build, update, and invoke Lambda functions.

Build a function:

```bash
$ env GOOS=linux GOARCH=amd64 go build -o uploadimg ./uploadfunc
$ zip -j ./uploadimg.zip uploadimg
```

Create the function and make sure that you do `export AWS_ACCOUNT_ID=...` (get your ID from the [console](https://console.aws.amazon.com/billing/home?#/account)):

```bash
$ aws lambda create-function \
  --function-name UploadImg \
  --zip-file fileb://uploadimg.zip \
  --runtime go1.x \
  --role arn:aws:iam::$AWS_ACCOUNT_ID:role/imgn-lambda \
  --handler uploadimg \
  --region eu-west-1
```

Update a function (assuming above build has been done as well as that you've created the function):

```bash
$ aws lambda update-function-code \
  --function-name UploadImg \
  --zip-file fileb://uploadimg.zip \
  --region eu-west-1
```

Directly invoke the Lamdba function like so:

```bash
$ aws lambda invoke --function-name UploadImg --region eu-west-1 uploadimg.json
```

## HTTP API integration

Set up triggers for Lambda functions via a HTTP API in the API Gateway.

Create the trigger for `/upload` like so (note that you must have `REST_API_ID` and `ROOT_PATH_ID` from the preparation above):

```bash
$ aws apigateway create-resource \
  --rest-api-id $REST_API_ID \
  --parent-id $ROOT_PATH_ID \
  --path-part upload \
  --region eu-west-1
```

From above output capture the resource ID via `export RES_ID=...` representing the path `upload/` and then set the allowed HTTP methods like so:

```bash
$ aws apigateway put-method \
  --rest-api-id $REST_API_ID \
  --resource-id $RES_ID \
  --http-method ANY \
  --authorization-type NONE \
  --region eu-west-1
```

Now you can define the integration with the `UploadImg` Lambda function on path `/upload` as follows:

```bash
$ aws apigateway put-integration \
  --rest-api-id $REST_API_ID \
  --resource-id $RES_ID \
  --http-method ANY \
  --type AWS_PROXY \
  --integration-http-method POST \
  --uri arn:aws:apigateway:eu-west-1:lambda:path/2015-03-31/functions/arn:aws:lambda:eu-west-1:$AWS_ACCOUNT_ID:function:UploadImg/invocations \
  --region eu-west-1
```

Make sure to fix permissions so that API Gateway is allowed to execute the Lambda function:

```bash
$ aws lambda add-permission \
  --function-name UploadImg \
  --statement-id mh9-uploadfunc \
  --action lambda:InvokeFunction \
  --principal apigateway.amazonaws.com \
  --source-arn arn:aws:execute-api:eu-west-1:$AWS_ACCOUNT_ID:$REST_API_ID/*/*/* \
  --region eu-west-1
```

For testing, you can trigger the function by invoking the path `upload/` (which `RES_ID` refers to) directly:

```bash
$ aws apigateway test-invoke-method \
  --rest-api-id $REST_API_ID \
  --resource-id $RES_ID \
  --http-method "GET" \
  --region eu-west-1
```

It's time to deploy the HTTP API (we're using the stage name `dev` here but whatever):

```bash
$ aws apigateway create-deployment \
  --rest-api-id $REST_API_ID \
  --stage-name dev
```

Now you can finally call the Lamdba function via the HTTP API using the stage `dev` and path `upload` like so:

```bash
$ curl https://$REST_API_ID.execute-api.eu-west-1.amazonaws.com/dev/upload
```

# imgn serverless

Clone this repo and work in the `functions/` subdirectory.

## Preparation

Create S3 buckets, one for the UI (`imgn-static`) and one for storing the uploaded images (`imgn-gallery`):

```bash
$ aws s3api create-bucket --bucket imgn-static --region eu-west-1 --create-bucket-configuration LocationConstraint=eu-west-1
$ aws s3api put-bucket-policy --bucket imgn-static --policy file://s3-ui-bucket-policy.json --region eu-west-1
$ aws s3 website s3://imgn-static/ --index-document index.html

$ aws s3api create-bucket --bucket imgn-gallery --region eu-west-1 --create-bucket-configuration LocationConstraint=eu-west-1
```

## UI

Deploy and update the UI as a static HTML site:

```bash
$ aws s3 sync ui/ s3://imgn-static --exclude ".DS_Store" --region eu-west-1
```

Now the UI is available via http://imgn-static.s3-website-eu-west-1.amazonaws.com/

## Lambda functions

How to build and deploy the Lambda functions and set up triggers via API Gateway.

### Preparation

Define role and permissions:

```bash
$ aws iam create-role --role-name imgn-lambda \
--assume-role-policy-document file://lambda-policy.json

$ aws iam attach-role-policy --role-name imgn-lambda \
--policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
```

### Build and deploy

Build a function:

```bash
$ env GOOS=linux GOARCH=amd64 go build -o uploadimg ./uploadfunc
$ zip -j ./uploadimg.zip uploadimg
```

Create the function and make sure that you do `export AWS_ACCOUNT_ID=...` (get our ID from the [console](https://console.aws.amazon.com/billing/home?#/account)):

```bash
$ aws lambda create-function \
 --function-name UploadImg \
 --zip-file fileb://uploadimg.zip \
 --runtime go1.x \
 --role arn:aws:iam::$AWS_ACCOUNT_ID:role/imgn-lambda \
 --handler uploadimg \
 --region eu-west-1
 ```

Directly invoke like so:

```bash
$ aws lambda invoke --function-name UploadImg --region eu-west-1 uploadimg.json
```

Set up triggers via a HTTP API in the API Gateway:

```bash
$ aws apigateway create-rest-api --name imgn --region eu-west-1
```

From the response, we capture the API ID via `export REST_API_ID=...` and then:

```bash
$ aws apigateway get-resources --rest-api-id $REST_API_ID --region eu-west-1
```

Same for the root path ID via `export ROOT_PATH_ID=...` and then:

```bash
$ aws apigateway create-resource \
 --rest-api-id $REST_API_ID \
 --parent-id $ROOT_PATH_ID \
 --path-part upload \
 --region eu-west-1
```

Same for the resource ID via `export RES_ID=...` and then:

```bash
$ aws apigateway put-method \
 --rest-api-id $REST_API_ID \
 --resource-id $RES_ID \
 --http-method ANY \
 --authorization-type NONE \
 --region eu-west-1
```

Define integration with the `upload` Lambda function:

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

Fix permissions so that API Gateway is allowed to execute `upload` Lambda function:

```bash
$ aws lambda add-permission \
 --function-name UploadImg \
 --statement-id mh9-uploadfunc \
 --action lambda:InvokeFunction \
 --principal apigateway.amazonaws.com \
 --source-arn arn:aws:execute-api:eu-west-1:$AWS_ACCOUNT_ID:$REST_API_ID/*/*/* \
 --region eu-west-1
```

Update a function:

```bash
$ env GOOS=linux GOARCH=amd64 go build -o uploadimg ./uploadfunc
$ zip -j ./uploadimg.zip uploadimg
$ aws lambda update-function-code \
 --function-name UploadImg \
 --zip-file fileb://uploadimg.zip \
 --region eu-west-1
```

Invoke directly to test:

```bash
$ aws apigateway test-invoke-method \
  --rest-api-id $REST_API_ID \
  --resource-id $RES_ID \
  --http-method "GET" \
  --region eu-west-1
```

Deploy the HTTP API:

```bash
$ aws apigateway create-deployment \
 --rest-api-id $REST_API_ID \
 --stage-name dev
```

Call the HTTP API:

```bash
$ curl https://$REST_API_ID.execute-api.eu-west-1.amazonaws.com/dev/upload
```

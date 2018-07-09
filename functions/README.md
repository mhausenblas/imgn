# imgn serverless

Clone repo and `cd functions`:

Deploy the UI as a static HTML site (replace `imgn-static` with your own bucket):

```bash
$ aws s3 sync ui/ s3://imgn-static --exclude ".DS_Store" --region eu-west-1
$ aws s3api put-bucket-policy --bucket imgn-static --policy file://s3-ui-bucket-policy.json --region eu-west-1
$ aws s3 website s3://imgn-static/ --index-document index.html
```

Now the UI is available via http://imgn-static.s3-website-eu-west-1.amazonaws.com/

Lambda functions:

Build:

```bash
$ env GOOS=linux GOARCH=amd64 go build -o uploadimg ./uploadfunc
$ zip -j ./uploadimg.zip uploadimg
```

Define role and permissions:

```bash
$ aws iam create-role --role-name imgn-lambda-uploadfunc \
--assume-role-policy-document file://uploadfunc/lambda-policy.json

$ aws iam attach-role-policy --role-name imgn-lambda-uploadfunc \
--policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
```

Make sure that you do `export AWS_ACCOUNT_ID=...` (get our ID from the [console](https://console.aws.amazon.com/billing/home?#/account)):

```bash
$ aws lambda create-function \
 --function-name UploadImg \
 --zip-file fileb://uploadimg.zip \
 --runtime go1.x \
 --role arn:aws:iam::$AWS_ACCOUNT_ID:role/imgn-lambda-uploadfunc \
 --handler uploadimg \
 --region eu-west-1
 ```

Invoke:

```bash
$ aws lambda invoke --function-name UploadImg --region eu-west-1 uploadimg.json
```
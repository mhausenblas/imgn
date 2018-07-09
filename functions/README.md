# imgn serverless

Clone repo and `cd functions`:

Deploy the UI as a static HTML site (replace `imgn-static` with your own bucket):

```bash
$ aws s3 sync ui/ s3://imgn-static --exclude ".DS_Store" --region eu-west-1
$ aws s3api put-bucket-policy --bucket imgn-static --policy file://s3-ui-bucket-policy.json --region eu-west-1
$ aws s3 website s3://imgn-static/ --index-document index.html
```

Now the UI is available via http://imgn-static.s3-website-eu-west-1.amazonaws.com/
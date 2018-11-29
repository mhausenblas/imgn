.PHONY: deps clean build deploy

deps:
	go get -u ./...

clean: 
	rm -rf ./uploadimg/uploadimg
	rm -rf ./listimgs/listimgs
	rm -rf ./extractmeta/extractmeta
	rm imgn-stack.yaml
	
build:
	GOOS=linux GOARCH=amd64 go build -o uploadimg/uploadimg ./uploadimg
	GOOS=linux GOARCH=amd64 go build -o listimgs/listimgs ./listimgs
	GOOS=linux GOARCH=amd64 go build -o extractmeta/extractmeta ./extractmeta

deploy:
	sam package --template-file template.yaml --output-template-file imgn-stack.yaml --s3-bucket imgn-app
	sam deploy --template-file imgn-stack.yaml --stack-name imgnstack --capabilities CAPABILITY_IAM
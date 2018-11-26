.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./uploadimg/uploadimg
	
build:
	GOOS=linux GOARCH=amd64 go build -o uploadimg/uploadimg ./uploadimg
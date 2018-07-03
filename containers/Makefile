frontend_version := 0.1
imgproc_version := 0.1

.PHONY: build build-frontend build-imgproc clean

build : build-frontend build-imgproc

build-frontend :
	GOOS=linux GOARCH=amd64 go build -o ./frontend/frontend frontend/main.go
	@docker build -t quay.io/mhausenblas/imgn-frontend:$(frontend_version) frontend/.
	@docker push quay.io/mhausenblas/imgn-frontend:$(frontend_version)

build-imgproc :
	GOOS=linux GOARCH=amd64 go build -o ./imgproc/imgproc imgproc/main.go
	@docker build -t quay.io/mhausenblas/imgn-imgproc:$(imgproc_version) imgproc/.
	@docker push quay.io/mhausenblas/imgn-imgproc:$(imgproc_version)

clean :
	@rm frontend/frontend
	@rm imgproc/imgproc
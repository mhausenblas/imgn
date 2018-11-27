package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

type UploadResponse struct {
	Status string `json:"status"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	imgf := strings.NewReader(request.Body)
	imgfname := "test.png"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("imgn-gallery"),
		Key:    aws.String(imgfname),
		Body:   imgf,
	})
	if err != nil {
		return serverError(err)
	}
	ur := &UploadResponse{
		Status: "Successfully uploaded image",
	}
	js, err := json.Marshal(ur)
	if err != nil {
		return serverError(err)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(js),
	}, nil

}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       fmt.Sprintf("%v", err.Error()),
	}, nil
}

func main() {
	lambda.Start(handler)
}

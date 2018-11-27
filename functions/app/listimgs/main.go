package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	s3client := s3.New(sess)
	resp, err := s3client.ListObjects(&s3.ListObjectsInput{Bucket: aws.String("imgn-gallery")})
	if err != nil {
		return serverError(err)
	}
	type Entry struct {
		Source string `json:"src"`
		Meta   string `json:"meta"`
	}
	flist := []Entry{}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
		flist = append(flist, Entry{
			Source: "/" + *item.Key,
			Meta:   "meta",
		})
	}
	js, err := json.Marshal(flist)
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

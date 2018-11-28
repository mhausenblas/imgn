package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	fmt.Println(err.Error())
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		Body: fmt.Sprintf("%v", err.Error()),
	}, nil
}

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
		flist = append(flist, Entry{
			Source: "http://imgn-gallery.s3-website-eu-west-1.amazonaws.com/" + *item.Key,
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

func main() {
	lambda.Start(handler)
}

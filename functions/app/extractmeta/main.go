package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handler() error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	s3client := s3.New(sess)
	resp, err := s3client.ListObjects(&s3.ListObjectsInput{Bucket: aws.String("imgn-gallery")})
	if err != nil {
		return err
	}
	for _, item := range resp.Contents {
		fmt.Printf("DEBUG:: item: %v\n", item)
	}
	return nil
}

func main() {
	lambda.Start(handler)
}

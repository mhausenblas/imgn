package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

func getMeta(bucket, imgfile string) string {
	// download metadata file:
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	if err != nil {
		log.Printf("ERROR:: %v\n", err)
		return "can't get metadata"
	}
	tmpf, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		log.Printf("ERROR:: %v\n", err)
		return "can't get metadata"
	}
	downloader := s3manager.NewDownloader(s)
	_, err = downloader.Download(tmpf,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(imgfile + ".meta"),
		})
	metadata, err := ioutil.ReadAll(tmpf)
	if err != nil {
		log.Printf("ERROR:: %v\n", err)
		return "no metadata available yet"
	}
	return string(metadata)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	gbucket := "imgn-gallery"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	s3client := s3.New(sess)
	resp, err := s3client.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(gbucket)})
	if err != nil {
		return serverError(err)
	}
	type Entry struct {
		Source string `json:"src"`
		Meta   string `json:"meta"`
	}
	flist := []Entry{}
	for _, item := range resp.Contents {
		if strings.HasSuffix(*item.Key, ".meta") {
			continue
		}
		flist = append(flist, Entry{
			Source: "http://imgn-gallery.s3-website-eu-west-1.amazonaws.com/" + *item.Key,
			Meta:   getMeta(gbucket, *item.Key),
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

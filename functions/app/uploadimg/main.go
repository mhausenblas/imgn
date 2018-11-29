package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

func parseForm(key, mpheader string, body io.Reader) (string, string, io.Reader, error) {
	var buf []byte
	filename := "unknown"
	contentType := "*"
	mediaType, params, err := mime.ParseMediaType(mpheader)
	if err != nil {
		return "", "", nil, err
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(body, params["boundary"])
		fmt.Printf("DEBUG:: boundary: %v\n", params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				fmt.Printf("DEBUG:: parse EOF\n")
				break
			}
			if err != nil {
				return "", "", nil, err
			}
			buf, err = ioutil.ReadAll(p)
			if err != nil {
				return "", "", nil, err
			}
			fmt.Printf("DEBUG:: part: %v\n", p.Header)
			filename = p.FileName()
			contentType = p.Header.Get("Content-Type")
		}
	}
	fmt.Printf("DEBUG:: parsed file name: %v\n", filename)
	fmt.Printf("DEBUG:: parsed content type: %v\n", contentType)
	fmt.Printf("DEBUG:: parsed size in bytes: %v\n", len(buf))
	// fmt.Printf("DEBUG:: parsed content: %v\n", string(buf))
	return filename, contentType, strings.NewReader(string(buf)), nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("DEBUG:: IsBase64Encoded: %v\n", request.IsBase64Encoded)
	gallerybucket := "imgn-gallery"
	// parse the image file name and the data from multipart formdata request:
	imgname, imgct, imgfile, err := parseForm("file", request.Headers["content-type"], strings.NewReader(request.Body))
	if err != nil {
		return serverError(err)
	}
	// upload image file into S3 bucket:
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	if err != nil {
		return serverError(err)
	}
	uploader := s3manager.NewUploader(s, func(u *s3manager.Uploader) {
		u.PartSize = 64 * 1024 * 1024
		u.LeavePartsOnError = true
	})
	uo, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(gallerybucket),
		Key:         aws.String(imgname),
		ContentType: aws.String(imgct),
		Body:        imgfile,
	})
	if err != nil {
		return serverError(err)
	}
	fmt.Printf("DEBUG:: upload result: %v\n", uo)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: fmt.Sprintf("Successfully uploaded %v into S3 bucket %v", imgname, gallerybucket),
	}, nil

}

func main() {
	lambda.Start(handler)
}

package main

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
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

func parseFormData(request events.APIGatewayProxyRequest) (imgname, imgct string, imgfile multipart.File, err error) {
	referrer := request.Headers["Referer"]
	// create a stdlib HTTP request so that we can use the FormFile methods:
	r, err := http.NewRequest(request.HTTPMethod, referrer, strings.NewReader(request.Body))
	if err != nil {
		return "", "", nil, err
	}
	fmt.Printf("DEBUG:: HTTP headers: %v\n", request.Headers)
	// need to set this header in order for FormFile to work
	// note that the raw multipart header is in this format:
	// content-type:multipart/form-data; boundary=----WebKitFormBoundaryHW2069S2hMazyq4B
	mpct := request.Headers["content-type"]
	r.Header.Set("Content-Type", mpct)
	r.Header.Set("Content-Length", strconv.Itoa(len(request.Body)))
	fmt.Printf("DEBUG:: content type: %v\n", mpct)
	// now let's parse the multipart form data (the key/name used in upload.html is 'file'):
	mpfile, mpheader, err := r.FormFile("file")
	if err != nil {
		return "", "", nil, err
	}
	defer mpfile.Close()
	// the name of the image user selected for upload:
	imgname = mpheader.Filename
	fmt.Printf("DEBUG:: filename: %v\n", imgname)
	// the image data itself:
	imgfile = mpfile
	// the content type of the image (yeah, not a good practice, but ...)
	switch {
	case strings.HasSuffix(imgname, "jpg"), strings.HasSuffix(imgname, "jpeg"):
		imgct = "image/jpeg"
	case strings.HasSuffix(imgname, "png"):
		imgct = "image/png"
	case strings.HasSuffix(imgname, "txt"):
		imgct = "text/plain"
	default:
		imgct = ""
	}
	fmt.Printf("DEBUG:: content type: %v\n", imgct)
	return imgname, imgct, imgfile, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("DEBUG:: IsBase64Encoded: %v\n", request.IsBase64Encoded)
	gallerybucket := "imgn-gallery"
	// parse the image file name and the data from multipart formdata request:
	imgname, imgct, imgfile, err := parseFormData(request)
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
		// ContentEncoding: aws.String("base64"),
		Body: imgfile,
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

package main

import (
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func has(list []string, element string) bool {
	for _, e := range list {
		if e == element {
			return true
		}
	}
	return false
}

func upload(bucket, fname, content string) error {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(s, func(u *s3manager.Uploader) {
		u.PartSize = 64 * 1024 * 1024
		u.LeavePartsOnError = true
	})
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fname),
		Body:   strings.NewReader(content),
	})
	return err
}

func extractMetadata(imgfile string) (metadata string, err error) {
	// download image file:
	//TBD
	// decode image and extract dimensions (== metadata):
	// image, _, err := image.DecodeConfig(content)
	// if err != nil {
	// 	return err
	// }
	// metadata := fmt.Sprintf("%dx%d", image.Width, image.Height)
	metadata = "100x100"
	log.Printf("DEBUG:: extracted metadata: %v from %v\n", metadata, imgfile)
	return metadata, nil
}

func handler() error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	s3client := s3.New(sess)
	resp, err := s3client.ListObjects(&s3.ListObjectsInput{Bucket: aws.String("imgn-gallery")})
	if err != nil {
		log.Printf("ERROR:: %v\n", err)
		return err
	}
	bucketcontent := []string{}
	// get all objects in the bucket:
	for _, obj := range resp.Contents {
		// log.Printf("DEBUG:: item: %v\n", obj)
		bucketcontent = append(bucketcontent, *obj.Key)
	}
	// go through the bucket and check if for a given
	// image file the respective metadata file exists:
	for _, file := range bucketcontent {
		if strings.HasSuffix(file, ".meta") {
			continue
		}
		metafile := file + ".meta"
		switch {
		case has(bucketcontent, file) && !has(bucketcontent, metafile):
			metadata, err := extractMetadata(file)
			if err != nil {
				log.Printf("ERROR:: %v\n", err)
			}
			err = upload("imgn-gallery", metafile, metadata)
			if err != nil {
				log.Printf("ERROR:: %v\n", err)
			}
		default:
			// NOP
		}
	}
	return nil
}

func main() {
	lambda.Start(handler)
}

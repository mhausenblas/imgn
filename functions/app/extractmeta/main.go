package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
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

func extractMetadata(bucket, imgfile string) (metadata string, err error) {
	// download image file:
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	if err != nil {
		return "", err
	}
	tmpf, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return "", err
	}
	downloader := s3manager.NewDownloader(s)
	_, err = downloader.Download(tmpf,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(imgfile),
		})
	// decode image and extract dimensions (== metadata):
	image, _, err := image.DecodeConfig(tmpf)
	if err != nil {
		return "", err
	}
	metadata = fmt.Sprintf("%dx%d", image.Width, image.Height)
	log.Printf("DEBUG:: extracted metadata: %v from %v\n", metadata, imgfile)
	return metadata, nil
}

func handler() error {
	gbucket := "imgn-gallery"
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
			metadata, err := extractMetadata(gbucket, file)
			if err != nil {
				log.Printf("ERROR:: %v\n", err)
			}
			err = upload(gbucket, metafile, metadata)
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

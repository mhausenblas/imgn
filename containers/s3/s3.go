package s3

import (
	"fmt"
	"strings"

	minio "github.com/minio/minio-go"
)

// newClient create a client that can talk to Minio
func newClient() (*minio.Client, error) {
	endpoint := "minio:9000"
	accessKeyID := "minio"
	secretAccessKey := "supersecret"
	useSSL := true
	mc, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return &minio.Client{}, err
	}
	return mc, nil
}

// CreateBucket creates a bucket.
func CreateBucket(name string) error {
	mc, err := newClient()
	if err != nil {
		return err
	}
	err = mc.MakeBucket(name, "us-east-1")
	return err
}

// ListBucket lists the objects in the bucket.
// If filter is non-empty, only returns objects
// with the specified suffix.
func ListBucket(name, filter string) ([]string, error) {
	var objects []string
	mc, err := newClient()
	if err != nil {
		return objects, err
	}
	exists, err := mc.BucketExists(name)
	if err != nil {
		return objects, fmt.Errorf(fmt.Sprintf("%s", err))
	}
	if !exists {
		return objects, fmt.Errorf(fmt.Sprintf("Bucket %s does not exist", name))
	}
	done := make(chan struct{})
	defer close(done)
	recursive := false
	for msg := range mc.ListObjects(name, "", recursive, done) {
		object := msg.Key
		switch {
		case filter != "":
			if strings.HasSuffix(object, filter) {
				objects = append(objects, object)
			}
		default:
			objects = append(objects, object)
		}
	}
	return objects, err
}

// UploadToBucket stores object in bucket.
func UploadToBucket(name, object, contentType string) error {
	mc, err := newClient()
	if err != nil {
		return err
	}
	_, err = mc.FPutObject(name, object, object, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}
	return nil
}

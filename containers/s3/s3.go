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
	useSSL := false
	mc, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return &minio.Client{}, err
	}
	return mc, nil
}

// CreateBucket creates a bucket if it doesn't exist yet.
func CreateBucket(name string) error {
	mc, err := newClient()
	if err != nil {
		return err
	}
	exists, err := mc.BucketExists(name)
	if err != nil {
		return err
	}
	if !exists {
		err = mc.MakeBucket(name, "us-east-1")
		if err != nil {
			return err
		}
	}
	return nil
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

// UploadToBucket stores file in path under object in bucket.
func UploadToBucket(name, object, path, contentType string) error {
	mc, err := newClient()
	if err != nil {
		return err
	}
	_, err = mc.FPutObject(name, object, path, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}
	return nil
}

// DownloadFromBucket retrieves object from bucket and stores it in path.
func DownloadFromBucket(name, object, path string) error {
	mc, err := newClient()
	if err != nil {
		return err
	}
	err = mc.FGetObject(name, object, path, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

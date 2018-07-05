package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mhausenblas/imgn/containers/s3"
)

func main() {
	bucketName := "gallery"
	galleryPath := "/app/gallery"
	if _, err := os.Stat(galleryPath); os.IsNotExist(err) {
		derr := os.Mkdir(galleryPath, os.ModePerm)
		if derr != nil {
			log.Printf("Can't create local gallery: %v", err)
		}
	}
	err := s3.CreateBucket(bucketName)
	if err != nil {
		log.Printf("Can't create bucket: %v", err)
	}
	for {
		syncBucket(bucketName)
		files, err := ioutil.ReadDir(galleryPath)
		if err != nil {
			log.Printf("Can't read local gallery: %v", err)
			return
		}
		for _, f := range files {
			if strings.HasSuffix(f.Name(), "jpg") ||
				strings.HasSuffix(f.Name(), "jpeg") ||
				strings.HasSuffix(f.Name(), "png") {
				extractMetadata(filepath.Join(galleryPath, f.Name()))
			}
		}
		time.Sleep(10 * time.Second)
	}
}

// syncBucket syncs the local gallery with the bucket
func syncBucket(name string) {
	objects, err := s3.ListBucket(name, "")
	if err != nil {
		log.Printf("Can't list objects in bucket: %v", err)
	}
	// download each file that we don't have yet locally:
	for _, object := range objects {
		localfile := filepath.Join(galleryPath, object)
		if _, err := os.Stat(localfile); err != nil {
			err = s3.DownloadFromBucket("gallery", object, localfile)
			if err != nil {
				log.Printf("Can't download file %v: %v", object, err)
			}
		}
	}
	log.Printf("Done downloading missing files from bucket %v", name)

	// upload metadata files that are not yet remote:
	files, err := ioutil.ReadDir(galleryPath)
	if err != nil {
		log.Printf("Can't scan local gallery for metadata: %v", err)
		return
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "meta") {
			if !has(objects, f.Name()) {
				err = s3.UploadToBucket("gallery", f.Name(), filepath.Join(galleryPath, f.Name()), "text/plain")
				if err != nil {
					log.Printf("Can't upload metadata file %v: %v", f.Name(), err)
				}
			}
		}
	}
}

func extractMetadata(imgfile string) {
	imgmetafile := imgfile + ".meta"
	if _, err := os.Stat(imgmetafile); err == nil {
		return
	}
	content, err := os.Open(imgfile)
	if err != nil {
		log.Printf("Can't open %s for metadata extraction: %v", imgfile, err)
		return
	}
	image, _, err := image.DecodeConfig(content)
	if err != nil {
		log.Printf("Can't parse metadata from %s: %v", imgfile, err)
		return
	}

	metafile, err := os.Create(imgmetafile)
	if err != nil {
		log.Printf("Can't create metadata file %s: %v", imgmetafile, err)
		return
	}
	defer metafile.Close()
	metadata := fmt.Sprintf("%dx%d", image.Width, image.Height)
	_, err = metafile.WriteString(metadata)
	if err != nil {
		log.Printf("Can't write metadata to %s: %v", imgmetafile, err)
	}
	log.Printf("Added metadata for: %s", imgfile)
}

func has(list []string, element string) bool {
	for _, e := range list {
		if e == element {
			return true
		}
	}
	return false
}

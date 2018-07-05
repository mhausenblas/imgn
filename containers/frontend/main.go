package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mhausenblas/imgn/containers/frontend/handlers"
	"github.com/mhausenblas/imgn/containers/s3"
)

var galleryPath = "/app/ui/gallery"

func main() {
	http.Handle("/", http.FileServer(http.Dir("./ui")))
	http.HandleFunc("/upload", handlers.UploadFile)
	http.HandleFunc("/gallery", handlers.ListFiles)
	log.Println("imgn server running")
	go syncBucket("gallery")
	http.ListenAndServe(":8080", nil)
}

// syncBucket syncs the local gallery with the bucket
func syncBucket(name string) {
	err := s3.CreateBucket(name)
	if err != nil {
		log.Printf("Can't create bucket: %v", err)
	}
	for {
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

		// upload image files that are not yet remote:
		files, err := ioutil.ReadDir(galleryPath)
		if err != nil {
			log.Printf("Can't scan local gallery for image files: %v", err)
			return
		}
		for _, f := range files {
			switch {
			case strings.HasSuffix(f.Name(), "jpg"), strings.HasSuffix(f.Name(), "jpeg"):
				if !has(objects, f.Name()) {
					err = s3.UploadToBucket("gallery", f.Name(), filepath.Join(galleryPath, f.Name()), "image/jpeg")
					if err != nil {
						log.Printf("Can't upload metadata file %v: %v", f.Name(), err)
					}
				}
			case strings.HasSuffix(f.Name(), "png"):
				if !has(objects, f.Name()) {
					err = s3.UploadToBucket("gallery", f.Name(), filepath.Join(galleryPath, f.Name()), "image/png")
					if err != nil {
						log.Printf("Can't upload metadata file %v: %v", f.Name(), err)
					}
				}
			default:
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func has(list []string, element string) bool {
	for _, e := range list {
		if e == element {
			return true
		}
	}
	return false
}

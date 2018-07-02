package imgproc

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
)

// Extract reads image metadata and writes it into a .meta file
func Extract() {
	for {
		files, err := ioutil.ReadDir("./ui/gallery")
		if err != nil {
			log.Printf("Can't list gallery: %v", err)
			return
		}
		for _, f := range files {
			if strings.HasSuffix(f.Name(), "jpg") ||
				strings.HasSuffix(f.Name(), "jpeg") ||
				strings.HasSuffix(f.Name(), "png") {
				extractMetadata(filepath.Join("./ui/gallery/", f.Name()))
			}
		}
		time.Sleep(10 * time.Second)
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

package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// UploadFile uploads a file to the server
func UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	file, handle, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	defer file.Close()

	mimeType := handle.Header.Get("Content-Type")
	switch mimeType {
	case "image/jpeg":
		saveFile(w, file, handle)
	case "image/png":
		saveFile(w, file, handle)
	default:
		jsonResponse(w, http.StatusBadRequest, "Sorry, I only support JPEG and PNG files.")
	}
}

// ListFiles creates a HTML snippet for all uploaded files
func ListFiles(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("./ui/gallery")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		jsonResponse(w, http.StatusInternalServerError, "Meh, data corruption :(")
		return
	}
	type Entry struct {
		Source string `json:"src"`
		Meta   string `json:"meta"`
	}
	flist := []Entry{}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "jpg") ||
			strings.HasSuffix(f.Name(), "jpeg") ||
			strings.HasSuffix(f.Name(), "png") {
			flist = append(flist, Entry{
				Source: filepath.Join("gallery/", f.Name()),
				Meta:   getMeta(filepath.Join("./ui/gallery/", f.Name())),
			})
		}
	}
	_ = json.NewEncoder(w).Encode(flist)
}

func getMeta(imgfile string) string {
	imgmetafile := imgfile + ".meta"
	if _, err := os.Stat(imgmetafile); err != nil {
		return "no metadata available yet"
	}
	metadata, err := ioutil.ReadFile(imgmetafile)
	if err != nil {
		return "no metadata available yet"
	}
	return string(metadata)
}

func saveFile(w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	err = ioutil.WriteFile("./ui/gallery/"+handle.Filename, data, 0666)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	jsonResponse(w, http.StatusCreated, "File uploaded successfully! :)")
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}

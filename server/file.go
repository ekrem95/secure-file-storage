package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ekrem95/secure-file-storage/database"
)

var baseDir = "/tmp/"

func fileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listFiles(w, r)
	case http.MethodPost:
		upload(w, r)
	default:
		errorHandler(w, r, http.StatusNotFound, "")
	}
}

func listFiles(w http.ResponseWriter, r *http.Request) {
	uid, err := strconv.Atoi(r.Header.Get("user_id"))
	if err != nil {
		internalError(w, r)
		return
	}

	var files database.Files

	err = files.Find(uid)
	if err != nil {
		internalError(w, r)
		return
	}

	response(w, r, files)
}
func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	suid := r.Header.Get("user_id")
	uid, err := strconv.Atoi(suid)
	if err != nil {
		internalError(w, r)
		return
	}

	fpath := baseDir + suid + "-" + handler.Filename

	f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	extension := filepath.Ext(handler.Filename)

	f2 := &database.File{Name: handler.Filename, Path: fpath, Ext: extension, Algorithms: "AES", UserID: uid}
	err = f2.Save()
	if err != nil {
		fmt.Println(err)
		internalError(w, r)
		return
	}

	res := struct {
		Message string `json:"message"`
	}{"File uploaded successfully"}

	response(w, r, res)
}

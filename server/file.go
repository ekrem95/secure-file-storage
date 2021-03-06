package server

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ekrem95/secure-file-storage/database"
	"github.com/ekrem95/secure-file-storage/encryption"
)

const perm os.FileMode = 0666

var (
	baseDir = "/tmp/"
	desKey  = []byte("des1key6")
	aesKey  = []byte("aes aes aes aes ")
)

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

	queryFileID := r.URL.Query()["file_id"]

	// list all files of a user
	if len(queryFileID) == 0 {
		var files database.Files

		if err = files.Find(uid); err != nil {
			internalError(w, r)
			return
		}

		response(w, r, files)
		return
	}

	fileID, err := strconv.Atoi(queryFileID[0])
	if err != nil {
		internalError(w, r)
		return
	}

	result, err := database.QueryRow(`SELECT path, algorithms FROM files WHERE id = $1`, fileID)
	if err != nil {
		internalError(w, r)
		return
	}

	var path, algorithms string
	if err = result.Scan(&path, &algorithms); err != nil {
		if err == sql.ErrNoRows {
			errorHandler(w, r, http.StatusBadRequest, "File not found")
			return
		}

		internalError(w, r)
		return
	}

	data, err := readFromFile(path)
	if err != nil {
		internalError(w, r)
		return
	}

	methods := strings.Split(algorithms, ",")

	for _, v := range methods {
		if v == "AES" {
			// data will be updated
			if err = encryption.AesDecrypt(&data, aesKey); err != nil {
				internalError(w, r)
				return
			}
		}

		if v == "DES" {
			// data will be updated
			if err = encryption.DesDecrypt(&data, desKey); err != nil {
				internalError(w, r)
				return
			}
		}
	}

	fmt.Fprint(w, string(data))
}

func readFromFile(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest, err.Error())
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

	f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		internalError(w, r)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	bytes, err := readFromFile(fpath)
	if err != nil {
		internalError(w, r)
		return
	}

	bytes, err = encryption.AesEncrypt(bytes, aesKey)
	if err != nil {
		internalError(w, r)
		return
	}

	bytes, err = encryption.DesEncrypt(bytes, desKey)
	if err != nil {
		internalError(w, r)
		return
	}

	if err = writeToFile(bytes, fpath); err != nil {
		internalError(w, r)
		return
	}

	extension := filepath.Ext(handler.Filename)
	// reverse order encryption methods in order to be able to decrypt easier
	alg := "DES,AES"

	fileInfo := &database.File{Name: handler.Filename, Path: fpath, Ext: extension, Algorithms: alg, UserID: uid}
	if err = fileInfo.Save(); err != nil {
		internalError(w, r)
		return
	}

	res := struct {
		Message string `json:"message"`
	}{"File uploaded successfully"}

	response(w, r, res)
}

func writeToFile(data []byte, file string) error {
	return ioutil.WriteFile(file, data, perm)
}

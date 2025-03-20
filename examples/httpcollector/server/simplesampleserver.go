package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const uploadDirectory = "uploads"

func main() {
	http.HandleFunc("/upload", uploadFileHandler)

	port := "8080"
	if envPort := os.Getenv("SIMPLESAMPLESERVER_PORT"); envPort != "" {
		if _, err := strconv.Atoi(envPort); err != nil {
			fmt.Printf("Error parsing port from environment. Invalid integer value: %q\n", envPort)
		} else {
			port = envPort
		}
	}
	if len(os.Args) > 1 {
		if _, err := strconv.Atoi(os.Args[1]); err != nil {
			fmt.Printf("Error parsing port from argument. Invalid integer value: %q\n", os.Args[1])
		} else {
			port = os.Args[1]
		}
	}

	fmt.Printf("Server started on localhost:%s\n", port)
	_ = http.ListenAndServe(":"+port, nil)
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Limit the size of the memory to 30MB
	err := r.ParseMultipartForm(30 << 20) // 30 MB
	if err != nil {
		fmt.Printf("Error parsing multipart form data: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	multipartFieldnameSha256 := "sha256"
	multipartFieldnameFile := "file"

	// Retrieve the file from the form data
	file, fileHeader, err := r.FormFile(multipartFieldnameFile)
	if err != nil {
		fmt.Printf("Error retrieving the file: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get full filename
	filename := fileHeader.Filename // Usually has only the basename
	if contentDisposition := fileHeader.Header.Get("Content-Disposition"); contentDisposition != "" {
		if _, params, err := mime.ParseMediaType(contentDisposition); err == nil {
			if name, ok := params["filename"]; ok {
				filename = name
			}
		}
	}

	// Retrieve file hash from the form data, or calculate it
	fileHash := r.FormValue(multipartFieldnameSha256)
	if fileHash == "" {
		fmt.Println("Notice sha256 value missing. Calculating hash")
		h := sha256.New()
		if _, err := io.Copy(h, file); err != nil {
			fmt.Printf("Error calculating hash: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileHash = fmt.Sprintf("%x", h.Sum(nil))
	}

	// Create a file to write the uploaded content
	dstDir := filepath.Join(uploadDirectory, fileHash[:2])
	dstPath := filepath.Join(dstDir, fileHash+".data")
	filenamePath := filepath.Join(dstDir, fileHash+".filename.txt")
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			fmt.Printf("Error creating destination directory: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dst, err := os.Create(dstPath)
		if err != nil {
			fmt.Printf("Error creating destination file: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file to the created file on the filesystem
		if _, err := io.Copy(dst, file); err != nil {
			fmt.Printf("Error copying data into file: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	filenameFile, err := os.OpenFile(filenamePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// Allow to fail, it's not critical
	if err != nil {
		fmt.Printf("Error opening filename file for appending: %v\n", err)
	} else {
		if _, err := filenameFile.Write([]byte(filename + "\n")); err != nil {
			fmt.Printf("Error writing filename to file: %v\n", err)
		}
	}

	fmt.Printf("Successfully received sample %s: filename=%q size=%d header=%+v\n", fileHash, filename, fileHeader.Size, fileHeader.Header)
}

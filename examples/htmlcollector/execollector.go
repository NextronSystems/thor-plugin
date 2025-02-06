package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/NextronSystems/jsonlog/thorlog/v3"
	"github.com/NextronSystems/thor-plugin"
)

const sampleServerUrl = "http://localhost:8084/upload"

func Init(config thor.Configuration, logger thor.Logger, actions thor.RegisterActions) {
	actions.AddPostProcessingHook(uploadSample)
	logger.Info("HTMLCollector plugin loaded!")
}

func uploadSample(logger thor.Logger, object thor.MatchedObject) {
	file, isFile := object.Object.(*thorlog.File)

	// Skip if the object is not a file
	if !isFile {
		return
	}

	// Select only EXE files
	if file.MagicHeader != "EXE" {
		return
	}

	// Note: The following implementation is fully backed by memory; in particular, the
	// file data is read into memory, too. If a small memory footprint is required,
	// consider a different approach, e.g., using io.Pipe's writer as multipart writer and
	// passing io.Pipe's reader to the HTTP request/client.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	multipartFieldnameSha256 := "sha256"
	multipartFieldnameFile := "file"

	hashPart, err := writer.CreateFormField(multipartFieldnameSha256)
	if err != nil {
		logger.Error("Failed to create multipart form field for sha256 hash", "error", err)
		return
	}
	_, err = hashPart.Write([]byte(file.Hashes.Sha256))
	if err != nil {
		logger.Error("Failed to write sha256 hash of sample to multipart section", "error", err)
		return
	}

	filePart, err := writer.CreateFormFile(multipartFieldnameFile, file.Path)
	if err != nil {
		logger.Error("Failed to create multipart file section for sample", "error", err)
		return
	}
	_, err = io.Copy(filePart, object.Content)
	if err != nil {
		logger.Error("Failed to write sample data into multipart section", "error", err)
		return
	}

	err = writer.Close()
	if err != nil {
		logger.Error("Failed to finish multipart message", "error", err)
		return
	}

	request, err := http.NewRequest("POST", sampleServerUrl, body)
	if err != nil {
		logger.Error("Failed to create HTTP POST message", "error", err)
		return
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logger.Error("Failed to execute the POST request to the server", "error", err)
		return
	}
	defer response.Body.Close()

	logger.Info("Uploaded sample", "object", object.Object)
}

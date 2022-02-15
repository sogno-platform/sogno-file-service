// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sogno-platform/file-service/api"
)

func addFileRequest(contents string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.csv")
	io.Copy(part, bytes.NewBufferString(contents))
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/files", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req
}

func TestAddFile(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	justNow := time.Now()
	origFileContents := "a|b\n1|2\n"
	req := addFileRequest(origFileContents)
	router.ServeHTTP(w, req)

	// Assert
	var resBody *api.ResponseFile
	json.Unmarshal([]byte(w.Body.String()), &resBody)

	assert.Equal(t, 200, w.Code)
	fileID := resBody.Data.FileID
	assert.NotEqual(t, "", fileID)

	lastModified := resBody.Data.LastModified
	assert.True(t, lastModified.After(justNow))

	url := resBody.Data.URL
	res, _ := http.Get(url)
	actualFileContents, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, origFileContents, string(actualFileContents))

	// Let's add the same file again and make sure the IDs are not the same
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var newResBody *api.ResponseFile
	json.Unmarshal([]byte(w.Body.String()), &newResBody)

	assert.Equal(t, 200, w.Code)
	newFileID := newResBody.Data.FileID
	assert.NotEqual(t, fileID, newFileID)
}

func TestGetFile(t *testing.T) {
	// Add a file
	router := setupRouter()
	w := httptest.NewRecorder()
	origFileContents := "a|b\n1|2\n"
	req := addFileRequest(origFileContents)
	router.ServeHTTP(w, req)

	var addFileRes *api.ResponseFile
	json.Unmarshal([]byte(w.Body.String()), &addFileRes)
	fileID := addFileRes.Data.FileID

	// Get that file
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/files/" + fileID, nil)
	router.ServeHTTP(w, req)

	// Assert
	var getFileRes *api.ResponseFile
	json.Unmarshal([]byte(w.Body.String()), &getFileRes)

	assert.Equal(t, 200, w.Code)

	url := getFileRes.Data.URL
	res, _ := http.Get(url)
	actualFileContents, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, origFileContents, string(actualFileContents))
}

func TestUpdateFile(t *testing.T) {
	// Add a file
	router := setupRouter()
	w := httptest.NewRecorder()
	origFileContents := "a|b\n1|2\n"
	req := addFileRequest(origFileContents)
	router.ServeHTTP(w, req)

	var addFileRes *api.ResponseFile
	json.Unmarshal([]byte(w.Body.String()), &addFileRes)
	fileID := addFileRes.Data.FileID

	// Update that file
	w = httptest.NewRecorder()
	newFileContents := "c|d\n3|4\n"
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.csv")
	io.Copy(part, bytes.NewBufferString(newFileContents))
	writer.Close()

	req, _ = http.NewRequest("PUT", "/api/files/" + fileID, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	var updateFileRes *api.ResponseFile
	json.Unmarshal([]byte(w.Body.String()), &updateFileRes)

	assert.Equal(t, 200, w.Code)
	url := updateFileRes.Data.URL
	res, _ := http.Get(url)
	actualFileContents, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, newFileContents, string(actualFileContents))
}

func TestDeleteFile(t *testing.T) {
	// Add a file
	router := setupRouter()
	w := httptest.NewRecorder()
	origFileContents := "a|b\n1|2\n"
	req := addFileRequest(origFileContents)
	router.ServeHTTP(w, req)

	var addFileRes *api.ResponseFile
	json.Unmarshal([]byte(w.Body.String()), &addFileRes)
	fileID := addFileRes.Data.FileID

	// Delete that file
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/files/" + fileID, nil)
	router.ServeHTTP(w, req)

	// Try to get it again
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/files/" + fileID, nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, 404, w.Code)
}

func TestListFiles(t *testing.T) {
	// Add a file
	router := setupRouter()
	w := httptest.NewRecorder()
	origFileContents := "a|b\n1|2\n"
	req := addFileRequest(origFileContents)
	router.ServeHTTP(w, req)

	// List files
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/files", nil)
	router.ServeHTTP(w, req)

	// Assert
	var resBody *api.ResponseFiles
	json.Unmarshal([]byte(w.Body.String()), &resBody)

	assert.Equal(t, 200, w.Code)
	assert.True(t, len(resBody.Data) >= 1)
}

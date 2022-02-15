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
)

func TestAddFile(t *testing.T) {
	justNow := time.Now()

	router := setupRouter()
	w := httptest.NewRecorder()

	originalBody := "a|b\n1|2\n"
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.csv")
	io.Copy(part, bytes.NewBufferString(originalBody))
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/files", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	var resBody map[string]map[string]string
	json.Unmarshal([]byte(w.Body.String()), &resBody)

	assert.Equal(t, 200, w.Code)
	data, contains := resBody["data"]
	assert.True(t, contains)

	fileID, _ := data["fileID"]
	assert.NotEqual(t, "", fileID)

	lastModified, _ := data["lastModified"]
	lastModifiedDT, _ := time.Parse(time.RFC3339, lastModified)
	assert.True(t, lastModifiedDT.After(justNow))

	url, _ := data["url"]
	res, _ := http.Get(url)
	actualBody, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, originalBody, string(actualBody))

	// Let's add the same file again and make sure the IDs are not the same
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var newResBody map[string]map[string]string
	json.Unmarshal([]byte(w.Body.String()), &newResBody)

	assert.Equal(t, 200, w.Code)
	data, contains = newResBody["data"]
	assert.True(t, contains)

	newFileID, _ := data["fileID"]
	assert.NotEqual(t, fileID, newFileID)
}

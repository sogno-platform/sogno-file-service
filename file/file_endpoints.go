// SPDX-License-Identifier: Apache-2.0

package file

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/sogno-platform/file-service/api"
)

func RegisterFileEndpoints(r *gin.RouterGroup) {
	r.GET("", getFiles)
	r.POST("", addFile)
	r.GET("/:fileID", getFile)
	r.PUT("/:fileID", updateFile)
	r.DELETE("/:fileID", deleteFile)
}

// addFile godoc
// @Summary Add file
// @ID addFile
// @Tags files
// @Produce json
// @Accept multipart/form-data
// @Success 200 {object} api.ResponseFile "File that was added"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param file formData file true "File to be uploaded"
// @Router /files [post]
func addFile(c *gin.Context) {

	fileID := uuid.New().String()
	fileHeader, err := c.FormFile("file")
	if err != nil {
		api.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	contentType := fileHeader.Header.Get("Content-Type")
	contentSize := fileHeader.Size
	content, err := fileHeader.Open()
	if err != nil {
		api.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	content.Close()

	putObject(fileID, content, contentSize, contentType)

	url, err := getObjectUrl(fileID)
	if err != nil {
		api.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	info, err := statObject(fileID)
	if err != nil {
		api.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	c.PureJSON(http.StatusOK, api.ResponseFile{
		Data: api.ResponseFileData{
			FileID:       fileID,
			LastModified: info.LastModified,
			URL:          url.String(),
		},
	})
}

// getFile godoc
// @Summary Get file info
// @ID getFile
// @Tags files
// @Produce json
// @Success 200 {object} api.ResponseFile "File info"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "File not found"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param fileID path string true "ID of file"
// @Router /files/{fileID} [get]
func getFile(c *gin.Context) {

	fileID := c.Param("fileID")
	info, err := statObject(fileID)
	var noSuchKeyError *NoSuchKeyError
	if errors.As(err, &noSuchKeyError) {
		api.ErrorJSON(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		api.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	url, err := getObjectUrl(fileID)
	if err != nil {
		api.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	c.PureJSON(http.StatusOK, api.ResponseFile{
		Data: api.ResponseFileData{
			FileID:       fileID,
			LastModified: info.LastModified,
			URL:          url.String(),
		},
	})
}

// updateFile godoc
// @Summary Update file
// @ID updateFile
// @Tags files
// @Produce json
// @Accept multipart/form-data
// @Success 200 {object} api.ResponseFile "File that was updated"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "File not found"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param fileID path string true "ID of file"
// @Param file formData file true "File to be uploaded"
// @Router /files/{fileID} [put]
func updateFile(c *gin.Context) {

	fileID := c.Param("fileID")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		api.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	// Check if the file exists
	info, err := statObject(fileID)
	var noSuchKeyError *NoSuchKeyError
	if errors.As(err, &noSuchKeyError) {
		api.ErrorJSON(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		api.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}

	contentType := fileHeader.Header.Get("Content-Type")
	contentSize := fileHeader.Size
	content, err := fileHeader.Open()
	if err != nil {
		api.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	content.Close()

	putObject(fileID, content, contentSize, contentType)

	url, err := getObjectUrl(fileID)
	if err != nil {
		api.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	info, err = statObject(fileID)
	if err != nil {
		api.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	c.PureJSON(http.StatusOK, api.ResponseFile{
		Data: api.ResponseFileData{
			FileID:       fileID,
			LastModified: info.LastModified,
			URL:          url.String(),
		},
	})
}

// deleteFile godoc
// @Summary Delete file
// @ID deleteFile
// @Tags files
// @Produce json
// @Success 200 {object} api.ResponseEmpty "Succeeds whether the file exists or not"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param fileID path string true "ID of file"
// @Router /files/{fileID} [delete]
func deleteFile(c *gin.Context) {

	err := deleteObject(c.Param("fileID"))

	if err != nil {
		api.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	c.PureJSON(http.StatusOK, api.ResponseEmpty{})
}

// getFiles godoc
// @Summary Get all files on the server
// @ID getFiles
// @Tags files
// @Produce json
// @Success 200 {object} api.ResponseFiles "Files available"
// @Failure 500 {object} api.ResponseFiles "Internal server error"
// @Router /files [get]
func getFiles(c *gin.Context) {

	var files []api.ResponseFileData

	objInfoChan, err := listObjects()
	if err != nil {
		api.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	for objInfo := range objInfoChan {
		err := objInfo.Err

		if err != nil {
			api.ErrorJSON(c, http.StatusInternalServerError, err)
			return
		} else {
			files = append(files, api.ResponseFileData{
				FileID:       objInfo.Key,
				LastModified: objInfo.LastModified,
			})
		}
	}
	c.PureJSON(http.StatusOK, api.ResponseFiles{Data: files})
}

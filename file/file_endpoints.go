// SPDX-License-Identifier: Apache-2.0

package file

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/sogno-platform/file-service/api"
	"github.com/sogno-platform/file-service/config"
)

func RegisterFileEndpoints(r *gin.RouterGroup) {
	controller, err := NewFileController()
	if err != nil {
		log.Fatalln(err)
	}
	r.GET("", controller.GetFiles)
	r.POST("", controller.AddFile)
	r.GET("/:fileID", controller.GetFile)
	r.PUT("/:fileID", controller.UpdateFile)
	r.DELETE("/:fileID", controller.DeleteFile)
}

type FileController struct {
	Bucket   string
	ObjStore *MinIOClient
}

func NewFileController() (*FileController, error) {
	endpoint := config.GlobalConfig.MinIOEndpoint
	bucket := config.GlobalConfig.MinIOBucket
	client, err := NewMinIOClient(endpoint)
	return &FileController{Bucket: bucket, ObjStore: client}, err
}

// AddFile godoc
// @Summary Add file
// @ID AddFile
// @Tags files
// @Produce json
// @Accept multipart/form-data
// @Success 200 {object} api.ResponseFile "File that was added"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param file formData file true "File to be uploaded"
// @Router /files [post]
func (f *FileController) AddFile(c *gin.Context) {

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

	f.ObjStore.PutObject(f.Bucket, fileID, content, contentSize, contentType)

	url, err := f.ObjStore.GetObjectUrl(f.Bucket, fileID)
	if err != nil {
		api.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	info, err := f.ObjStore.StatObject(f.Bucket, fileID)
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

// GetFile godoc
// @Summary Get file info
// @ID GetFile
// @Tags files
// @Produce json
// @Success 200 {object} api.ResponseFile "File info"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "File not found"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param fileID path string true "ID of file"
// @Router /files/{fileID} [get]
func (f *FileController) GetFile(c *gin.Context) {

	fileID := c.Param("fileID")
	info, err := f.ObjStore.StatObject(f.Bucket, fileID)
	var noSuchKeyError *NoSuchKeyError
	if errors.As(err, &noSuchKeyError) {
		api.ErrorJSON(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		api.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	url, err := f.ObjStore.GetObjectUrl(f.Bucket, fileID)
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

// UpdateFile godoc
// @Summary Update file
// @ID UpdateFile
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
func (f *FileController) UpdateFile(c *gin.Context) {

	fileID := c.Param("fileID")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		api.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	// Check if the file exists
	info, err := f.ObjStore.StatObject(f.Bucket, fileID)
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

	f.ObjStore.PutObject(f.Bucket, fileID, content, contentSize, contentType)

	url, err := f.ObjStore.GetObjectUrl(f.Bucket, fileID)
	if err != nil {
		api.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	info, err = f.ObjStore.StatObject(f.Bucket, fileID)
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

// DeleteFile godoc
// @Summary Delete file
// @ID DeleteFile
// @Tags files
// @Produce json
// @Success 200 {object} api.ResponseEmpty "Succeeds whether the file exists or not"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param fileID path string true "ID of file"
// @Router /files/{fileID} [delete]
func (f *FileController) DeleteFile(c *gin.Context) {

	err := f.ObjStore.DeleteObject(f.Bucket, c.Param("fileID"))

	if err != nil {
		api.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	c.PureJSON(http.StatusOK, api.ResponseEmpty{})
}

// GetFiles godoc
// @Summary Get all files on the server
// @ID GetFiles
// @Tags files
// @Produce json
// @Success 200 {object} api.ResponseFiles "Files available"
// @Failure 500 {object} api.ResponseFiles "Internal server error"
// @Router /files [get]
func (f *FileController) GetFiles(c *gin.Context) {

	var files []api.ResponseFileData

	objInfoChan, err := f.ObjStore.ListObjects(f.Bucket)
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

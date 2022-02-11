// SPDX-License-Identifier: Apache-2.0

package file

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterFileEndpoints(r *gin.RouterGroup) {
	r.GET("", getFiles)
	r.POST("", addFile)
	r.GET("/:fileID", getFile)
	//r.PUT("/:fileID", updateFile)
	//r.DELETE("/:fileID", deleteFile)
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			},
		})
		return
	}

	contentType := fileHeader.Header.Get("Content-Type")
	contentSize := fileHeader.Size
	content, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error(),
			},
		})
		return
	}
	content.Close()

	putObject(fileID, content, contentSize, contentType)
	url, err := getObjectUrl(fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"fileID":       fileID,
			"lastModified": "TODO",
			"url":          url.String(),
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
// @Router /files [post]
func getFile(c *gin.Context) {

	fileID := c.Param("fileID")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    http.StatusBadRequest,
				"message": "fileID required",
			},
		})
		return
	}
	url, err := getObjectUrl(fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"fileID":       fileID,
			"lastModified": "TODO",
			// FIXME: Some encoding is happening which breaks the URL
			"url":          url.String(),
		},
	})
}

// getFiles godoc
// @Summary Get all files on the server
// @ID getFiles
// @Tags files
// @Produce json
// @Success 200 {array} api.ResponseFile "Files available"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Router /files [get]
func getFiles(c *gin.Context) {

	var files []gin.H

	for objInfo := range listObjects() {
		var url *url.URL
		err := objInfo.Err

		if err == nil {
			url, err = getObjectUrl(objInfo.Key)
		}
		if err != nil {
			files = append(files, gin.H{
				"error":        err.Error(),
				"fileID":       "",
				"lastModified": "",
				"url":          "",
			})
		} else {
			files = append(files, gin.H{
				"error":        err,
				"fileID":       objInfo.Key,
				"lastModified": objInfo.LastModified,
				"url":          url.String(),
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"data": files,
	})
}

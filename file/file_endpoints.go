// SPDX-License-Identifier: Apache-2.0

package file

import (
	"github.com/gin-gonic/gin"
)

func RegisterFileEndpoints(r *gin.RouterGroup) {
	r.GET("", getFiles)
	//r.POST("", addFile)
	//r.GET("/:fileID", getFile)
	//r.PUT("/:fileID", updateFile)
	//r.DELETE("/:fileID", deleteFile)
}

// getFiles godoc
// @Summary Get all files on the server
// @ID getFiles
// @Tags files
// @Produce json
// @Success 200 {object} api.ResponseFiles "Files available"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Router /files [get]
func getFiles(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

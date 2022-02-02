// SPDX-License-Identifier: Apache-2.0

package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/sogno-platform/file-service/file"
)

func RegisterEndpoints(r *gin.RouterGroup) {

	file.RegisterFileEndpoints(r.Group("/files"))

}

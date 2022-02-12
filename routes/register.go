// SPDX-License-Identifier: Apache-2.0

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"github.com/sogno-platform/file-service/docs"
	"github.com/sogno-platform/file-service/file"
)

func RegisterEndpoints(r *gin.RouterGroup) {

	docs.SwaggerInfo_swagger.BasePath = "/api"
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	file.RegisterFileEndpoints(r.Group("/files"))

}

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sogno-platform/file-service/routes"
)

func main() {
	r := gin.Default()

	r.GET("", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/api/docs/index.html")
	})
	api := r.Group("/api")
	routes.RegisterEndpoints(api)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// Package docs is the package for the docs rest handler
package docs

import (
	_ "codematic/docs"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// New creates a new instance of the docs rest handler
func New(r *gin.RouterGroup) {
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

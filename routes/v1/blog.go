package v1

import (
	"github.com/NonNtp/gin-gorm-advanced/controllers"
	"github.com/NonNtp/gin-gorm-advanced/middlewares"
	"github.com/gin-gonic/gin"
)

func InitBlogRoutes(rg *gin.RouterGroup) {
	routerGroup := rg.Group("/blogs")

	routerGroup.POST("/", middlewares.AuthJWT(), controllers.CreateBlog)
}

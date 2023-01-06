package controllers

import (
	"net/http"

	"github.com/NonNtp/gin-gorm-advanced/db"
	"github.com/NonNtp/gin-gorm-advanced/dto"
	"github.com/NonNtp/gin-gorm-advanced/models"
	"github.com/gin-gonic/gin"
)

func CreateBlog(ctx *gin.Context) {
	var blog models.Blog

	var form dto.BlogRequest
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blog = models.Blog{
		Topic:  form.Topic,
		Detail: form.Detail,
		UserID: form.UserID,
	}

	if err := db.Conn.Create(&blog).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.CreateBlogResponse{
		ID:     blog.ID,
		Topic:  blog.Topic,
		Detail: blog.Detail,
		UserID: blog.UserID,
	})

}

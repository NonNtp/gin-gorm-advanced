package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/NonNtp/gin-gorm-advanced/db"
	"github.com/NonNtp/gin-gorm-advanced/dto"
	"github.com/NonNtp/gin-gorm-advanced/models"
	"github.com/NonNtp/gin-gorm-advanced/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/matthewhartstonge/argon2"
)

func GetAll(ctx *gin.Context) {
	var users []models.User
	query := db.Conn.Preload("Blogs").Order("created_at").Scopes(utils.Paginate(ctx)).Find(&users)
	//db.Conn.Raw("select * from users order by created_at").Scan(&users)

	if query.RowsAffected < 1 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found "})
		return
	}

	var result []dto.UserAndBlogResponse
	for _, user := range users {
		userResult := dto.UserAndBlogResponse{
			ID:       user.ID,
			Fullname: user.Fullname,
			Email:    user.Email,
		}

		var blogs []dto.BlogResponse
		for _, blog := range user.Blogs {
			blogs = append(blogs, dto.BlogResponse{
				ID:     blog.ID,
				Topic:  blog.Topic,
				Detail: blog.Detail,
			})
		}
		userResult.Blogs = blogs
		result = append(result, userResult)
	}

	ctx.JSON(http.StatusOK, result)
}

func SearchUsers(ctx *gin.Context) {
	fullname := ctx.Query("fullname")
	email := ctx.Query("email")

	var users []models.User
	query := db.Conn.Preload("Blogs")

	if fullname != "" {
		query = query.Where("fullname LIKE ?", "%"+fullname+"%")
	}
	if email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}
	query.Scopes(utils.Paginate(ctx)).Find(&users)

	if query.RowsAffected < 1 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found "})
		return
	}

	var result []dto.UserAndBlogResponse
	for _, user := range users {
		userResult := dto.UserAndBlogResponse{
			ID:       user.ID,
			Fullname: user.Fullname,
			Email:    user.Email,
		}

		var blogs []dto.BlogResponse
		for _, blog := range user.Blogs {
			blogs = append(blogs, dto.BlogResponse{
				ID:     blog.ID,
				Topic:  blog.Topic,
				Detail: blog.Detail,
			})
		}
		userResult.Blogs = blogs
		result = append(result, userResult)
	}

	ctx.JSON(http.StatusOK, result)
}

func GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User

	query := db.Conn.Preload("Blogs").First(&user, id)
	if query.RowsAffected <= 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found "})
		return
	}

	result := dto.UserAndBlogResponse{
		ID:       user.ID,
		Fullname: user.Fullname,
		Email:    user.Email,
	}

	var blogs []dto.BlogResponse
	for _, blog := range user.Blogs {
		blogs = append(blogs, dto.BlogResponse{
			ID:     blog.ID,
			Topic:  blog.Topic,
			Detail: blog.Detail,
		})
	}
	result.Blogs = blogs

	ctx.JSON(http.StatusOK, result)
}

func GetProfile(ctx *gin.Context) {
	user := ctx.MustGet("user")

	ctx.JSON(http.StatusOK, user)
}

func Register(ctx *gin.Context) {
	var form dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Fullname: form.Fullname,
		Email:    form.Email,
		Password: form.Password,
	}

	// check user existing
	userExist := db.Conn.Where("email = ?", form.Email).First(&user)
	if userExist.RowsAffected == 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email existing"})
		return
	}

	if err := db.Conn.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "You are Successfully Registered !"})
}

func Login(ctx *gin.Context) {
	var form dto.LoginRequest

	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:    form.Email,
		Password: form.Password,
	}

	// check user account
	userAccount := db.Conn.Where("email = ?", form.Email).First(&user)
	if userAccount.RowsAffected < 1 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "This email does not exist in the system."})
		return
	}

	// compare password between database and form with argon2
	ok, _ := argon2.VerifyEncoded([]byte(form.Password), []byte(user.Password))
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Password is incorrect"})
		return
	}

	//Initial token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 2).Unix(),
	})

	//Create token
	jwtSecret := os.Getenv("JWT_SECRET")
	accessToken, _ := token.SignedString([]byte(jwtSecret))

	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
	})
}

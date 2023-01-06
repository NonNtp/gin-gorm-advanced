package dto

type BlogRequest struct {
	Topic  string `json:"topic" binding:"required"`
	Detail string `json:"detail" binding:"required"`
	UserID uint   `json:"user_id" binding:"required"`
}

type BlogResponse struct {
	ID     uint   `json:"id"`
	Topic  string `json:"topic"`
	Detail string `json:"detail"`
}

type CreateBlogResponse struct {
	ID     uint   `json:"id"`
	Topic  string `json:"topic"`
	Detail string `json:"detail"`
	UserID uint   `json:"user_id"`
}

package models

import (
	"time"
)

type UserUpdateRequest struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Login    string `json:"login"`
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CreatePostRequest struct {
	PostName        string   `json:"post_name"`
	PostDescription string   `json:"post_description"`
	Tags            []string `json:"tags"`
	SecurityFlag    bool     `json:"security_flag"`
}

type PostID struct {
	PostID int `json:"post_id"`
}

type UpdatePostRequest struct {
	PostID          int      `json:"post_id"`
	PostName        string   `json:"post_name"`
	PostDescription string   `json:"post_description"`
	Tags            []string `json:"tags"`
	SecurityFlag    bool     `json:"security_flag"`
}

type GetPostResponse struct {
	PostID          int       `json:"post_id"`
	PostName        string    `json:"post_name"`
	PostDescription string    `json:"post_description"`
	SecurityFlag    bool      `json:"security_flag"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Tags            []string  `json:"tags"`
	UserID          int       `json:"user_id"`
}

type GetPostListResponse struct {
	Posts []*GetPostResponse `json:"posts"`
}

type PostCommentRequest struct {
	PostID      int    `json:"post_id"`
	Description string `json:"description"`
}

type GetCommentResponse struct {
	CommentID   int    `json:"comment_id"`
	UserID      int    `json:"user_id"`
	PostID      int    `json:"post_id"`
	Description string `json:"description"`
}

type GetCommentListResponse struct {
	Comments []*GetCommentResponse `json:"comments"`
}

type DynamicListResponse struct {
	Dynamic []*DynamicResponse `json:"dynamic"`
}

type DynamicResponse struct {
	Count int       `json:"count"`
	Date  time.Time `json:"date"`
}

type TopParameter struct {
	Parameter string `json:"parameter"`
}

type TopTenResponse struct {
	Top []int `json:"top"`
}

type CountResponse struct {
	Count int32 `json:"count"`
}

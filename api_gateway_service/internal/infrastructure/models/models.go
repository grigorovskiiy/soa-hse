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

type GetPostListRequest struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"page_size"`
}

type CreatePostRequest struct {
	PostName        string   `json:"post_name"`
	PostDescription string   `json:"post_description"`
	Tags            []string `json:"tags"`
	SecurityFlag    bool     `json:"security_flag"`
}

type GetDeletePostRequest struct {
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
	PostId          int       `json:"post_id"`
	PostName        string    `json:"post_name"`
	PostDescription string    `json:"post_description"`
	SecurityFlag    bool      `json:"security_flag"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Tags            []string  `json:"tags"`
	UserId          int       `json:"user_id"`
}

type GetPostListResponse struct {
	Posts []*GetPostResponse `json:"posts"`
}

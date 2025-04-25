package models

import (
	"github.com/uptrace/bun"
	"time"
)

type DbPost struct {
	bun.BaseModel `bun:"table:posts,select:posts"`
	Id            int       `bun:"id,pk,autoincrement" json:"id"`
	Name          string    `bun:"name" json:"name"`
	Description   string    `bun:"description" json:"description"`
	UserId        int       `bun:"user_id" json:"user_id"`
	SecurityFlag  bool      `bun:"security_flag" json:"security_flag"`
	CreatedAt     time.Time `bun:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at" json:"updated_at"`
	Tags          []string  `bun:"tags" json:"tags"`
	Views         int       `bun:"views" json:"views"`
}

type DbComment struct {
	bun.BaseModel `bun:"table:comments,select:comments"`
	Id            int    `bun:"id,pk,autoincrement" json:"id"`
	PostId        int    `bun:"post_id" json:"post_id"`
	UserId        int    `bun:"user_id" json:"user_id"`
	Description   string `bun:"description" json:"description"`
}

type DbLike struct {
	bun.BaseModel `bun:"table:likes,select:likes"`
	Id            int `bun:"id,pk,autoincrement" json:"id"`
	UserId        int `bun:"user_id" json:"user_id"`
	PostId        int `bun:"post_id" json:"post_id"`
}

type DbView struct {
	bun.BaseModel `bun:"table:views,select:views"`
	Id            int `bun:"id,pk,autoincrement" json:"id"`
	PostId        int `bun:"post_id" json:"post_id"`
	UserId        int `bun:"user_id" json:"user_id"`
}

type LikeViewCommentUpdate struct {
	UserId int       `bun:"user_id" json:"user_id"`
	PostId int       `bun:"post_id" json:"post_id"`
	Time   time.Time `bun:"time" json:"time"`
}

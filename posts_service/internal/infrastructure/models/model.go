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
}

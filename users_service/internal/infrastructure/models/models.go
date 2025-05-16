package models

import (
	"github.com/uptrace/bun"
	"time"
)

type DbUser struct {
	bun.BaseModel `bun:"table:users,select:users"`
	Id            int       `bun:"id,pk,autoincrement" json:"id"`
	Name          string    `bun:"name" json:"name"`
	Surname       string    `bun:"surname" json:"surname"`
	Email         string    `bun:"email" json:"email"`
	Password      string    `bun:"password" json:"password"`
	Login         string    `bun:"login" json:"login"`
	CreatedAt     time.Time `bun:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at" json:"updated_at"`
}

type UserUpdateRequest struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
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

type ClientUpdate struct {
	UserId int       `bun:"user_id" json:"user_id"`
	Time   time.Time `bun:"time" json:"time"`
}

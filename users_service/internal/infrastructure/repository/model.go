package repository

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

type UserUpdate struct {
	bun.BaseModel `bun:"table:users,select:users"`
	Name          string `bun:"name" json:"name"`
	Surname       string `bun:"surname" json:"surname"`
	Email         string `bun:"email" json:"email"`
	Password      string `bun:"password" json:"password"`
	Login         string `bun:"login" json:"login"`
}

type UserGetRegisterLogin struct {
	Password string `bun:"password" json:"password"`
	Login    string `bun:"login" json:"login"`
	Email    string `bun:"email" json:"email"`
}

package repository

type UserUpdate struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Login    string `json:"login"`
}

type UserGetRegisterLogin struct {
	Password string `json:"password"`
	Login    string `json:"login"`
	Email    string `json:"email"`
}

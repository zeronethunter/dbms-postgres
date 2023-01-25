package models

type User struct {
	Nickname string `json:"nickname" db:"nickname"`
	Email    string `json:"email" db:"email"`
	FullName string `json:"fullname" db:"fullname"`
	About    string `json:"about" db:"about"`
}

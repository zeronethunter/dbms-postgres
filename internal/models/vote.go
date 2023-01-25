package models

type Vote struct {
	Nickname string `json:"nickname" db:"nickname"`
	Voice    int64  `json:"voice" db:"voice"`
}

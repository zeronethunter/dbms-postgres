package models

import "time"

type Thread struct {
	ID      uint64    `json:"id"`
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
}

type ThreadResponse struct {
	ID      uint64    `json:"id" db:"id"`
	Author  string    `json:"author" db:"author_nickname"`
	Created time.Time `json:"created" db:"created"`
	Forum   string    `json:"forum" db:"forum"`
	Message string    `json:"message" db:"message"`
	Slug    string    `json:"slug" db:"slug"`
	Title   string    `json:"title" db:"title"`
	Votes   int64     `json:"votes" db:"votes"`
}

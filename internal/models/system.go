package models

type System struct {
	UserCount   uint64 `json:"user" db:"user"`
	ForumCount  uint64 `json:"forum" db:"forum"`
	ThreadCount uint64 `json:"thread" db:"thread"`
	PostCount   uint64 `json:"post" db:"post"`
}

package models

type Forum struct {
	Title string `json:"title" db:"title"`
	User  string `json:"user" db:"author_nickname"`
	Slug  string `json:"slug" db:"slug"`
}

type ForumResponse struct {
	Title   string `json:"title" db:"title"`
	User    string `json:"user" db:"author_nickname"`
	Slug    string `json:"slug" db:"slug"`
	Threads uint64 `json:"threads" db:"threads"`
	Posts   uint64 `json:"posts" db:"posts"`
}

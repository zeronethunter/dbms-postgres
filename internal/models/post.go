package models

type Post struct {
	ID       uint64 `json:"id" db:"id"`
	Author   string `json:"author" db:"author_nickname"`
	Created  string `json:"created" db:"created"`
	Forum    string `json:"forum" db:"forum_slug"`
	IsEdited bool   `json:"isEdited" db:"is_edited"`
	Message  string `json:"message" db:"message"`
	Parent   uint64 `json:"parent" db:"parent_id"`
	Thread   uint64 `json:"thread" db:"thread_id"`
}

type FullPost struct {
	Post   *Post           `json:"post"`
	Author *User           `json:"author"`
	Forum  *ForumResponse  `json:"forum"`
	Thread *ThreadResponse `json:"thread"`
}

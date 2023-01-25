package internalErrors

import "errors"

var (
	ErrAlreadyExist                  = errors.New("already exists")
	ErrNoRows                        = errors.New("no rows in result set")
	ErrPostWasCreatedInAnotherThread = errors.New("post was created in another thread")
	ErrPostAuthorNotFound            = errors.New("post author not found")
	ErrNoRowsBySlug                  = errors.New("no rows by slug")
	ErrNoRowsByID                    = errors.New("no rows by id")
	ErrConflictEmail                 = errors.New("conflict email")
	ErrConflictNickname              = errors.New("conflict nickname")
	ErrUserNotFound                  = errors.New("user not found")
	ErrSlugAlreadyExist              = errors.New("slug already exists")
	ErrWrongForumSlug                = errors.New("wrong forum slug")
	ErrNoParentPost                  = errors.New("no parent post")
)

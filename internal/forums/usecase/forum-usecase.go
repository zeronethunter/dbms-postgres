package forumUsecase

import (
	"database/sql"
	"time"

	internalErrors "technopark-dbms-forum/internal"
	forumRepository "technopark-dbms-forum/internal/forums/repository"
	"technopark-dbms-forum/internal/models"
)

type ForumUsecase struct {
	r *forumRepository.Postgres
}

func NewForumUsecase(repo *forumRepository.Postgres) *ForumUsecase {
	return &ForumUsecase{r: repo}
}

func (f *ForumUsecase) Create(forum *models.Forum) (interface{}, error) {
	res, err := f.r.Create(forum)
	if err == internalErrors.ErrAlreadyExist {
		fullRes, err := f.GetFullBySlug(forum.Slug)
		if err != nil {
			return nil, err
		}

		return fullRes, internalErrors.ErrAlreadyExist
	}
	return res, err
}

func (f *ForumUsecase) GetBySlug(slug string) (*models.Forum, error) {
	res, err := f.r.GetBySlug(slug)
	if err == sql.ErrNoRows {
		return nil, internalErrors.ErrNoRows
	}
	return res, err
}

func (f *ForumUsecase) GetFullBySlug(slug string) (*models.ForumResponse, error) {
	res, err := f.r.GetFullBySlug(slug)
	if err == sql.ErrNoRows {
		return nil, internalErrors.ErrNoRows
	}
	return res, err
}

func (f *ForumUsecase) GetThreadsBySlug(slug string, limit int64, since string, desc bool) ([]*models.ThreadResponse, error) {
	_, err := f.r.GetBySlug(slug)
	if err == sql.ErrNoRows {
		return nil, internalErrors.ErrNoRows
	}

	sinceTime := time.Time{}
	if since != "" {
		sinceTime, err = time.Parse(time.RFC3339, since)
		if err != nil {
			return nil, err
		}
	}
	res, err := f.r.GetThreadsBySlug(slug, limit, sinceTime, desc)
	if err == sql.ErrNoRows {
		return nil, internalErrors.ErrNoRows
	}
	return res, err
}

func (f *ForumUsecase) GetUsersBySlug(slug string, limit int64, since string, desc bool) ([]*models.User, error) {
	forum, err := f.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	slug = forum.Slug

	res, err := f.r.GetUsersBySlug(slug, limit, since, desc)
	if err == sql.ErrNoRows {
		return nil, internalErrors.ErrNoRows
	}
	return res, err
}

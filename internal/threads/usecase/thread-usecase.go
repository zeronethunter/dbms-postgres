package threadUsecase

import (
	"strconv"
	"time"

	postRepository "technopark-dbms-forum/internal/posts/repository"

	internalErrors "technopark-dbms-forum/internal"

	"technopark-dbms-forum/internal/models"
	threadRepository "technopark-dbms-forum/internal/threads/repository"
)

type ThreadUsecase struct {
	threadRepo *threadRepository.Postgres
	postsRepo  *postRepository.Postgres
}

func NewThreadUsecase(threadRepo *threadRepository.Postgres, postsRepo *postRepository.Postgres) *ThreadUsecase {
	return &ThreadUsecase{
		threadRepo: threadRepo,
		postsRepo:  postsRepo,
	}
}

func (t *ThreadUsecase) Create(thread *models.Thread) (*models.ThreadResponse, error) {
	return t.threadRepo.Create(thread)
}

func (t *ThreadUsecase) GetBySlugOrID(slugOrID string) (*models.ThreadResponse, error) {
	id, err := strconv.ParseUint(slugOrID, 10, 64)
	if err != nil {
		thread, err := t.threadRepo.GetBySlug(slugOrID)
		if err == internalErrors.ErrNoRows {
			return thread, internalErrors.ErrNoRowsBySlug
		}
		return thread, err
	}

	thread, err := t.threadRepo.GetByID(id)
	if err == internalErrors.ErrNoRows {
		return thread, internalErrors.ErrNoRowsByID
	}
	return thread, err
}

func (t *ThreadUsecase) Update(slugOrID, message, title string) (*models.ThreadResponse, error) {
	id, err := strconv.ParseUint(slugOrID, 10, 64)
	if err != nil {
		return t.threadRepo.UpdateBySlug(&models.Thread{
			Slug:    slugOrID,
			Message: message,
			Title:   title,
		})
	}

	return t.threadRepo.UpdateByID(&models.Thread{
		ID:      id,
		Message: message,
		Title:   title,
	})
}

func (t *ThreadUsecase) Vote(slugOrID string, vote *models.Vote) (*models.ThreadResponse, error) {
	id, err := strconv.ParseUint(slugOrID, 10, 64)
	if err != nil {
		return t.threadRepo.VoteBySlug(slugOrID, vote)
	}

	return t.threadRepo.VoteByID(id, vote)
}

func (t *ThreadUsecase) CreatePosts(slugOrID string, posts []*models.Post) ([]*models.Post, error) {
	id, err := strconv.ParseUint(slugOrID, 10, 64)
	var thread *models.ThreadResponse
	if err != nil {
		thread, err = t.threadRepo.GetBySlug(slugOrID)
		if err == internalErrors.ErrNoRows {
			return nil, internalErrors.ErrNoRowsBySlug
		} else if err != nil {
			return nil, err
		}
	} else {
		thread, err = t.threadRepo.GetByID(id)
		if err == internalErrors.ErrNoRows {
			return nil, internalErrors.ErrNoRowsByID
		} else if err != nil {
			return nil, err
		}
	}

	timeNow := time.Now().Format(time.RFC3339)
	for index := range posts {
		if posts[index].Parent != 0 {
			parent, err := t.postsRepo.GetByID(posts[index].Parent)
			if err == internalErrors.ErrNoRows {
				return nil, internalErrors.ErrPostWasCreatedInAnotherThread
			}

			if parent.Forum != thread.Forum {
				return nil, internalErrors.ErrPostWasCreatedInAnotherThread
			}
		}

		posts[index].Thread = thread.ID
		posts[index].Forum = thread.Forum
		posts[index].Created = timeNow
	}

	return t.threadRepo.CreatePosts(posts)
}

func (t *ThreadUsecase) GetPosts(slugOrID string, limit, since uint64, sort string, desc bool) ([]*models.Post, error) {
	id, err := strconv.ParseUint(slugOrID, 10, 64)
	var thread *models.ThreadResponse
	if err != nil {
		thread, err = t.threadRepo.GetBySlug(slugOrID)
		if err == internalErrors.ErrNoRows {
			return nil, internalErrors.ErrNoRowsBySlug
		} else if err != nil {
			return nil, err
		}
	} else {
		thread, err = t.threadRepo.GetByID(id)
		if err == internalErrors.ErrNoRows {
			return nil, internalErrors.ErrNoRowsByID
		} else if err != nil {
			return nil, err
		}
	}

	return t.threadRepo.GetPostsByID(thread.ID, limit, since, sort, desc)
}

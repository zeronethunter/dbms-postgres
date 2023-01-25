package postUsecase

import (
	"technopark-dbms-forum/internal/models"
	postRepository "technopark-dbms-forum/internal/posts/repository"
)

type PostUsecase struct {
	r *postRepository.Postgres
}

func NewPostUsecase(repo *postRepository.Postgres) *PostUsecase {
	return &PostUsecase{r: repo}
}

func (p *PostUsecase) GetByID(id uint64) (*models.Post, error) {
	return p.r.GetByID(id)
}

func (p *PostUsecase) Update(post *models.Post) (*models.Post, error) {
	_, err := p.r.Update(post)
	if err != nil {
		return nil, err
	}
	return p.r.GetByID(post.ID)
}

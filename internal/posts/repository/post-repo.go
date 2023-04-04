package postRepository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	internalErrors "technopark-dbms-forum/internal"
	"technopark-dbms-forum/internal/models"
)

type Postgres struct {
	sqlx *sqlx.DB
}

func NewPostgres(url string) (*Postgres, error) {
	newSQLX, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}

	if err = newSQLX.Ping(); err != nil {
		return nil, err
	}

	return &Postgres{sqlx: newSQLX}, nil
}

func (p *Postgres) Close() error {
	return p.sqlx.Close()
}

func (p *Postgres) GetByID(id uint64) (*models.Post, error) {
	post := models.Post{}
	err := p.sqlx.Get(
		&post,
		`
			SELECT id, author_nickname, forum_slug, message, thread_id, parent_id, is_edited, created
			FROM posts
			WHERE id = $1
		`,
		id,
	)
	if err == sql.ErrNoRows {
		return nil, internalErrors.ErrNoRows
	} else if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == "23503" {
			return nil, internalErrors.ErrUserNotFound
		}
		return nil, err
	}

	return &post, nil
}

func (p *Postgres) Update(newPost *models.Post) (*models.Post, error) {
	_, err := p.sqlx.Exec(
		`
			UPDATE posts 
			SET message = COALESCE(NULLIF($1, ''), message), 
			    is_edited = CASE WHEN message = COALESCE(NULLIF($1, ''), message) THEN is_edited ELSE true END
			WHERE id = $2
		`,
		newPost.Message,
		newPost.ID,
	)
	if err != nil {
		return nil, err
	}

	return newPost, nil
}

package forumRepository

import (
	"database/sql"
	"time"

	"github.com/lib/pq"

	internalErrors "technopark-dbms-forum/internal"

	"github.com/jmoiron/sqlx"
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

func (p *Postgres) Create(f *models.Forum) (*models.Forum, error) {
	forum := models.Forum{}
	err := p.sqlx.Get(
		&forum,
		"SELECT title, author_nickname, slug FROM forums WHERE slug = $1",
		f.Slug,
	)
	if err != sql.ErrNoRows {
		return &forum, internalErrors.ErrAlreadyExist
	} else if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if _, err = p.sqlx.Exec(
		"INSERT INTO forums (title, author_nickname, slug) VALUES ($1, $2, $3)",
		f.Title,
		f.User,
		f.Slug,
	); err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == "23503" {
			return nil, internalErrors.ErrUserNotFound
		}

		return nil, err
	}

	return f, nil
}

func (p *Postgres) GetBySlug(slug string) (*models.Forum, error) {
	forum := models.Forum{}
	err := p.sqlx.Get(
		&forum,
		"SELECT title, author_nickname, slug FROM forums WHERE slug = $1",
		slug,
	)
	if err != nil {
		return nil, err
	}

	return &forum, nil
}

func (p *Postgres) GetFullBySlug(slug string) (*models.ForumResponse, error) {
	forum := models.ForumResponse{}
	err := p.sqlx.Get(
		&forum,
		`
			SELECT f.title, f.author_nickname, f.slug, f.posts, f.threads 
			FROM forums as f 
			WHERE f.slug = $1
		`,
		slug,
	)
	if err != nil {
		return nil, err
	}

	return &forum, nil
}

func (p *Postgres) GetUsersBySlug(slug string, limit int64, since string, desc bool) ([]*models.User, error) {
	users := make([]*models.User, 0)

	if desc {
		if since != "" {
			err := p.sqlx.Select(
				&users,
				`
				SELECT u.nickname, u.fullname, u.about, u.email
				FROM users as u
				JOIN user_forum as uf ON uf.nickname = u.nickname
				WHERE uf.forum_slug = $1 and u.nickname < $2
				ORDER BY u.nickname desc 
				LIMIT $3
			`,
				slug,
				since,
				limit,
			)
			if err != nil {
				return nil, err
			}
		} else {
			err := p.sqlx.Select(
				&users,
				`
				SELECT u.nickname, u.fullname, u.about, u.email
				FROM users as u
				JOIN user_forum as uf ON uf.nickname = u.nickname
				WHERE uf.forum_slug = $1
				ORDER BY u.nickname desc 
				LIMIT $2
			`,
				slug,
				limit,
			)
			if err != nil {
				return nil, err
			}
		}
	} else {
		if since != "" {
			err := p.sqlx.Select(
				&users,
				`
				SELECT u.nickname, u.fullname, u.about, u.email
				FROM users as u
				JOIN user_forum as uf ON uf.nickname = u.nickname
				WHERE uf.forum_slug = $1 and u.nickname > $2
				ORDER BY u.nickname 
				LIMIT $3
			`,
				slug,
				since,
				limit,
			)
			if err != nil {
				return nil, err
			}
		} else {
			err := p.sqlx.Select(
				&users,
				`
				SELECT u.nickname, u.fullname, u.about, u.email
				FROM users as u
				JOIN user_forum as uf ON uf.nickname = u.nickname
				WHERE uf.forum_slug = $1
				ORDER BY u.nickname 
				LIMIT $2
			`,
				slug,
				limit,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	return users, nil
}

func (p *Postgres) GetThreadsBySlug(slug string, limit int64, since time.Time, desc bool) ([]*models.ThreadResponse, error) {
	threads := make([]*models.ThreadResponse, 0)

	if since != (time.Time{}) {
		if desc {
			err := p.sqlx.Select(
				&threads,
				`
					SELECT t.id, t.author_nickname, t.forum, t.message, t.slug, t.title, t.created, t.votes
					FROM threads t
					WHERE t.forum = $1 AND t.created <= $2
					ORDER BY t.created desc
					LIMIT $3
				`,
				slug,
				since,
				limit,
			)
			if err != nil {
				return nil, err
			}
		} else {
			err := p.sqlx.Select(
				&threads,
				`
					SELECT t.id, t.author_nickname, t.forum, t.message, t.slug, t.title, t.created, t.votes
					FROM threads t
					WHERE t.forum = $1 AND t.created >= $2
					ORDER BY t.created
					LIMIT $3
				`,
				slug,
				since,
				limit,
			)
			if err != nil {
				return nil, err
			}
		}
	} else {
		if desc {
			err := p.sqlx.Select(
				&threads,
				`
					SELECT t.id, t.author_nickname, t.forum, t.message, t.slug, t.title, t.created, t.votes
					FROM threads t
					WHERE t.forum = $1
					ORDER BY t.created desc
					LIMIT $2
				`,
				slug,
				limit,
			)
			if err != nil {
				return nil, err
			}
		} else {
			err := p.sqlx.Select(
				&threads,
				`
					SELECT t.id, t.author_nickname, t.forum, t.message, t.slug, t.title, t.created, t.votes
					FROM threads t
					WHERE t.forum = $1
					ORDER BY t.created
					LIMIT $2
				`,
				slug,
				limit,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	return threads, nil
}

package threadRepository

import (
	"database/sql"
	"fmt"

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

	newSQLX.SetMaxOpenConns(100)

	if err = newSQLX.Ping(); err != nil {
		return nil, err
	}

	return &Postgres{sqlx: newSQLX}, nil
}

func (p *Postgres) Close() error {
	return p.sqlx.Close()
}

func (p *Postgres) Create(t *models.Thread) (*models.ThreadResponse, error) {
	if t.Slug != "" {
		if thread, err := p.GetBySlug(t.Slug); err != internalErrors.ErrNoRows {
			return thread, internalErrors.ErrSlugAlreadyExist
		}
	}

	thread := models.ThreadResponse{
		ID:      t.ID,
		Author:  t.Author,
		Created: t.Created,
		Forum:   t.Forum,
		Message: t.Message,
		Slug:    t.Slug,
		Title:   t.Title,
		Votes:   0,
	}

	if err := p.sqlx.QueryRowx(
		`
			INSERT INTO threads (author_nickname, created, forum, message, slug, title) 
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id;
		`,
		t.Author,
		t.Created,
		t.Forum,
		t.Message,
		t.Slug,
		t.Title,
	).Scan(&thread.ID); err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code == "23503" {
				return nil, internalErrors.ErrUserNotFound
			}
		}
		return nil, err
	}

	return &thread, nil
}

func (p *Postgres) GetBySlug(slug string) (*models.ThreadResponse, error) {
	thread := models.ThreadResponse{}
	err := p.sqlx.Get(
		&thread,
		`
			SELECT id, author_nickname, created, forum, message, slug, title, votes
			FROM threads
			WHERE slug = $1
		`,
		slug,
	)
	if err == sql.ErrNoRows {
		return nil, internalErrors.ErrNoRows
	} else if err != nil {
		return nil, err
	}

	return &thread, nil
}

func (p *Postgres) GetByID(id uint64) (*models.ThreadResponse, error) {
	thread := models.ThreadResponse{}
	err := p.sqlx.Get(
		&thread,
		`
			SELECT id, author_nickname, created, forum, message, slug, title, votes
			FROM threads
			WHERE id = $1
		`,
		id,
	)
	if err == sql.ErrNoRows {
		return nil, internalErrors.ErrNoRows
	} else if err != nil {
		return nil, err
	}

	return &thread, nil
}

func (p *Postgres) UpdateByID(t *models.Thread) (*models.ThreadResponse, error) {
	thread := models.ThreadResponse{}
	err := p.sqlx.Get(
		&thread,
		`
			SELECT id, author_nickname, created, forum, message, slug, title, votes
			FROM threads
			WHERE id = $1
		`,
		t.ID,
	)
	if err == sql.ErrNoRows {
		return nil, internalErrors.ErrNoRowsByID
	} else if err != nil {
		return nil, err
	}

	if t.Message != "" {
		thread.Message = t.Message
	}
	if t.Title != "" {
		thread.Title = t.Title
	}

	_, err = p.sqlx.Exec(
		`
			UPDATE threads
			SET message = $1, title = $2
			WHERE id = $3
		`,
		thread.Message,
		thread.Title,
		thread.ID,
	)
	if err != nil {
		return nil, err
	}

	return &thread, nil
}

func (p *Postgres) UpdateBySlug(t *models.Thread) (*models.ThreadResponse, error) {
	thread := models.ThreadResponse{
		ID:      t.ID,
		Author:  t.Author,
		Created: t.Created,
		Forum:   t.Forum,
		Message: t.Message,
		Slug:    t.Slug,
		Title:   t.Title,
		Votes:   0,
	}
	err := p.sqlx.Get(
		&thread,
		`
			SELECT id, author_nickname, created, forum, message, slug, title, votes
			FROM threads
			WHERE slug = $1
		`,
		t.Slug,
	)
	if err == sql.ErrNoRows {
		return &thread, internalErrors.ErrNoRowsBySlug
	} else if err != nil {
		return &thread, err
	}

	if t.Message != "" {
		thread.Message = t.Message
	}
	if t.Title != "" {
		thread.Title = t.Title
	}

	_, err = p.sqlx.Exec(
		`
			UPDATE threads
			SET message = $1, title = $2
			WHERE id = $3
		`,
		thread.Message,
		thread.Title,
		thread.ID,
	)
	if err != nil {
		return nil, err
	}

	return &thread, nil
}

func (p *Postgres) VoteBySlug(slug string, v *models.Vote) (*models.ThreadResponse, error) {
	thread := models.ThreadResponse{
		Slug: slug,
	}

	err := p.sqlx.Get(
		&thread,
		`
			SELECT id, author_nickname, created, forum, message, slug, title, votes
			FROM threads
			WHERE slug = $1
		`,
		slug,
	)
	if err == sql.ErrNoRows {
		return &thread, internalErrors.ErrNoRows
	} else if err != nil {
		return &thread, err
	}

	_, err = p.sqlx.Exec(
		`
			INSERT INTO votes (nickname, thread_id, voice)
			VALUES ($1, $2, $3)
			ON CONFLICT (nickname, thread_id) DO UPDATE
			SET voice = $3
		`,
		v.Nickname,
		thread.ID,
		v.Voice,
	)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == "23503" {
			return nil, internalErrors.ErrUserNotFound
		}
		return nil, err
	}

	err = p.sqlx.Get(
		&thread,
		`
			SELECT id, author_nickname, created, forum, message, slug, title, votes
			FROM threads
			WHERE slug = $1
		`,
		slug,
	)
	if err != nil {
		return &thread, err
	}

	return &thread, nil
}

func (p *Postgres) VoteByID(id uint64, v *models.Vote) (*models.ThreadResponse, error) {
	thread := models.ThreadResponse{
		ID: id,
	}

	err := p.sqlx.Get(
		&thread,
		`
			SELECT id, author_nickname, created, forum, message, slug, title, votes
			FROM threads
			WHERE id = $1
		`,
		id,
	)
	if err == sql.ErrNoRows {
		return &thread, internalErrors.ErrNoRows
	} else if err != nil {
		return &thread, err
	}

	_, err = p.sqlx.Exec(
		`
			INSERT INTO votes (nickname, thread_id, voice)
			VALUES ($1, $2, $3)
			ON CONFLICT (nickname, thread_id) DO UPDATE
			SET voice = $3
		`,
		v.Nickname,
		thread.ID,
		v.Voice,
	)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == "23503" {
			return nil, internalErrors.ErrUserNotFound
		}
		return nil, err
	}

	err = p.sqlx.Get(
		&thread,
		`
			SELECT id, author_nickname, created, forum, message, slug, title, votes
			FROM threads
			WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return &thread, err
	}

	return &thread, nil
}

func (p *Postgres) getPostsByIDFlat(id uint64, limit uint64, since uint64, desc bool) ([]*models.Post, error) {
	posts := make([]*models.Post, 0)

	query := `
		SELECT p.id, p.author_nickname, p.created, p.forum_slug, p.is_edited, p.message, p.parent_id, p.thread_id
		FROM posts p
		WHERE p.thread_id = $1
	`

	if since != 0 {
		if desc {
			query += fmt.Sprintf(" AND p.id < %d", since)
		} else {
			query += fmt.Sprintf(" AND p.id > %d", since)
		}
	}

	if desc {
		query += " ORDER BY created DESC, p.id DESC"
	} else {
		query += " ORDER BY created, p.id"
	}

	if limit != 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	} else {
		query += " LIMIT 100"
	}

	err := p.sqlx.Select(
		&posts,
		query,
		id,
	)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *Postgres) getPostsByIDTree(id uint64, limit uint64, since uint64, desc bool) ([]*models.Post, error) {
	posts := make([]*models.Post, 0)

	query := `
		SELECT p.id, p.author_nickname, p.created, p.forum_slug, p.is_edited, p.message, p.parent_id, p.thread_id
		FROM posts p
		WHERE p.thread_id = $1
	`

	if since != 0 {
		if desc {
			query += fmt.Sprintf(" AND p.path < (SELECT path FROM posts WHERE id = %d)", since)
		} else {
			query += fmt.Sprintf(" AND p.path > (SELECT path FROM posts WHERE id = %d)", since)
		}
	}

	if desc {
		query += " ORDER BY p.path DESC"
	} else {
		query += " ORDER BY p.path"
	}

	if limit != 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	} else {
		query += " LIMIT 100"
	}

	err := p.sqlx.Select(
		&posts,
		query,
		id,
	)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *Postgres) getPostsByIDParentTree(id uint64, limit uint64, since uint64, desc bool) ([]*models.Post, error) {
	var err error
	posts := make([]*models.Post, 0)

	if limit == 0 {
		limit = 100
	}

	if since == 0 {
		if desc {
			err = p.sqlx.Select(
				&posts,
				`
					SELECT id, parent_id, author_nickname, message, is_edited, forum_slug, created, thread_id
					FROM posts
					WHERE path[1] IN (
					      	SELECT id 
					      	FROM posts 
					      	WHERE thread_id = $1 AND parent_id = 0 
					      	ORDER BY id DESC 
					      	LIMIT $2
					      )
					ORDER BY path[1] DESC, path;
				`,
				id,
				limit,
			)
		} else {
			err = p.sqlx.Select(
				&posts,
				`
					SELECT id, parent_id, author_nickname, message, is_edited, forum_slug, created, thread_id 
					FROM posts
					WHERE path[1] IN (
					      	SELECT id 
					      	FROM posts 
					      	WHERE thread_id = $1 AND parent_id = 0 
					      	ORDER BY id 
					      	LIMIT $2
					      )
					ORDER BY path;
				`,
				id,
				limit,
			)
		}
	} else {
		if desc {
			err = p.sqlx.Select(
				&posts,
				`
					SELECT id, parent_id, author_nickname, message, is_edited, forum_slug, created, thread_id 
					FROM posts
					WHERE path[1] IN (
					      	SELECT id 
					      	FROM posts 
					      	WHERE thread_id = $1 AND parent_id = 0 AND path[1] < (SELECT path[1] FROM posts WHERE id = $2)
					      	ORDER BY id DESC 
					      	LIMIT $3
					      )
					ORDER BY path[1] DESC, path;
				`,
				id,
				since,
				limit,
			)
		} else {
			err = p.sqlx.Select(
				&posts,
				`
					SELECT id, parent_id, author_nickname, message, is_edited, forum_slug, created, thread_id 
					FROM posts p 
					WHERE p.path[1] IN (
					      	SELECT id 
					      	FROM posts 
					      	WHERE thread_id = $1 AND parent_id = 0 AND path[1] > (SELECT path[1] FROM posts WHERE id = $2)
					      	ORDER BY id 
					      	LIMIT $3
					      )
					ORDER BY path;
				`,
				id,
				since,
				limit,
			)
		}
	}
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *Postgres) GetPostsByID(id uint64, limit uint64, since uint64, sort string, desc bool) ([]*models.Post, error) {
	posts := make([]*models.Post, 0)
	var err error
	switch sort {
	case "flat":
		posts, err = p.getPostsByIDFlat(id, limit, since, desc)
		if err == sql.ErrNoRows {
			return posts, internalErrors.ErrNoRows
		} else if err != nil {
			return nil, err
		}
	case "tree":
		posts, err = p.getPostsByIDTree(id, limit, since, desc)
		if err == sql.ErrNoRows {
			return posts, internalErrors.ErrNoRows
		} else if err != nil {
			return nil, err
		}
	case "parent_tree":
		posts, err = p.getPostsByIDParentTree(id, limit, since, desc)
		if err == sql.ErrNoRows {
			return posts, internalErrors.ErrNoRows
		} else if err != nil {
			return nil, err
		}
	default:
		posts, err = p.getPostsByIDFlat(id, limit, since, desc)
		if err == sql.ErrNoRows {
			return posts, internalErrors.ErrNoRows
		} else if err != nil {
			return nil, err
		}
	}

	return posts, nil
}

func (p *Postgres) CreatePosts(posts []*models.Post) ([]*models.Post, error) {
	tx, err := p.sqlx.Beginx()
	if err != nil {
		return nil, err
	}

	for index, post := range posts {
		err = tx.QueryRowx(`
			INSERT INTO posts (author_nickname, created, forum_slug, message, parent_id, thread_id)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`,
			post.Author,
			post.Created,
			post.Forum,
			post.Message,
			post.Parent,
			post.Thread,
		).Scan(&posts[index].ID)
		if err != nil {
			tx.Rollback()

			pgErr, ok := err.(*pq.Error)
			if ok {
				switch pgErr.Code {
				case "23503":
					return nil, internalErrors.ErrPostAuthorNotFound
				case "23505":
					return nil, internalErrors.ErrPostWasCreatedInAnotherThread
				}
			}
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

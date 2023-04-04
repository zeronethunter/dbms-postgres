package systemRepository

import (
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

func (p *Postgres) ClearAll() error {
	_, err := p.sqlx.Exec(
		`
			TRUNCATE forums CASCADE;
			TRUNCATE votes CASCADE;
			TRUNCATE posts CASCADE;
			TRUNCATE threads CASCADE;
			TRUNCATE users CASCADE;
			TRUNCATE user_forum CASCADE;
		`,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetInfo() (*models.System, error) {
	var sys models.System

	err := p.sqlx.Get(
		&sys,
		`
			SELECT COUNT(*) as forum
			FROM forums
		`,
	)
	if err != nil {
		return nil, err
	}

	err = p.sqlx.Get(
		&sys,
		`
			SELECT COUNT(*) as "user"
			FROM users
		`,
	)
	if err != nil {
		return nil, err
	}

	if err = p.sqlx.Get(
		&sys,
		`
			SELECT COUNT(*) as thread
			FROM threads
		`,
	); err != nil {
		return nil, err
	}

	if err = p.sqlx.Get(
		&sys,
		`
			SELECT COUNT(*) as post
			FROM posts
		`,
	); err != nil {
		return nil, err
	}

	return &sys, nil
}

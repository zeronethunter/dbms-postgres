package userRepository

import (
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

func (p *Postgres) Create(u *models.User) ([]*models.User, error) {
	var user []*models.User
	err := p.sqlx.Select(
		&user,
		"SELECT nickname, email, fullname, about FROM users WHERE nickname = $1 OR email = $2",
		u.Nickname,
		u.Email,
	)
	if len(user) == 0 {
		if _, err = p.sqlx.Exec(
			"INSERT INTO users (nickname, email, fullname, about) VALUES ($1, $2, $3, $4)",
			u.Nickname,
			u.Email,
			u.FullName,
			u.About,
		); err != nil {
			return nil, err
		}

		user = append(user, u)
	} else {
		return user, internalErrors.ErrAlreadyExist
	}

	return user, nil
}

func (p *Postgres) GetByNickname(nickname string) (*models.User, error) {
	var user models.User
	if err := p.sqlx.Get(
		&user,
		"SELECT nickname, email, fullname, about FROM users WHERE nickname = $1",
		nickname,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *Postgres) Update(u *models.User) error {
	_, err := p.sqlx.Exec(
		`
			UPDATE users SET email = $1, fullname = $2, about = $3 WHERE nickname = $4
		`,
		u.Email,
		u.FullName,
		u.About,
		u.Nickname,
	)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code == "23505" {
				if pgErr.Constraint == "users_email_key" {
					return internalErrors.ErrConflictEmail
				} else {
					return internalErrors.ErrConflictNickname
				}
			}
		}
		return err
	}

	return nil
}

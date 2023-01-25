package userUsecase

import (
	"database/sql"

	"github.com/jinzhu/copier"

	internalErrors "technopark-dbms-forum/internal"
	"technopark-dbms-forum/internal/models"
	userRepository "technopark-dbms-forum/internal/users/repository"
)

type UserUsecase struct {
	r *userRepository.Postgres
}

func NewUserUsecase(repo *userRepository.Postgres) *UserUsecase {
	return &UserUsecase{r: repo}
}

func (u *UserUsecase) Create(user *models.User) (interface{}, error) {
	users, err := u.r.Create(user)
	if err == internalErrors.ErrAlreadyExist {
		return users, err
	} else if err != nil {
		return nil, err
	}

	if len(users) == 1 {
		return users[0], nil
	}

	return users, nil
}

func (u *UserUsecase) GetByNickname(nickname string) (*models.User, error) {
	user, err := u.r.GetByNickname(nickname)
	if err == sql.ErrNoRows {
		return nil, internalErrors.ErrNoRows
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) Update(user *models.User) error {
	oldUser, err := u.GetByNickname(user.Nickname)
	if err != nil {
		return err
	}

	if err = copier.CopyWithOption(oldUser, user, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}

	*user = *oldUser

	if err = u.r.Update(user); err == sql.ErrNoRows {
		return internalErrors.ErrNoRows
	}

	return err
}

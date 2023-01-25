package userDelivery

import (
	"fmt"
	"net/http"

	internalErrors "technopark-dbms-forum/internal"

	"github.com/labstack/echo/v4"

	"technopark-dbms-forum/internal/models"
	userUsecase "technopark-dbms-forum/internal/users/usecase"
)

type Handler struct {
	u *userUsecase.UserUsecase
}

func NewHandler(userUsecase *userUsecase.UserUsecase) *Handler {
	return &Handler{u: userUsecase}
}

func (h *Handler) Create(c echo.Context) error {
	nickname := c.Param("nickname")

	user := models.User{
		Nickname: nickname,
	}

	err := c.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	response, err := h.u.Create(&user)
	if err == internalErrors.ErrAlreadyExist {
		return c.JSON(http.StatusConflict, response)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) Get(c echo.Context) error {
	nickname := c.Param("nickname")

	user, err := h.u.GetByNickname(nickname)
	if err == internalErrors.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find user by nickname: %s", nickname))
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) Update(c echo.Context) error {
	nickname := c.Param("nickname")

	user := models.User{
		Nickname: nickname,
	}

	err := c.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = h.u.Update(&user)
	if err == internalErrors.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find user by nickname: %s", nickname))
	} else if err == internalErrors.ErrConflictEmail {
		return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("This email is already registered by user: %s", nickname))
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, user)
}

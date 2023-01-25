package forumDelivery

import (
	"fmt"
	"net/http"
	"strconv"

	userUsecase "technopark-dbms-forum/internal/users/usecase"

	internalErrors "technopark-dbms-forum/internal"

	"github.com/labstack/echo/v4"

	forumUsecase "technopark-dbms-forum/internal/forums/usecase"
	"technopark-dbms-forum/internal/models"
)

type Handler struct {
	forumUsecase *forumUsecase.ForumUsecase
	userUsecase  *userUsecase.UserUsecase
}

func NewHandler(forumUsecase *forumUsecase.ForumUsecase, userUsecase *userUsecase.UserUsecase) *Handler {
	return &Handler{
		forumUsecase: forumUsecase,
		userUsecase:  userUsecase,
	}
}

func (h *Handler) Create(c echo.Context) error {
	slug := c.Param("slug")

	forum := models.Forum{
		Slug: slug,
	}

	err := c.Bind(&forum)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	user, err := h.userUsecase.GetByNickname(forum.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find user with nickname: %s", forum.User))
	}
	forum.User = user.Nickname

	response, err := h.forumUsecase.Create(&forum)
	if err == internalErrors.ErrAlreadyExist {
		return c.JSON(http.StatusConflict, response)
	} else if err == internalErrors.ErrUserNotFound {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find user with nickname: %s", forum.User))
	} else if err == internalErrors.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find user with nickname: %s", forum.User))
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) GetDetails(c echo.Context) error {
	slug := c.Param("slug")

	forum, err := h.forumUsecase.GetFullBySlug(slug)
	if err == internalErrors.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find forum with slug: %s", slug))
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, forum)
}

func (h *Handler) GetThreads(c echo.Context) error {
	slug := c.Param("slug")

	var limit int64 = 100
	desc := false
	since := ""
	var err error

	if c.QueryParam("limit") != "" {
		limit, err = strconv.ParseInt(c.QueryParam("limit"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	if c.QueryParam("desc") != "" {
		desc, err = strconv.ParseBool(c.QueryParam("desc"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	if c.QueryParam("since") != "" {
		since = c.QueryParam("since")
	}

	threads, err := h.forumUsecase.GetThreadsBySlug(slug, limit, since, desc)
	if err == internalErrors.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find forum with slug: %s", slug))
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, threads)
}

func (h *Handler) GetUsers(c echo.Context) error {
	slug := c.Param("slug")

	var limit int64 = 100
	desc := false
	since := ""
	var err error

	if c.QueryParam("limit") != "" {
		limit, err = strconv.ParseInt(c.QueryParam("limit"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	if c.QueryParam("desc") != "" {
		desc, err = strconv.ParseBool(c.QueryParam("desc"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	if c.QueryParam("since") != "" {
		since = c.QueryParam("since")
	}

	users, err := h.forumUsecase.GetUsersBySlug(slug, limit, since, desc)
	if err == internalErrors.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find forum by slug: %s", slug))
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, users)
}

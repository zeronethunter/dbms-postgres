package postDelivery

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"technopark-dbms-forum/internal/models"

	threadUsecase "technopark-dbms-forum/internal/threads/usecase"

	forumUsecase "technopark-dbms-forum/internal/forums/usecase"

	userUsecase "technopark-dbms-forum/internal/users/usecase"

	internalErrors "technopark-dbms-forum/internal"

	"github.com/labstack/echo/v4"

	postUsecase "technopark-dbms-forum/internal/posts/usecase"
)

type Handler struct {
	postUsecase   *postUsecase.PostUsecase
	userUsecase   *userUsecase.UserUsecase
	forumUsecase  *forumUsecase.ForumUsecase
	threadUsecase *threadUsecase.ThreadUsecase
}

func NewHandler(
	postUsecase *postUsecase.PostUsecase,
	userUsecase *userUsecase.UserUsecase,
	forumUsecase *forumUsecase.ForumUsecase,
	threadUsecase *threadUsecase.ThreadUsecase,
) *Handler {
	return &Handler{
		postUsecase:   postUsecase,
		userUsecase:   userUsecase,
		forumUsecase:  forumUsecase,
		threadUsecase: threadUsecase,
	}
}

func (h *Handler) GetInfo(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	related := strings.Split(c.QueryParam("related"), ",")

	post, err := h.postUsecase.GetByID(id)
	if err == internalErrors.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find post with id: %d", id))
	}

	fullInfo := &models.FullPost{
		Post: post,
	}

	for _, elem := range related {
		if elem == "user" {
			user, err := h.userUsecase.GetByNickname(post.Author)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			fullInfo.Author = user
		}
		if elem == "forum" {
			forum, err := h.forumUsecase.GetFullBySlug(post.Forum)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			fullInfo.Forum = forum
		}
		if elem == "thread" {
			thread, err := h.threadUsecase.GetBySlugOrID(strconv.FormatUint(post.Thread, 10))
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			fullInfo.Thread = thread
		}
	}

	return c.JSON(http.StatusOK, fullInfo)
}

func (h *Handler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	post := models.Post{
		ID: id,
	}

	err = c.Bind(&post)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	updatedPost, err := h.postUsecase.Update(&post)
	if err == internalErrors.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find post with id: %d", id))
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, updatedPost)
}

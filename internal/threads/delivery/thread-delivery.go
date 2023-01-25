package threadDelivery

import (
	"fmt"
	"net/http"
	"strconv"

	forumUsecase "technopark-dbms-forum/internal/forums/usecase"

	internalErrors "technopark-dbms-forum/internal"

	threadUsecase "technopark-dbms-forum/internal/threads/usecase"

	"github.com/labstack/echo/v4"

	"technopark-dbms-forum/internal/models"
)

type Handler struct {
	threadUsecase *threadUsecase.ThreadUsecase
	forumUsecase  *forumUsecase.ForumUsecase
}

func NewHandler(threadUsecase *threadUsecase.ThreadUsecase, forumUsecase *forumUsecase.ForumUsecase) *Handler {
	return &Handler{
		threadUsecase: threadUsecase,
		forumUsecase:  forumUsecase,
	}
}

func (h *Handler) Create(c echo.Context) error {
	slug := c.Param("slug")

	t := models.Thread{}

	if err := c.Bind(&t); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	forum, err := h.forumUsecase.GetFullBySlug(slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread forum by slug: %s", slug))
	}

	t.Forum = forum.Slug

	response, err := h.threadUsecase.Create(&t)
	if err == internalErrors.ErrAlreadyExist {
		return c.JSON(http.StatusConflict, response)
	} else if err == internalErrors.ErrSlugAlreadyExist {
		return c.JSON(http.StatusConflict, response)
	} else if err == internalErrors.ErrUserNotFound {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread author by nickname: %s", t.Author))
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) GetDetails(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")

	thread, err := h.threadUsecase.GetBySlugOrID(slugOrID)
	if err == internalErrors.ErrNoRowsBySlug {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread with slug: %s", slugOrID))
	} else if err == internalErrors.ErrNoRowsByID {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread with id: %s", slugOrID))
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, thread)
}

func (h *Handler) Update(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")

	thread := models.Thread{}

	if err := c.Bind(&thread); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	response, err := h.threadUsecase.Update(slugOrID, thread.Message, thread.Title)
	if err == internalErrors.ErrNoRowsBySlug {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread with slug: %s", slugOrID))
	} else if err == internalErrors.ErrNoRowsByID {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread with id: %s", slugOrID))
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) Vote(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")

	vote := models.Vote{}

	if err := c.Bind(&vote); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	response, err := h.threadUsecase.Vote(slugOrID, &vote)
	if err == internalErrors.ErrNoRows {
		if response.ID == 0 {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread with slug: %s", slugOrID))
		} else if response.Slug == "" {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread with id: %d", response.ID))
		}

		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread with slug: %s", slugOrID))
	} else if err == internalErrors.ErrUserNotFound {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find user by nickname: %s", vote.Nickname))
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) CreatePosts(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")

	var posts []*models.Post

	if err := c.Bind(&posts); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	response, err := h.threadUsecase.CreatePosts(slugOrID, posts)
	if err == internalErrors.ErrNoRowsBySlug {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find post thread by slug: %s", slugOrID))
	} else if err == internalErrors.ErrNoRowsByID {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find post thread by id: %s", slugOrID))
	} else if err == internalErrors.ErrPostWasCreatedInAnotherThread {
		return echo.NewHTTPError(http.StatusConflict, "Post was created in another thread")
	} else if err == internalErrors.ErrPostAuthorNotFound {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find post author by nickname: %s", posts[0].Author))
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) GetPosts(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")

	limit, err := strconv.ParseUint(c.QueryParam("limit"), 10, 64)
	if err != nil {
		limit = 100
	}
	since, err := strconv.ParseUint(c.QueryParam("since"), 10, 64)
	if err != nil {
		since = 0
	}
	sort := c.QueryParam("sort")
	if sort == "" {
		sort = "flat"
	}
	desc, err := strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		desc = false
	}

	posts, err := h.threadUsecase.GetPosts(slugOrID, limit, since, sort, desc)
	if err == internalErrors.ErrNoRowsBySlug {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread with slug: %s", slugOrID))
	} else if err == internalErrors.ErrNoRowsByID {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find thread with id: %s", slugOrID))
	} else if err == internalErrors.ErrWrongForumSlug {
		return echo.NewHTTPError(http.StatusConflict, "Parent post was created in another thread")
	} else if err == internalErrors.ErrNoParentPost {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find parent post")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, posts)
}

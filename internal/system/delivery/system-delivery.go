package systemDelivery

import (
	"net/http"

	"github.com/labstack/echo/v4"

	systemRepository "technopark-dbms-forum/internal/system/repository"
)

type Handler struct {
	systemRepo *systemRepository.Postgres
}

func NewHandler(repo *systemRepository.Postgres) *Handler {
	return &Handler{
		systemRepo: repo,
	}
}

func (h *Handler) GetInfo(c echo.Context) error {
	info, err := h.systemRepo.GetInfo()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(200, info)
}

func (h *Handler) Clear(c echo.Context) error {
	err := h.systemRepo.ClearAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(200, "OK")
}

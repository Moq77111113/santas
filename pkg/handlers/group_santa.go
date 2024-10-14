package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Group) Santa(ctx echo.Context) error {
	id := ctx.Get("id").(int)

	mms, err := h.ExclusionRepo.GenerateSanta(ctx.Request().Context(), id)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "members not found")
	}
	return ctx.JSON(http.StatusOK, mms)
}

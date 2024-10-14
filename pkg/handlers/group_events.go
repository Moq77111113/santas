package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/moq77111113/chmoly-santas/pkg/services"
)

func (h *Group) Subscribe(ctx echo.Context) error {

	id := ctx.Get("id").(int)

	_, err := h.getGroupByID(ctx)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "group not found")
	}

	h.SSE.AddClient(ctx, fmt.Sprintf("%s:%d", channelBase, id))
	return nil
}

func (h *Group) broadcast(ctx echo.Context, id int) error {

	exc, err := h.ExclusionRepo.MembersWithExclusions(ctx.Request().Context(), id)

	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(exc)

	if err != nil {
		return err
	}
	h.SSE.Broadcast(fmt.Sprintf("%s:%d", channelBase, id), services.Message{
		Data: string(jsonData),
		Type: "exclusions",
	})

	return nil
}

func (h *Group) broadcastConfig(ctx echo.Context, id int) error {
	conf, err := h.getConfig(ctx)
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	h.SSE.Broadcast(fmt.Sprintf("%s:%d", channelBase, id), services.Message{
		Data: string(jsonData),
		Type: "config",
	})
	return nil
}

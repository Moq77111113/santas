package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Join a group
func (h *Group) Join(ctx echo.Context) error {

	id := ctx.Get("id").(int)

	me, err := getCurrentUser(ctx)
	if err != nil {
		return err
	}

	err = h.GroupRepo.AddMember(ctx.Request().Context(), id, me.ID)

	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "unable to add member")
	}
	h.broadcast(ctx, id)
	return ctx.JSON(http.StatusCreated, nil)
}

// Returns a group members
func (h *Group) GetMembers(ctx echo.Context) error {
	id := ctx.Get("id").(int)

	mms, err := h.GroupRepo.GetMembers(ctx.Request().Context(), id)

	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "members not found")
	}

	return ctx.JSON(http.StatusOK, mms)
}

// Removes a member from a group
func (h *Group) RemoveMember(ctx echo.Context) error {
	id := ctx.Get("id").(int)
	memberId := ctx.Get("memberId").(int)

	me, err := getCurrentUser(ctx)
	if err != nil {
		return err
	}

	gr, err := h.getGroupByID(ctx)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "group not found")
	}

	if gr.Owner.ID != me.ID && memberId != me.ID {
		return echo.NewHTTPError(http.StatusForbidden, "You are not allowed to remove this member")
	}

	_, err = h.GroupRepo.RemoveMember(ctx.Request().Context(), id, memberId)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "Member not found")
	}

	h.broadcast(ctx, id)
	return ctx.NoContent(http.StatusNoContent)
}

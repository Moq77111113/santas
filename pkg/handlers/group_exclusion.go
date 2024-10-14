package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/moq77111113/chmoly-santas/pkg/services"
)

// Adds an exclusion to a member in a group
func (h *Group) AddExclusion(ctx echo.Context) error {
	var form addExclusionForm

	if err := form.BindAndValidate(ctx, &form); err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, form.Errors())
	}

	id := ctx.Get("id").(int)
	memberId := ctx.Get("memberId").(int)

	err := h.checkExclusion(ctx, id, memberId, form.ExcludeId)

	if err != nil {
		return err
	}

	_, err = h.ExclusionRepo.AddExclusion(ctx.Request().Context(), services.AddExclusion{
		GroupId:   id,
		MemberId:  memberId,
		ExcludeId: form.ExcludeId,
	})

	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "unable to add exclusion")
	}

	h.broadcast(ctx, id)
	return ctx.NoContent(http.StatusCreated)
}

func (h *Group) RemoveExclusion(ctx echo.Context) error {
	id := ctx.Get("id").(int)
	memberId := ctx.Get("memberId").(int)
	excludeId := ctx.Get("excludeId").(int)

	err := h.ExclusionRepo.RemoveExclusion(ctx.Request().Context(), id, memberId, excludeId)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "Member not found")
	}

	h.broadcast(ctx, id)
	return ctx.NoContent(http.StatusNoContent)
}

// Returns a group exclusions
func (h *Group) MembersWithExclusions(ctx echo.Context) error {
	id := ctx.Get("id").(int)

	exc, err := h.ExclusionRepo.MembersWithExclusions(ctx.Request().Context(), id)

	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "members not found")
	}

	return ctx.JSON(http.StatusOK, exc)
}

func (h *Group) checkExclusion(ctx echo.Context, groupId, memberId, excludeId int) error {
	if memberId == excludeId {
		return echo.NewHTTPError(http.StatusBadRequest, "member cannot exclude themselves")
	}

	maxAllowedMembers, err := h.getMaxExclusions(ctx, groupId)
	if err != nil {
		return err
	}

	mexc, err := h.ExclusionRepo.MemberWithExclusions(ctx.Request().Context(), groupId, memberId)

	if err != nil {
		return err
	}

	// Check if the member is excluding more than the permitted percentage
	if len(mexc.ExcludedMembers) >= maxAllowedMembers {
		return echo.NewHTTPError(http.StatusBadRequest, "Member is excluding more than the permitted percentage")
	}

	// Check if the member is the excluded member is excluded more than the permitted percentage
	eexc, err := h.ExclusionRepo.MembersExcludedBy(ctx.Request().Context(), groupId, excludeId)

	if err != nil {
		return err
	}

	if len(eexc.ExcludedBy) >= maxAllowedMembers {
		return echo.NewHTTPError(http.StatusBadRequest, "Member is excluded more than the permitted percentage")
	}

	return nil
}

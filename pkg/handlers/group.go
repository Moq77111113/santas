package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/moq77111113/chmoly-santas/pkg/form"
	"github.com/moq77111113/chmoly-santas/pkg/middleware"
	"github.com/moq77111113/chmoly-santas/pkg/services"
)

type (
	Group struct {
		GroupRepo     *services.GroupRepo
		ExclusionRepo *services.ExclusionRepo
		SSE           *services.SSEClient
		Auth          *services.AuthClient
	}

	createGroupForm struct {
		Name string `form:"name" validate:"required"`
		form.Form
	}

	addExclusionForm struct {
		ExcludeId int `form:"memberId" validate:"required"`
		form.Form
	}
)

const (
	channelBase = "group"
)

func init() {
	Register(new(Group))
}

func (h *Group) Init(c *services.Container) error {
	h.GroupRepo = c.Repositories.Group
	h.ExclusionRepo = c.Repositories.Exclusion
	h.SSE = c.SSE
	h.Auth = c.Auth
	return nil
}

func (h *Group) Routes(g *echo.Group) {

	groups := g.Group("/group", middleware.WithAuthentication(func(c echo.Context) error {
		return c.JSON(401, "Unauthorized")
	}))

	groups.GET("", h.List)
	groups.POST("", h.CreateGroup)

	withId := groups.Group("/:id", checkParamMw("id"))
	withId.GET("", h.GetGroup)
	withId.POST("/join", h.Join)
	withId.GET("/events", h.Subsribe)
	withId.GET("/member", h.GetMembers)

	withId.GET("/exclusion", h.MembersWithExclusions)
	withId.DELETE("/member/:memberId", h.RemoveMember, checkParamMw("memberId"))
	withId.POST("/member/:memberId/exclusion", h.AddExclusion, checkParamMw("memberId"))
	withId.DELETE("/member/:memberId/exclusion/:excludeId", h.RemoveExclusion, checkParamMw("memberId"), checkParamMw("excludeId"))
}

func (h *Group) IsApi() bool {
	return true
}

func checkParamMw(path string, types ...string) echo.MiddlewareFunc {
	paramType := "int"
	if len(types) > 0 {
		paramType = types[0]
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			param := ctx.Param(path)
			if param == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid "+path)
			}

			switch paramType {
			case "int":
				id, err := strconv.Atoi(param)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%v must be an integer", path))
				}
				ctx.Set(path, id)
			default: // unhandled other

			}

			return next(ctx)
		}
	}
}

func (h *Group) List(ctx echo.Context) error {
	grp, err := h.GroupRepo.List(ctx.Request().Context())
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusOK, []string{})
	}

	return ctx.JSON(http.StatusOK, grp)
}

// Returns a group
func (h *Group) GetGroup(ctx echo.Context) error {

	gr, err := h.getGroupByID(ctx)

	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "group not found")
	}

	return ctx.JSON(http.StatusOK, gr)
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

// Creates a group
func (h *Group) CreateGroup(ctx echo.Context) error {

	me, err := h.Auth.GetAuthenticatedUser(ctx)
	if err != nil {
		return err
	}

	var form createGroupForm
	if err := form.BindAndValidate(ctx, &form); err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, form.Errors())
	}

	gr, err := h.GroupRepo.CreateGroup(ctx.Request().Context(), me, form.Name)

	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "unable to create group")
	}

	return ctx.JSON(http.StatusCreated, gr)
}

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

// Adds an exclusion to a member in a group
func (h *Group) AddExclusion(ctx echo.Context) error {
	var form addExclusionForm

	if err := form.BindAndValidate(ctx, &form); err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, form.Errors())
	}

	id := ctx.Get("id").(int)
	memberId := ctx.Get("memberId").(int)

	_, err := h.ExclusionRepo.AddExclusion(ctx.Request().Context(), services.AddExclusion{
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

func (h *Group) Subsribe(ctx echo.Context) error {

	id := ctx.Get("id").(int)

	_, err := h.getGroupByID(ctx)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "group not found")
	}

	h.SSE.AddClient(ctx, fmt.Sprintf("%s:%d", channelBase, id))
	return nil
}

func (h *Group) getGroupByID(ctx echo.Context) (*services.EnrichedWithOwner, error) {
	id := ctx.Get("id").(int)
	return h.GroupRepo.Get(ctx.Request().Context(), id)
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
	h.SSE.Broadcast(fmt.Sprintf("%s:%d", channelBase, id), string(jsonData))

	return nil
}

package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/moq77111113/chmoly-santas/config"
	"github.com/moq77111113/chmoly-santas/pkg/form"
	"github.com/moq77111113/chmoly-santas/pkg/middleware"
	"github.com/moq77111113/chmoly-santas/pkg/services"
)

type (
	Group struct {
		Config        *config.Config
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

	GroupConfig struct {
		MaxMemberExclusions int `json:"maxMemberExclusions"`
	}
)

const (
	channelBase = "group"
)

func init() {
	Register(new(Group))
}

func (h *Group) Init(c *services.Container) error {
	h.Config = c.Config
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
	groups.POST("", h.Create)

	withId := groups.Group("/:id", checkParamMw("id"))
	withId.GET("", h.GetOne)
	withId.POST("/join", h.Join)
	withId.GET("/events", h.Subscribe)
	withId.GET("/member", h.GetMembers)
	withId.GET("/config", h.GetConfig)
	withId.GET("/exclusion", h.MembersWithExclusions)
	withId.DELETE("/member/:memberId", h.RemoveMember, checkParamMw("memberId"))
	withId.POST("/member/:memberId/exclusion", h.AddExclusion, checkParamMw("memberId"))
	withId.DELETE("/member/:memberId/exclusion/:excludeId", h.RemoveExclusion, checkParamMw("memberId"), checkParamMw("excludeId"))
	withId.GET("/santas", h.Santa)
}

func (h *Group) IsApi() bool {
	return true
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
func (h *Group) GetOne(ctx echo.Context) error {

	gr, err := h.getGroupByID(ctx)

	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "group not found")
	}

	return ctx.JSON(http.StatusOK, gr)
}

// Creates a group
func (h *Group) Create(ctx echo.Context) error {

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

// Get Config
func (h *Group) GetConfig(ctx echo.Context) error {

	maxAllowedMembers, err := h.getMaxExclusions(ctx, ctx.Get("id").(int))
	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "group not found")
	}

	c := &GroupConfig{
		MaxMemberExclusions: maxAllowedMembers,
	}
	return ctx.JSON(http.StatusOK, c)
}

func (h *Group) getGroupByID(ctx echo.Context) (*services.EnrichedWithOwner, error) {
	id := ctx.Get("id").(int)
	return h.GroupRepo.Get(ctx.Request().Context(), id)
}

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
	h.SSE.Broadcast(fmt.Sprintf("%s:%d", channelBase, id), string(jsonData))

	return nil
}

func (h *Group) getMaxExclusions(ctx echo.Context, groupId int) (int, error) {
	exclusionPercentage := h.Config.App.ExclusionPercentage // Between 0 and 1
	roundUp := h.Config.App.RoundUp

	mbs, err := h.GroupRepo.GetMembers(ctx.Request().Context(), groupId)

	if err != nil {
		return 0, err
	}

	maxExclusions := float64(len(mbs)) * float64(exclusionPercentage)
	var maxAllowedMembers int
	if roundUp {
		maxAllowedMembers = int(math.Ceil(maxExclusions))
	} else {
		maxAllowedMembers = int(math.Floor(maxExclusions))
	}

	return maxAllowedMembers, nil
}

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

	registerForm struct {
		GroupID int    `form:"groupId" validate:"required"`
		Name    string `form:"name" validate:"required"`
		form.Form
	}

	addMemberForm struct {
		MemberName string `form:"name" validate:"required"`
		form.Form
	}

	addExclusionForm struct {
		ExcludeName string `form:"name" validate:"required"`
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

	groups := g.Group("/group")

	groups.POST("", h.CreateGroup)
	groups.POST("/register", h.Register)
	// 	,
	// 	 middleware.WithoutAuthentication(func(c echo.Context) string {
	// 		return ("/group")
	// 	})
	// )

	withId := groups.Group("/:id", checkParamMw("id"),
		// Grab id from path & send to middleware.WithAuthentication
		middleware.WithAuthentication(func(c echo.Context) string {
			return fmt.Sprintf("/register?groupId=%s", c.Param("id"))
		}))
	withId.GET("", h.GetGroup)
	withId.POST("/member", h.AddMember)
	withId.GET("/events", h.RegisterChannel)
	withId.GET("/member", h.GetMembers)
	withId.GET("/exclusion", h.MembersWithExclusions)
	withId.POST("/member/:memberId/exclusion", h.AddExclusion, checkParamMw("memberId"))
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

// Returns a group
func (h *Group) GetGroup(ctx echo.Context) error {
	id := ctx.Get("id").(int)

	gr, err := h.GroupRepo.Get(ctx.Request().Context(), id)

	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "group not found")
	}

	return ctx.JSON(http.StatusOK, gr)
}

// Register a member and set the authenticated user
func (h *Group) Register(ctx echo.Context) error {
	var form registerForm

	if err := form.BindAndValidate(ctx, &form); err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, form.Errors())
	}

	mm, err := h.GroupRepo.AddMember(ctx.Request().Context(), form.GroupID, form.Name)

	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "unable to add member")
	}

	err = h.Auth.SetAuthenticatedUser(ctx, mm)
	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "unable to register")
	}

	return ctx.Redirect(http.StatusFound, "/foobar")
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

	var form createGroupForm
	if err := form.BindAndValidate(ctx, &form); err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, form.Errors())
	}

	gr, err := h.GroupRepo.CreateGroup(ctx.Request().Context(), form.Name)

	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "unable to create group")
	}

	return ctx.JSON(http.StatusCreated, gr)
}

// Adds a member to a group
func (h *Group) AddMember(ctx echo.Context) error {
	var form addMemberForm
	if err := form.BindAndValidate(ctx, &form); err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, form.Errors())
	}

	id := ctx.Get("id").(int)

	mm, err := h.GroupRepo.AddMember(ctx.Request().Context(), id, form.MemberName)

	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "unable to add member")
	}
	h.broadcast(ctx, id)
	return ctx.JSON(http.StatusCreated, mm)
}

// Removes a member from a group
func (h *Group) RemoveMember(ctx echo.Context) error {
	id := ctx.Get("id").(int)
	memberId := ctx.Get("memberId").(int)

	_, err := h.GroupRepo.RemoveMember(ctx.Request().Context(), id, memberId)
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
	mms, err := h.GroupRepo.GetMembers(ctx.Request().Context(), id)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "Member not found")
	}
	var excludeId int
	for _, mm := range mms {
		// Compare insensitive name
		if strings.EqualFold(mm.Name, form.ExcludeName) {
			excludeId = mm.ID
			break
		}
	}
	if excludeId == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Member not found")
	}

	_, err = h.ExclusionRepo.AddExclusion(ctx.Request().Context(), services.AddExclusion{
		GroupId:   id,
		MemberId:  memberId,
		ExcludeId: excludeId,
	})

	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "unable to add exclusion")
	}

	h.broadcast(ctx, id)
	return ctx.NoContent(http.StatusCreated)
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

func (h *Group) RegisterChannel(ctx echo.Context) error {

	id := ctx.Get("id").(int)

	_, err := h.GroupRepo.Get(ctx.Request().Context(), id)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "group not found")
	}

	h.SSE.AddClient(ctx, fmt.Sprintf("%s:%d", channelBase, id))
	return nil
}

func (h *Group) BroadcastGroup(id int, message string) {
	h.SSE.Broadcast(fmt.Sprintf("%s:%d", channelBase, id), message)
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

	h.BroadcastGroup(id, string(jsonData))
	return nil
}

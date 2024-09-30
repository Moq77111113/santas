package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/moq77111113/chmoly-santas/internal/core"
)

type (
	Group struct {
		Participants []Participant `json:"participants"`
	}

	Participant struct {
		Name       string   `json:"name"`
		Exclusions []string `json:"exclusions"`
	}
)

type groupApi struct {
	app core.App
}

var groups = map[string]Group{}

func (api *groupApi) notify() error {
	data, err := json.Marshal(groups)
	if err != nil {
		return err
	}

	fmt.Sprintln("%v", string(data))

	api.app.Notifier.Notify(string(data))
	return nil
}
func bindGroupApi(app core.App, rg *echo.Group) {
	api := groupApi{app: app}

	group := rg.Group("/group")
	group.GET("", api.list)
	group.POST("", api.create)
	group.POST("/:id/participant", api.addParticipant)
	group.DELETE("/:id/participant/:name", api.removeParticipant)

}

func (api *groupApi) list(c echo.Context) error {
	return c.JSON(200, groups)
}

func (api *groupApi) create(c echo.Context) error {

	var group Group
	if err := c.Bind(&group); err != nil {
		return err
	}

	groups[fmt.Sprintf("%v", len(groups)+1)] = group

	api.notify()
	return c.JSON(200, group)
}

func (api *groupApi) addParticipant(c echo.Context) error {

	if _, ok := groups[c.Param("id")]; !ok {
		return echo.ErrNotFound
	}
	group := groups[c.Param("id")]
	var participant Participant
	if err := c.Bind(&participant); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid participant")
	}

	if len(participant.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid participant")
	}

	group.Participants = append(groups[c.Param("id")].Participants, participant)
	groups[c.Param("id")] = group

	api.notify()
	return c.JSON(200, group)
}

func (api *groupApi) removeParticipant(c echo.Context) error {

	if _, ok := groups[c.Param("id")]; !ok {
		return echo.ErrNotFound
	}
	group := groups[c.Param("id")]

	if len(group.Participants) == 0 {
		return echo.ErrNotFound
	}

	trimmedName := strings.TrimSpace(c.Param("name"))
	for i, p := range group.Participants {

		if strings.EqualFold((trimmedName), strings.TrimSpace(p.Name)) {
			group.Participants = append(group.Participants[:i], group.Participants[i+1:]...)
			groups[c.Param("id")] = group
			api.notify()
			return c.JSON(200, group)
		}
	}

	return echo.NewHTTPError(http.StatusNotFound, "Participant not found")
}

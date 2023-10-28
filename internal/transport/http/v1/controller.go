package v1

import (
	"github.com/labstack/echo/v4"
)

type Controller interface {
	InitRoutes(e *echo.Group)
}

type controller struct {
	People Controller
}

func NewController(peopleContoller Controller) *controller {
	return &controller{People: peopleContoller}
}

func (h *controller) InitRoutes(api *echo.Group) {
	v1 := api.Group("/v1")

	h.People.InitRoutes(v1)
}

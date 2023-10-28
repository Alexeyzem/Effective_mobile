package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gitlab.com/petly-app/backend/pkg/errors"
	"local/EffectiveMobile/pkg/logger"
	"net/http"
)

type errResponse struct {
	Error string `json:"error"`
}

type ResponsePresenter struct {
}

func (p *ResponsePresenter) CreatedResponse(c echo.Context, body interface{}, createdEntityID string) error {
	loc := fmt.Sprintf("%s/%s", c.Path(), createdEntityID)
	c.Response().Header().Set("Location", loc)
	return c.JSON(http.StatusCreated, body)
}

func (p *ResponsePresenter) OkResponse(c echo.Context, body interface{}) error {
	if body == nil {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, body)
}

func (p *ResponsePresenter) FailResponse(c echo.Context, err error, log logger.Logger) error {
	rlog := logger.WithPrefix(log, "RequestID", c.Response().Header().Get(echo.HeaderXRequestID))
	rlog.Error(err)

	responseErr := err
	statusCode := 0
	switch {
	case errors.Is(err, errors.New("not found")):
		statusCode = http.StatusNotFound
	case errors.Is(err, errors.New("bad request")):
		statusCode = http.StatusBadRequest
	case errors.Is(err, errors.New("internal server")):
		statusCode = http.StatusInternalServerError
	default:
		statusCode = http.StatusInternalServerError
		responseErr = errors.MaskWrapErr(errors.New("internal server"), err)
	}

	return c.JSON(statusCode, errResponse{Error: responseErr.Error()})
}

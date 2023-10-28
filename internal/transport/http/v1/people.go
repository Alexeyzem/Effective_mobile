package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"local/EffectiveMobile/internal/domain"
	"local/EffectiveMobile/internal/transport/http"
	v1 "local/EffectiveMobile/internal/transport/http/v1/model"
	"local/EffectiveMobile/pkg/logger"
	http2 "net/http"
	"strconv"
	"time"
)

type PeopleService interface {
	Create(ctx context.Context, people domain.People) (domain.People, error)
	People(ctx context.Context, id string) (domain.People, error)
	PeopleList(ctx context.Context, count int) ([]domain.People, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, people domain.People) (domain.People, error)
}

type PeopleController struct {
	peopleService PeopleService
	validator     *validator.Validate
	log           logger.Logger
	http.ResponsePresenter
}

func NewPeopleConroller(service PeopleService, log logger.Logger, validator *validator.Validate) *PeopleController {
	return &PeopleController{
		peopleService: service,
		log:           log,
		validator:     validator,
	}
}

func (pc *PeopleController) InitRoutes(v1 *echo.Group) {
	peopleGroup := v1.Group("/people")
	{
		peopleGroup.GET("", pc.handleGetPeopleList)
		peopleGroup.GET("/:id", pc.handleGetPeople)
		peopleGroup.DELETE("/:id", pc.handleDeletePeople)
		peopleGroup.POST("", pc.handleCreatePeople)
		peopleGroup.PUT("/:id", pc.handleUpdatePeople)
	}

}

func (pc *PeopleController) handleGetPeopleList(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 1*time.Minute)
	defer cancel()

	r := new(v1.GetPeopleListRequest)
	if err := c.Bind(r); err != nil {
		return pc.FailResponse(c, fmt.Errorf("%w: %v", err, errors.New("bad request")), pc.log)
	}

	count := c.QueryParam("count")
	if count == "" {
		count = "0"
	}
	log := logger.WithPrefix(pc.log, "Count", count)
	cnt, err := strconv.Atoi(count)
	if err != nil {
		return pc.FailResponse(c, fmt.Errorf("%w: %v", err, errors.New("bad request")), log)
	}

	people, err := pc.peopleService.PeopleList(ctx, cnt)
	if err != nil {
		return pc.FailResponse(c, errors.Wrap(err, "people service PeopleList"), log)
	}

	return pc.OkResponse(c, v1.PeopleListFromDomain(people))
}

func (pc *PeopleController) handleGetPeople(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 1*time.Minute)
	defer cancel()

	id := c.Param("id")
	log := logger.WithPrefix(pc.log, "PeopleID", id)

	people, err := pc.peopleService.People(ctx, id)
	if err != nil {
		return pc.FailResponse(c, errors.Wrap(err, "people service People"), log)
	}

	return pc.OkResponse(c, people)
}

func (pc *PeopleController) handleCreatePeople(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 1*time.Minute)
	defer cancel()

	people := &v1.CommonPeopleRequest{}
	if err := c.Bind(people); err != nil {
		return pc.FailResponse(c, fmt.Errorf("%w: %v", err, errors.New("bad request")), pc.log)
	}

	err := enrichment(fmt.Sprintf("https://api.agify.io/?name=%s", people.FirstName), people)
	if err != nil {
		return pc.FailResponse(c, errors.Wrap(err, "people service create"), pc.log)
	}

	err = enrichment(fmt.Sprintf("https://api.genderize.io/?name=%s", people.FirstName), people)
	if err != nil {
		return pc.FailResponse(c, errors.Wrap(err, "people service create"), pc.log)
	}

	err = enrichmentNation(people)
	if err != nil {
		return pc.FailResponse(c, errors.Wrap(err, "people service create"), pc.log)
	}

	CreatedPeople, err := pc.peopleService.Create(ctx, people.ToDomain())
	if err != nil {
		return pc.FailResponse(c, errors.Wrap(err, "people service create"), pc.log)
	}
	return pc.CreatedResponse(c, CreatedPeople, CreatedPeople.ID)
}

func enrichment(url string, people *v1.CommonPeopleRequest) error {
	r, err := http2.Get(url)
	if err != nil {
		return err
	}
	bs := make([]byte, 1024)
	n, err := r.Body.Read(bs)
	if err != nil {
		return err
	}
	bs = bs[:n]
	if err = json.Unmarshal(bs, &people); err != nil {
		return err
	}
	return nil
}

type ForParseCountryId struct {
	CountryId string `json:"country_id"`
}

type forParseCountry struct {
	Country []ForParseCountryId `json:"country"`
}

func enrichmentNation(people *v1.CommonPeopleRequest) error {
	url := fmt.Sprintf("https://api.nationalize.io/?name=%s", people.FirstName)
	r, err := http2.Get(url)
	if err != nil {
		return err
	}
	bs := make([]byte, 1024)
	n, err := r.Body.Read(bs)
	if err != nil {
		return err
	}
	bs = bs[:n]
	var nation forParseCountry
	if err = json.Unmarshal(bs, &nation); err != nil {
		return err
	}
	if len(nation.Country) < 1 {
		return errors.New("people service nation")
	}
	people.Nationality = nation.Country[0].CountryId
	return nil
}

func (pc *PeopleController) handleUpdatePeople(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 1*time.Minute)
	defer cancel()

	id := c.Param("id")
	log := logger.WithPrefix(pc.log, "PeopleId", id)

	people := new(v1.CommonPeopleRequest)
	if err := c.Bind(people); err != nil {
		return pc.FailResponse(c, errors.Wrap(err, "people service update"), log)
	}

	err := enrichment(fmt.Sprintf("https://api.agify.io/?name=%s", people.FirstName), people)
	if err != nil {
		return pc.FailResponse(c, errors.Wrap(err, "people service create"), pc.log)
	}

	err = enrichment(fmt.Sprintf("https://api.genderize.io/?name=%s", people.FirstName), people)
	if err != nil {
		return pc.FailResponse(c, errors.Wrap(err, "people service create"), pc.log)
	}

	err = enrichmentNation(people)
	if err != nil {
		return pc.FailResponse(c, errors.Wrap(err, "people service create"), pc.log)
	}

	updatedPeople, err := pc.peopleService.Update(ctx, people.ToDomain())

	if err != nil {
		return pc.FailResponse(c, errors.Wrap(err, "people service update"), log)
	}

	return pc.OkResponse(c, updatedPeople)
}

func (pc *PeopleController) handleDeletePeople(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 1*time.Minute)
	defer cancel()

	id := c.Param("id")
	log := logger.WithPrefix(pc.log, "PeopleId", id)

	if err := pc.peopleService.Delete(ctx, id); err != nil {
		return pc.FailResponse(c, errors.Wrap(err, "People service delete"), log)
	}

	return pc.OkResponse(c, nil)
}

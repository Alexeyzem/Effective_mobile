package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"local/EffectiveMobile/internal/domain"
)

type PeopleRepository interface {
	Create(ctx context.Context, people domain.People) (domain.People, error)
	Find(ctx context.Context, id string) (domain.People, error)
	PeopleList(ctx context.Context, count int) ([]domain.People, error)
	Update(ctx context.Context, people domain.People) (domain.People, error)
	Delete(ctx context.Context, id string) error
}

type People struct {
	repo PeopleRepository
}

func NewPeople(repo PeopleRepository) *People {
	return &People{repo: repo}
}

func (p *People) People(ctx context.Context, id string) (domain.People, error) {
	people, err := p.repo.Find(ctx, id)
	if err != nil {
		return domain.People{}, errors.Wrap(err, "repository find people")
	}

	return people, nil
}

func (p *People) Create(ctx context.Context, people domain.People) (domain.People, error) {
	people.ID = uuid.New().String()
	createdPeople, err := p.repo.Create(ctx, people)
	if err != nil {
		return domain.People{}, errors.Wrap(err, "repository create people")
	}

	return createdPeople, nil
}

func (p *People) PeopleList(ctx context.Context, count int) ([]domain.People, error) {
	peopleList, err := p.repo.PeopleList(ctx, count)
	if err != nil {
		return nil, errors.Wrap(err, "repository get people list")
	}

	return peopleList, nil
}

func (p *People) Update(ctx context.Context, people domain.People) (domain.People, error) {
	updatedPeople, err := p.repo.Update(ctx, people)
	if err != nil {
		return domain.People{}, errors.Wrap(err, "repository update people")
	}

	return updatedPeople, nil
}

func (p *People) Delete(ctx context.Context, id string) error {
	if err := p.repo.Delete(ctx, id); err != nil {
		return errors.Wrap(err, "repository delete people")
	}

	return nil
}

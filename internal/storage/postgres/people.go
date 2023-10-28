package storage

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"local/EffectiveMobile/internal/domain"
	"local/EffectiveMobile/internal/storage/dbo"
	"time"
)

type PeoplePostgres struct {
	db *gorm.DB
}

func NewPeoplePostgres(db *gorm.DB) *PeoplePostgres {
	return &PeoplePostgres{db: db}
}

func (r *PeoplePostgres) Create(ctx context.Context, people domain.People) (domain.People, error) {
	dboPeople := dbo.PeopleFromDomain(people)
	err := r.db.WithContext(ctx).Create(&dboPeople).Error
	if err != nil {
		return domain.People{}, errors.Wrap(err, "postgres gorm create")
	}

	return dboPeople.ToDomain(), nil
}

func (r *PeoplePostgres) Find(ctx context.Context, id string) (domain.People, error) {
	var people dbo.People

	err := r.db.WithContext(ctx).First(&people, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.People{}, errors.New("not found")
	}

	if err != nil {
		return domain.People{}, errors.Wrap(err, "postgres gorm first")
	}

	return people.ToDomain(), nil
}

func (r *PeoplePostgres) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&dbo.People{}).Error
	if err != nil {
		return errors.Wrap(err, "postgres gorm delete")
	}

	return nil
}

func (r *PeoplePostgres) Update(ctx context.Context, people domain.People) (domain.People, error) {
	dboPeople := dbo.PeopleFromDomain(people)
	updates := map[string]interface{}{
		"first_name":  dboPeople.FirstName,
		"last_name":   dboPeople.LastName,
		"middle_name": dboPeople.MiddleName,
		"age":         dboPeople.Age,
		"gender":      dboPeople.Gender,
		"nationality": dboPeople.Nationality,
		"updated_at":  time.Now(),
	}
	result := r.db.WithContext(ctx).Clauses(clause.Returning{}).Model(&dboPeople).Updates(updates)

	if result.RowsAffected == 0 {
		return domain.People{}, errors.New("not found")
	}

	if result.Error != nil {
		return domain.People{}, errors.Wrap(result.Error, "postgres gorm update")
	}

	return dboPeople.ToDomain(), nil
}

func (r *PeoplePostgres) PeopleList(ctx context.Context, count int) ([]domain.People, error) {
	dboPeople := make([]dbo.People, 0)
	var result *gorm.DB
	if count == 0 {
		result = r.db.WithContext(ctx).Select("id, first_name, last_name, middle_name, gender, age, nationality").Find(&dboPeople)
	} else {
		result = r.db.WithContext(ctx).Select("id, first_name, last_name, middle_name, gender, age, nationality").Limit(count).Find(&dboPeople)
	}
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "postgres gorm people list")
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("not found")
	}

	people := make([]domain.People, len(dboPeople))
	for i, d := range dboPeople {
		people[i] = d.ToDomain()
	}

	return people, nil
}

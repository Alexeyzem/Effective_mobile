package dbo

import (
	"local/EffectiveMobile/internal/domain"
	"time"
)

const peopleTables = "people"

type People struct {
	ID          string    `gorm:"primaryKey"`
	FirstName   string    `gorm:"column:first_name"`
	LastName    string    `gorm:"column:last_name"`
	MiddleName  string    `gorm:"column:middle_name"`
	Age         int       `gorm:"column:age"`
	Gender      string    `gorm:"column:gender"`
	Nationality string    `gorm:"column:nationality"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (dbo *People) ToDomain() domain.People {
	return domain.People{
		ID:          dbo.ID,
		FirstName:   dbo.FirstName,
		LastName:    dbo.LastName,
		MiddleName:  dbo.MiddleName,
		Age:         dbo.Age,
		Gender:      dbo.Gender,
		Nationality: dbo.Nationality,
	}
}

func PeopleFromDomain(people domain.People) People {
	createTime := time.Now()
	return People{
		ID:          people.ID,
		FirstName:   people.FirstName,
		LastName:    people.LastName,
		MiddleName:  people.MiddleName,
		Age:         people.Age,
		Gender:      people.Gender,
		Nationality: people.Nationality,
		CreatedAt:   createTime,
		UpdatedAt:   createTime,
	}
}
func (dbo *People) TableName() string { return peopleTables }

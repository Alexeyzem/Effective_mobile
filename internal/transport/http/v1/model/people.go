package v1

import "local/EffectiveMobile/internal/domain"

type CommonPeopleRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	MiddleName  string `json:"middle_name"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}
type GetPeopleListRequest struct {
	Count string `json:"count"`
}

func (r *CommonPeopleRequest) ToDomain() domain.People {
	return domain.People{
		FirstName:   r.FirstName,
		MiddleName:  r.MiddleName,
		LastName:    r.LastName,
		Age:         r.Age,
		Gender:      r.Gender,
		Nationality: r.Nationality,
	}
}

type UpdatePeopleRequest struct {
	CommonPeopleRequest
}

func (r *UpdatePeopleRequest) ToDomain(id string) domain.People {
	return domain.People{
		ID:          id,
		FirstName:   r.FirstName,
		MiddleName:  r.MiddleName,
		LastName:    r.LastName,
		Age:         r.Age,
		Gender:      r.Gender,
		Nationality: r.Nationality,
	}
}

type PeopleList struct {
	People []CommonPeopleRequest `json:"people"`
}

func PeopleListFromDomain(people []domain.People) PeopleList {
	outPeople := make([]CommonPeopleRequest, len(people))
	for i, human := range people {
		outPeople[i].Age = human.Age
		outPeople[i].Gender = human.Gender
		outPeople[i].Nationality = human.Nationality
		outPeople[i].FirstName = human.FirstName
		outPeople[i].LastName = human.LastName
		outPeople[i].MiddleName = human.MiddleName
	}
	return PeopleList{outPeople}
}

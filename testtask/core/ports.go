package core

import (
	"context"
)

type APIClient interface {
	GetAge(string) (APIAgeResponse, error)
	GetNation(string) (APINationResponse, error)
	GetGender(string) (APIGenderResponse, error)
}

type DB interface {
	GetPerson(context.Context, string) (Person, error)
	CreatePerson(context.Context, Person) error
	GetPeople(context.Context, PersonFilters) ([]Person, error)
	UpdatePerson(context.Context, Person) error
	DeletePerson(context.Context, string) error
}

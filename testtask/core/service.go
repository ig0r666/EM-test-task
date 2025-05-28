package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type Service struct {
	db  DB
	api APIClient
	log *slog.Logger
}

func NewService(db DB, api APIClient, log *slog.Logger) *Service {
	return &Service{
		db:  db,
		api: api,
		log: log,
	}
}

func (s *Service) CreatePerson(ctx context.Context, req PersonRequest) (*Person, error) {
	age, err := s.api.GetAge(req.Name)
	if err != nil {
		s.log.Error("failed to get age", "error", err)
		return nil, ErrAPIFailed
	}

	gender, err := s.api.GetGender(req.Name)
	if err != nil {
		s.log.Error("failed to get gender", "error", err)
		return nil, ErrAPIFailed
	}

	nation, err := s.api.GetNation(req.Name)
	if err != nil {
		s.log.Error("failed to get nation", "error", err)
		return nil, ErrAPIFailed
	}

	id, err := uuid.NewRandom()
	if err != nil {
		s.log.Error("failed to generate UUID", "error", err)
		return nil, fmt.Errorf("failed to generate ID")
	}

	person := Person{
		ID:          id.String(),
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         age.Age,
		Gender:      gender.Gender,
		Nationality: nation.Country[0].CountryID,
	}

	if err := s.db.CreatePerson(ctx, person); err != nil {
		s.log.Error("failed to create person", "error", err)
		return nil, err
	}

	return &person, nil
}

func (s *Service) GetPeople(ctx context.Context, filters PersonFilters) ([]Person, error) {
	people, err := s.db.GetPeople(ctx, filters)
	if err != nil {
		if errors.Is(err, ErrPersonNotFound) {
			return nil, err
		}
		s.log.Error("failed to get people", "error", err)
		return nil, err
	}
	return people, nil
}

func (s *Service) GetPerson(ctx context.Context, id string) (*Person, error) {
	person, err := s.db.GetPerson(ctx, id)
	if err != nil {
		if errors.Is(err, ErrPersonNotFound) {
			return nil, err
		}
		s.log.Error("failed to get person", "id", id, "error", err)
		return nil, err
	}
	return &person, nil
}

func (s *Service) UpdatePerson(ctx context.Context, person Person) error {
	if err := s.db.UpdatePerson(ctx, person); err != nil {
		if errors.Is(err, ErrPersonNotFound) {
			return err
		}
		s.log.Error("failed to update person", "error", err)
		return err
	}
	return nil
}

func (s *Service) DeletePerson(ctx context.Context, id string) error {
	if err := s.db.DeletePerson(ctx, id); err != nil {
		s.log.Error("failed to delete person", "id", id, "error", err)
		return err
	}
	return nil
}

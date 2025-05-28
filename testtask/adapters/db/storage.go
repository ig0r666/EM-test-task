package db

import (
	"EMtask/testtask/core"
	"context"
	"database/sql"
	"errors"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
}

func New(log *slog.Logger, address string) (*DB, error) {

	db, err := sqlx.Connect("pgx", address)
	if err != nil {
		log.Error("connection problem", "address", address, "error", err)
		return nil, err
	}

	return &DB{
		log:  log,
		conn: db,
	}, nil
}

func (d *DB) GetPerson(ctx context.Context, id string) (core.Person, error) {
	const query = `
		SELECT people_id, name, surname, patronymic, age, gender, nationality
		FROM people 
		WHERE people_id = $1
	`

	var person core.Person
	if err := d.conn.GetContext(ctx, &person, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.Person{}, core.ErrPersonNotFound
		}
		d.log.Error("failed to get person", "id", id, "error", err)
		return core.Person{}, err
	}

	return person, nil
}

func (d *DB) GetPeople(ctx context.Context, filters core.PersonFilters) ([]core.Person, error) {
	query := `
        SELECT people_id, name, surname, patronymic, age, gender, nationality
        FROM people 
        WHERE 1=1
    `
	args := []interface{}{}

	if filters.Age != "" {
		query += " AND age = ?"
		args = append(args, filters.Age)
	}
	if filters.Gender != "" {
		query += " AND gender = ?"
		args = append(args, filters.Gender)
	}
	if filters.Nationality != "" {
		query += " AND nationality = ?"
		args = append(args, filters.Nationality)
	}

	limit := "10"
	if filters.Limit != "" {
		limit = filters.Limit
	}
	offset := "0"
	if filters.Offset != "" {
		offset = filters.Offset
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	query = d.conn.Rebind(query)

	var people []core.Person
	if err := d.conn.SelectContext(ctx, &people, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrPersonNotFound
		}
		d.log.Error("failed to get people", "error", err, "filters", filters)
		return nil, err
	}

	return people, nil
}

func (d *DB) CreatePerson(ctx context.Context, person core.Person) error {
	const query = `
		INSERT INTO people 
		(people_id, name, surname, patronymic, age, gender, nationality)
		VALUES (:people_id, :name, :surname, :patronymic, :age, :gender, :nationality)
	`

	_, err := d.conn.NamedExecContext(ctx, query, person)
	if err != nil {
		d.log.Error("failed to create person", "error", err, "person", person)
		return err
	}

	return nil
}

func (d *DB) UpdatePerson(ctx context.Context, person core.Person) error {
	const query = `
		UPDATE people 
		SET 
			name = :name,
			surname = :surname,
			patronymic = :patronymic,
			age = :age,
			gender = :gender
		WHERE people_id = :people_id
	`

	result, err := d.conn.NamedExecContext(ctx, query, person)
	if err != nil {
		d.log.Error("failed to update person", "error", err, "person", person)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return core.ErrPersonNotFound
	}

	return nil
}

func (d *DB) DeletePerson(ctx context.Context, id string) error {
	const query = `DELETE FROM people WHERE people_id = $1`

	result, err := d.conn.ExecContext(ctx, query, id)
	if err != nil {
		d.log.Error("failed to delete person", "error", err, "id", id)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return core.ErrPersonNotFound
	}

	return nil
}

package repository

import (
	"context"
	"errors"
	"fmt"

	"neuroscan/internal/domain"
	"neuroscan/internal/toolshed"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactRepository interface {
	GetContactByID(ctx context.Context, uid string, timepoint int) (domain.Contact, error)
	GetContactByUID(ctx context.Context, uid string, timepoint int) (domain.Contact, error)
	ContactExists(ctx context.Context, uid string, timepoint int) (bool, error)
	SearchContacts(ctx context.Context, query domain.APIV1Request) ([]domain.Contact, error)
	CountContacts(ctx context.Context, query domain.APIV1Request) (int, error)
	CreateContact(ctx context.Context, uid string, filename string, timepoint int, color toolshed.Color) error
	DeleteContact(ctx context.Context, uid string, timepoint int) error
	IngestContact(ctx context.Context, contact domain.Contact, skipExisting bool, force bool) (bool, error)
}

type PostgresContactRepository struct {
	DB *pgxpool.Pool
}

func NewPostgresContactRepository(db *pgxpool.Pool) *PostgresContactRepository {
	return &PostgresContactRepository{
		DB: db,
	}
}

func (r *PostgresContactRepository) GetContactByID(ctx context.Context, uid string, timepoint int) (domain.Contact, error) {
	query := "SELECT id, uid, timepoint, filename, color FROM contacts WHERE id = $1 AND timepoint = $2"

	var contact domain.Contact
	err := r.DB.QueryRow(ctx, query, uid, timepoint).Scan(&contact.ID, &contact.UID, &contact.Timepoint, &contact.Filename, &contact.Color)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Contact{}, nil
		}

		return domain.Contact{}, err
	}

	return contact, nil
}

func (r *PostgresContactRepository) GetContactByUID(ctx context.Context, uid string, timepoint int) (domain.Contact, error) {
	query := "SELECT id, uid, timepoint, filename, color FROM contacts WHERE uid = $1 AND timepoint = $2"

	var contact domain.Contact
	err := r.DB.QueryRow(ctx, query, uid, timepoint).Scan(&contact.ID, &contact.UID, &contact.Timepoint, &contact.Filename, &contact.Color)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Contact{}, nil
		}

		return domain.Contact{}, err
	}

	return contact, nil
}

func (r *PostgresContactRepository) ContactExists(ctx context.Context, uid string, timepoint int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM contacts WHERE uid = $1 AND timepoint = $2)"

	var exists bool
	err := r.DB.QueryRow(ctx, query, uid, timepoint).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return exists, nil
}

func (r *PostgresContactRepository) SearchContacts(ctx context.Context, query domain.APIV1Request) ([]domain.Contact, error) {
	q := "SELECT id, uid, timepoint, filename, color FROM contacts "

	parsedQuery, args := query.ToPostgresQuery()

	q += parsedQuery

	rows, _ := r.DB.Query(ctx, q, args...)

	contacts, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Contact])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.Contact{}, nil
		}

		return nil, err
	}

	return contacts, err
}

func (r *PostgresContactRepository) CountContacts(ctx context.Context, query domain.APIV1Request) (int, error) {
	var count int

	q := "SELECT COUNT(*) FROM contacts "

	parsedQuery, args := query.ToPostgresQuery()

	q += parsedQuery

	err := r.DB.QueryRow(ctx, q, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresContactRepository) CreateContact(ctx context.Context, uid string, filename string, timepoint int, color toolshed.Color) error {
	exists, err := r.ContactExists(ctx, uid, timepoint)

	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("contact already exists")
	}

	query := "INSERT INTO contacts (uid, timepoint, filename, color) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING"

	_, err = r.DB.Exec(ctx, query, uid, timepoint, filename, color)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresContactRepository) DeleteContact(ctx context.Context, uid string, timepoint int) error {
	query := "DELETE FROM contacts WHERE uid = $1 AND timepoint = $2"

	_, err := r.DB.Exec(ctx, query, uid, timepoint)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresContactRepository) IngestContact(ctx context.Context, contact domain.Contact, skipExisting bool, force bool) (bool, error) {
	exists, err := r.ContactExists(ctx, contact.UID, contact.Timepoint)

	if err != nil {
		return false, err
	}

	if skipExisting && exists {
		return true, nil
	}

	if force && exists {
		err = r.DeleteContact(ctx, contact.UID, contact.Timepoint)
		if err != nil {
			return false, err
		}
	}

	err = r.CreateContact(ctx, contact.UID, contact.Filename, contact.Timepoint, contact.Color)
	if err != nil {
		return false, err
	}

	return true, nil
}

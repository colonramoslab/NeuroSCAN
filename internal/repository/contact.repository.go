package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"neuroscan/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactRepository interface {
	GetContactByUID(ctx context.Context, uid string, timepoint int) (domain.Contact, error)
	ContactExists(ctx context.Context, uid string, timepoint int) (bool, error)
	SearchContacts(ctx context.Context, query domain.APIV1Request) ([]domain.Contact, error)
	CountContacts(ctx context.Context, query domain.APIV1Request) (int, error)
	CreateContact(ctx context.Context, contact domain.Contact) error
	IngestContact(ctx context.Context, contact domain.Contact, skipExisting bool, force bool) (bool, error)
	TruncateContacts(ctx context.Context) error
}

type PostgresContactRepository struct {
	DB *pgxpool.Pool
}

func NewPostgresContactRepository(db *pgxpool.Pool) *PostgresContactRepository {
	return &PostgresContactRepository{
		DB: db,
	}
}

func (r *PostgresContactRepository) GetContactByUID(ctx context.Context, uid string, timepoint int) (domain.Contact, error) {
	query := "SELECT id, uid, ulid, timepoint, filename, color FROM contacts WHERE uid = $1 AND timepoint = $2"

	var contact domain.Contact
	err := r.DB.QueryRow(ctx, query, uid, timepoint).Scan(&contact.ID, &contact.UID, &contact.ULID, &contact.Timepoint, &contact.Filename, &contact.Color)
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
	q := "SELECT id, uid, ulid, timepoint, filename, color FROM contacts "

	parsedQuery, args := r.ParseContactAPIV1Request(ctx, query)

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

	parsedQuery, args := r.ParseContactAPIV1Request(ctx, query)

	q += parsedQuery

	err := r.DB.QueryRow(ctx, q, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresContactRepository) CreateContact(ctx context.Context, contact domain.Contact) error {
	exists, err := r.ContactExists(ctx, contact.UID, contact.Timepoint)

	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("contact already exists")
	}

	query := "INSERT INTO contacts (uid, ulid, timepoint, filename, color) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING"

	_, err = r.DB.Exec(ctx, query, contact.UID, contact.ULID, contact.Timepoint, contact.Filename, contact.Color)
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

	err = r.CreateContact(ctx, contact)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *PostgresContactRepository) TruncateContacts(ctx context.Context) error {
	query := "TRUNCATE TABLE contacts RESTART IDENTITY CASCADE"

	_, err := r.DB.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresContactRepository) ParseContactAPIV1Request(ctx context.Context, req domain.APIV1Request) (string, []interface{}) {

	queryParts := []string{"where 1=1"}
	args := []interface{}{}

	if req.Timepoint != nil {
		args = append(args, req.Timepoint)
		queryParts = append(queryParts, fmt.Sprintf("timepoint = $%d", len(args)))
	}

	if len(req.UIDs) > 0 {
		// we need to build a query where UID is like or, looping over the UIDs, wrapping them in % and adding them to the array[]
		uidArray := []string{}
		for _, uid := range req.UIDs {
			uidArray = append(uidArray, fmt.Sprintf("%%%s%%", strings.ToLower(uid)))
		}
		args = append(args, uidArray)
		queryParts = append(queryParts, fmt.Sprintf("LOWER(uid) ILIKE ANY($%d)", len(args)))
	}

	query := strings.Join(queryParts, " AND ")

	// if count is true, return the query and args before adding the sort and limit
	if req.Count {
		return query, args
	}

	if req.Sort != "" {
		// split by ":", first part is the field, second part is the direction
		parts := strings.Split(req.Sort, ":")

		if len(parts) == 2 {

			// if the second part is not asc or desc, default to asc
			if parts[1] != "asc" && parts[1] != "desc" {
				parts[1] = "asc"
			}

			query += fmt.Sprintf(" order by %s %s", parts[0], parts[1])
		}
	}

	if req.Limit > 0 {
		args = append(args, req.Limit)
		query += fmt.Sprintf(" limit $%d", len(args))
	} else {
		query += " limit 100"
	}

	if req.Offset > 0 {
		args = append(args, req.Offset)
		query += fmt.Sprintf(" offset $%d", len(args))
	}

	return query, args
}

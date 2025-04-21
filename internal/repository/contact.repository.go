package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"neuroscan/internal/cache"
	"neuroscan/internal/domain"
	"neuroscan/internal/toolshed"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactRepository interface {
	GetContactByULID(ctx context.Context, id string) (domain.Contact, error)
	GetContactByUID(ctx context.Context, uid string, timepoint int) (domain.Contact, error)
	ContactExists(ctx context.Context, uid string, timepoint int) (bool, error)
	SearchContacts(ctx context.Context, query domain.APIV1Request) ([]domain.Contact, error)
	CountContacts(ctx context.Context, query domain.APIV1Request) (int, error)
	CreateContact(ctx context.Context, contact domain.Contact) error
	UpdateContact(ctx context.Context, contact domain.Contact) error
	IngestContact(ctx context.Context, contact domain.Contact, skipExisting bool, force bool) (bool, error)
	TruncateContacts(ctx context.Context) error
}

type Contact struct {
	ID          int             `db:"id"`
	ULID        string          `db:"ulid"`
	UID         string          `db:"uid"`
	Timepoint   int             `db:"timepoint"`
	Filename    string          `db:"filename"`
	Color       toolshed.Color  `db:"color"`
	SurfaceArea sql.NullFloat64 `db:"surface_area"`
}

func (c *Contact) ToDomain(neuron *domain.Neuron, totalPatches *int, totalCellPatchSA *float64) domain.Contact {
	contact := domain.Contact{
		ID:         c.ID,
		ULID:       c.ULID,
		UID:        c.UID,
		Timepoint:  c.Timepoint,
		Filename:   c.Filename,
		Color:      c.Color,
		CellStats:  &domain.CellStats{},
		PatchStats: &domain.PatchStats{},
	}

	if c.SurfaceArea.Valid {
		contact.PatchStats.PatchSurfaceArea = &c.SurfaceArea.Float64
	}

	if totalPatches != nil {
		contact.PatchStats.TotalCount = totalPatches
	}

	if totalCellPatchSA != nil {
		contact.PatchStats.TotalCellPatchSurfaceArea = totalCellPatchSA
	}

	if neuron != nil {
		if neuron.CellStats.Volume != nil {
			contact.CellStats.Volume = neuron.CellStats.Volume
		}

		if neuron.CellStats.SurfaceArea != nil {
			contact.CellStats.SurfaceArea = neuron.CellStats.SurfaceArea
		}
	}

	return contact
}

type PostgresContactRepository struct {
	cache cache.Cache
	DB    *pgxpool.Pool
}

func NewPostgresContactRepository(db *pgxpool.Pool, c cache.Cache) *PostgresContactRepository {
	return &PostgresContactRepository{
		cache: c,
		DB:    db,
	}
}

func (r *PostgresContactRepository) GetContactByULID(ctx context.Context, id string) (domain.Contact, error) {
	query := "SELECT * FROM contacts WHERE ulid = $1"

	var contact Contact
	err := r.DB.QueryRow(ctx, query, id).Scan(&contact.ID, &contact.ULID, &contact.UID, &contact.Timepoint, &contact.Filename, &contact.Color, &contact.SurfaceArea)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Contact{}, nil
		}

		return domain.Contact{}, err
	}

	neuron, err := r.ContactNeuron(ctx, contact.UID, contact.Timepoint)
	if err != nil {
		return domain.Contact{}, err
	}

	totalPatches, err := r.CellPatchCount(ctx, contact.UID, contact.Timepoint)
	if err != nil {
		return domain.Contact{}, err
	}

	totalCellPatchSA, err := r.CellContactSurfaceArea(ctx, contact.UID, contact.Timepoint)
	if err != nil {
		return domain.Contact{}, err
	}

	return contact.ToDomain(&neuron, &totalPatches, &totalCellPatchSA), nil
}

func (r *PostgresContactRepository) GetContactByUID(ctx context.Context, uid string, timepoint int) (domain.Contact, error) {
	query := "SELECT * FROM contacts WHERE uid = $1 AND timepoint = $2"

	var contact Contact
	err := r.DB.QueryRow(ctx, query, uid, timepoint).Scan(&contact.ID, &contact.ULID, &contact.UID, &contact.Timepoint, &contact.Filename, &contact.Color, &contact.SurfaceArea)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Contact{}, nil
		}

		return domain.Contact{}, err
	}

	neuron, err := r.ContactNeuron(ctx, contact.UID, contact.Timepoint)
	if err != nil {
		return domain.Contact{}, err
	}

	totalPatches, err := r.CellPatchCount(ctx, contact.UID, contact.Timepoint)
	if err != nil {
		return domain.Contact{}, err
	}

	totalCellPatchSA, err := r.CellContactSurfaceArea(ctx, contact.UID, contact.Timepoint)
	if err != nil {
		return domain.Contact{}, err
	}

	return contact.ToDomain(&neuron, &totalPatches, &totalCellPatchSA), nil
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
	q := "SELECT * FROM contacts "

	parsedQuery, args := r.ParseContactAPIV1Request(ctx, query)

	q += parsedQuery

	rows, _ := r.DB.Query(ctx, q, args...)

	contacts, err := pgx.CollectRows(rows, pgx.RowToStructByName[Contact])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.Contact{}, nil
		}

		return nil, err
	}

	domainContacts := make([]domain.Contact, len(contacts))

	for i := range contacts {
		domainContacts[i] = contacts[i].ToDomain(nil, nil, nil)
	}

	return domainContacts, err
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

// UpdateContact takes a contact and updates the fields accordingly
func (r *PostgresContactRepository) UpdateContact(ctx context.Context, contact domain.Contact) error {
	query := `UPDATE contacts SET `
	var args []any
	args = append(args, contact.ULID)

	if contact.PatchStats != nil {
		if contact.PatchStats.PatchSurfaceArea != nil {
			args = append(args, *contact.PatchStats.PatchSurfaceArea)
			query += fmt.Sprintf("surface_area = $%d, ", len(args))
		}
	}

	if len(args) == 1 {
		return nil
	}

	query = strings.TrimSuffix(query, ", ")
	query = strings.TrimSuffix(query, ",")
	if !strings.HasSuffix(query, " ") {
		query += " "
	}

	query += `where ulid = $1`

	_, err := r.DB.Exec(ctx, query, args...)
	if err != nil {
		fmt.Printf("%v", err.Error())
		return err
	}

	return nil
}

func (r *PostgresContactRepository) ContactNeuron(ctx context.Context, uid string, timepoint int) (domain.Neuron, error) {
	// we need to get the cell by splitting the contact by the word "by"
	parts := strings.Split(uid, "by")

	if len(parts) < 1 {
		return domain.Neuron{}, errors.New("invalid uid")
	}

	cellUID := parts[0]

	if cellUID == "" {
		return domain.Neuron{}, errors.New("invalid uid")
	}

	query := "SELECT * FROM neurons WHERE uid = $1 and timepoint = $2"

	var neuron Neuron
	err := r.DB.QueryRow(ctx, query, cellUID, timepoint).Scan(&neuron.ID, &neuron.ULID, &neuron.UID, &neuron.Timepoint, &neuron.Filename, &neuron.Color, &neuron.Volume, &neuron.SurfaceArea)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Neuron{}, nil
		}
		return domain.Neuron{}, err
	}

	return neuron.ToDomain(), nil
}

func (r *PostgresContactRepository) ContactSurfaceArea(ctx context.Context, uid string, timepoint int) (float64, error) {
	cacheKey := fmt.Sprintf("neuron:contact_surface_area:%s:%d", uid, timepoint)

	if cachedPSA, found := r.cache.Get(cacheKey); found {
		if cached, ok := cachedPSA.(float64); ok {
			return cached, nil
		}
	}

	query := "SELECT sum(surface_area) FROM contacts WHERE uid like $1 AND timepoint = $2;"
	like := fmt.Sprintf("%s%%", uid)

	var total sql.NullFloat64
	err := r.DB.QueryRow(ctx, query, like, timepoint).Scan(&total)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	if total.Valid {
		count := total.Float64

		r.cache.Set(cacheKey, count)

		return count, nil
	}

	return 0, nil
}

func (r *PostgresContactRepository) CellPatchCount(ctx context.Context, contactUID string, timepoint int) (int, error) {
	parts := strings.Split(contactUID, "by")

	if len(parts) < 1 {
		return 0, errors.New("invalid contact uid")
	}

	cellUID := parts[0]

	if cellUID == "" {
		return 0, errors.New("invalid cell uid")
	}

	like := fmt.Sprintf("%sby%%", cellUID)

	query := "SELECT count(*) FROM contacts WHERE uid LIKE $1 AND timepoint = $2;"

	var count int
	err := r.DB.QueryRow(ctx, query, like, timepoint).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}

func (r *PostgresContactRepository) CellContactSurfaceArea(ctx context.Context, contactUID string, timepoint int) (float64, error) {
	parts := strings.Split(contactUID, "by")

	if len(parts) < 1 {
		return 0, errors.New("invalid contact uid")
	}

	cellUID := parts[0]

	if cellUID == "" {
		return 0, errors.New("invalid cell uid")
	}

	like := fmt.Sprintf("%sby%%", cellUID)

	query := "SELECT sum(surface_area) FROM contacts WHERE uid LIKE $1 AND timepoint = $2;"

	var total sql.NullFloat64
	err := r.DB.QueryRow(ctx, query, like, timepoint).Scan(&total)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	if total.Valid {
		count := total.Float64
		return count, nil
	}

	return 0, nil
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

func (r *PostgresContactRepository) ParseContactAPIV1Request(ctx context.Context, req domain.APIV1Request) (string, []any) {
	queryParts := []string{"where 1=1"}
	args := []any{}

	if req.Timepoint != nil {
		args = append(args, req.Timepoint)
		queryParts = append(queryParts, fmt.Sprintf("timepoint = $%d", len(args)))
	}

	if len(req.UIDs) > 0 {
		// we need to build a query where UID is like or, looping over the UIDs, wrapping them in % and adding them to the array[]
		uidArray := []string{}
		for _, uid := range req.UIDs {
			uidArray = append(uidArray, fmt.Sprintf("%%%sby%%", strings.ToLower(uid)))
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

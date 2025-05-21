package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"neuroscan/internal/domain"
	"neuroscan/internal/repository"
)

type ContactService interface {
	GetContactByULID(ctx context.Context, id string) (domain.Contact, error)
	GetContactByUID(ctx context.Context, uid string, timepoint int) (domain.Contact, error)
	ContactExists(ctx context.Context, uid string, timepoint int) (bool, error)
	SearchContacts(ctx context.Context, query domain.APIV1Request) ([]domain.Contact, error)
	CountContacts(ctx context.Context, query domain.APIV1Request) (int, error)
	CreateContact(ctx context.Context, contact domain.Contact) error
	UpdateContact(ctx context.Context, contact domain.Contact) error
	IngestContact(ctx context.Context, contact domain.Contact, skipExisting bool, force bool) (bool, error)
	TruncateContacts(ctx context.Context) error
	ParseMeta(ctx context.Context, row []string, timepoint int, dataType string) error
}

type contactService struct {
	repo repository.ContactRepository
}

func NewContactService(repo repository.ContactRepository) ContactService {
	return &contactService{
		repo: repo,
	}
}

func (s *contactService) GetContactByULID(ctx context.Context, id string) (domain.Contact, error) {
	return s.repo.GetContactByULID(ctx, id)
}

func (s *contactService) GetContactByUID(ctx context.Context, uid string, timepoint int) (domain.Contact, error) {
	return s.repo.GetContactByUID(ctx, uid, timepoint)
}

func (s *contactService) ContactExists(ctx context.Context, uid string, timepoint int) (bool, error) {
	return s.repo.ContactExists(ctx, uid, timepoint)
}

func (s *contactService) SearchContacts(ctx context.Context, query domain.APIV1Request) ([]domain.Contact, error) {
	return s.repo.SearchContacts(ctx, query)
}

func (s *contactService) CountContacts(ctx context.Context, query domain.APIV1Request) (int, error) {
	return s.repo.CountContacts(ctx, query)
}

func (s *contactService) CreateContact(ctx context.Context, contact domain.Contact) error {
	return s.repo.CreateContact(ctx, contact)
}

func (s *contactService) UpdateContact(ctx context.Context, contact domain.Contact) error {
	return s.repo.UpdateContact(ctx, contact)
}

func (s *contactService) IngestContact(ctx context.Context, contact domain.Contact, skipExisting bool, force bool) (bool, error) {
	return s.repo.IngestContact(ctx, contact, skipExisting, force)
}

func (s *contactService) TruncateContacts(ctx context.Context) error {
	return s.repo.TruncateContacts(ctx)
}

func (s *contactService) ParseMeta(ctx context.Context, row []string, timepoint int, dataType string) error {
	uid := row[0]
	val := row[1]

	contact, err := s.repo.GetContactByUID(ctx, uid, timepoint)
	if err != nil {
		errorString := fmt.Sprintf("unable to find contact from meta uid %s and timepoint %d: %s", uid, timepoint, err.Error())
		return errors.New(errorString)
	}

	contact.PatchStats = &domain.PatchStats{}

	switch dataType {
	case "surface_area":
		val = strings.TrimSpace(val)
		value, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		contact.PatchStats.PatchSurfaceArea = &value
	default:
		return errors.New("unknown data type")
	}

	return s.repo.UpdateContact(ctx, contact)
}

package service

import (
	"context"
	"errors"

	"neuroscan/internal/domain"
	"neuroscan/internal/repository"
	"neuroscan/internal/toolshed"
)

type ContactService interface {
	GetContactByID(ctx context.Context, uid string, timepoint int) (domain.Contact, error)
	GetContactByUID(ctx context.Context, uid string, timepoint int) (domain.Contact, error)
	ContactExists(ctx context.Context, uid string, timepoint int) (bool, error)
	SearchContacts(ctx context.Context, query domain.APIV1Request) ([]domain.Contact, error)
	CountContacts(ctx context.Context, query domain.APIV1Request) (int, error)
	CreateContact(ctx context.Context, uid string, filename string, timepoint int, color toolshed.Color) error
	IngestContact(ctx context.Context, contact domain.Contact, skipExisting bool, force bool) (bool, error)
	ParseContact(ctx context.Context, filePath string) (domain.Contact, error)
}

type contactService struct {
	repo repository.ContactRepository
}

func NewContactService(repo repository.ContactRepository) ContactService {
	return &contactService{
		repo: repo,
	}
}

func (s *contactService) GetContactByID(ctx context.Context, uid string, timepoint int) (domain.Contact, error) {
	return s.repo.GetContactByID(ctx, uid, timepoint)
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

func (s *contactService) CreateContact(ctx context.Context, uid string, filename string, timepoint int, color toolshed.Color) error {
	return s.repo.CreateContact(ctx, uid, filename, timepoint, color)
}

func (s *contactService) IngestContact(ctx context.Context, contact domain.Contact, skipExisting bool, force bool) (bool, error) {
	return s.repo.IngestContact(ctx, contact, skipExisting, force)
}

func (s *contactService) ParseContact(ctx context.Context, filePath string) (domain.Contact, error) {
	fileMetas, err := toolshed.FilePathParse(filePath)

	if err != nil {
		return domain.Contact{}, errors.New("error parsing contact file path: " + err.Error())
	}

	fileMeta := fileMetas[0]

	contact := domain.Contact{
		UID:       fileMeta.UID,
		Filename:  fileMeta.Filename,
		Timepoint: fileMeta.Timepoint,
		Color:     fileMeta.Color,
	}

	return contact, nil
}

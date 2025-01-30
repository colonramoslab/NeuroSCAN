package service

import (
	"context"

	"neuroscan/internal/domain"
	"neuroscan/internal/repository"
)

type ContactService interface {
	GetContactByUID(ctx context.Context, uid string, timepoint int) (domain.Contact, error)
	ContactExists(ctx context.Context, uid string, timepoint int) (bool, error)
	SearchContacts(ctx context.Context, query domain.APIV1Request) ([]domain.Contact, error)
	CountContacts(ctx context.Context, query domain.APIV1Request) (int, error)
	CreateContact(ctx context.Context, contact domain.Contact) error
	IngestContact(ctx context.Context, contact domain.Contact, skipExisting bool, force bool) (bool, error)
	TruncateContacts(ctx context.Context) error
}

type contactService struct {
	repo repository.ContactRepository
}

func NewContactService(repo repository.ContactRepository) ContactService {
	return &contactService{
		repo: repo,
	}
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

func (s *contactService) IngestContact(ctx context.Context, contact domain.Contact, skipExisting bool, force bool) (bool, error) {
	return s.repo.IngestContact(ctx, contact, skipExisting, force)
}

func (s *contactService) TruncateContacts(ctx context.Context) error {
	return s.repo.TruncateContacts(ctx)
}

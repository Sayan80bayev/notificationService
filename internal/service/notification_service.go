package service

import (
	"context"
	"errors"
	"notificationService/internal/model"
	"notificationService/internal/repository"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidID      = errors.New("invalid notification id")
	ErrInvalidUserID  = errors.New("invalid user id")
	ErrMissingTitle   = errors.New("notification title is required")
	ErrMissingMessage = errors.New("notification message is required")
)

type NotificationService interface {
	CreateNotification(ctx context.Context, n *model.Notification) (*model.Notification, error)
	GetNotificationByID(ctx context.Context, id uuid.UUID) (*model.Notification, error)
	GetNotificationsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Notification, error)
	MarkNotificationAsRead(ctx context.Context, id uuid.UUID) error
	DeleteNotification(ctx context.Context, id uuid.UUID) error
}

type notificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) NotificationService {
	return &notificationService{repo: repo}
}

// CreateNotification adds a new notification securely
func (s *notificationService) CreateNotification(ctx context.Context, n *model.Notification) (*model.Notification, error) {
	if n == nil {
		return nil, errors.New("notification cannot be nil")
	}
	if n.UserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	if n.Title == "" {
		return nil, ErrMissingTitle
	}
	if n.Message == "" {
		return nil, ErrMissingMessage
	}
	if n.Type == "" {
		n.Type = "system" // default type
	}
	n.ID = uuid.New()
	n.CreatedAt = time.Now().UTC()
	n.IsRead = false
	n.ReadAt = nil

	if err := s.repo.Create(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

// GetNotificationByID fetches a single notification
func (s *notificationService) GetNotificationByID(ctx context.Context, id uuid.UUID) (*model.Notification, error) {
	if id == uuid.Nil {
		return nil, ErrInvalidID
	}
	return s.repo.FindByID(ctx, id)
}

// GetNotificationsByUser fetches notifications with pagination
func (s *notificationService) GetNotificationsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Notification, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.FindByUserID(ctx, userID, limit, offset)
}

// MarkNotificationAsRead marks a notification as read
func (s *notificationService) MarkNotificationAsRead(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrInvalidID
	}
	return s.repo.MarkAsRead(ctx, id, time.Now().UTC())
}

// DeleteNotification removes a notification
func (s *notificationService) DeleteNotification(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrInvalidID
	}
	return s.repo.Delete(ctx, id)
}

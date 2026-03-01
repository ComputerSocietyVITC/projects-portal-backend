package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrInviteNotFound      = errors.New("invite not found")
	ErrInviteExpired       = errors.New("invite has expired")
	ErrInviteUsed          = errors.New("invite has already been used")
	ErrInviteEmailMismatch = errors.New("email does not match invite")
)

type InviteService interface {
	CreateInvite(ctx context.Context, email, role string, createdBy uuid.UUID) (*models.Invite, error)
	ValidateAndUseInvite(ctx context.Context, token, email string) error
}

type inviteService struct {
	inviteRepo repository.InviteRepository
}

func NewInviteService(inviteRepo repository.InviteRepository) InviteService {
	return &inviteService{inviteRepo: inviteRepo}
}

func (s *inviteService) CreateInvite(ctx context.Context, email, role string, createdBy uuid.UUID) (*models.Invite, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	invite := &models.Invite{
		Email:     email,
		Role:      role,
		Token:     token,
		ExpiresAt: time.Now().Add(72 * time.Hour).Unix(),
		CreatedBy: createdBy,
	}

	if err := s.inviteRepo.Create(ctx, invite); err != nil {
		return nil, err
	}

	return invite, nil
}

func (s *inviteService) ValidateAndUseInvite(ctx context.Context, token, email string) error {
	invite, err := s.inviteRepo.GetByToken(ctx, token)
	if err != nil {
		return ErrInviteNotFound
	}

	if invite.Used {
		return ErrInviteUsed
	}

	if time.Now().Unix() > invite.ExpiresAt {
		return ErrInviteExpired
	}

	if invite.Email != email {
		return ErrInviteEmailMismatch
	}

	return s.inviteRepo.MarkAsUsed(ctx, invite.ID)
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

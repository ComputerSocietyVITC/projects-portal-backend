package service

import (
	"errors"
	"time"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type InviteService struct {
	Repo *repository.InviteRepo
}

func NewInviteService(r *repository.InviteRepo) *InviteService {
	return &InviteService{Repo: r}
}

func (svc *InviteService) ValidateInviteToken(email, token string) (string, error) {
	invite, err := svc.Repo.GetInviteByToken(token)
	if err != nil {
		return "", err
	}

	if time.Now().Unix() > invite.ExpiresAt {
		return "", errors.New("invite expired")
	}

	if invite.Used {
		return "", errors.New("invite already used")
	}

	if invite.Email != email {
		return "", errors.New("invite email mismatch")
	}

	return invite.Role, nil
}

func (svc *InviteService) UseInviteToken(token string) error {
	invite, err := svc.Repo.GetInviteByToken(token)
	if err != nil {
		return err
	}

	invite.Used = true
	if err := svc.Repo.UpdateInvite(invite); err != nil {
		return err
	}

	return nil
}

func (svc *InviteService) CreateInvite(email, role string, adminUUID uuid.UUID) (*models.Invite, error) {
	token := uuid.NewString()

	invite := &models.Invite{
		Email:     email,
		Role:      role,
		Token:     token,
		ExpiresAt: time.Now().Add(48 * time.Hour).Unix(),
		CreatedBy: adminUUID,
	}

	if err := svc.Repo.CreateInvite(invite); err != nil {
		return nil, err
	}

	svc.Repo.Logger.Info("Created new invite", zap.String("email", email), zap.String("role", role), zap.String("token", token), zap.String("created_by", adminUUID.String()))

	return invite, nil
}

func (svc *InviteService) DeleteInvite(inviteUUID uuid.UUID) error {
	return svc.Repo.DeleteInvite(inviteUUID)
}

func (svc *InviteService) GetInvites() ([]*models.Invite, error) {
	return svc.Repo.GetInvites()
}

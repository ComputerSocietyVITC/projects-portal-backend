package repository

import (
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type InviteRepo struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (r *InviteRepo) CreateInvite(invite *models.Invite) error {
	return r.DB.Create(invite).Error
}

func (r *InviteRepo) DeleteInvite(inviteUUID uuid.UUID) error {
	return r.DB.Delete(&models.Invite{ID: inviteUUID}).Error
}

func (r *InviteRepo) GetInviteByID(inviteUUID uuid.UUID) (*models.Invite, error) {
	var invite = models.Invite{ID: inviteUUID}

	if err := r.DB.First(&invite).Error; err != nil {
		return nil, err
	}

	return &invite, nil
}

func (r *InviteRepo) GetInviteByToken(inviteToken string) (*models.Invite, error) {
	var invite = models.Invite{Token: inviteToken}

	if err := r.DB.First(&invite).Error; err != nil {
		return nil, err
	}

	return &invite, nil
}

func (r *InviteRepo) GetInvites() ([]*models.Invite, error) {
	var invites []*models.Invite
	if err := r.DB.Find(&invites).Error; err != nil {
		return nil, err
	}

	return invites, nil
}

func (r *InviteRepo) UpdateInvite(invite *models.Invite) error {
	return r.DB.Save(invite).Error
}

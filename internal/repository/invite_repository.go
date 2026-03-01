package repository

import (
	"context"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InviteRepository interface {
	GetByToken(ctx context.Context, token string) (*models.Invite, error)
	Create(ctx context.Context, invite *models.Invite) error
	MarkAsUsed(ctx context.Context, id uuid.UUID) error
}

type inviteRepository struct {
	db *gorm.DB
}

func NewInviteRepository(db *gorm.DB) InviteRepository {
	return &inviteRepository{db: db}
}

func (r *inviteRepository) GetByToken(ctx context.Context, token string) (*models.Invite, error) {
	var invite models.Invite
	if err := r.db.WithContext(ctx).First(&invite, "token = ?", token).Error; err != nil {
		return nil, err
	}
	return &invite, nil
}

func (r *inviteRepository) Create(ctx context.Context, invite *models.Invite) error {
	return r.db.WithContext(ctx).Create(invite).Error
}

func (r *inviteRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Invite{}).Where("id = ?", id).Update("used", true).Error
}

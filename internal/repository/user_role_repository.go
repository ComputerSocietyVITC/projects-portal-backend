package repository

import (
	"context"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRoleRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.UserRole, error)
	Create(ctx context.Context, userRole *models.UserRole) error
	Delete(ctx context.Context, userID, roleID uuid.UUID) error
}

type userRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) UserRoleRepository {
	return &userRoleRepository{db: db}
}

func (r *userRoleRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.UserRole, error) {
	var userRoles []models.UserRole
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return nil, err
	}
	return userRoles, nil
}

func (r *userRoleRepository) Create(ctx context.Context, userRole *models.UserRole) error {
	return r.db.WithContext(ctx).Create(userRole).Error
}

func (r *userRoleRepository) Delete(ctx context.Context, userID, roleID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.UserRole{}).Error
}

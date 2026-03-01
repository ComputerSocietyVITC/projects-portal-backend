package repository

import (
	"context"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectMemberRepository interface {
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.ProjectMember, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.ProjectMember, error)
	Create(ctx context.Context, member *models.ProjectMember) error
	Delete(ctx context.Context, projectID, userID uuid.UUID) error
	Exists(ctx context.Context, projectID, userID uuid.UUID) (bool, error)
}

type projectMemberRepository struct {
	db *gorm.DB
}

func NewProjectMemberRepository(db *gorm.DB) ProjectMemberRepository {
	return &projectMemberRepository{db: db}
}

func (r *projectMemberRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.ProjectMember, error) {
	var members []models.ProjectMember
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *projectMemberRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.ProjectMember, error) {
	var members []models.ProjectMember
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *projectMemberRepository) Create(ctx context.Context, member *models.ProjectMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *projectMemberRepository) Delete(ctx context.Context, projectID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("project_id = ? AND user_id = ?", projectID, userID).Delete(&models.ProjectMember{}).Error
}

func (r *projectMemberRepository) Exists(ctx context.Context, projectID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error
	return count > 0, err
}

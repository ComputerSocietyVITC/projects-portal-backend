package repository

import (
	"context"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	GetAllProjects(ctx context.Context) ([]models.Project, error)
	GetProjectByID(ctx context.Context, id uint) (*models.Project, error)
	CreateProject(ctx context.Context, project *models.Project) (*models.Project, error)
	UpdateProject(ctx context.Context, id uint, updates map[string]interface{}) (*models.Project, error)
	DeleteProject(ctx context.Context, project *models.Project) error
}

type projectRepository struct {
	DB *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{DB: db}
}

func (r *projectRepository) GetAllProjects(ctx context.Context) ([]models.Project, error) {
	var projects []models.Project
	if err := r.DB.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) GetProjectByID(ctx context.Context, id uint) (*models.Project, error) {
	var project models.Project
	if err := r.DB.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) CreateProject(ctx context.Context, project *models.Project) (*models.Project, error) {
	if err := r.DB.Create(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}

func (r *projectRepository) UpdateProject(ctx context.Context, id uint, updates map[string]interface{}) (*models.Project, error) {
	var project models.Project
	if err := r.DB.First(&project, id).Error; err != nil {
		return nil, err
	}
	if err := r.DB.Model(&project).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) DeleteProject(ctx context.Context, project *models.Project) error {
	if err := r.DB.Delete(project.ID).Error; err != nil {
		return err
	}
	return nil
}

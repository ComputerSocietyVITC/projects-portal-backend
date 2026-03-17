package service

import (
	"context"
	"errors"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrInvalidID       = errors.New("invalid input")
)

type ProjectService interface {
	GetAllProjects(ctx context.Context) ([]models.Project, error)
	GetProjectByID(ctx context.Context, id uint) (*models.Project, error)
	CreateProject(ctx context.Context, project *models.Project) (*models.Project, error)
	UpdateProject(ctx context.Context, id uint, updates map[string]interface{}) (*models.Project, error)
	DeleteProject(ctx context.Context, id uint) error
}

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &projectService{repo: repo}
}

func (s *projectService) GetAllProjects(ctx context.Context) ([]models.Project, error) {
	return s.repo.GetAllProjects(ctx)
}

func (s *projectService) GetProjectByID(ctx context.Context, id uint) (*models.Project, error) {
	project, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	return project, nil
}

func (s *projectService) CreateProject(ctx context.Context, project *models.Project) (*models.Project, error) {
	if project.Title == "" {
		return nil, ErrInvalidID
	}
	return s.repo.CreateProject(ctx, project)
}

func (s *projectService) UpdateProject(ctx context.Context, id uint, updates map[string]interface{}) (*models.Project, error) {
	project, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	if title, ok := updates["title"].(string); ok {
		project.Title = title
	}
	if description, ok := updates["description"].(string); ok {
		project.Description = description
	}
	return s.repo.UpdateProject(ctx, id, updates)
}

func (s *projectService) DeleteProject(ctx context.Context, id uint) error {
	project, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return err
	}
	if err := s.repo.DeleteProject(ctx, project); err != nil {
		return err
	}
	return nil
}

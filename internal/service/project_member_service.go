package service

import (
	"context"
	"errors"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrMemberAlreadyExists = errors.New("user is already a member of this project")
	ErrMemberNotFound      = errors.New("user is not a member of this project")
)

type ProjectMemberService interface {
	AddUserToProject(ctx context.Context, projectID, userID uuid.UUID, role string) error
	RemoveUserFromProject(ctx context.Context, projectID, userID uuid.UUID) error
	GetProjectMembers(ctx context.Context, projectID uuid.UUID) ([]models.ProjectMember, error)
}

type projectMemberService struct {
	repo repository.ProjectMemberRepository
}

func NewProjectMemberService(repo repository.ProjectMemberRepository) ProjectMemberService {
	return &projectMemberService{repo: repo}
}

func (s *projectMemberService) AddUserToProject(ctx context.Context, projectID, userID uuid.UUID, role string) error {
	exists, err := s.repo.Exists(ctx, projectID, userID)
	if err != nil {
		return err
	}
	if exists {
		return ErrMemberAlreadyExists
	}

	member := &models.ProjectMember{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
	}

	return s.repo.Create(ctx, member)
}

func (s *projectMemberService) RemoveUserFromProject(ctx context.Context, projectID, userID uuid.UUID) error {
	exists, err := s.repo.Exists(ctx, projectID, userID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrMemberNotFound
	}

	return s.repo.Delete(ctx, projectID, userID)
}

func (s *projectMemberService) GetProjectMembers(ctx context.Context, projectID uuid.UUID) ([]models.ProjectMember, error) {
	return s.repo.GetByProjectID(ctx, projectID)
}

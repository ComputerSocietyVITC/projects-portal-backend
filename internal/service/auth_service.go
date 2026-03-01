package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, email, password, name, token string) (*models.User, error)
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type authService struct {
	userRepo      repository.UserRepository
	userRoleRepo  repository.UserRoleRepository
	roleRepo      repository.RoleRepository
	inviteService InviteService
}

func NewAuthService(userRepo repository.UserRepository, userRoleRepo repository.UserRoleRepository, roleRepo repository.RoleRepository, inviteService InviteService) AuthService {
	return &authService{
		userRepo:      userRepo,
		userRoleRepo:  userRoleRepo,
		roleRepo:      roleRepo,
		inviteService: inviteService,
	}
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	roles, err := s.GetUserRoles(ctx, user.ID)
	if err != nil {
		return "", err
	}

	token, err := generateJWT(user.ID, user.Email, roles)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) Register(ctx context.Context, email, password, name, token string) (*models.User, error) {
	if err := s.inviteService.ValidateAndUseInvite(ctx, token, email); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Name:         name,
		Status:       "active",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error) {
	userRoles, err := s.userRoleRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var roleNames []string
	for _, ur := range userRoles {
		role, err := s.roleRepo.GetByID(ctx, ur.RoleID)
		if err != nil {
			continue
		}
		roleNames = append(roleNames, role.Name)
	}

	return roleNames, nil
}

func generateJWT(userID uuid.UUID, email string, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"roles":   roles,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}

	return token.SignedString([]byte(secret))
}

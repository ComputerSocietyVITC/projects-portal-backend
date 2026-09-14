package service

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo          *repository.AuthRepo
	RedisRepo     *repository.RedisRepo
	InviteService *InviteService
}

func NewAuthService(r *repository.AuthRepo, redis *repository.RedisRepo, inviteSvc *InviteService) *AuthService {
	return &AuthService{Repo: r, RedisRepo: redis, InviteService: inviteSvc}
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET")), nil
}

func (svc *AuthService) CreateRefreshToken(userid uuid.UUID) (string, error) {
	var refreshToken = models.UserRefreshToken{UserID: userid.String()}

	claims := &models.JwtUserRefreshToken{
		UserID: userid.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTok, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	refreshToken.RefreshToken = refreshTok

	if err := svc.RedisRepo.SaveRefreshToken(&refreshToken); err != nil {
		return "", err
	}
	return refreshTok, nil
}

func (svc *AuthService) GetUserFromRefreshToken(refreshToken string) (*models.User, error) {
	parsedRefreshToken, _ := jwt.ParseWithClaims(refreshToken, &models.JwtUserRefreshToken{}, jwtKeyFunc)

	claim, ok := parsedRefreshToken.Claims.(*models.JwtUserRefreshToken)
	if !ok || !parsedRefreshToken.Valid {
		return nil, jwt.ErrTokenExpired
	}

	if modelRefreshToken, err := svc.RedisRepo.GetRefreshToken(claim.UserID); err != nil || modelRefreshToken.RefreshToken != refreshToken {
		return nil, errors.New("invalid refresh token")
	}

	userUUID := uuid.MustParse(claim.UserID)
	user, err := svc.Repo.GetUserByID(userUUID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *AuthService) CreateAccessToken(user *models.User) (string, error) {
	var userid = user.ID

	role, err := svc.Repo.GetUserRole(userid)
	if err != nil {
		return "", err
	}

	claims := &models.JwtUserAccessToken{
		UserID: userid.String(),
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return accessToken, nil
}

func (svc *AuthService) Register(email, password, name, token string) (*models.User, error) {
	roleName, err := svc.InviteService.ValidateInviteToken(email, token)
	if err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Email:        email,
		PasswordHash: string(passwordHash),
		Name:         name,
	}

	if err := svc.Repo.CreateUser(user); err != nil {
		return nil, err
	}

	if err := svc.Repo.CreateRoleReference(user.ID, roleName); err != nil {
		return nil, err
	}

	if err := svc.InviteService.UseInviteToken(token); err != nil {
		return nil, err
	}

	svc.Repo.Logger.Info("Registered new user", zap.String("email", email), zap.String("name", name), zap.String("role", roleName), zap.String("token", token))
	return user, nil
}

func (svc *AuthService) Login(email, password string) (*models.User, error) {
	user, err := svc.Repo.GetUserByEmail(email)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(strings.TrimSpace(user.PasswordHash)), []byte(strings.TrimSpace(password)))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *AuthService) Logout(userID string) error {
	return svc.RedisRepo.DeleteRefreshToken(userID)
}

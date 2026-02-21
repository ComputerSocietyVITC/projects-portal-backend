package repository

import (
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthRepo struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (r *AuthRepo) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *AuthRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepo) GetUserByID(userUUID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("id = ?", userUUID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepo) GetUserRole(userUUID uuid.UUID) (string, error) {
	var roleID = models.UserRole{UserID: userUUID}
	if err := r.DB.Where("user_id = ?", userUUID).First(&roleID).Error; err != nil {
		return "", err
	}

	var role = models.Role{ID: roleID.RoleID}
	if err := r.DB.Where("id = ?", roleID.RoleID).First(&role).Error; err != nil {
		return "", err
	}
	r.Logger.Log(zap.InfoLevel, "got role for user", zap.String("user_uuid", userUUID.String()), zap.String("role_name", role.Name), zap.String("role_id", role.ID.String()))
	return role.Name, nil
}

func (r *AuthRepo) CreateRoleReference(userUUID uuid.UUID, roleName string) error {
	var role models.Role
	if err := r.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}

	var userRole = models.UserRole{UserID: userUUID, RoleID: role.ID}
	return r.DB.Create(userRole).Error
}

package main

import (
	"log"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/config"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("=== Seeding RBAC Data ===")

	// Create roles
	roles := []models.Role{
		{Name: "group_head"},
		{Name: "member"},
		{Name: "viewer"},
	}

	for _, role := range roles {
		var existing models.Role
		if err := db.Where("name = ?", role.Name).First(&existing).Error; err != nil {
			if err := db.Create(&role).Error; err != nil {
				log.Printf("Failed to create role %s: %v", role.Name, err)
			} else {
				log.Printf("✓ Created role: %s (ID: %s)", role.Name, role.ID)
			}
		} else {
			log.Printf("⊘ Role already exists: %s", role.Name)
			role.ID = existing.ID
		}
	}

	// Create test users with proper password hashing
	hashedPassword1, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashedPassword2, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	users := []models.User{
		{
			Email:        "grouphead@example.com",
			PasswordHash: string(hashedPassword1),
			Name:         "Group Head User",
			Status:       "active",
		},
		{
			Email:        "member@example.com",
			PasswordHash: string(hashedPassword2),
			Name:         "Regular Member",
			Status:       "active",
		},
	}

	for i, user := range users {
		var existing models.User
		if err := db.Where("email = ?", user.Email).First(&existing).Error; err != nil {
			if err := db.Create(&user).Error; err != nil {
				log.Printf("Failed to create user %s: %v", user.Email, err)
			} else {
				log.Printf("✓ Created user: %s (ID: %s)", user.Name, user.ID)
				users[i].ID = user.ID
			}
		} else {
			log.Printf("⊘ User already exists: %s", user.Email)
			users[i].ID = existing.ID
		}
	}

	// Assign roles to users
	// Group head gets group_head role
	var groupHeadRole models.Role
	db.Where("name = ?", "group_head").First(&groupHeadRole)

	userRole1 := models.UserRole{
		UserID: users[0].ID,
		RoleID: groupHeadRole.ID,
	}

	var existingUR1 models.UserRole
	if err := db.Where("user_id = ? AND role_id = ?", userRole1.UserID, userRole1.RoleID).First(&existingUR1).Error; err != nil {
		if err := db.Create(&userRole1).Error; err != nil {
			log.Printf("Failed to assign group_head role: %v", err)
		} else {
			log.Printf("✓ Assigned group_head role to %s", users[0].Email)
		}
	} else {
		log.Printf("⊘ Role already assigned to %s", users[0].Email)
	}

	// Regular member gets member role
	var memberRole models.Role
	db.Where("name = ?", "member").First(&memberRole)

	userRole2 := models.UserRole{
		UserID: users[1].ID,
		RoleID: memberRole.ID,
	}

	var existingUR2 models.UserRole
	if err := db.Where("user_id = ? AND role_id = ?", userRole2.UserID, userRole2.RoleID).First(&existingUR2).Error; err != nil {
		if err := db.Create(&userRole2).Error; err != nil {
			log.Printf("Failed to assign member role: %v", err)
		} else {
			log.Printf("✓ Assigned member role to %s", users[1].Email)
		}
	} else {
		log.Printf("⊘ Role already assigned to %s", users[1].Email)
	}

	log.Println("\n=== RBAC Seeding Complete ===")
	log.Println("\nTest Credentials:")
	log.Println("Group Head: grouphead@example.com / password123")
	log.Println("Member: member@example.com / password123")
}

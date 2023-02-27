package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/meshachdamilare/auth-with-gorm-postgres/config"
	"log"
	"time"
)

func init() {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}
	config.ConnectDB(&conf)
	config.DB.AutoMigrate(User{})
	fmt.Println(" Migration complete.")
}

type User struct {
	ID                 uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name               string    `gorm:"type:varchar(255);not null"`
	Email              string    `gorm:"uniqueIndex;not null"`
	Password           string    `gorm:"not null"`
	Role               string    `gorm:"type:varchar(255);not null"`
	Provider           string    `gorm:"not null"`
	Photo              string    `gorm:"not null"`
	VerificationCode   string
	PasswordResetToken string
	Verified           bool `gorm:"not null"`
	PasswordResetAt    time.Time
	CreatedAt          time.Time `gorm:"not null"`
	UpdatedAt          time.Time `gorm:"not null"`
}

type SignUp struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
	Photo           string `json:"photo" binding:"required"`
}

type SignIn struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	Photo     string    `json:"photo,omitempty"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordInput struct {
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

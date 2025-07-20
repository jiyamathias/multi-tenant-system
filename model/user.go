package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	// Password string representation of the user password value
	Password string

	// User schema
	User struct {
		ID        uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
		TenantID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenantId"`
		Tenant    *Tenant        `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE;" json:"-"`
		FirstName string         `gorm:"size:50;not null" json:"firstName"`
		LastName  string         `gorm:"size:50;not null" json:"lastName"`
		Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
		Password  Password       `gorm:"not null" json:"-"`
		IsActive  bool           `gorm:"default:true" json:"isActive"`
		CreatedAt time.Time      `json:"createdAt"`
		UpdatedAt time.Time      `json:"updatedAt"`
		DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	}

	// PublicUser schema
	PublicUser struct {
		Email       string `json:"email"`
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastName"`
		PhoneNumber string `json:"phoneNumber"`
	}

	// AuthResponse schema
	AuthResponse struct {
		Token              *string `json:"token"`
		Refresh            *string `json:"refresh"`
		User               *User   `json:"user"`
		AccessTokenExpiry  *string `json:"accessTokenExpiry"`
		RefreshTokenExpiry *string `json:"refreshTokenExpiry"`
	}
)

const (
// EmailSignUp for email signups

)

// String representation of a password either as encrypted or not
func (p Password) String() string {
	return string(p)
}

// Encrypt securely encrypts a Password object
func (p Password) Encrypt() Password {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(p.String()), 14)
	return Password(string(bytes))
}

// Check compared password and returns if they match
func (p Password) Check(password Password) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p), []byte(password))
	return err == nil
}

// PublicUser method helps us not to expose sensitive user datas
func (u *User) PublicUser() *User {
	return &User{
		Email:     u.Email,
		FirstName: u.FirstName,
	}
}

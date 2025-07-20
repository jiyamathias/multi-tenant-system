package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Tenant struct {
		ID           uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
		BusinessName string         `gorm:"size:100;not null" json:"businessName"`
		Email        string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
		CreatedAt    time.Time      `json:"createdAt"`
		UpdatedAt    time.Time      `json:"updatedAt"`
		DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
		Users        []User         `gorm:"foreignKey:TenantID" json:"-"`
	}
)

package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	// Balance schema
	Balance struct {
		ID              uuid.UUID       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
		UserID          uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
		TransactionType TransactionType `gorm:"type:varchar(50);not null" json:"transaction_type"`
		TransactionID   uuid.UUID       `gorm:"type:uuid;not null" json:"transaction_id"`
		BalanceBefore   float64         `gorm:"not null" json:"balance_before"`
		BalanceAfter    float64         `gorm:"not null" json:"balance_after"`
		CreatedAt       time.Time       `gorm:"default:now()" json:"created_at"`
		UpdatedAt       *time.Time      `json:"updated_at,omitempty"`
		DeletedAt       gorm.DeletedAt  `gorm:"index" json:"-"`
	}
)

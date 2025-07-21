package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	// Wallet schema
	Wallet struct {
		ID              uuid.UUID        `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
		UserID          uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
		User            *User            `gorm:"foreignKey:UserID;references:ID"`
		TransactionID   *uuid.UUID       `json:"transaction_id"`
		Transaction     *Transaction     `gorm:"foreignKey:TransactionID;references:ID"`
		TransactionType *TransactionType `gorm:"type:varchar(50)" json:"transaction_type"`
		BalanceID       *uuid.UUID       `gorm:"type:uuid" json:"balance_id"`
		BalanceBefore   float64          `gorm:"not null" json:"balance_before"`
		BalanceAfter    float64          `gorm:"not null" json:"balance_after"`
		CreatedAt       time.Time        `gorm:"default:now()" json:"created_at"`
		UpdatedAt       *time.Time       `json:"updated_at"`
		DeletedAt       gorm.DeletedAt   `gorm:"index" json:"-"`
	}
)

const (
	// DebitTransaction represents a debit transaction
	DebitTransaction TransactionType = "debit"
	// CreditTransaction represents a credit transaction
	CreditTransaction TransactionType = "credit"
)

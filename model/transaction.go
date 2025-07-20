package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	// TransactionTypeDebit represents a debit transaction
	TransactionTypeDebit TransactionType = "debit"

	// TransactionTypeCredit represents a credit transaction
	TransactionTypeCredit TransactionType = "credit"

	// TransactionStatusSuccessful represent represents a successful transaction
	TransactionStatusSuccessful TransactionStatus = "successful"
	//  TransactionStatusPending represents a pending transaction
	TransactionStatusPending TransactionStatus = "pending"
	// TransactionStatusFailed represents a failed transaction
	TransactionStatusFailed TransactionStatus = "failed"
	// TransactionStatusCanceled represents a canceled transaction
	TransactionStatusCanceled TransactionStatus = "canceled"
	// TransactionStatusRefunded represents a refunded transaction
	TransactionStatusRefunded TransactionStatus = "refunded"

	// TransactionFlowRevenue represents a revenue transaction
	TransactionFlowRevenue TransactionFlow = "revenue"
	// TransactionFlowWithdrawal represents a withdrawal transaction
	TransactionFlowWithdrawal TransactionFlow = "withdrawal"
)

type (
	// TransactionType credit or debit transaction
	TransactionType string

	// TransactionStatus  of type string
	TransactionStatus string

	// TransactionFlow string
	TransactionFlow string

	// Transaction schema
	Transaction struct {
		ID              uuid.UUID         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
		UserID          uuid.UUID         `gorm:"type:uuid;not null;index" json:"user_id" validate:"required"`
		User            *User             `gorm:"foreignKey:UserID;references:ID"`
		Amount          float64           `gorm:"not null" json:"amount" validate:"required"`
		Charges         float64           `json:"charges"`
		MetaData        *map[string]any   `gorm:"type:jsonb" json:"meta_data"`
		Currency        string            `json:"currency"`
		TransactionType TransactionType   `gorm:"type:varchar(50);not null" json:"transaction_type"`
		Status          TransactionStatus `gorm:"type:varchar(50);not null" json:"status"`
		TransactionFlow TransactionFlow   `gorm:"type:varchar(50)" json:"transaction_flow"`
		CreatedAt       time.Time         `gorm:"default:now()" json:"created_at"`
		UpdatedAt       *time.Time        `json:"updated_at,omitempty"`
		DeletedAt       gorm.DeletedAt    `gorm:"index" json:"-"`
	}
)

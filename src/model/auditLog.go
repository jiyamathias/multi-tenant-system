package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	// AuditLogAction type
	AuditLogAction string

	// Actor type
	Actor string

	// AuditLog schema
	AuditLog struct {
		ID            uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
		TenantID      *uuid.UUID     `json:"tenantId"`
		UserID        *uuid.UUID     `json:"userId"`
		TransactionID *uuid.UUID     `json:"transaction_id"`
		ActionDone    AuditLogAction `gorm:"type:varchar(50);not null" json:"action_done"`
		Actor         Actor          `gorm:"type:varchar(100);not null" json:"actor"`
		Messages      string         `gorm:"type:text" json:"messages"`
		CreatedAt     time.Time      `gorm:"default:now()" json:"created_at"`
		UpdatedAt     *time.Time     `gorm:"default:null" json:"updated_at"`
		DeletedAt     gorm.DeletedAt `gorm:"ndex" json:"-"`
	}
)

const (
	// ActorTenant when tenant makes the action
	ActorTenant Actor = "tenant"
	// ActorUser when user makes the action
	ActorUser Actor = "user"

	// ActionCreated is the action when the transaction is created
	ActionCreated AuditLogAction = "created"
	// ActionSuccess is the action when the transaction is successful
	ActionSuccess AuditLogAction = "success"
	// ActionPending is the action when the transaction is pending
	ActionPending AuditLogAction = "pending"
	//  ActionFailed is the action when the transaction is failed after payment
	ActionFailed AuditLogAction = "failed"
	// ActionInDispute is the action when the transaction is being disputed
	ActionInDispute AuditLogAction = "in_dispute"
	// ActionResolved is the action when the transaction dispute is resolved
	ActionResolved AuditLogAction = "resolved"
)

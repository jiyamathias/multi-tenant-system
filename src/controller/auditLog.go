package controller

import (
	"context"

	"github.com/google/uuid"

	"codematic/model"
	"codematic/model/pagination"
	"codematic/pkg/helper"
)

// CreateAuditLog add a new invoice into the AuditLog table
func (c *Controller) CreateAuditLog(ctx context.Context, AuditLog model.AuditLog) (model.AuditLog, error) {
	return c.auditLogStorage.CreateAuditLog(ctx, AuditLog)
}

// GetAllAuditLogsByTransactionID returns all the logs of actions on a transaction
func (c *Controller) GetAllAuditLogsByTransactionID(ctx context.Context, txID uuid.UUID, page pagination.Page) ([]*model.AuditLog, pagination.PageInfo, error) {
	if txID == uuid.Nil {
		return nil, pagination.PageInfo{}, helper.ErrIDMissing
	}

	return c.auditLogStorage.GetAllAuditLogsByTransactionID(ctx, txID, page)
}

func (c *Controller) GetAuditLogByID(ctx context.Context, id uuid.UUID) (model.AuditLog, error) {
	if id == uuid.Nil {
		return model.AuditLog{}, helper.ErrIDMissing
	}

	return c.auditLogStorage.GetAuditLogByID(ctx, id)
}

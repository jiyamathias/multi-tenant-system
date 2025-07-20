package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"codematic/model"
	"codematic/model/pagination"
	"codematic/pkg/helper"
)

// AuditLogDatabase shows every methods available under the audit log database
type AuditLogDatabase interface {
	CreateAuditLog(ctx context.Context, AuditLog model.AuditLog) (model.AuditLog, error)
	GetAllAuditLogsByTransactionID(ctx context.Context, txID uuid.UUID, page pagination.Page) ([]*model.AuditLog, pagination.PageInfo, error)
	GetAuditLogByID(ctx context.Context, id uuid.UUID) (model.AuditLog, error)
}

// AuditLog object
type AuditLog struct {
	logger  zerolog.Logger
	storage *Storage
}

// NewAuditLog creates a new instance of the audit log
func NewAuditLog(s *Storage) *AuditLogDatabase {
	s.Logger.With().Str(helper.LogStrKeyLevel, "auditLog").Logger()

	a := &AuditLog{
		storage: s,
	}

	auditLog := AuditLogDatabase(a)
	return &auditLog
}

// CreateAuditLog add a new invoice into the AuditLog table
func (a *AuditLog) CreateAuditLog(ctx context.Context, auditLog model.AuditLog) (model.AuditLog, error) {
	db := a.storage.DB.WithContext(ctx).Model(&model.AuditLog{}).Create(&auditLog)
	if db.Error != nil {
		a.storage.Logger.Err(db.Error).Msgf("AuditLog:: AuditLog creation error: %v, (%v)", ErrRecordCreatingFailed, db.Error)
		return model.AuditLog{}, ErrRecordCreatingFailed
	}

	return auditLog, nil
}

// GetAllAuditLogsByTransactionID returns all the logs of actions on a transaction
func (a *AuditLog) GetAllAuditLogsByTransactionID(ctx context.Context, txID uuid.UUID, page pagination.Page) ([]*model.AuditLog, pagination.PageInfo, error) {
	var auditLogs []*model.AuditLog

	offset := 0
	// load defaults
	if page.Number == nil {
		tmpPageNumber := pagination.PageDefaultNumber
		page.Number = &tmpPageNumber
	}
	if page.Size == nil {
		tmpPageSize := pagination.PageDefaultSize
		page.Size = &tmpPageSize
	}
	if page.SortBy == nil {
		tmpPageSortBy := pagination.PageDefaultSortBy
		page.SortBy = &tmpPageSortBy
	}
	if page.SortDirectionDesc == nil {
		tmpPageSortDirectionDesc := pagination.PageDefaultSortDirectionDesc
		page.SortDirectionDesc = &tmpPageSortDirectionDesc
	}

	if *page.Number > 1 {
		offset = *page.Size * (*page.Number - 1)
	}
	sortDirection := pagination.PageSortDirectionDescending
	if !*page.SortDirectionDesc {
		sortDirection = pagination.PageSortDirectionAscending
	}

	queryDraft := a.storage.DB.WithContext(ctx).Where("transaction_id = ?", txID)

	var count int64
	queryDraft.Model(model.AuditLog{}).Count(&count)

	db := queryDraft.Debug().Offset(offset).Limit(*page.Size).
		Order(fmt.Sprintf("%s %s", *page.SortBy, sortDirection)).
		Find(&auditLogs)

	if db.Error != nil {
		a.logger.Err(db.Error).Msgf("GetAllAuditLogsByTransactionID:: error: %v, (%v)", ErrEmptyResult, db.Error)
		return nil, pagination.PageInfo{}, ErrEmptyResult
	}

	return auditLogs, pagination.PageInfo{
		Page:            *page.Number,
		Size:            *page.Size,
		HasNextPage:     int64(offset+*page.Size) < count,
		HasPreviousPage: *page.Number > 1,
		TotalCount:      count,
	}, nil
}

// GetAuditLogByID returns the details of a log bases on the id
func (a *AuditLog) GetAuditLogByID(ctx context.Context, id uuid.UUID) (model.AuditLog, error) {
	var auditLog model.AuditLog

	db := a.storage.DB.WithContext(ctx).Where("id = ?", id).First(&auditLog).Order("created_at desc")
	if db.Error != nil || strings.EqualFold(auditLog.ID.String(), helper.ZeroUUID) {
		a.storage.Logger.Err(db.Error).Msgf("GetAuditLogByID error: %v (%v)", ErrRecordNotFound, db.Error)
		return auditLog, ErrRecordNotFound
	}

	return auditLog, nil
}

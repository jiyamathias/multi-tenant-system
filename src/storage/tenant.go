package storage

import (
	"context"
	"strings"

	"codematic/model"
	"codematic/pkg/helper"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// TenantDatabase shows every methods available under the tenant to interact with the database
type TenantDatabase interface {
	CreateTenant(ctx context.Context, tenant model.Tenant) (model.Tenant, error)
	GetTenantByID(ctx context.Context, id uuid.UUID) (model.Tenant, error)
	UpdateTenantByID(ctx context.Context, tenant model.Tenant) error
	GetTenantByEmail(ctx context.Context, email string) (model.Tenant, error)
}

// Tenant object
type Tenant struct {
	logger  zerolog.Logger
	storage *Storage
}

// NewTenant creates a new reference to the Tenant storage entity
func NewTenant(s *Storage) *TenantDatabase {
	l := s.Logger.With().Str(helper.LogStrKeyLevel, "tenant").Logger()
	tenant := &Tenant{
		logger:  l,
		storage: s,
	}

	tenantDatabase := TenantDatabase(tenant)
	return &tenantDatabase
}

// CreateTenant creates a new role in the tenants table
func (t *Tenant) CreateTenant(ctx context.Context, tenant model.Tenant) (model.Tenant, error) {
	db := t.storage.DB.WithContext(ctx).Model(&model.Tenant{}).Create(&tenant)
	if db.Error != nil {
		t.logger.Err(db.Error).Msgf("CreateTenant error: %v, (%v)", ErrRecordCreatingFailed, db.Error)
		return model.Tenant{}, ErrRecordCreatingFailed
	}

	return tenant, nil
}

// GetTenantByID returns a tenant matching the ID
func (t *Tenant) GetTenantByID(ctx context.Context, id uuid.UUID) (model.Tenant, error) {
	var tenant model.Tenant
	db := t.storage.DB.WithContext(ctx).Where("id = ?", id).First(&tenant)
	if db.Error != nil || strings.EqualFold(tenant.ID.String(), helper.ZeroUUID) {
		t.logger.Err(db.Error).Msgf("GetTenantByID error: %v (%v)", ErrRecordNotFound, db.Error)
		return tenant, ErrRecordNotFound
	}

	return tenant, nil
}

// UpdateTenantByID should update the tenant record in the storage
func (t *Tenant) UpdateTenantByID(ctx context.Context, tenant model.Tenant) error {
	db := t.storage.DB.WithContext(ctx).Model(model.Tenant{
		ID: tenant.ID,
	}).UpdateColumns(tenant)

	if db.Error != nil {
		t.logger.Err(db.Error).Msgf("UpdateTenantByID error: %v, (%v)", ErrRecordUpdateFailed, db.Error)
		return ErrRecordUpdateFailed
	}

	return nil
}

// GetTenantByEmail returns a tenant matching the email
func (t *Tenant) GetTenantByEmail(ctx context.Context, email string) (model.Tenant, error) {
	var tenant model.Tenant

	db := t.storage.DB.WithContext(ctx).Where("email = ?", email).First(&tenant)
	if db.Error != nil || strings.EqualFold(tenant.ID.String(), helper.ZeroUUID) {
		t.logger.Err(db.Error).Msgf("GetTenantByEmail error: %v (%v)", ErrRecordNotFound, db.Error)
		return tenant, ErrRecordNotFound
	}

	return tenant, nil
}

package storage

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"codematic/model"
	"codematic/pkg/helper"
)

// BalanceDatabase is the interface to interact with the balance table
type BalanceDatabase interface {
	CreateBalance(ctx context.Context, balance model.Balance) (model.Balance, error)
	UpdateUserBalance(ctx context.Context, balance model.Balance) (model.Balance, error)
	GetLastBalanceByUserID(ctx context.Context, userID uuid.UUID) (model.Balance, error)
}

// Balance object
type Balance struct {
	logger  zerolog.Logger
	storage *Storage
}

// NewBalance creates a new reference to a Balance storage entity
func NewBalance(s *Storage) *BalanceDatabase {
	l := s.Logger.With().Str(helper.LogStrKeyLevel, "balance").Logger()
	balance := &Balance{
		logger:  l,
		storage: s,
	}
	balanceDatabase := BalanceDatabase(balance)
	return &balanceDatabase
}

// CreateBalance create a new balaance row on the table
func (b *Balance) CreateBalance(ctx context.Context, balance model.Balance) (model.Balance, error) {
	db := b.storage.DB.WithContext(ctx).Model(&model.Balance{}).Create(&balance)
	if db.Error != nil {
		b.storage.Logger.Err(db.Error).Msgf("Balance:: balance creation error: %v, (%v)", ErrRecordCreatingFailed, db.Error)

		if strings.Contains(db.Error.Error(), "duplicate key value") {
			return model.Balance{}, ErrDuplicateRecord
		}
		return model.Balance{}, ErrRecordCreatingFailed
	}

	return balance, nil
}

// UpdateUserBalance updates user balance
func (b *Balance) UpdateUserBalance(ctx context.Context, balance model.Balance) (model.Balance, error) {
	db := b.storage.DB.WithContext(ctx).Model(model.Balance{
		ID: balance.ID,
	}).UpdateColumns(balance)

	if db.Error != nil {
		b.storage.Logger.Err(db.Error).Msgf("UpdateUserBalance error: %v, (%v)", ErrRecordUpdateFailed, db.Error)
		return balance, ErrRecordUpdateFailed
	}

	return balance, nil
}

// GetLastBalanceByUserID returns the latest balance for a user, or 0.00 if none exists.
func (b *Balance) GetLastBalanceByUserID(ctx context.Context, userID uuid.UUID) (model.Balance, error) {
	var balance model.Balance

	err := b.storage.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&balance).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// No balance record exists, return zero balance
		return model.Balance{
			UserID:        userID,
			BalanceBefore: 0.00,
			BalanceAfter:  0.00,
		}, nil
	}

	if err != nil {
		return model.Balance{}, err
	}

	return balance, nil
}

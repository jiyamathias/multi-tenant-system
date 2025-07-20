package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"codematic/model"
	"codematic/pkg/helper"
)

// WalletDatabase shows every methods available under the wallet to interact with the database
type WalletDatabase interface {
	CreateWallet(ctx context.Context, wallet model.Wallet) (model.Wallet, error)
	GetWalletByUserID(ctx context.Context, userID uuid.UUID) (model.Wallet, error)
	UpdateWalletByID(ctx context.Context, wallet model.Wallet) error
}

// Wallet object
type Wallet struct {
	storage *Storage
}

// NewWallet allowing the wallet object to implement all of the methods in the NewWalletDatabase interface
func NewWallet(s *Storage) *WalletDatabase {
	s.Logger.With().Str(helper.LogStrKeyLevel, "wallet").Logger()

	w := &Wallet{
		storage: s,
	}

	wallet := WalletDatabase(w)
	return &wallet
}

// CreateWallet add a new wallet into the wallet table
func (w *Wallet) CreateWallet(ctx context.Context, wallet model.Wallet) (model.Wallet, error) {
	db := w.storage.DB.WithContext(ctx).Model(&model.Wallet{}).Create(&wallet)
	if db.Error != nil {
		w.storage.Logger.Err(db.Error).Msgf("Wallet ::: Wallet creation error: %v, (%v)", ErrRecordCreatingFailed, db.Error)
		return model.Wallet{}, ErrRecordCreatingFailed
	}

	return wallet, nil
}

// GetWalletByUserID gets a wallet by it the user's id from the wallet table
func (w *Wallet) GetWalletByUserID(ctx context.Context, userID uuid.UUID) (model.Wallet, error) {
	var wallet model.Wallet

	db := w.storage.DB.WithContext(ctx).Where("user_id = ?", userID).First(&wallet)
	if db.Error != nil || strings.EqualFold(wallet.ID.String(), helper.ZeroUUID) {
		w.storage.Logger.Err(db.Error).Msgf("GetWalletByUserID ::: error: %v (%v)", ErrRecordNotFound, db.Error)
		fmt.Println(db.Error)
		return wallet, db.Error
	}

	return wallet, nil
}

// UpdateWalletByID updates users wallet record in the wallet table using the wallet ID
func (w *Wallet) UpdateWalletByID(ctx context.Context, wallet model.Wallet) error {
	db := w.storage.DB.WithContext(ctx).Model(model.Wallet{
		ID: wallet.ID,
	}).UpdateColumns(wallet)

	if db.Error != nil {
		w.storage.Logger.Err(db.Error).Msgf("UpdateWalletByID ::: error: %v, (%v)", ErrRecordUpdateFailed, db.Error)
		return ErrRecordUpdateFailed
	}

	return nil
}

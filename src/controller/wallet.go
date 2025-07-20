package controller

import (
	"context"

	"github.com/google/uuid"

	"codematic/model"
)

// CreateWallet add a new wallet into the wallet table
func (c *Controller) CreateWallet(ctx context.Context, wallet model.Wallet) (model.Wallet, error) {
	return c.walletStorage.CreateWallet(ctx, wallet)
}

// GetWalletByUserID gets a wallet by it the user's id from the wallet table
func (c *Controller) GetWalletByUserID(ctx context.Context, userID uuid.UUID) (model.Wallet, error) {
	return c.walletStorage.GetWalletByUserID(ctx, userID)
}

// UpdateWalletByID updates users wallet record in the wallet table using the wallet ID
func (c *Controller) UpdateWalletByID(ctx context.Context, wallet model.Wallet) error {
	return c.walletStorage.UpdateWalletByID(ctx, wallet)
}

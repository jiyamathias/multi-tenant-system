package controller

import (
	"context"

	"github.com/google/uuid"

	"codematic/model"
	"codematic/model/pagination"
)

// CreateTransaction method inserts a new transaction record into the transactions table
func (c *Controller) CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	return c.transactionStorage.CreateTransaction(ctx, transaction)
}

// GetTransactionsByUserID retrieves all transactions for a specific user
func (c *Controller) GetTransactionsByUserID(ctx context.Context, userID uuid.UUID, transactionFlow *model.TransactionFlow, page pagination.Page) ([]model.Transaction, pagination.PageInfo, error) {
	return c.transactionStorage.GetTransactionsByUserID(ctx, userID, transactionFlow, page)
}

// GetTransactionByID retrieves a specific transaction by its ID
func (c *Controller) GetTransactionByID(ctx context.Context, transactionID uuid.UUID) (model.Transaction, error) {
	return c.transactionStorage.GetTransactionByID(ctx, transactionID)
}

// UpdateTransactionByID updates a transaction by transaction ID
func (c *Controller) UpdateTransactionByID(ctx context.Context, transaction model.Transaction) error {
	return c.transactionStorage.UpdateTransactionByID(ctx, transaction)
}

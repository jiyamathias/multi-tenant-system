package storage

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"codematic/model"
	"codematic/model/pagination"
	"codematic/pkg/helper"
)

// TransactionDatabase enlists all possible operations to be performed with the transaction storage
type TransactionDatabase interface {
	CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error)
	GetTransactionsByUserID(ctx context.Context, userID uuid.UUID, transactionFlow *model.TransactionFlow, page pagination.Page) ([]model.Transaction, pagination.PageInfo, error)
	GetTransactionByID(ctx context.Context, transactionID uuid.UUID) (model.Transaction, error)
	UpdateTransactionByID(ctx context.Context, transaction model.Transaction) error
}

// Transaction config object
type Transaction struct {
	logger  zerolog.Logger
	storage *Storage
}

// NewTransaction creates a new reference to a transaction storage entity
func NewTransaction(s *Storage) *TransactionDatabase {
	l := s.Logger.With().Str(helper.LogStrKeyLevel, "transaction").Logger()
	Transaction := &Transaction{
		logger:  l,
		storage: s,
	}
	TransactionDatabase := TransactionDatabase(Transaction)
	return &TransactionDatabase
}

// CreateTransaction method inserts a new transaction record into the transactions table
func (tx *Transaction) CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	db := tx.storage.DB.WithContext(ctx).Model(&model.Transaction{}).Create(&transaction)
	if db.Error != nil {
		tx.logger.Err(db.Error).Msgf("TransactionService:: Transaction creation error: %v", db.Error)
		if strings.Contains(db.Error.Error(), "duplicate key value") {
			return model.Transaction{}, ErrDuplicateRecord
		}
		return model.Transaction{}, ErrRecordCreatingFailed
	}

	return transaction, nil
}

// GetTransactionsByUserID retrieves all transactions for a specific user
func (tx *Transaction) GetTransactionsByUserID(ctx context.Context, userID uuid.UUID, transactionFlow *model.TransactionFlow, page pagination.Page) ([]model.Transaction, pagination.PageInfo, error) {
	var transactions []model.Transaction

	offset := 0
	if page.Number == nil {
		tmpPageNumber := pagination.PageDefaultNumber
		page.Number = &tmpPageNumber
	}
	if page.Size == nil {
		tmpPageSize := pagination.PageDefaultSize
		page.Size = &tmpPageSize
	}

	if *page.Number > 1 {
		offset = *page.Size * (*page.Number - 1)
	}

	query := tx.storage.DB.WithContext(ctx).Where("user_id = ?", userID)

	if transactionFlow != nil {
		query = tx.storage.DB.WithContext(ctx).Where("user_id = ? AND transaction_flow = ?", userID, transactionFlow)
	}

	var count int64
	query.Model(model.Transaction{}).Count(&count)

	db := query.Offset(offset).Limit(*page.Size).Find(&transactions)
	if db.Error != nil {
		tx.logger.Err(db.Error).Msgf("TransactionService:: Error fetching transactions for user %v: %v", userID, db.Error)
		return nil, pagination.PageInfo{}, ErrRecordNotFound
	}

	return transactions, pagination.PageInfo{
		Page:            *page.Number,
		Size:            *page.Size,
		HasNextPage:     int64(offset+*page.Size) < count,
		HasPreviousPage: *page.Number > 1,
		TotalCount:      count,
	}, nil
}

// GetTransactionByID retrieves a specific transaction by its ID
func (tx *Transaction) GetTransactionByID(ctx context.Context, transactionID uuid.UUID) (model.Transaction, error) {
	var transaction model.Transaction

	db := tx.storage.DB.WithContext(ctx).Where("id = ?", transactionID).First(&transaction)
	if db.Error != nil {
		tx.logger.Err(db.Error).Msgf("TransactionService:: Error fetching transaction by ID %v: %v", transactionID, db.Error)
		return model.Transaction{}, ErrRecordNotFound
	}

	return transaction, nil
}

// UpdateTransactionByID updates a transaction by transaction ID
func (tx *Transaction) UpdateTransactionByID(ctx context.Context, transaction model.Transaction) error {
	db := tx.storage.DB.WithContext(ctx).Model(model.Transaction{}).Where("id = ?", transaction.ID).
		Updates(transaction)
	if db.Error != nil {
		tx.logger.Err(db.Error).Msgf("UpdateTransactionByID ::: %s", db.Error)
		return ErrRecordUpdateFailed
	}

	return nil
}

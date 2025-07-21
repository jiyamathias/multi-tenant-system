package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"codematic/model"
)

func (c *Controller) ProcessPaymentWebhook(ctx context.Context, payload model.PaymentWebhook) error {
	txID := strings.Split(payload.Data.Reference, "_")
	tx, err := c.GetTransactionByID(ctx, uuid.MustParse(txID[1]))
	if err != nil {
		c.logger.Err(err).Msgf("GetTransactionByID ===> error getting transaction by ID %v", err)
		return ErrTransactionID
	}

	user, err := c.GetUserByID(ctx, tx.UserID)
	if err != nil {
		c.logger.Err(err).Msgf("GetUserByID ===> error getting user by ID %v", err)
		return err
	}

	// TODO save the ID into redis to avoid updating the transaction matching ID several times
	setStatus, err := c.GetStringValue(ctx, fmt.Sprintf("wbk_%s_%s", txID[0], tx.ID.String()))
	if err != nil {
		if err == redis.Nil {
			setStatus = ""
		} else {
			return err
		}
	}

	if setStatus == string(model.TransactionStatusSuccessful) {
		// do not do anything, since the transaction has already been updated. Just return without an error
		return nil
	}

	tx.MetaData = &payload.Data.Metadata
	tx.Charges = payload.Data.Fees
	tx.Amount = payload.Data.Amount
	tx.Currency = payload.Data.Currency

	// check if imcoming transaction status is successful
	switch payload.Data.Status {
	case "success":
		// previous balance
		previousBal, err := c.GetLastBalanceByUserID(ctx, tx.UserID)
		if err != nil {
			c.logger.Err(err).Msgf("error getting user last balance ===> %v", err)
			return err
		}

		// create new balance
		balance := model.Balance{
			ID:              uuid.New(),
			UserID:          tx.UserID,
			TransactionType: model.TransactionTypeCredit,
			TransactionID:   tx.ID,
			BalanceBefore:   previousBal.BalanceAfter,
			BalanceAfter:    previousBal.BalanceAfter + payload.Data.Amount,
		}

		// if the transaction type is withdrawal
		if txID[0] == "dbt" {
			balance.TransactionType = model.DebitTransaction
			balance.BalanceAfter = previousBal.BalanceAfter - payload.Data.Amount
		}

		newBal, err := c.CreateBalance(ctx, balance)
		if err != nil {
			c.logger.Err(err).Msgf("error creating user balance ===> %v", err)
			return err
		}

		// get user wallet balance
		wallet, err := c.GetWalletByUserID(ctx, tx.UserID)
		if err != nil {
			c.logger.Err(err).Msgf("error getting wallet by userID ===> %v", err)
			return err
		}

		txType := model.CreditTransaction

		// update wallet balance
		wallet.BalanceBefore = wallet.BalanceAfter
		wallet.BalanceAfter = wallet.BalanceAfter + payload.Data.Amount
		wallet.TransactionID = &tx.ID
		wallet.TransactionType = &txType
		wallet.BalanceID = &newBal.ID

		// if the transaction type is withdrawal
		if txID[0] == "dbt" {
			txType = model.DebitTransaction
			wallet.TransactionType = &txType
			wallet.BalanceAfter = previousBal.BalanceAfter - payload.Data.Amount
		}

		if err := c.UpdateWalletByID(ctx, wallet); err != nil {
			c.logger.Err(err).Msgf("error updating wallet by ID ===> %v", err)
			return err
		}

		// create a audit log
		auditLog := model.AuditLog{
			ID:            uuid.New(),
			TransactionID: &tx.ID,
			UserID:        &tx.UserID,
			TenantID:      &user.TenantID,
			Actor:         model.ActorUser,
			ActionDone:    model.ActionSuccess,
			Messages:      "reveived a credit to wallet",
		}

		if txID[0] == "dbt" {
			auditLog.Messages = "reveived a debit to wallet"
		}

		if _, err = c.CreateAuditLog(ctx, auditLog); err != nil {
			c.logger.Err(err).Msgf("error creating audit log")
			return err
		}

		// save that key to redis for 7 days, it is safe to assume that webhooks would not keep on being sent for 7 straight days :D.
		// And as such, safely delete the key + value after 7 days
		if err := c.SetValue(ctx, fmt.Sprintf("wbk_%s_%s", txID[0], tx.ID.String()), model.TransactionStatusSuccessful, 7*24*time.Hour); err != nil {
			return err
		}
	case "failed":
		tx.Status = model.TransactionStatusFailed

		// create a audit log
		auditLog := model.AuditLog{
			ID:            uuid.New(),
			TransactionID: &tx.ID,
			UserID:        &tx.UserID,
			TenantID:      &user.TenantID,
			Actor:         model.ActorUser,
			ActionDone:    model.ActionFailed,
			Messages:      "transaction failed",
		}

		if _, err = c.CreateAuditLog(ctx, auditLog); err != nil {
			c.logger.Err(err).Msgf("error creating audit log")
			return err
		}
	case "pending":
		tx.Status = model.TransactionStatusPending

		// create a audit log
		auditLog := model.AuditLog{
			ID:            uuid.New(),
			TransactionID: &tx.ID,
			UserID:        &tx.UserID,
			TenantID:      &user.TenantID,
			Actor:         model.ActorUser,
			ActionDone:    model.ActionPending,
			Messages:      "transaction pending",
		}

		if _, err = c.CreateAuditLog(ctx, auditLog); err != nil {
			c.logger.Err(err).Msgf("error creating audit log")
			return err
		}
	}

	if err := c.UpdateTransactionByID(ctx, tx); err != nil {
		c.logger.Err(err).Msgf("error updating transaction by ID ===> %v", err)
		return err
	}

	return nil
}

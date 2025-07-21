package controller

import (
	"codematic/model"
	"context"

	"github.com/google/uuid"
)

func (c *Controller) VirtualAccount(ctx context.Context, userID uuid.UUID, fullName, bankName string) (model.VirtualAccount, error) {
	user, err := c.GetUserByID(ctx, userID)
	if err != nil {
		c.logger.Err(err).Msgf("error getting user by ID ::: %v", err)
		return model.VirtualAccount{}, err
	}

	payload := model.InitiateTransaction{
		FullName: fullName,
		BankName: bankName,
	}

	virtualAccount, err := c.paymentService.InitiateTransaction(model.PaymentProviderFlutterwave, model.PaymentActionVirtualAccount, payload)
	if err != nil {
		c.logger.Err(err).Msgf("VirtualAccount ::: InitiateTransaction ===> %v", err)
		return model.VirtualAccount{}, err
	}

	auditLog := model.AuditLog{
		ID:       uuid.New(),
		TenantID: &user.TenantID,
		// TransactionID: &transaction.ID,
		UserID:     &user.ID,
		Actor:      model.ActorUser,
		ActionDone: model.ActionCreated,
		Messages:   "user logged in",
	}

	if _, err = c.CreateAuditLog(ctx, auditLog); err != nil {
		c.logger.Err(err).Msgf("error creating audit log")
		return model.VirtualAccount{}, err
	}

	return virtualAccount, nil
}

// Deposit is a method used to add funds to once wallet. We would be assumming that we already have the users card details.
// And as such, all we need to for the user to pass in the amount they would want to depoist into their wallet, and it gets processed
// and their wallet gets deposited if no errors occures
func (c *Controller) Deposit(ctx context.Context, userID uuid.UUID, amount float64) error {
	user, err := c.GetUserByID(ctx, userID)
	if err != nil {
		c.logger.Err(err).Msgf("error getting user by ID ::: %v", err)
		return err
	}

	payload := model.InitiateTransaction{
		Amount: amount,
	}

	_, err = c.paymentService.InitiateTransaction(model.PaymentProviderFlutterwave, model.PaymentActionDeposit, payload)
	if err != nil {
		c.logger.Err(err).Msgf("Deposit ::: InitiateTransaction ===> %v", err)
		return err
	}

	// create a transaction history
	transaction := model.Transaction{
		ID:              uuid.New(),
		UserID:          user.ID,
		Amount:          amount,
		Charges:         (amount * 10) / 100, // assumming processing charges is 10%
		TransactionType: model.CreditTransaction,
		TransactionFlow: model.TransactionFlowRevenue,
	}

	if _, err := c.CreateTransaction(ctx, transaction); err != nil {
		c.logger.Err(err).Msgf("Deposit ::: CreateTransaction ::: error creating transaction history ===> %v", err)
		return err
	}

	// create a audit log
	auditLog := model.AuditLog{
		ID:            uuid.New(),
		TenantID:      &user.TenantID,
		TransactionID: &transaction.ID,
		UserID:        &user.ID,
		Actor:         model.ActorUser,
		ActionDone:    model.ActionCreated,
		Messages:      "deposit created",
	}

	if _, err = c.CreateAuditLog(ctx, auditLog); err != nil {
		c.logger.Err(err).Msgf("error creating audit log")
		return err
	}

	return nil
}

func (c *Controller) Transfer(ctx context.Context, userID uuid.UUID, bankNumber, accountNumber string, amount float64) error {
	user, err := c.GetUserByID(ctx, userID)
	if err != nil {
		c.logger.Err(err).Msgf("error getting user by ID ::: %v", err)
		return err
	}

	payload := model.InitiateTransaction{
		BankNumber:    bankNumber,
		AccountNumber: accountNumber,
		Amount:        amount,
	}

	_, err = c.paymentService.InitiateTransaction(model.PaymentProviderFlutterwave, model.PaymentActionTransfer, payload)
	if err != nil {
		c.logger.Err(err).Msgf("Transfer ::: InitiateTransaction ===> %v", err)
		return err
	}

	// create a transaction history
	transaction := model.Transaction{
		ID:              uuid.New(),
		UserID:          user.ID,
		Amount:          amount,
		Charges:         (amount * 10) / 100, // assumming processing charges is 10%
		TransactionType: model.DebitTransaction,
		TransactionFlow: model.TransactionFlowWithdrawal,
	}

	if _, err := c.CreateTransaction(ctx, transaction); err != nil {
		c.logger.Err(err).Msgf("Deposit ::: CreateTransaction ::: error creating transaction history ===> %v", err)
		return err
	}

	auditLog := model.AuditLog{
		ID:            uuid.New(),
		TenantID:      &user.TenantID,
		TransactionID: &transaction.ID,
		UserID:        &user.ID,
		Actor:         model.ActorUser,
		ActionDone:    model.ActionCreated,
		Messages:      "withdrawal created",
	}

	if _, err = c.CreateAuditLog(ctx, auditLog); err != nil {
		c.logger.Err(err).Msgf("error creating audit log")
		return err
	}

	return nil
}

package controller

import "codematic/model"

func (c *Controller) VirtualAccount(fullName, bankName string) error {
	payload := model.InitiateTransaction{
		FullName: fullName,
		BankName: bankName,
	}

	_, err := c.paymentService.InitiateTransaction(model.PaymentProviderFlutterwave, model.PaymentActionVirtualAccount, payload)
	if err != nil {
		c.logger.Err(err).Msgf("VirtualAccount ::: InitiateTransaction ===> %v", err)
		return err
	}

	return nil
}

func (c *Controller) Deposit(amount float64) error {
	payload := model.InitiateTransaction{
		Amount: amount,
	}

	_, err := c.paymentService.InitiateTransaction(model.PaymentProviderFlutterwave, model.PaymentActionDeposit, payload)
	if err != nil {
		c.logger.Err(err).Msgf("Deposit ::: InitiateTransaction ===> %v", err)
		return err
	}

	return nil
}

func (c *Controller) Transfer(bankNumber, accountNumber string, amount float64) error {
	payload := model.InitiateTransaction{
		BankNumber:    bankNumber,
		AccountNumber: accountNumber,
		Amount:        amount,
	}
	_, err := c.paymentService.InitiateTransaction(model.PaymentProviderFlutterwave, model.PaymentActionTransfer, payload)
	if err != nil {
		c.logger.Err(err).Msgf("Transfer ::: InitiateTransaction ===> %v", err)
		return err
	}

	return nil
}

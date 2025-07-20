package payment

import (
	"fmt"

	"codematic/model"
	"codematic/pkg/helper"
)

const (
	flutterwaveBaseURL string = "flutter-wave-base-url"
	paystackBaseURL    string = "pay-stack-base-url"
)

// VirtualAccount create a virtual account for a user
func (f *flutterwaveProvider) VirtualAccount(fullName, bankName string) (model.VirtualAccount, error) {
	acctNum, err := helper.GenerateRandomDigits(10)
	if err != nil {
		return model.VirtualAccount{}, helper.ErrCreatingAcctNumber
	}

	return model.VirtualAccount{
		BankName:      bankName,
		AccountNumber: acctNum,
	}, nil
}

// Transfer make external transfers to with the users bank or someone else's bank account
func (f *flutterwaveProvider) Transfer(bankNumber string, accountNumber string, amount float64) error {
	fmt.Printf("Flutterwave Withdraw: %s -> %s : %.2f\n", bankNumber, accountNumber, amount)
	return nil
}

// Deposit simulate making an API call to flutterwave assumming we already have
// the users card details, so we just need an amount to be passed in
func (f *flutterwaveProvider) Deposit(amount float64) error {
	fmt.Printf("Flutterwave Deposit: %.2f\n", amount)
	return nil
}

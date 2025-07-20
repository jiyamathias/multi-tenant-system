package model

type (
	PaymentAction   string
	PaymentProvider string
)

const (
	PaymentActionTransfer       PaymentAction = "transfer"
	PaymentActionDeposit        PaymentAction = "deposit"
	PaymentActionVirtualAccount PaymentAction = "virtual_account"

	PaymentProviderPaystack    PaymentProvider = "paystack"
	PaymentProviderFlutterwave PaymentProvider = "flutterwave"
)

type (
	VirtualAccount struct {
		BankName      string `json:"bankName"`
		AccountNumber string `json:"accountNumber"`
	}

	InitiateTransaction struct {
		AccountNumber string
		BankName      string
		Amount        float64
		BankNumber    string
		FullName      string
	}
)

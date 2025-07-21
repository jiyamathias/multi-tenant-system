package payment

type (
	depositRequest struct {
		Amount float64 `json:"amount" validate:"required"`
	}

	makeTransferRequest struct {
		BankNumber    string  `json:"bankNumber" validate:"required"`
		AccountNumber string  `json:"accountNumber" validate:"required"`
		Amount        float64 `json:"amount" validate:"required"`
	}

	bankTransferRequest struct {
		FulName  string `json:"fullName" validate:"required"`
		BankName string `json:"bankName" validate:"required"`
	}
)

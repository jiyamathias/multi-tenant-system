package model

type (
	PaymentWebhook struct {
		Event string `json:"event"` // success, failed, pending
		Data  struct {
			Status    string         `json:"status"`
			Reference string         `json:"reference"`
			Amount    float64        `json:"amount"`
			Currency  string         `json:"currency"`
			Metadata  map[string]any `json:"metadata"`
			Fees      float64        `json:"fees"`
		} `json:"data"`
	}
)

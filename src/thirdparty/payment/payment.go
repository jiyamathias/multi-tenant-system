// Package payment contains the setup for a payment gateway
package payment

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"

	"codematic/model"
	"codematic/pkg/environment"
	"codematic/pkg/helper"
	"codematic/storage"
)

const (
	packageName = "payment"
)

// PaymentProvider defines the contract each payment provider must implement
type PaymentProvider interface {
	VirtualAccount(fullName, bankName string) (model.VirtualAccount, error)
	Transfer(bankNumber string, accountNumber string, amount float64) error
	Deposit(amount float64) error
}

// PaymentService provides access to all registered payment providers
type (
	PaymentService struct {
		providers map[model.PaymentProvider]PaymentProvider // e.g., "paystack", "flutterwave"
		logger    zerolog.Logger
		env       *environment.Env
		storage   *storage.Storage
	}
	// flutterwaveProvider implements PaymentProvider
	flutterwaveProvider struct {
		APIKey string
	}

	// paystackProvider implements PaymentProvider
	paystackProvider struct {
		APIKey string
	}
)

// NewFlutterwaveProvider initializes a new Flutterwave provider
func NewFlutterwaveProvider(apiKey string) *flutterwaveProvider {
	return &flutterwaveProvider{APIKey: apiKey}
}

// NewPaystackProvider initializes a new Paystack provider
func NewPaystackProvider(apiKey string) *paystackProvider {
	return &paystackProvider{APIKey: apiKey}
}

// New initializes the PaymentService with all supported providers
func New(z zerolog.Logger, ev *environment.Env, s *storage.Storage) *PaymentService {
	l := z.With().Str(helper.LogStrKeyLevel, packageName).Logger()

	paystack := NewPaystackProvider(ev.Get("PAYSTACK_API_KEY"))
	flutterwave := NewFlutterwaveProvider(ev.Get("FLUTTERWAVE_API_KEY"))

	providers := map[model.PaymentProvider]PaymentProvider{
		model.PaymentProviderPaystack:    paystack,
		model.PaymentProviderFlutterwave: flutterwave,
	}

	return &PaymentService{
		providers: providers,
		logger:    l,
		env:       ev,
		storage:   s,
	}
}

// GetProvider returns the payment provider by key (e.g., "paystack")
func (ps *PaymentService) GetProvider(providerKey model.PaymentProvider) (PaymentProvider, bool) {
	p, ok := ps.providers[providerKey]
	return p, ok
}

func (s *PaymentService) InitiateTransaction(provider model.PaymentProvider, action model.PaymentAction, payload model.InitiateTransaction) (model.VirtualAccount, error) {
	p, ok := s.providers[provider]
	if !ok {
		return model.VirtualAccount{}, fmt.Errorf("unsupported provider: %s", provider)
	}

	switch action {
	case model.PaymentActionVirtualAccount:
		return p.VirtualAccount(payload.FullName, payload.BankName)
	case model.PaymentActionTransfer:
		return model.VirtualAccount{}, p.Transfer(payload.BankNumber, payload.AccountNumber, payload.Amount)
	case model.PaymentActionDeposit:
		return model.VirtualAccount{}, p.Deposit(payload.Amount)
	default:
		return model.VirtualAccount{}, errors.New("unsupported action")
	}
}

func (p *PaymentService) IsIdempotencyKeyUsed(ctx context.Context, key string) (bool, error) {
	// TODO check DB idempotency key
	return false, nil
}

func (p *PaymentService) MarkIdempotencyKeyUsed(ctx context.Context, key string) error {
	// TODO store idempotency key
	return nil
}

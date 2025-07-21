// Package controller defines implementation that exposes logics of the app
package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"codematic/model"
	"codematic/model/pagination"
	"codematic/pkg/environment"
	"codematic/pkg/helper"
	"codematic/pkg/middleware"
	"codematic/storage"
	"codematic/storage/redis"
	"codematic/thirdparty/payment"
)

const packageName = "controller"

// Operations enlist all possible operations for this controller across all modules
//
//go:generate mockgen -source controller.go -destination ./mock/mock_controller.go -package mock Operations
type Operations interface {
	Middleware() *middleware.Middleware

	CreateUser(ctx context.Context, u model.User) (model.User, error)
	AuthenticateUser(ctx context.Context, email, password string) (model.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (model.User, error)
	UpdateUserByID(ctx context.Context, userID uuid.UUID, u model.User) (model.User, error)

	GetAllAuditLogsByTransactionID(ctx context.Context, txID uuid.UUID, page pagination.Page) ([]*model.AuditLog, pagination.PageInfo, error)
	GetAuditLogByID(ctx context.Context, id uuid.UUID) (model.AuditLog, error)

	CreateBalance(ctx context.Context, balance model.Balance) (model.Balance, error)
	GetLastBalanceByUserID(ctx context.Context, userID uuid.UUID) (model.Balance, error)

	GetWalletByUserID(ctx context.Context, userID uuid.UUID) (model.Wallet, error)

	CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error)
	GetTransactionsByUserID(ctx context.Context, userID uuid.UUID, transactionFlow *model.TransactionFlow, page pagination.Page) ([]model.Transaction, pagination.PageInfo, error)
	GetTransactionByID(ctx context.Context, transactionID uuid.UUID) (model.Transaction, error)
	UpdateTransactionByID(ctx context.Context, transaction model.Transaction) error

	ProcessPaymentWebhook(ctx context.Context, payload model.PaymentWebhook) error

	CreateTenant(ctx context.Context, tenant model.Tenant) (model.Tenant, error)
	GetAllUsersByTenantID(ctx context.Context, tenantId uuid.UUID, page pagination.Page) ([]*model.User, pagination.PageInfo, error)
	AuthenticateTenant(ctx context.Context, email, password string) (model.Tenant, error)

	VirtualAccount(ctx context.Context, userID uuid.UUID, fullName, bankName string) (model.VirtualAccount, error)
	Deposit(ctx context.Context, userID uuid.UUID, amount float64) error
	Transfer(ctx context.Context, userID uuid.UUID, bankNumber, accountNumber string, amount float64) error
}

// Controller object to hold necessary reference to other dependencies
type Controller struct {
	logger     zerolog.Logger
	env        *environment.Env
	middleware *middleware.Middleware

	// storage layers
	userStorage        storage.UserDatabase
	auditLogStorage    storage.AuditLogDatabase
	balanceStorage     storage.BalanceDatabase
	walletStorage      storage.WalletDatabase
	transactionStorage storage.TransactionDatabase
	tenantStorage      storage.TenantDatabase

	redis redis.KvStore
	// third party services
	paymentService payment.PaymentService
}

// New creates a new instance of Controller
func New(z zerolog.Logger, s *storage.Storage, m *middleware.Middleware) *Operations {
	l := z.With().Str(helper.LogStrKeyModule, packageName).Logger()

	// init all storage layer under here
	user := storage.NewUser(s)

	auditLog := storage.NewAuditLog(s)
	balance := storage.NewBalance(s)
	wallet := storage.NewWallet(s)
	transaction := storage.NewTransaction(s)
	tenant := storage.NewTenant(s)

	newRedis := redis.NewRedis(s.Env, z, s.Env.Get("REDIS_SERVER_ADDRESS"))

	ctrl := &Controller{
		logger:      l,
		env:         s.Env,
		middleware:  m,
		userStorage: *user,

		auditLogStorage:    *auditLog,
		balanceStorage:     *balance,
		walletStorage:      *wallet,
		transactionStorage: *transaction,
		tenantStorage:      *tenant,

		redis: *newRedis,
	}

	op := Operations(ctrl)
	return &op
}

// Middleware returns the middleware object exposed by this app
func (c *Controller) Middleware() *middleware.Middleware {
	return c.middleware
}

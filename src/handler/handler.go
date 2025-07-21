// Package handler contains the all endpoints in this application
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"codematic/controller"
	auditLog "codematic/handler/auditLog"
	"codematic/handler/auth"
	"codematic/handler/docs"
	"codematic/handler/payment"
	"codematic/handler/tenant"
	"codematic/handler/transaction"
	"codematic/handler/wallet"
	"codematic/handler/webhook"
	"codematic/pkg/environment"
	"codematic/pkg/helper"
)

// Handler object
type Handler struct {
	application controller.Operations
	logger      *zerolog.Logger
	env         *environment.Env
	api         *gin.RouterGroup
}

// New creates a new instance of Handler
func New(z zerolog.Logger, ev *environment.Env, engine *gin.Engine, a controller.Operations) *Handler {
	log := z.With().Str(helper.LogStrPartnerLevel, "handler").Logger()
	apiGroup := engine.Group("/api")
	return &Handler{
		application: a,
		logger:      &log,
		env:         ev,
		api:         apiGroup,
	}
}

// Build setups the APi endpoints
func (h *Handler) Build() {
	v1 := h.api.Group("/v1")

	auth.New(v1, *h.logger, h.application, h.env)
	auditLog.New(v1, *h.logger, h.application, h.env)
	transaction.New(v1, *h.logger, h.application, h.env)
	webhook.New(v1, *h.logger, h.application, h.env)
	wallet.New(v1, *h.logger, h.application, h.env)
	tenant.New(v1, *h.logger, h.application, h.env)
	payment.New(v1, *h.logger, h.application, h.env)
	docs.New(v1)
}

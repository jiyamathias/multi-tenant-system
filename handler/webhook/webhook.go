// Package webhook handles webhooks reqponse from the payment gateway
package webhook

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"codematic/controller"
	restModel "codematic/handler/model"
	"codematic/model"
	"codematic/pkg/environment"
)

type webhookHandler struct {
	logger      zerolog.Logger
	controller  controller.Operations
	environment *environment.Env
}

// New creates a new instance of the subscription rest handler
func New(r *gin.RouterGroup, l zerolog.Logger, c controller.Operations, env *environment.Env) {
	webhk := webhookHandler{
		logger:      l,
		controller:  c,
		environment: env,
	}

	webhookGroup := r.Group("/webhook")
	webhookGroup.POST("/payment", webhk.processPayment())
}

func (w webhookHandler) processPayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request model.PaymentWebhook

		if err := c.ShouldBindJSON(&request); err != nil {
			w.logger.Error().Msgf("%v", err)
			return
		}

		if request.Data.Status != "success" && request.Data.Status != "failed" && request.Data.Status != "pending" {
			restModel.ErrorResponse(c, http.StatusBadRequest, "status can either be success, failed, pending")
			return
		}

		if err := w.controller.ProcessPaymentWebhook(context.Background(), request); err != nil {
			w.logger.Error().Msgf("error performing update on transaction %v", err)
			return
		}

		c.JSON(200, "success")
		c.Abort()
	}
}

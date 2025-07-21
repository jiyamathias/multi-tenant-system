// Package webhook handles webhooks reqponse from the payment gateway
package webhook

import (
	"context"
	"net/http"
	"strings"

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

		// set a dunny authentication to prevent webhook from updating records if it is not coming from the set payment gateway
		paymentAuth := c.GetHeader("auth")
		if paymentAuth != "payment" {
			restModel.ErrorResponse(c, http.StatusBadRequest, "this webhook is not from our payment gateway")
			return
		}

		// add some validation since we are simulation the webhook just to ensure everything works as it should.
		if !strings.Contains(request.Data.Reference, "_") {
			restModel.ErrorResponse(c, http.StatusBadRequest, "the reference must contain an underscore")
			return
		}

		txID := strings.Split(request.Data.Reference, "_")
		if txID[0] != "dbt" && txID[0] != "crt" {
			restModel.ErrorResponse(c, http.StatusBadRequest, "the reference format is dbt... or crt...")
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

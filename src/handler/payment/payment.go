package payment

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"codematic/controller"
	restModel "codematic/handler/model"
	"codematic/pkg/environment"
	"codematic/pkg/middleware"
)

type paymentHandler struct {
	logger      zerolog.Logger
	controller  controller.Operations
	environment *environment.Env
}

// New creates a new instance of the tenant rest handler
func New(r *gin.RouterGroup, l zerolog.Logger, c controller.Operations, env *environment.Env) {

	payment := paymentHandler{
		logger:      l,
		controller:  c,
		environment: env,
	}
	paymentGroup := r.Group("/payment")

	paymentGroup.POST("/deposit", payment.controller.Middleware().AuthMiddleware(), payment.makeDeposit())
	paymentGroup.POST("/transfer", payment.controller.Middleware().AuthMiddleware(), payment.makeTransfer())
	paymentGroup.POST("/bank-transfer", payment.controller.Middleware().AuthMiddleware(), payment.bankTransfer())
}

// makeDeposit 	godoc
//
//	@Summary		makeDeposit
//	@Description	this endpoint is used to make a deposit
//	@Tags			payment
//	@Accept			json
//	@Produce		json
//	@Param			depositRequest	body		depositRequest				true	"deposit request body"
//	@Success		200				{object}	restModel.GenericResponse	"wallet top up successful"
//	@Router			/payment/deposit [post]
func (p *paymentHandler) makeDeposit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request depositRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			p.logger.Error().Msgf("%v", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, "incomplete details please fill out the missing details")
			return
		}

		err := restModel.ValidateRequest(request)
		if err != nil {
			p.logger.Error().Msgf("%v", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, restModel.ErrIncompleteDetails.Error())
			return
		}

		userID, err := uuid.Parse(c.GetString(middleware.ActorIDInContext))
		if err != nil {
			p.logger.Err(err).Msgf("makeDeposit ::: error parsing uuid ==> %s", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		if err := p.controller.Deposit(context.Background(), userID, request.Amount); err != nil {
			p.logger.Error().Msgf("makeDeposit ::: %v", err)
			restModel.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		restModel.OkResponse(c, http.StatusOK, "wallet top up successful", nil)
	}
}

// makeTransfer 	godoc
//
//	@Summary		makeTransfer
//	@Description	this endpoint is used to make transfer
//	@Tags			payment
//	@Accept			json
//	@Produce		json
//	@Param			makeTransferRequest	body		makeTransferRequest				true	"make transfer request body"
//	@Success		200				{object}	restModel.GenericResponse	"transfer successful"
//	@Router			/payment/transfer [post]
func (p *paymentHandler) makeTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request makeTransferRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			p.logger.Error().Msgf("%v", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, "incomplete details please fill out the missing details")
			return
		}

		err := restModel.ValidateRequest(request)
		if err != nil {
			p.logger.Error().Msgf("%v", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, restModel.ErrIncompleteDetails.Error())
			return
		}

		userID, err := uuid.Parse(c.GetString(middleware.ActorIDInContext))
		if err != nil {
			p.logger.Err(err).Msgf("makeTransfer ::: error parsing uuid ==> %s", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		if err := p.controller.Transfer(context.Background(), userID, request.BankNumber, request.AccountNumber, request.Amount); err != nil {
			p.logger.Error().Msgf("makeTransfer ::: %v", err)
			restModel.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		restModel.OkResponse(c, http.StatusOK, "transfer successful", nil)
	}
}

// bankTransfer 	godoc
//
//	@Summary		bankTransfer
//	@Description	this endpoint is used to get a one time virtual account that is to be used top up once wallet
//	@Tags			payment
//	@Accept			json
//	@Produce		json
//	@Param			bankTransferRequest	body		bankTransferRequest				true	"bank transfer request body"
//	@Success		200				{object}	restModel.GenericResponse	"bank account created successful"
//	@Router			/payment/bank-transfer [post]
func (p *paymentHandler) bankTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request bankTransferRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			p.logger.Error().Msgf("%v", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, "incomplete details please fill out the missing details")
			return
		}

		err := restModel.ValidateRequest(request)
		if err != nil {
			p.logger.Error().Msgf("%v", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, restModel.ErrIncompleteDetails.Error())
			return
		}

		userID, err := uuid.Parse(c.GetString(middleware.ActorIDInContext))
		if err != nil {
			p.logger.Err(err).Msgf("bankTransfer ::: error parsing uuid ==> %s", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		virtualAccount, err := p.controller.VirtualAccount(context.Background(), userID, request.FulName, request.BankNumber)
		if err != nil {
			p.logger.Error().Msgf("bankTransfer ::: %v", err)
			restModel.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		restModel.OkResponse(c, http.StatusOK, "bank account created successful", virtualAccount)
	}
}

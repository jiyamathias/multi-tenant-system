package transaction

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"codematic/controller"
	restModel "codematic/handler/model"
	"codematic/model"
	"codematic/pkg/environment"
	"codematic/pkg/helper"
	"codematic/pkg/middleware"
)

type tsHandler struct {
	logger      zerolog.Logger
	controller  controller.Operations
	environment *environment.Env
}

// New creates a new instance of the transaction rest handler
func New(r *gin.RouterGroup, l zerolog.Logger, c controller.Operations, env *environment.Env) {
	ts := tsHandler{
		logger:      l,
		controller:  c,
		environment: env,
	}

	tsGroup := r.Group("/transaction")
	tsGroup.GET("/:id", ts.controller.Middleware().AuthMiddleware(), ts.getTransactionByID())
	tsGroup.GET("", ts.controller.Middleware().AuthMiddleware(), ts.getTransactionsByUserID())
	tsGroup.GET("/flow", ts.controller.Middleware().AuthMiddleware(), ts.getAllTransactionsByFlow())
}

// getTransactionByID 	godoc
//
//	@Summary		getTransactionByID
//	@Description	this endpoint gets a transaction by it ID
//	@Tags			transaction
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						false	"transactionID"
//	@Success		200	{object}	restModel.GenericResponse	"transaction details fetched successfully"
//	@Router			/transaction/{id} [get]
func (ts *tsHandler) getTransactionByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		params := c.Param("id")
		if len(params) == 0 {
			ts.logger.Err(helper.ErrSomeFieldsMissing).Msgf("GetTransactionByID :::  ==> %s", helper.ErrSomeFieldsMissing)
			restModel.ErrorResponse(c, http.StatusBadRequest, helper.ErrSomeFieldsMissing.Error())
			return
		}

		tsID, err := uuid.Parse(params)
		if err != nil {
			ts.logger.Err(err).Msgf("getTransactionByID :::  ==> %s", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		transaction, err := ts.controller.GetTransactionByID(context.Background(), tsID)
		if err != nil {
			ts.logger.Err(err).Msgf("getTransactionByID :::  ==> %s", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		restModel.OkResponse(c, http.StatusOK, "transaction details fetched successfully", transaction)
	}
}

// getTransactionsByUserID 	godoc
//
//	@Summary		getTransactionsByUserID
//	@Description	this endpoint is used to get all transactions belonging to a particular user
//	@Tags			transaction
//	@Accept			json
//	@Produce		json
//	@Param			page				query		string						false	"page"
//	@Param			size				query		string						false	"size"
//	@Param			sort_by				query		string						true	"sort_by"
//	@Param			sort_direction_desc	query		string						true	"sort_direction_desc"
//	@Success		200					{object}	restModel.GenericResponse	"transactions fetched successfully"
//	@Router			/transaction [get]
func (ts *tsHandler) getTransactionsByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.ActorIDInContext)

		transactions, pagination, err := ts.controller.GetTransactionsByUserID(context.Background(), uuid.MustParse(userID), nil, helper.ParsePageParams(c))
		if err != nil {
			ts.logger.Err(err).Msgf("getTransactionsByUserID :::  ==> %s", err)
			restModel.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		restModel.OkPaginatedResponse(c, http.StatusOK, "transactions fetched successfully", transactions, pagination)
	}
}

// getAllTransactionsByFlow 	godoc
//
//	@Summary		getAllTransactionsByFlow
//	@Description	this endpoint is used to get all transactions belonging to a particular user based on the flow of transaction flow
//	@Tags			transaction
//	@Accept			json
//	@Produce		json
//	@Param			page				query		string						false	"page"
//	@Param			size				query		string						false	"size"
//	@Param			sort_by				query		string						true	"sort_by"
//	@Param			sort_direction_desc	query		string						true	"sort_direction_desc"
//	@Param			flow				query		string						true	"flow"
//	@Success		200					{object}	restModel.GenericResponse	"transactions fetched successfully"
//	@Router			/transaction/flow [get]
func (ts *tsHandler) getAllTransactionsByFlow() gin.HandlerFunc {
	return func(c *gin.Context) {
		flow := c.Query("flow")
		if len(flow) == 0 {
			ts.logger.Error().Msgf("getAllTransactionsByFlow ::: %v", helper.ErrSomeFieldsMissing)
			restModel.ErrorResponse(c, http.StatusBadRequest, helper.ErrSomeFieldsMissing.Error())
			return
		}

		userID := c.GetString(middleware.ActorIDInContext)
		user, err := ts.controller.GetUserByID(context.Background(), uuid.MustParse(userID))
		if err != nil {
			ts.logger.Error().Msgf("getAllTransactionsByFlow ::: %v", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		if flow != string(model.TransactionFlowRevenue) && flow != string(model.TransactionFlowWithdrawal) {
			restModel.ErrorResponse(c, http.StatusBadRequest, "flow can either be: revenue or withdrawal")
			return
		}

		txFlow := model.TransactionFlow(flow)
		transactions, pageInfo, err := ts.controller.GetTransactionsByUserID(context.Background(), user.ID, &txFlow, helper.ParsePageParams(c))
		if err != nil {
			ts.logger.Err(err).Msgf("getAllTransactionsByFlow :::  ==> %s", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		restModel.OkPaginatedResponse(c, http.StatusOK, "transactions fetched successfully", transactions, pageInfo)
	}
}

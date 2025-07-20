// Package auditlog contains every request that pertains to the transactions acions
package auditlog

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"codematic/controller"
	restModel "codematic/handler/model"
	"codematic/pkg/environment"
	"codematic/pkg/helper"
)

type auditLogHandler struct {
	logger      zerolog.Logger
	controller  controller.Operations
	environment *environment.Env
}

// New creates a new instance of the audit log rest handler
func New(r *gin.RouterGroup, l zerolog.Logger, c controller.Operations, env *environment.Env) {

	auditLog := auditLogHandler{
		logger:      l,
		controller:  c,
		environment: env,
	}
	auditLogGroup := r.Group("/auditLog")

	auditLogGroup.GET("/transaction/:id", auditLog.controller.Middleware().AuthMiddleware(), auditLog.getAllAuditLogsByTransactionID())
	auditLogGroup.GET("/:id", auditLog.controller.Middleware().AuthMiddleware(), auditLog.getAuditLogByID())

}

// getAllAuditLogsByTransactionID 	godoc
//
//	@Summary		getAllAuditLogsByTransactionID
//	@Description	this endpoint gets all audit logs by the transaction ID
//	@Tags			audit-log
//	@Param			Authorization	header	string	true	"Bearer <token>"
//	@Accept			json
//	@Produce		json
//	@Param			id					path		string						false	"transaction id"
//	@Param			page				query		string						false	"page"
//	@Param			size				query		string						false	"size"
//	@Param			sort_by				query		string						false	"sort_by"
//	@Param			sort_direction_desc	query		string						false	"sort_direction_desc"
//	@Success		200					{object}	restModel.GenericResponse	"audit logs fetched successfully"
//	@Router			/auditLog/transaction/{id} [get]
func (a *auditLogHandler) getAllAuditLogsByTransactionID() gin.HandlerFunc {
	return func(c *gin.Context) {
		txID := c.Param("id")
		if len(txID) == 0 {
			a.logger.Error().Msgf("getAllAuditLogsByTransactionID ::: %v", helper.ErrIDMissing)
			restModel.ErrorResponse(c, http.StatusBadRequest, helper.ErrIDMissing.Error())
			return
		}

		auditlogs, pagination, err := a.controller.GetAllAuditLogsByTransactionID(context.Background(), uuid.MustParse(txID), helper.ParsePageParams(c))
		if err != nil {
			a.logger.Error().Msgf("getAllAuditLogsByTransactionID ::: %v", err)
			restModel.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		restModel.OkPaginatedResponse(c, http.StatusOK, "audit logs fetched successfully", auditlogs, pagination)
	}
}

// getAuditLogByID 	godoc
//
//	@Summary		getAuditLogByID
//	@Description	this endpoint gets an audit log by ID
//	@Tags			audit-log
//	@Param			Authorization	header	string	true	"Bearer <token>"
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						false	"audit id"
//	@Success		200	{object}	restModel.GenericResponse	"audit log fetched successfully"
//	@Router			/auditLog/{id} [get]
func (a *auditLogHandler) getAuditLogByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		auditLogID := c.Param("id")
		if len(auditLogID) == 0 {
			a.logger.Error().Msgf("getAuditLogByID ::: %v", helper.ErrIDMissing)
			restModel.ErrorResponse(c, http.StatusBadRequest, helper.ErrIDMissing.Error())
			return
		}

		auditLog, err := a.controller.GetAuditLogByID(context.Background(), uuid.MustParse(auditLogID))
		if err != nil {
			a.logger.Error().Msgf("getAuditLogByID ::: %v", err)
			restModel.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		restModel.OkResponse(c, http.StatusOK, "audit log fetched successfully", auditLog)
	}
}

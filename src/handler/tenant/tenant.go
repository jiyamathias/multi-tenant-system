package tenant

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
	"codematic/pkg/middleware"
)

type tenantHandler struct {
	logger      zerolog.Logger
	controller  controller.Operations
	environment *environment.Env
}

// New creates a new instance of the tenant rest handler
func New(r *gin.RouterGroup, l zerolog.Logger, c controller.Operations, env *environment.Env) {

	tenant := tenantHandler{
		logger:      l,
		controller:  c,
		environment: env,
	}
	auditLogGroup := r.Group("/tenant")

	auditLogGroup.POST("", tenant.createTenant())
	auditLogGroup.GET("", tenant.controller.Middleware().TenantAuthMiddleware(), tenant.getAllUsersByTenantID())

}

// createTenant 	godoc
//
//	@Summary		createTenant
//	@Description	this endpoint create a new tenent
//	@Tags			tenant
//	@Accept			json
//	@Produce		json
//	@Param			tenantRequest	body		tenantRequest				true	"tenant request body"
//	@Success		201				{object}	restModel.GenericResponse	"tenant created successfully"
//	@Router			/tenant [post]
func (t *tenantHandler) createTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request tenantRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			t.logger.Error().Msgf("%v", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, "incomplete details please fill out the missing details")
			return
		}

		tenant, err := t.controller.CreateTenant(context.Background(), request.toModel())
		if err != nil {
			t.logger.Error().Msgf("CreateTenant ::: %v", err)
			restModel.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		restModel.OkResponse(c, http.StatusCreated, "tenant created successfully", tenant)
	}
}

// getAllUsersByTenantID 	godoc
//
//	@Summary		getAllUsersByTenantID
//	@Description	this endpoint gets all users under a tenant
//	@Tags			tenant
//	@Param			Authorization	header	string	true	"Bearer <token>"
//	@Accept			json
//	@Produce		json
//	@Param			page				query		string						false	"page"
//	@Param			size				query		string						false	"size"
//	@Param			sort_by				query		string						false	"sort_by"
//	@Param			sort_direction_desc	query		string						false	"sort_direction_desc"
//	@Success		200	{object}	restModel.GenericResponse	"audit log fetched successfully"
//	@Router			/tenant [get]
func (t *tenantHandler) getAllUsersByTenantID() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, err := uuid.Parse(c.GetString(middleware.ActorIDInContext))
		if err != nil {
			t.logger.Err(err).Msgf("getWalletByUserID ::: error parsing uuid ==> %s", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		users, pagination, err := t.controller.GetAllUsersByTenantID(context.Background(), tenantID, helper.ParsePageParams(c))
		if err != nil {
			t.logger.Error().Msgf("getAuditLogByID ::: %v", err)
			restModel.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		restModel.OkPaginatedResponse(c, http.StatusOK, "users fetched successfully", users, pagination)
	}
}

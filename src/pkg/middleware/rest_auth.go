package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	restModel "codematic/handler/model"
	"codematic/model"
	"codematic/pkg/helper"
)

// AuthMiddleware authenticates a restful api call and inject the userID and userType into to context
func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		if len(bearerToken) == 0 {
			restModel.ErrorResponse(c, http.StatusUnauthorized, "Bearer token is missing")
			return
		}
		if !strings.HasPrefix(bearerToken, "Bearer ") {
			restModel.ErrorResponse(c, http.StatusUnauthorized, "Token type must be bearer")
			return
		}

		actorID, err := ParseToken(&m.env, strings.TrimPrefix(bearerToken, "Bearer "))
		if err != nil {
			restModel.ErrorResponse(c, http.StatusBadRequest, "unable to parse token")
			return
		}

		var actorType model.ActorType
		user := model.User{}
		ctx := context.WithValue(c.Request.Context(), helper.GinContextKey, c)

		// check if the actor is a user
		db := m.storage.DB.WithContext(ctx).Where("id = ?", actorID).First(&user)
		if db.Error != nil || strings.EqualFold(actorID, helper.ZeroUUID) {
			restModel.ErrorResponse(c, http.StatusBadRequest, ErrInvalidToken.Error())
			return
		}

		//validate the tenant
		tenant := model.Tenant{}
		db = m.storage.DB.WithContext(ctx).Where("id = ?", actorID).First(&tenant)
		if db.Error != nil || strings.EqualFold(actorID, helper.ZeroUUID) {
			restModel.ErrorResponse(c, http.StatusBadRequest, ErrInvalidTenant.Error())
			return
		}

		actorID = user.ID.String()
		actorType = model.ActorTypeUser
		tenantID := user.TenantID.String()

		c.Set(ActorIDInContext, actorID)
		c.Set(ActorTypeInContext, actorType)
		c.Set(UserInContext, &user)
		c.Set(TenantIDInContext, tenantID)

		c.Next()
	}
}

// TenantAuthMiddleware authenticates a restful api call and inject the tenant into to context
func (m *Middleware) TenantAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		if len(bearerToken) == 0 {
			restModel.ErrorResponse(c, http.StatusUnauthorized, "Bearer token is missing")
			return
		}
		if !strings.HasPrefix(bearerToken, "Bearer ") {
			restModel.ErrorResponse(c, http.StatusUnauthorized, "Token type must be bearer")
			return
		}

		actorID, err := TenantParseToken(&m.env, strings.TrimPrefix(bearerToken, "Bearer "))
		if err != nil {
			restModel.ErrorResponse(c, http.StatusBadRequest, "unable to parse token")
			return
		}

		var actorType model.ActorType
		tenant := model.Tenant{}
		ctx := context.WithValue(c.Request.Context(), helper.GinContextKey, c)

		// check if the actor is a tenant
		db := m.storage.DB.WithContext(ctx).Where("id = ?", actorID).First(&tenant)
		if db.Error != nil || strings.EqualFold(actorID, helper.ZeroUUID) {
			restModel.ErrorResponse(c, http.StatusBadRequest, ErrInvalidToken.Error())
			return
		}

		actorID = tenant.ID.String()
		actorType = model.ActorTypeTenant

		c.Set(ActorIDInContext, actorID)
		c.Set(ActorTypeInContext, actorType)

		c.Next()
	}
}

// CheckAuthMiddleware authenticates a restful api call and inject the userID and userType into to context if it exists
func (m *Middleware) CheckAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID := ""
		actorType := ""

		bearerToken := c.Request.Header.Get("Authorization")
		if len(bearerToken) < 10 {
			c.Set(ActorIDInContext, actorID)
			c.Set(ActorTypeInContext, actorType)
			c.Next()
			return
		}

		actorID, err := ParseToken(&m.env, strings.TrimPrefix(bearerToken, "Bearer "))
		if err != nil {
			restModel.ErrorResponse(c, http.StatusBadRequest, "unable to parse token")
			return
		}

		user := model.User{}
		ctx := context.WithValue(c.Request.Context(), helper.GinContextKey, c)

		// check if the actor is a user
		db := m.storage.DB.WithContext(ctx).Where("id = ?", actorID).First(&user)
		if db.Error != nil || strings.EqualFold(actorID, helper.ZeroUUID) {
			// might be an admin, check the admin table
			restModel.ErrorResponse(c, http.StatusBadRequest, ErrInvalidToken.Error())
			return
		}

		actorID = user.ID.String()
		actorType = string(model.ActorTypeUser)
		tenantID := user.TenantID.String()

		c.Set(ActorIDInContext, actorID)
		c.Set(ActorTypeInContext, actorType)

		c.Set(TenantIDInContext, tenantID)
		c.Next()
	}
}

func (m *Middleware) CheckTenantAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID := ""
		var actorType model.ActorType

		bearerToken := c.Request.Header.Get("Authorization")
		if len(bearerToken) < 10 {
			c.Set(ActorIDInContext, actorID)
			c.Set(ActorTypeInContext, actorType)
			c.Next()
			return
		}

		actorID, err := TenantParseToken(&m.env, strings.TrimPrefix(bearerToken, "Bearer "))
		if err != nil {
			restModel.ErrorResponse(c, http.StatusBadRequest, "unable to parse token")
			return
		}

		tenant := model.Tenant{}
		ctx := context.WithValue(c.Request.Context(), helper.GinContextKey, c)

		// check if the actor is a tenant
		db := m.storage.DB.WithContext(ctx).Where("id = ?", actorID).First(&tenant)
		if db.Error != nil || strings.EqualFold(actorID, helper.ZeroUUID) {
			restModel.ErrorResponse(c, http.StatusBadRequest, ErrInvalidToken.Error())
			return
		}

		actorID = tenant.ID.String()
		actorType = model.ActorTypeTenant

		c.Set(ActorIDInContext, actorID)
		c.Set(ActorTypeInContext, actorType)

		c.Next()
	}
}

// CorsMiddleware adds a CORS check for api rest endpoints allowing only list of origins defined in env
func (m *Middleware) CorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.DefaultConfig())
}

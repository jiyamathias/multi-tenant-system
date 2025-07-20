// Package auth contains every request that pertains to the users auth including the request object
package auth

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

type authHandler struct {
	logger      zerolog.Logger
	controller  controller.Operations
	environment *environment.Env
}

// New creates a new instance of the auth rest handler
func New(r *gin.RouterGroup, l zerolog.Logger, c controller.Operations, env *environment.Env) {
	auth := authHandler{
		logger:      l,
		controller:  c,
		environment: env,
	}

	authGroup := r.Group("/auth")

	authGroup.POST("/signup", auth.signup())
	authGroup.POST("/login", auth.login())
	authGroup.GET("/user/:id", auth.controller.Middleware().AuthMiddleware(), auth.getUserByID())
	authGroup.PATCH("/user", auth.controller.Middleware().AuthMiddleware(), auth.updateUserByID())

}

// signup 	godoc
//
//	@Summary		user signup
//	@Description	this endpoint signs up a new user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			signupRequest	body		signupRequest				true	"signup request body"
//	@Success		201				{object}	restModel.GenericResponse	"signup successful"
//	@Router			/auth/signup [post]
func (a *authHandler) signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request signupRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			a.logger.Error().Msgf("%v", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, "incomplete details please fill out the missing details")
			return
		}

		if err := restModel.ValidateRequest(request); err != nil {
			a.logger.Error().Msgf("%v", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, restModel.ErrIncompleteDetails.Error())
			return
		}

		newUser, err := a.controller.CreateUser(context.Background(), request.toUserModel())
		if err != nil {
			a.logger.Error().Msgf("%v", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		tokenDetails, err := a.controller.Middleware().CreateToken(c, a.environment, &newUser)
		if err != nil {
			a.logger.Err(err).Msgf("signup ::: Unable to generate token ==> %s", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		response := loginResponse{
			User:               newUser,
			AccessToken:        tokenDetails.AccessToken,
			AccessTokenExpiry:  tokenDetails.AccessTokenExpiry,
			RefreshToken:       tokenDetails.RefreshToken,
			RefreshTokenExpiry: tokenDetails.RefreshTokenExpiry,
		}

		restModel.OkResponse(c, http.StatusOK, "signup successful", response)
	}
}

// login 	godoc
//
//	@Summary		login
//	@Description	this endpoint is used to log a user in
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			loginRequest	body		loginRequest				true	"login request body"
//	@Success		200				{object}	restModel.GenericResponse	"user logged in successfully"
//	@Router			/auth/login [post]
func (a *authHandler) login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest

		// run the validation first
		if err := c.ShouldBindJSON(&req); err != nil {
			a.logger.Error().Msgf("%v", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		err := restModel.ValidateRequest(req)
		if err != nil {
			a.logger.Error().Msgf("%v", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, restModel.ErrIncompleteLoginDetails.Error())
			return
		}

		user, err := a.controller.AuthenticateUser(context.Background(), req.Email, req.Password)
		if err != nil {
			a.logger.Error().Msgf("%v", err)
			restModel.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		// create token
		tokenDetails, err := a.controller.Middleware().CreateToken(c, a.environment, &user)
		if err != nil {
			a.logger.Err(err).Msgf("Login ::: Unable to generate token ==> %s", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		response := loginResponse{
			User:               user,
			AccessToken:        tokenDetails.AccessToken,
			AccessTokenExpiry:  tokenDetails.AccessTokenExpiry,
			RefreshToken:       tokenDetails.RefreshToken,
			RefreshTokenExpiry: tokenDetails.RefreshTokenExpiry,
		}

		restModel.OkResponse(c, http.StatusOK, "user logged in successfully", response)
	}
}

// getUserByID 	godoc
//
//	@Summary		getUserByID
//	@Description	this endpoint gets a user by ID
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						false	"userID"
//	@Success		200	{object}	restModel.GenericResponse	"user details fetched successfully"
//	@Router			/auth/user/{id} [get]
func (a *authHandler) getUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		params := c.Param("id")
		if len(params) == 0 {
			a.logger.Err(helper.ErrUserIDParamsMissing).Msgf("getUserByID :::  ==> %s", helper.ErrUserIDParamsMissing)
			restModel.ErrorResponse(c, http.StatusBadRequest, helper.ErrUserIDParamsMissing.Error())
			return
		}

		userID, err := uuid.Parse(params)
		if err != nil {
			a.logger.Err(err).Msgf("getUserByID :::  ==> %s", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		user, err := a.controller.GetUserByID(context.Background(), userID)
		if err != nil {
			a.logger.Err(err).Msgf("getUserByID :::  ==> %s", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		restModel.OkResponse(c, http.StatusOK, "user details fetched successfully", user)
	}
}

// updateUserByID 	godoc
//
//	@Summary		updateUserByID
//	@Description	this endpoint is used to update any of the users record
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			updateUserRequest	body		updateUserRequest			true	"update user request body"
//	@Success		201					{object}	restModel.GenericResponse	"user updated successfully"
//	@Router			/auth/user [patch]
func (a *authHandler) updateUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input updateUserRequest

		// run the validation first
		if err := c.ShouldBindJSON(&input); err != nil {
			a.logger.Error().Msgf("%v", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		err := restModel.ValidateRequest(input)
		if err != nil {
			a.logger.Error().Msgf("%v", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		userID, err := uuid.Parse(c.GetString(middleware.ActorIDInContext))
		if err != nil {
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		userData := input.toUserModel(userID)

		updatedUser, err := a.controller.UpdateUserByID(context.Background(), userID, userData)
		if err != nil {
			a.logger.Err(err).Msgf("UpdateUserByID ::: Unable to update user ==> %s", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		restModel.OkResponse(c, http.StatusCreated, "user updated successfully", updatedUser)
	}
}

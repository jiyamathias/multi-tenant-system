package wallet

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

type walletHandler struct {
	logger      zerolog.Logger
	controller  controller.Operations
	environment *environment.Env
}

// New creates a new instance of the wallet rest handler
func New(r *gin.RouterGroup, l zerolog.Logger, c controller.Operations, env *environment.Env) {
	wallet := walletHandler{
		logger:      l,
		controller:  c,
		environment: env,
	}

	walletGroup := r.Group("/wallet")

	walletGroup.POST("", wallet.controller.Middleware().AuthMiddleware(), wallet.getWalletByUserID())
}

// getWalletByUserID 	godoc
//
//	@Summary		getWalletByUserID
//	@Description	this endpoint gets gets a users wallet balance
//	@Tags			wallet
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	restModel.GenericResponse	"wallet balance fetched successfully"
//	@Router			/wallet [get]
func (w *walletHandler) getWalletByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := uuid.Parse(c.GetString(middleware.ActorIDInContext))
		if err != nil {
			w.logger.Err(err).Msgf("getWalletByUserID ::: error parsing uuid ==> %s", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		wallet, err := w.controller.GetWalletByUserID(context.Background(), userID)
		if err != nil {
			w.logger.Err(err).Msgf("getWalletByUserID :::  ==> %s", err)
			restModel.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		restModel.OkResponse(c, http.StatusOK, "wallet balance fetched successfully", wallet)
	}
}

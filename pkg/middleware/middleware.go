// Package middleware contains methods that allows for authorization
package middleware

import (
	"crypto/rsa"

	ginJwt "github.com/appleboy/gin-jwt/v2"
	"github.com/rs/zerolog"

	"codematic/model"
	"codematic/pkg/environment"
	"codematic/pkg/helper"
	"codematic/storage"
)

const (
	// RequestBodyInContext context key holder
	RequestBodyInContext = "request_body_in_context"
	// ActorIDInContext context key holder
	ActorIDInContext = "actor_id_in_context"
	// TenantIDInContext context key holder
	TenantIDInContext = "tenant_id_in_context"
	// ActorTypeInContext context key holder
	ActorTypeInContext = "actor_type_in_context"
	// UserInContext context key holder
	UserInContext = "user_in_context"
	// packageName name of this package
	packageName = "middleware"
)

type (
	// ActorType is a type string
	ActorType string

	// Middleware object
	Middleware struct {
		logger  zerolog.Logger
		env     environment.Env
		jwt     *ginJwt.GinJWTMiddleware
		pKey    *rsa.PrivateKey
		storage *storage.Storage
	}
)

// NewMiddleware new instance of our custom ginJwt middleware
func NewMiddleware(z zerolog.Logger, env environment.Env, s *storage.Storage) *Middleware {
	mWare, _ := jwtMiddleware(&env, env.Get("JWT_ACCESS_TOKEN_SECRET"))
	l := z.With().Str(helper.LogStrKeyModule, packageName).Logger()
	return &Middleware{
		logger:  l,
		env:     env,
		jwt:     mWare,
		storage: s,
	}
}

// jwtMiddleware generates a JWT token
func jwtMiddleware(env *environment.Env, secretKey string) (*ginJwt.GinJWTMiddleware, error) {
	return ginJwt.New(&ginJwt.GinJWTMiddleware{
		Realm:      realm,
		Key:        []byte(secretKey),
		MaxRefresh: jwtAccessTokenExpiry(env),
		PayloadFunc: func(data any) ginJwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return ginJwt.MapClaims{
					identityKey: v.ID,
				}
			}
			return ginJwt.MapClaims{}
		},
		IdentityKey: identityKey,
		Timeout:     jwtAccessTokenExpiry(env),
	})
}

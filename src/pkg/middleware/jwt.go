package middleware

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	ginJwt "github.com/appleboy/gin-jwt/v2"
	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"codematic/model"
	"codematic/pkg/environment"
)

type (
	// Tokens object
	Tokens struct {
		AccessToken        string
		RefreshToken       string
		AccessTokenExpiry  string
		RefreshTokenExpiry string
	}
)

var (
	identityKey                = "id"
	realm                      = "codematic"
	claimsID                   = "id"
	claimsExpiry               = "exp"
	claimsCreatedAt            = "created_at"
	tenantID                   = "tenant_id"
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("token is invalid")
	ErrInvalidTenant           = errors.New("invalid tenant")
)

func jwtAccessTokenExpiry(env *environment.Env) time.Duration {
	ttl, err := strconv.Atoi(env.Get("JWT_ACCESS_TOKEN_EXPIRY"))
	if err != nil {
		return time.Hour * 4
	}
	return time.Minute * time.Duration(ttl)
}

func jwtRefreshTokenExpiry(env *environment.Env) time.Duration {
	ttl, err := strconv.Atoi(env.Get("JWT_REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		return time.Hour * 24
	}

	return time.Hour * time.Duration(ttl)
}

// CreateToken creates a new JWT token
func (m *Middleware) CreateToken(c *gin.Context, env *environment.Env, user *model.User) (*Tokens, error) {
	accessToken := jwtGo.New(jwtGo.GetSigningMethod(m.jwt.SigningAlgorithm))
	accessClaims := accessToken.Claims.(jwtGo.MapClaims)

	refreshToken := jwtGo.New(jwtGo.GetSigningMethod(m.jwt.SigningAlgorithm))
	refreshClaims := refreshToken.Claims.(jwtGo.MapClaims)

	if m.jwt.PayloadFunc != nil {
		for key, value := range m.jwt.PayloadFunc(user) {
			accessClaims[key] = value
			refreshClaims[key] = value
		}
	}

	accessExpire := time.Now().Add(jwtAccessTokenExpiry(env))
	refreshExpire := time.Now().Add(jwtRefreshTokenExpiry(env))

	accessClaims[claimsID] = user.ID
	accessClaims[tenantID] = user.TenantID
	accessClaims[claimsExpiry] = accessExpire.Unix()
	accessClaims[claimsCreatedAt] = m.jwt.TimeFunc().Unix()

	refreshClaims[claimsID] = user.ID
	refreshClaims[tenantID] = user.TenantID
	refreshClaims[claimsExpiry] = refreshExpire.Unix()
	refreshClaims[claimsCreatedAt] = m.jwt.TimeFunc().Unix()

	accessTokenString, err := m.signedString(accessToken)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := m.signedString(refreshToken)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:        accessTokenString,
		RefreshToken:       refreshTokenString,
		AccessTokenExpiry:  accessExpire.String(),
		RefreshTokenExpiry: refreshExpire.String(),
	}, err
}

// CreateTenantToken creates a new JWT token
func (m *Middleware) CreateTenantToken(c *gin.Context, env *environment.Env, tenant *model.Tenant) (*Tokens, error) {
	accessToken := jwtGo.New(jwtGo.GetSigningMethod(m.jwt.SigningAlgorithm))
	accessClaims := accessToken.Claims.(jwtGo.MapClaims)

	refreshToken := jwtGo.New(jwtGo.GetSigningMethod(m.jwt.SigningAlgorithm))
	refreshClaims := refreshToken.Claims.(jwtGo.MapClaims)

	if m.jwt.PayloadFunc != nil {
		for key, value := range m.jwt.PayloadFunc(tenant) {
			accessClaims[key] = value
			refreshClaims[key] = value
		}
	}

	accessExpire := time.Now().Add(jwtAccessTokenExpiry(env))
	refreshExpire := time.Now().Add(jwtRefreshTokenExpiry(env))

	accessClaims[tenantID] = tenant.ID
	accessClaims[claimsExpiry] = accessExpire.Unix()
	accessClaims[claimsCreatedAt] = m.jwt.TimeFunc().Unix()

	refreshClaims[tenantID] = tenant.ID
	refreshClaims[claimsExpiry] = refreshExpire.Unix()
	refreshClaims[claimsCreatedAt] = m.jwt.TimeFunc().Unix()

	accessTokenString, err := m.signedString(accessToken)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := m.signedString(refreshToken)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:        accessTokenString,
		RefreshToken:       refreshTokenString,
		AccessTokenExpiry:  accessExpire.String(),
		RefreshTokenExpiry: refreshExpire.String(),
	}, err
}

// GetGinJWTMiddleware returns GinJWTMiddleware
func (m *Middleware) GetGinJWTMiddleware() *ginJwt.GinJWTMiddleware {
	return m.jwt
}

func (m *Middleware) signedString(token *jwtGo.Token) (string, error) {
	var tokenString string
	var err error
	if m.usingPublicKeyAlgo() {
		tokenString, err = token.SignedString(m.pKey)
	} else {
		tokenString, err = token.SignedString(m.jwt.Key)
	}
	return tokenString, err
}

func (m *Middleware) usingPublicKeyAlgo() bool {
	switch m.jwt.SigningAlgorithm {
	case "RS256", "RS512", "RS384":
		return true
	}
	return false
}

// ValidateRefreshToken validates refresh token
func (m *Middleware) ValidateRefreshToken(z zerolog.Logger, c *gin.Context, env *environment.Env, token string) (*uuid.UUID, error) {
	tokenGotten, err := jwtGo.Parse(token, func(token *jwtGo.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwtGo.SigningMethodHMAC); !ok {
			z.Error().Msgf("RefreshToken unexpected signing method: (%v)", token.Header["alg"])

			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(env.Get("JWT_REFRESH_TOKEN_SECRET")), nil
	})

	//any error may be due to token expiration
	if err != nil {
		z.Err(err).Msgf("RefreshToken error: %v", err)
		return nil, err
	}

	//is token valid?
	if err = tokenGotten.Claims.Valid(); err != nil {
		z.Err(err).Msgf("RefreshToken error: %v", err)
		return nil, err
	}

	claims, ok := tokenGotten.Claims.(jwtGo.MapClaims)
	claimsUUID := claims[claimsID].(string)
	//get the last refresh token for this user/customer
	refreshTokenCookie, err := c.Cookie(claimsUUID)
	//error may be due to cookie expiration OR a new refresh token has been generated
	if err != nil || refreshTokenCookie != token {
		z.Err(err).Msgf("RefreshToken:Cookie error: %v", err)
		return nil, ErrInvalidToken
	}

	if ok && tokenGotten.Valid {
		//convert the interface to uuid.UUID
		parsedUUID, err := uuid.Parse(claimsUUID)
		if err != nil {
			z.Err(err).Msgf("RefreshToken: Invalid user (%v)", err)
			return nil, err
		}

		return &parsedUUID, nil
	}

	return nil, ErrInvalidToken
}

func (m *Middleware) ValidateRefreshTenantToken(z zerolog.Logger, c *gin.Context, env *environment.Env, token string) (*uuid.UUID, error) {
	tokenGotten, err := jwtGo.Parse(token, func(token *jwtGo.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwtGo.SigningMethodHMAC); !ok {
			z.Error().Msgf("RefreshToken unexpected signing method: (%v)", token.Header["alg"])

			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(env.Get("JWT_REFRESH_TOKEN_SECRET")), nil
	})

	//any error may be due to token expiration
	if err != nil {
		z.Err(err).Msgf("RefreshToken error: %v", err)
		return nil, err
	}

	//is token valid?
	if err = tokenGotten.Claims.Valid(); err != nil {
		z.Err(err).Msgf("RefreshToken error: %v", err)
		return nil, err
	}

	claims, ok := tokenGotten.Claims.(jwtGo.MapClaims)
	claimsUUID := claims[tenantID].(string)
	//get the last refresh token for this user/customer
	refreshTokenCookie, err := c.Cookie(claimsUUID)
	//error may be due to cookie expiration OR a new refresh token has been generated
	if err != nil || refreshTokenCookie != token {
		z.Err(err).Msgf("RefreshToken:Cookie error: %v", err)
		return nil, ErrInvalidToken
	}

	if ok && tokenGotten.Valid {
		//convert the interface to uuid.UUID
		parsedUUID, err := uuid.Parse(claimsUUID)
		if err != nil {
			z.Err(err).Msgf("RefreshToken: Invalid user (%v)", err)
			return nil, err
		}

		return &parsedUUID, nil
	}

	return nil, ErrInvalidToken
}

// ParseToken parses the JWT token
func ParseToken(env *environment.Env, tokenStr string) (jwtGo.MapClaims, error) {
	token, err := jwtGo.Parse(tokenStr, func(token *jwtGo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtGo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(env.Get("JWT_ACCESS_TOKEN_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwtGo.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// TenantParseToken parses the JWT token
func TenantParseToken(env *environment.Env, tokenStr string) (string, error) {
	token, err := jwtGo.Parse(tokenStr, func(token *jwtGo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtGo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(env.Get("JWT_ACCESS_TOKEN_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwtGo.MapClaims); ok && token.Valid {
		tenantID := claims[tenantID].(string)
		return tenantID, nil
	}

	return "", err
}

// JwtAuthorization retrieves the user ID from a JWT claims
func (m *Middleware) JwtAuthorization(c *gin.Context) (*uuid.UUID, error) {
	claims, err := m.jwt.GetClaimsFromJWT(c)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(claims[claimsID].(string))
	if err != nil {
		return nil, err
	}

	return &userID, nil
}

func (m *Middleware) TenantJwtAuthorization(c *gin.Context) (*uuid.UUID, error) {
	claims, err := m.jwt.GetClaimsFromJWT(c)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(claims[tenantID].(string))
	if err != nil {
		return nil, err
	}

	return &userID, nil
}

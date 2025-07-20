package auth

import (
	"time"

	"github.com/google/uuid"

	"codematic/model"
)

type (
	signupRequest struct {
		FirstName string `json:"firstName" validate:"required"`
		LastName  string `json:"lastName" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required"`
	}

	updateUserRequest struct {
		FirstName *string `json:"firstName"`
		LastName  *string `json:"lastName"`
	}

	loginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	loginResponse struct {
		User               model.User `json:"user"`
		AccessToken        string     `json:"accessToken"`
		AccessTokenExpiry  string     `json:"accessTokenExpiry"`
		RefreshToken       string     `json:"refreshToken"`
		RefreshTokenExpiry string     `json:"refreshTokenExpiry"`
	}
)

func (s *signupRequest) toUserModel() model.User {
	password := model.Password(s.Password)

	u := model.User{
		Email:     s.Email,
		Password:  password,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	return u
}

func (u *updateUserRequest) toUserModel(userID uuid.UUID) model.User {
	var user model.User

	if u.FirstName != nil {
		user.FirstName = *u.FirstName
	}
	if u.LastName != nil {
		user.LastName = *u.LastName
	}

	user.ID = userID

	return user
}

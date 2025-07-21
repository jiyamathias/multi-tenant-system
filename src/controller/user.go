package controller

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"codematic/model"
	"codematic/model/pagination"
	"codematic/pkg/helper"
)

// CreateUser creates a new user in the database
func (c *Controller) CreateUser(ctx context.Context, u model.User) (model.User, error) {
	u.Email = strings.ToLower(u.Email)

	encryptedPass := u.Password.Encrypt()
	u.Password = encryptedPass

	// validate if the tenantID was passed in
	if u.TenantID == uuid.Nil {
		return model.User{}, helper.ErrTenantIDMissing
	}

	_, err := c.userStorage.GetUserByEmail(ctx, u.Email)
	if err == nil {
		c.logger.Err(err).Msgf("CreateUser::: Email already taken %s", err)
		return model.User{}, ErrEmailAlreadyExists
	}

	newUser, err := c.userStorage.CreateUser(ctx, u)
	if err != nil {
		c.logger.Err(err).Msgf("CreateUser::: Unable to insert user into db %s", err)
		return model.User{}, err
	}

	// create a wallet for every newly created users
	wallet := model.Wallet{
		ID:            uuid.New(),
		UserID:        newUser.ID,
		BalanceBefore: 0.0,
		BalanceAfter:  0.0,
	}

	_, err = c.CreateWallet(ctx, wallet)
	if err != nil {
		c.logger.Err(err).Msgf("CreateWallet::: Unable to create wallet %s", err)
		return model.User{}, err
	}

	auditLog := model.AuditLog{
		ID:         uuid.New(),
		TenantID:   &newUser.TenantID,
		UserID:     &newUser.ID,
		Actor:      model.ActorTenant,
		ActionDone: model.ActionSuccess,
		Messages:   "tenant added a new user",
	}

	if _, err = c.CreateAuditLog(ctx, auditLog); err != nil {
		c.logger.Err(err).Msgf("error creating audit log")
		return model.User{}, err
	}

	return newUser, nil
}

// GetUserByEmail returns a model.User object associated with the email address specified
func (c *Controller) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	user, err := c.userStorage.GetUserByEmail(ctx, email)
	if err != nil {
		c.logger.Err(err).Msgf("GetUserByEmail::: Unable to fetch user by given email %s ", err)
		return model.User{}, err
	}
	return user, nil
}

// AuthenticateUser returns a model.User object of matching email and password(not hashed password) else returns an error
func (c *Controller) AuthenticateUser(ctx context.Context, email, password string) (model.User, error) {
	user, err := c.userStorage.GetUserByEmail(ctx, strings.ToLower(email))
	if err != nil {
		c.logger.Err(err).Msgf("AuthenticateUser::: Unable to fetch user details %s", err)
		return model.User{}, ErrIncorrectLoginDetails
	}

	// check password hash
	if ok := user.Password.Check(model.Password(password)); !ok {
		return model.User{}, ErrIncorrectLoginDetails
	}

	auditLog := model.AuditLog{
		ID:         uuid.New(),
		TenantID:   &user.TenantID,
		UserID:     &user.ID,
		Actor:      model.ActorUser,
		ActionDone: model.ActionSuccess,
		Messages:   "user logged in",
	}

	if _, err = c.CreateAuditLog(ctx, auditLog); err != nil {
		c.logger.Err(err).Msgf("error creating audit log")
		return model.User{}, err
	}

	return user, nil
}

// GetUserByID fetch a user matching the ID
func (c *Controller) GetUserByID(ctx context.Context, userID uuid.UUID) (model.User, error) {
	// get user by ID
	return c.userStorage.GetUserByID(ctx, userID)
}

// UpdateUserByID updates a users record
func (c *Controller) UpdateUserByID(ctx context.Context, userID uuid.UUID, u model.User) (model.User, error) {
	// get user by ID
	user, err := c.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		c.logger.Err(err).Msgf("GetUserById::: could not update user %s", err)
		return model.User{}, err
	}

	u.UpdatedAt = time.Now()

	// update user by ID
	user, err = c.userStorage.UpdateUserByID(ctx, userID, u)
	if err != nil {
		c.logger.Err(err).Msgf("UpdateUser::: Could not get user details %s", err)
		return model.User{}, nil
	}

	user, err = c.userStorage.GetUserByID(ctx, u.ID)
	if err != nil {
		c.logger.Err(err).Msgf("UpdateUser::: Could not get user details %s", err)
		return model.User{}, nil
	}

	return user, nil
}

// GetAllUsersByTenantID gets all users signed up under a particular tenant
func (c *Controller) GetAllUsersByTenantID(ctx context.Context, tenantId uuid.UUID, page pagination.Page) ([]*model.User, pagination.PageInfo, error) {
	return c.userStorage.GetAllUsersByTenantID(ctx, tenantId, page)
}

package controller

import (
	"context"

	"github.com/google/uuid"

	"codematic/model"
)

// CreateBalance creates a new balance information
func (c *Controller) CreateBalance(ctx context.Context, balance model.Balance) (model.Balance, error) {
	newBalance, err := c.balanceStorage.CreateBalance(ctx, balance)
	if err != nil {
		c.logger.Err(err).Msgf("CreateBalance:: unable to insert new balance information in db %s", err)
		return model.Balance{}, err
	}

	return newBalance, nil
}

func (c *Controller) GetLastBalanceByUserID(ctx context.Context, userID uuid.UUID) (model.Balance, error) {
	user, err := c.GetUserByID(ctx, userID)
	if err != nil {
		c.logger.Error().Msgf("GetTransactionByTypeOfFlow ::: error retrieving user => %v", err)
		return model.Balance{}, err
	}

	balance, err := c.balanceStorage.GetLastBalanceByUserID(ctx, user.ID)
	if err != nil {
		c.logger.Err(err).Msgf("GetBalanceByID failed: unable to fetch record %s", err)
		return balance, err
	}

	return balance, err
}

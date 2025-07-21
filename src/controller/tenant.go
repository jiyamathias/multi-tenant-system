package controller

import (
	"context"
	"strings"

	"codematic/model"
)

func (c *Controller) CreateTenant(ctx context.Context, tenant model.Tenant) (model.Tenant, error) {
	tenant.Email = strings.ToLower(tenant.Email)

	encryptedPass := tenant.Password.Encrypt()
	tenant.Password = encryptedPass

	_, err := c.tenantStorage.GetTenantByEmail(ctx, tenant.Email)
	if err == nil {
		c.logger.Err(err).Msgf("GetTenantByEmail::: Email already taken %s", err)
		return model.Tenant{}, ErrEmailAlreadyExists
	}

	newTenant, err := c.tenantStorage.CreateTenant(ctx, tenant)
	if err != nil {
		c.logger.Err(err).Msgf("CreateTenant::: Unable to insert tenant into db %s", err)
		return model.Tenant{}, err
	}

	return newTenant, nil
}

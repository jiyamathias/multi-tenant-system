package tenant

import (
	"codematic/model"

	"github.com/google/uuid"
)

type (
	tenantRequest struct {
		BusinessName string `json:"businessName" validate:"required"`
		Email        string `json:"email" validate:"required"`
		Password     string `json:"password" validate:"required"`
	}
)

func (t *tenantRequest) toModel() model.Tenant {
	password := model.Password(t.Password)

	return model.Tenant{
		ID:           uuid.New(),
		BusinessName: t.BusinessName,
		Email:        t.Email,
		Password:     password,
	}
}

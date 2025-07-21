package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"codematic/model"
	"codematic/model/pagination"
	"codematic/pkg/helper"
)

// UserDatabase enlist all possible storage operations for the User
//
//go:generate mockgen -source user.go -destination ./mock/mock_user.go -package mock UserDatabase
type UserDatabase interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	UpdateUserByID(ctx context.Context, userID uuid.UUID, user model.User) (model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	UpdateLastLoggedIn(ctx context.Context, email string, when time.Time) error

	GetAllUsersByTenantID(ctx context.Context, tenantId uuid.UUID, page pagination.Page) ([]*model.User, pagination.PageInfo, error)
}

// User object
type User struct {
	logger  zerolog.Logger
	storage *Storage
}

// NewUser creates a new instance of the user
func NewUser(s *Storage) *UserDatabase {
	l := s.Logger.With().Str(helper.LogStrKeyLevel, "user").Logger()
	user := &User{
		logger:  l,
		storage: s,
	}
	userDatabase := UserDatabase(user)
	return &userDatabase
}

// CreateUser adds a new row into the user table
func (u *User) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	db := u.storage.DB.WithContext(ctx).Model(&model.User{}).Create(&user)
	if db.Error != nil {
		u.logger.Err(db.Error).Msgf("User::CreateUser error: %v, (%v)", ErrRecordCreatingFailed, db.Error)
		if strings.Contains(db.Error.Error(), "duplicate key value") {
			return model.User{}, ErrDuplicateRecord
		}
		return model.User{}, ErrRecordCreatingFailed
	}
	return user, nil
}

// UpdateUserByID should update the user record in the storage
func (u *User) UpdateUserByID(ctx context.Context, userID uuid.UUID, user model.User) (model.User, error) {
	db := u.storage.DB.WithContext(ctx).Model(model.User{
		ID: userID,
	}).UpdateColumns(user)

	if db.Error != nil {
		u.logger.Err(db.Error).Msgf("UserStorage ::: UpdateByID error: %v, (%v)", ErrRecordUpdateFailed, db.Error)
		return user, ErrRecordUpdateFailed
	}

	return user, nil
}

// GetUserByID returns a user matching the ID
func (u *User) GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	var user model.User
	db := u.storage.DB.WithContext(ctx).Where("id = ?", id).First(&user)
	if db.Error != nil || strings.EqualFold(user.ID.String(), helper.ZeroUUID) {
		u.logger.Err(db.Error).Msgf("User::GetUserByID error: %v (%v)", ErrRecordNotFound, db.Error)
		return user, ErrRecordNotFound
	}
	return user, nil
}

// GetUserByEmail returns a user matching the email address
func (u *User) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	db := u.storage.DB.WithContext(ctx).Where("email = ?", email).First(&user)
	if db.Error != nil || strings.EqualFold(user.ID.String(), helper.ZeroUUID) {
		u.logger.Err(db.Error).Msgf("User::GetUserByEmail error: %v (%v)", ErrRecordNotFound, db.Error)
		return user, ErrRecordNotFound
	}
	return user, nil
}

// UpdateLastLoggedIn updated the user's last login
func (u *User) UpdateLastLoggedIn(ctx context.Context, email string, when time.Time) error {
	db := u.storage.DB.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).
		Updates(map[string]interface{}{
			"last_logged_in_at": when,
			"updated_at":        when,
		})
	if db.Error != nil {
		u.logger.Err(db.Error).Msgf("User::UpdateLastLoggedIn error: %v (%v)", ErrRecordUpdateFailed, db.Error)
		return ErrRecordUpdateFailed
	}

	return nil
}

// GetAllUsersByTenantID gets all users signed up under a particular tenant
func (u *User) GetAllUsersByTenantID(ctx context.Context, tenantId uuid.UUID, page pagination.Page) ([]*model.User, pagination.PageInfo, error) {
	var users []*model.User

	offset := 0
	// load defaults
	if page.Number == nil {
		tmpPageNumber := pagination.PageDefaultNumber
		page.Number = &tmpPageNumber
	}
	if page.Size == nil {
		tmpPageSize := pagination.PageDefaultSize
		page.Size = &tmpPageSize
	}
	if page.SortBy == nil {
		tmpPageSortBy := pagination.PageDefaultSortBy
		page.SortBy = &tmpPageSortBy
	}

	if page.SortDirectionDesc == nil {
		tmpPageSortDirectionDesc := pagination.PageDefaultSortDirectionDesc
		page.SortDirectionDesc = &tmpPageSortDirectionDesc
	}

	if *page.Number > 1 {
		offset = *page.Size * (*page.Number - 1)
	}
	sortDirection := pagination.PageSortDirectionDescending
	if !*page.SortDirectionDesc {
		sortDirection = pagination.PageSortDirectionAscending
	}

	queryDraft := u.storage.DB.WithContext(ctx).Where("tenant_id = ?", tenantId)

	var count int64
	queryDraft.Model(model.User{}).Count(&count)

	db := queryDraft.Debug().Offset(offset).Limit(*page.Size).
		Order(fmt.Sprintf("%s %s", *page.SortBy, sortDirection)).
		Find(&users)

	if db.Error != nil {
		u.logger.Err(db.Error).Msgf("GetAllUsersByTenantID:: error: %v, (%v)", ErrEmptyResult, db.Error)
		return nil, pagination.PageInfo{}, ErrEmptyResult
	}

	return users, pagination.PageInfo{
		Page:            *page.Number,
		Size:            *page.Size,
		HasNextPage:     int64(offset+*page.Size) < count,
		HasPreviousPage: *page.Number > 1,
		TotalCount:      count,
	}, nil
}

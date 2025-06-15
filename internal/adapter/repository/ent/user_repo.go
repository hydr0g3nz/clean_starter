package repository

import (
	"context"
	"time"

	"github.com/hydr0g3nz/wallet_topup_system/internal/adapter/ent"
	"github.com/hydr0g3nz/wallet_topup_system/internal/adapter/ent/user"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/entity"
	i "github.com/hydr0g3nz/wallet_topup_system/internal/domain/repository"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/vo"
)

// userRepository implements UserRepository interface
type userRepository struct {
	client *ent.Client
}

// NewUserRepository creates a new user repository
func NewUserRepository(client *ent.Client) i.UserRepository {
	return &userRepository{
		client: client,
	}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, userEntity *entity.User) (*entity.User, error) {
	entUser, err := r.client.User.
		Create().
		SetEmail(userEntity.Email).
		SetPasswordHash(userEntity.PasswordHash).
		SetRole(user.Role(userEntity.Role)).
		SetIsActive(userEntity.IsActive).
		SetEmailVerified(userEntity.EmailVerified).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return r.entToEntity(entUser), nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id int) (*entity.User, error) {
	entUser, err := r.client.User.
		Query().
		Where(user.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil // or custom not found error
		}
		return nil, err
	}

	return r.entToEntity(entUser), nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	entUser, err := r.client.User.
		Query().
		Where(user.Email(email)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil // or custom not found error
		}
		return nil, err
	}

	return r.entToEntity(entUser), nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, userEntity *entity.User) (*entity.User, error) {
	updateQuery := r.client.User.
		UpdateOneID(userEntity.ID).
		SetEmail(userEntity.Email).
		SetPasswordHash(userEntity.PasswordHash).
		SetRole(user.Role(userEntity.Role)).
		SetIsActive(userEntity.IsActive).
		SetEmailVerified(userEntity.EmailVerified).
		SetUpdatedAt(time.Now())

	// Set last login if provided
	if userEntity.LastLoginAt != nil {
		updateQuery = updateQuery.SetLastLoginAt(*userEntity.LastLoginAt)
	}

	entUser, err := updateQuery.Save(ctx)
	if err != nil {
		return nil, err
	}

	return r.entToEntity(entUser), nil
}

// Delete deletes a user by ID
func (r *userRepository) Delete(ctx context.Context, id int) error {
	return r.client.User.
		DeleteOneID(id).
		Exec(ctx)
}

// List retrieves users with pagination
func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	entUsers, err := r.client.User.
		Query().
		Order(ent.Desc(user.FieldCreatedAt)).
		Limit(limit).
		Offset(offset).
		All(ctx)

	if err != nil {
		return nil, err
	}

	return r.entUsersToEntities(entUsers), nil
}

// ListByRole retrieves users by role with pagination
func (r *userRepository) ListByRole(ctx context.Context, role string, limit, offset int) ([]*entity.User, error) {
	entUsers, err := r.client.User.
		Query().
		Where(user.RoleEQ(user.Role(role))).
		Order(ent.Desc(user.FieldCreatedAt)).
		Limit(limit).
		Offset(offset).
		All(ctx)

	if err != nil {
		return nil, err
	}

	return r.entUsersToEntities(entUsers), nil
}

// UpdateLastLogin updates the last login timestamp
func (r *userRepository) UpdateLastLogin(ctx context.Context, id int) error {
	return r.client.User.
		UpdateOneID(id).
		SetLastLoginAt(time.Now()).
		Exec(ctx)
}

// entToEntity converts ent.User to domain entity.User
func (r *userRepository) entToEntity(entUser *ent.User) *entity.User {
	userEntity := &entity.User{
		ID:            entUser.ID,
		Email:         entUser.Email,
		PasswordHash:  entUser.PasswordHash,
		Role:          vo.UserRole(entUser.Role),
		IsActive:      entUser.IsActive,
		EmailVerified: entUser.EmailVerified,
		CreatedAt:     entUser.CreatedAt,
		UpdatedAt:     entUser.UpdatedAt,
	}

	// Handle nullable LastLoginAt
	if entUser.LastLoginAt != nil {
		userEntity.LastLoginAt = entUser.LastLoginAt
	}

	return userEntity
}

// entityToEnt converts domain entity.User to ent.User (helper method if needed)
func (r *userRepository) entityToEnt(userEntity *entity.User) *ent.User {
	entUser := &ent.User{
		ID:            userEntity.ID,
		Email:         userEntity.Email,
		PasswordHash:  userEntity.PasswordHash,
		Role:          user.Role(userEntity.Role),
		IsActive:      userEntity.IsActive,
		EmailVerified: userEntity.EmailVerified,
		CreatedAt:     userEntity.CreatedAt,
		UpdatedAt:     userEntity.UpdatedAt,
	}

	// Handle nullable LastLoginAt
	if userEntity.LastLoginAt != nil {
		entUser.LastLoginAt = userEntity.LastLoginAt
	}

	return entUser
}

// entUsersToEntities converts slice of ent.User to slice of domain entity.User
func (r *userRepository) entUsersToEntities(entUsers []*ent.User) []*entity.User {
	entities := make([]*entity.User, len(entUsers))
	for i, entUser := range entUsers {
		entities[i] = r.entToEntity(entUser)
	}
	return entities
}

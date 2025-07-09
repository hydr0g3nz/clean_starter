package repository

import (
	"context"
	"database/sql"
	"errors"

	sqlc "github.com/hydr0g3nz/clean_stater/internal/adapter/sqlc/generated"
	"github.com/hydr0g3nz/clean_stater/internal/domain/entity"
	"github.com/hydr0g3nz/clean_stater/internal/domain/repository"
	"github.com/hydr0g3nz/clean_stater/internal/domain/vo"
	"github.com/hydr0g3nz/clean_stater/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &userRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	dbUser, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		Role:          sqlc.UserRole(user.Role.String()),
		IsActive:      utils.ConvertToPGBool(user.IsActive),
		EmailVerified: utils.ConvertToPGBool(user.EmailVerified),
	})
	if err != nil {
		return nil, err
	}

	return r.dbUserToEntity(dbUser)
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*entity.User, error) {
	dbUser, err := r.queries.GetUserByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.dbUserToEntity(dbUser)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.dbUserToEntity(dbUser)
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	dbUser, err := r.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:            int32(user.ID),
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		Role:          sqlc.UserRole(user.Role.String()),
		IsActive:      utils.ConvertToPGBool(user.IsActive),
		EmailVerified: utils.ConvertToPGBool(user.EmailVerified),
	})
	if err != nil {
		return nil, err
	}

	return r.dbUserToEntity(dbUser)
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	return r.queries.DeleteUser(ctx, int32(id))
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	dbUsers, err := r.queries.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	return r.dbUsersToEntities(dbUsers)
}

func (r *userRepository) ListByRole(ctx context.Context, role string, limit, offset int) ([]*entity.User, error) {
	dbUsers, err := r.queries.ListUsersByRole(ctx, sqlc.ListUsersByRoleParams{
		Role:   sqlc.UserRole(role),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	return r.dbUsersToEntities(dbUsers)
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id int) error {
	return r.queries.UpdateUserLastLogin(ctx, int32(id))
}

// Helper methods for conversion
func (r *userRepository) dbUserToEntity(dbUser *sqlc.User) (*entity.User, error) {
	role, err := vo.ParseUserRole(string(dbUser.Role))
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:            int(dbUser.ID),
		Email:         dbUser.Email,
		PasswordHash:  dbUser.PasswordHash,
		Role:          role,
		IsActive:      utils.ConvertToBool(dbUser.IsActive),
		EmailVerified: utils.ConvertToBool(dbUser.EmailVerified),
		CreatedAt:     dbUser.CreatedAt.Time,
		UpdatedAt:     dbUser.UpdatedAt.Time,
	}

	if dbUser.LastLoginAt.Valid {
		user.LastLoginAt = &dbUser.LastLoginAt.Time
	}

	return user, nil
}

func (r *userRepository) dbUsersToEntities(dbUsers []*sqlc.User) ([]*entity.User, error) {
	entities := make([]*entity.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		entity, err := r.dbUserToEntity(dbUser)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}

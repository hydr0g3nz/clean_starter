package usecase

import (
	"context"
	"fmt"

	"github.com/hydr0g3nz/wallet_topup_system/config"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/entity"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/infra"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/repository"
	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/vo"
	"golang.org/x/crypto/bcrypt"
)

// UserUsecase interface defines the contract for user business logic
type UserUsecase interface {
	Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	GetProfile(ctx context.Context, userID int) (*UserResponse, error)
	UpdateProfile(ctx context.Context, userID int, req *UpdateProfileRequest) (*UserResponse, error)
	ChangePassword(ctx context.Context, userID int, req *ChangePasswordRequest) error
	VerifyEmail(ctx context.Context, userID int) error
	DeactivateUser(ctx context.Context, userID int) error
	ActivateUser(ctx context.Context, userID int) error
	GetUsersByRole(ctx context.Context, role string, limit, offset int) ([]*UserResponse, error)
	DeleteUser(ctx context.Context, userID int) error
}

// userUsecase implements UserUsecase interface
type userUsecase struct {
	userRepo repository.UserRepository
	logger   infra.Logger
	config   *config.Config
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(
	userRepo repository.UserRepository,
	logger infra.Logger,
	config *config.Config,
) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		logger:   logger,
		config:   config,
	}
}

// Register creates a new user account
func (u *userUsecase) Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error) {
	u.logger.Info("Starting user registration", "email", req.Email, "role", req.Role)

	// Check if user already exists
	existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		u.logger.Error("Error checking existing user", "error", err, "email", req.Email)
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		u.logger.Warn("User already exists", "email", req.Email)
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := u.hashPassword(req.Password)
	if err != nil {
		u.logger.Error("Error hashing password", "error", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	role, err := vo.ParseUserRole(req.Role)
	if err != nil {
		u.logger.Error("Invalid user role", "error", err, "role", req.Role)
		return nil, err
	}

	// Create user entity
	user := &entity.User{
		Email:         req.Email,
		PasswordHash:  hashedPassword,
		Role:          role,
		IsActive:      true,
		EmailVerified: false,
	}

	// Save to database
	createdUser, err := u.userRepo.Create(ctx, user)
	if err != nil {
		u.logger.Error("Error creating user", "error", err, "email", req.Email)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	u.logger.Info("User registered successfully", "userID", createdUser.ID, "email", createdUser.Email)

	return u.toUserResponse(createdUser), nil
}

// Login authenticates a user and updates last login
func (u *userUsecase) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	u.logger.Info("User login attempt", "email", req.Email)

	// Get user by email
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		u.logger.Error("Error getting user", "error", err, "email", req.Email)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		u.logger.Warn("User not found", "email", req.Email)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		u.logger.Warn("Inactive user login attempt", "userID", user.ID, "email", req.Email)
		return nil, fmt.Errorf("user account is deactivated")
	}

	// Verify password
	if !u.verifyPassword(req.Password, user.PasswordHash) {
		u.logger.Warn("Invalid password", "userID", user.ID, "email", req.Email)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Update last login
	if err := u.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		u.logger.Error("Error updating last login", "error", err, "userID", user.ID)
		// Don't fail login for this error
	}

	u.logger.Info("User logged in successfully", "userID", user.ID, "email", user.Email)

	return &LoginResponse{
		User:  u.toUserResponse(user),
		Token: "jwt_token_placeholder", // You'll need to implement JWT generation
	}, nil
}

// GetProfile retrieves user profile
func (u *userUsecase) GetProfile(ctx context.Context, userID int) (*UserResponse, error) {
	u.logger.Debug("Getting user profile", "userID", userID)

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		u.logger.Error("Error getting user profile", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	if user == nil {
		u.logger.Warn("User not found", "userID", userID)
		return nil, fmt.Errorf("user not found")
	}

	return u.toUserResponse(user), nil
}

// UpdateProfile updates user profile information
func (u *userUsecase) UpdateProfile(ctx context.Context, userID int, req *UpdateProfileRequest) (*UserResponse, error) {
	u.logger.Info("Updating user profile", "userID", userID)

	// Get current user
	currentUser, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		u.logger.Error("Error getting current user", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if currentUser == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if email is being changed and if it's unique
	if req.Email != "" && req.Email != currentUser.Email {
		existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
		if err != nil {
			u.logger.Error("Error checking email uniqueness", "error", err, "email", req.Email)
			return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
		}
		if existingUser != nil {
			return nil, fmt.Errorf("email already in use")
		}
		currentUser.Email = req.Email
		currentUser.EmailVerified = false // Reset verification if email changed
	}

	// Update user
	updatedUser, err := u.userRepo.Update(ctx, currentUser)
	if err != nil {
		u.logger.Error("Error updating user", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	u.logger.Info("User profile updated successfully", "userID", userID)

	return u.toUserResponse(updatedUser), nil
}

// ChangePassword changes user password
func (u *userUsecase) ChangePassword(ctx context.Context, userID int, req *ChangePasswordRequest) error {
	u.logger.Info("Changing user password", "userID", userID)

	// Get current user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		u.logger.Error("Error getting user", "error", err, "userID", userID)
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Verify current password
	if !u.verifyPassword(req.CurrentPassword, user.PasswordHash) {
		u.logger.Warn("Invalid current password", "userID", userID)
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := u.hashPassword(req.NewPassword)
	if err != nil {
		u.logger.Error("Error hashing new password", "error", err, "userID", userID)
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.PasswordHash = hashedPassword
	_, err = u.userRepo.Update(ctx, user)
	if err != nil {
		u.logger.Error("Error updating password", "error", err, "userID", userID)
		return fmt.Errorf("failed to update password: %w", err)
	}

	u.logger.Info("Password changed successfully", "userID", userID)
	return nil
}

// VerifyEmail marks user email as verified
func (u *userUsecase) VerifyEmail(ctx context.Context, userID int) error {
	u.logger.Info("Verifying user email", "userID", userID)

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	user.EmailVerified = true
	_, err = u.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	u.logger.Info("Email verified successfully", "userID", userID)
	return nil
}

// DeactivateUser deactivates a user account
func (u *userUsecase) DeactivateUser(ctx context.Context, userID int) error {
	u.logger.Info("Deactivating user", "userID", userID)

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	user.IsActive = false
	_, err = u.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	u.logger.Info("User deactivated successfully", "userID", userID)
	return nil
}

// ActivateUser activates a user account
func (u *userUsecase) ActivateUser(ctx context.Context, userID int) error {
	u.logger.Info("Activating user", "userID", userID)

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	user.IsActive = true
	_, err = u.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}

	u.logger.Info("User activated successfully", "userID", userID)
	return nil
}

// GetUsersByRole retrieves users by role
func (u *userUsecase) GetUsersByRole(ctx context.Context, role string, limit, offset int) ([]*UserResponse, error) {
	u.logger.Debug("Getting users by role", "role", role, "limit", limit, "offset", offset)

	users, err := u.userRepo.ListByRole(ctx, role, limit, offset)
	if err != nil {
		u.logger.Error("Error getting users by role", "error", err, "role", role)
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}

	return u.toUserResponses(users), nil
}

// DeleteUser deletes a user account
func (u *userUsecase) DeleteUser(ctx context.Context, userID int) error {
	u.logger.Info("Deleting user", "userID", userID)

	// Get user first to check if it exists
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Delete user
	if err := u.userRepo.Delete(ctx, userID); err != nil {
		u.logger.Error("Error deleting user", "error", err, "userID", userID)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	u.logger.Info("User deleted successfully", "userID", userID)
	return nil
}

// Helper methods

// hashPassword hashes a password using bcrypt
func (u *userUsecase) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// verifyPassword verifies a password against its hash
func (u *userUsecase) verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// toUserResponse converts entity to response
func (u *userUsecase) toUserResponse(user *entity.User) *UserResponse {
	return &UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Role:          user.Role.String(),
		IsActive:      user.IsActive,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		LastLoginAt:   user.LastLoginAt,
	}
}

// toUserResponses converts slice of entities to responses
func (u *userUsecase) toUserResponses(users []*entity.User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = u.toUserResponse(user)
	}
	return responses
}

package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/wallet_topup_system/internal/adapter/dto"
	usecase "github.com/hydr0g3nz/wallet_topup_system/internal/application"
)

// UserController handles HTTP requests related to user operations
type UserController struct {
	userUseCase usecase.UserUsecase
}

// NewUserController creates a new instance of UserController
func NewUserController(userUseCase usecase.UserUsecase) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}

// Register handles user registration
func (c *UserController) Register(ctx *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err)
	}

	if req.Email == "" || req.Password == "" || req.Role == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Email, Password, and Role are required",
		})
	}

	response, err := c.userUseCase.Register(ctx.Context(), &usecase.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusCreated, "User registered successfully", response)
}

// Login handles user authentication
func (c *UserController) Login(ctx *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err)
	}

	if req.Email == "" || req.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Email and Password are required",
		})
	}

	response, err := c.userUseCase.Login(ctx.Context(), &usecase.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Login successful", response)
}

// GetProfile handles getting user profile
func (c *UserController) GetProfile(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("id")
	if userIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "User ID is required",
		})
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid User ID format",
		})
	}

	response, err := c.userUseCase.GetProfile(ctx.Context(), userID)
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Profile retrieved successfully", response)
}

// UpdateProfile handles updating user profile
func (c *UserController) UpdateProfile(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("id")
	if userIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "User ID is required",
		})
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid User ID format",
		})
	}

	var req dto.UpdateProfileRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err)
	}

	response, err := c.userUseCase.UpdateProfile(ctx.Context(), userID, &usecase.UpdateProfileRequest{
		Email: req.Email,
	})
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Profile updated successfully", response)
}

// ChangePassword handles changing user password
func (c *UserController) ChangePassword(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("id")
	if userIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "User ID is required",
		})
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid User ID format",
		})
	}

	var req dto.ChangePasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err)
	}

	if req.CurrentPassword == "" || req.NewPassword == "" || req.ConfirmPassword == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "CurrentPassword, NewPassword, and ConfirmPassword are required",
		})
	}

	err = c.userUseCase.ChangePassword(ctx.Context(), userID, &usecase.ChangePasswordRequest{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	})
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Password changed successfully", nil)
}

// VerifyEmail handles email verification
func (c *UserController) VerifyEmail(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("id")
	if userIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "User ID is required",
		})
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid User ID format",
		})
	}

	err = c.userUseCase.VerifyEmail(ctx.Context(), userID)
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Email verified successfully", nil)
}

// DeactivateUser handles user deactivation
func (c *UserController) DeactivateUser(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("id")
	if userIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "User ID is required",
		})
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid User ID format",
		})
	}

	err = c.userUseCase.DeactivateUser(ctx.Context(), userID)
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "User deactivated successfully", nil)
}

// ActivateUser handles user activation
func (c *UserController) ActivateUser(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("id")
	if userIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "User ID is required",
		})
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid User ID format",
		})
	}

	err = c.userUseCase.ActivateUser(ctx.Context(), userID)
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "User activated successfully", nil)
}

// GetUsersByRole handles getting users by role with pagination
func (c *UserController) GetUsersByRole(ctx *fiber.Ctx) error {
	role := ctx.Query("role")
	if role == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Role query parameter is required",
		})
	}

	limitStr := ctx.Query("limit", "10")
	offsetStr := ctx.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid limit parameter",
		})
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid offset parameter",
		})
	}

	response, err := c.userUseCase.GetUsersByRole(ctx.Context(), role, limit, offset)
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Users retrieved successfully", response)
}

// DeleteUser handles user deletion
func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("id")
	if userIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "User ID is required",
		})
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid User ID format",
		})
	}

	err = c.userUseCase.DeleteUser(ctx.Context(), userID)
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "User deleted successfully", nil)
}

// GetMe handles getting current user profile (from JWT token)
func (c *UserController) GetMe(ctx *fiber.Ctx) error {
	// Assuming you have middleware that sets user ID in context
	userID := ctx.Locals("userID")
	if userID == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "User not authenticated",
		})
	}

	userIDInt, ok := userID.(int)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Invalid user ID in context",
		})
	}

	response, err := c.userUseCase.GetProfile(ctx.Context(), userIDInt)
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Current user profile retrieved successfully", response)
}

// UpdateMe handles updating current user profile
func (c *UserController) UpdateMe(ctx *fiber.Ctx) error {
	// Assuming you have middleware that sets user ID in context
	userID := ctx.Locals("userID")
	if userID == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "User not authenticated",
		})
	}

	userIDInt, ok := userID.(int)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Invalid user ID in context",
		})
	}

	var req dto.UpdateProfileRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err)
	}

	response, err := c.userUseCase.UpdateProfile(ctx.Context(), userIDInt, &usecase.UpdateProfileRequest{
		Email: req.Email,
	})
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Profile updated successfully", response)
}

// ChangeMyPassword handles changing current user password
func (c *UserController) ChangeMyPassword(ctx *fiber.Ctx) error {
	// Assuming you have middleware that sets user ID in context
	userID := ctx.Locals("userID")
	if userID == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "User not authenticated",
		})
	}

	userIDInt, ok := userID.(int)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Invalid user ID in context",
		})
	}

	var req dto.ChangePasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err)
	}

	if req.CurrentPassword == "" || req.NewPassword == "" || req.ConfirmPassword == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "CurrentPassword, NewPassword, and ConfirmPassword are required",
		})
	}

	err := c.userUseCase.ChangePassword(ctx.Context(), userIDInt, &usecase.ChangePasswordRequest{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	})
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Password changed successfully", nil)
}

// RegisterRoutes registers the routes for the user controller
func (c *UserController) RegisterRoutes(router fiber.Router) {
	userGroup := router.Group("/users")

	// Public routes
	userGroup.Post("/register", c.Register)
	userGroup.Post("/login", c.Login)

	// Protected routes (require authentication)
	userGroup.Get("/me", c.GetMe)
	userGroup.Put("/me", c.UpdateMe)
	userGroup.Put("/me/password", c.ChangeMyPassword)

	// Admin routes (require admin role)
	userGroup.Get("/", c.GetUsersByRole)
	userGroup.Get("/:id", c.GetProfile)
	userGroup.Put("/:id", c.UpdateProfile)
	userGroup.Put("/:id/password", c.ChangePassword)
	userGroup.Put("/:id/verify-email", c.VerifyEmail)
	userGroup.Put("/:id/deactivate", c.DeactivateUser)
	userGroup.Put("/:id/activate", c.ActivateUser)
	userGroup.Delete("/:id", c.DeleteUser)
}

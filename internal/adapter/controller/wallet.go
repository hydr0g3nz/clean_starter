package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/wallet_topup_system/internal/adapter/dto"
	usecase "github.com/hydr0g3nz/wallet_topup_system/internal/application"
)

// WalletController handles HTTP requests related to wallet operations
type WalletController struct {
	walletUseCase usecase.WalletUsecase
}

// NewWalletController creates a new instance of WalletController
func NewWalletController(walletUseCase usecase.WalletUsecase) *WalletController {
	return &WalletController{
		walletUseCase: walletUseCase,
	}
}

// VerifyTopup handles the verification of a top-up request
func (c *WalletController) VerifyTopup(ctx *fiber.Ctx) error {
	var req dto.VerifyRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err)
	}

	if req.UserID == 0 || req.Amount <= 0 || req.PaymentMethod == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "UserID, Amount, and PaymentMethod are required and must be valid",
		})
	}

	response, err := c.walletUseCase.VerifyTopup(req.UserID, req.Amount, req.PaymentMethod)
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Top-up verified successfully", response)
}

// ConfirmTopup handles the confirmation of a top-up transaction
func (c *WalletController) ConfirmTopup(ctx *fiber.Ctx) error {
	var req dto.ConfirmRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err)
	}

	if req.TransactionID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Transaction ID is required",
		})
	}

	transaction, wallet, err := c.walletUseCase.ConfirmTopup(req.TransactionID)
	if err != nil {
		return HandleError(ctx, err)
	}

	response := dto.ConfirmResponse{
		TransactionID: transaction.ID,
		UserID:        transaction.UserID,
		Amount:        transaction.Amount.Amount(),
		Status:        transaction.Status.String(),
		Balance:       wallet.Balance.Amount(),
	}

	return SuccessResp(ctx, fiber.StatusOK, "Top-up confirmed successfully", response)
}

// RegisterRoutes registers the routes for the wallet controller
func (c *WalletController) RegisterRoutes(router fiber.Router) {
	walletGroup := router.Group("/wallet")
	walletGroup.Post("/verify", c.VerifyTopup)
	walletGroup.Post("/confirm", c.ConfirmTopup)
}

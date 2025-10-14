package controller

import (
	"github.com/MishraShardendu22/Scanner/util"
	"github.com/gofiber/fiber/v2"
)

// GetCurrentToken returns information about the current active token
func GetCurrentToken(c *fiber.Ctx) error {
	tokenManager := util.GetTokenManager()

	response := fiber.Map{
		"current_index": tokenManager.GetCurrentIndex(),
		"total_tokens":  tokenManager.GetTokenCount(),
		"current_token": maskToken(tokenManager.GetCurrentToken()),
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Current token info", response, "")
}

// RotateToken manually rotates to the next token
func RotateToken(c *fiber.Ctx) error {
	tokenManager := util.GetTokenManager()
	tokenManager.RotateToken()

	response := fiber.Map{
		"current_index": tokenManager.GetCurrentIndex(),
		"total_tokens":  tokenManager.GetTokenCount(),
		"current_token": maskToken(tokenManager.GetCurrentToken()),
		"message":       "Token rotated successfully",
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Token rotated", response, "")
}

// maskToken masks the token for security (shows first and last 4 characters)
func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "..." + token[len(token)-4:]
}

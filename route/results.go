package route

import (
	"github.com/MishraShardendu22/Scanner/controller"
	"github.com/gofiber/fiber/v2"
)

// SetupResultRoutes sets up routes for retrieving scan results
func SetupResultRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Get specific scan result by scan_id
	// GET /api/results/:scan_id
	api.Get("/results/:scan_id", controller.GetScanResult)

	// Get all scan results with pagination
	// GET /api/results?page=1&limit=10
	api.Get("/results", controller.GetAllResults)

	// Get dashboard with aggregated statistics
	// GET /api/dashboard
	api.Get("/dashboard", controller.GetDashboard)
}

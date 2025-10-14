package route

import (
	"github.com/MishraShardendu22/Scanner/controller"
	"github.com/gofiber/fiber/v2"
)

// SetupScanRoutes sets up routes for triggering security scans
func SetupScanRoutes(app *fiber.App) {
	api := app.Group("/api/scan")

	// Scan a specific request by request_id
	// POST /api/scan/:request_id
	api.Post("/:request_id", controller.ScanRequest)

	// Scan a request by MongoDB ID
	// POST /api/scan/by-id/:id
	api.Post("/by-id/:id", controller.ScanByID)

	// Scan all models for an organization
	// POST /api/scan/org/:org/models
	api.Post("/org/:org/models", controller.ScanOrgModels)

	// Scan all datasets for an organization
	// POST /api/scan/org/:org/datasets
	api.Post("/org/:org/datasets", controller.ScanOrgDatasets)

	// Scan all spaces for an organization
	// POST /api/scan/org/:org/spaces
	api.Post("/org/:org/spaces", controller.ScanOrgSpaces)
}

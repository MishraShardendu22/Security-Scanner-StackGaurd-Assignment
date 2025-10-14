package route

import (
	"github.com/MishraShardendu22/Scanner/controller"
	"github.com/gofiber/fiber/v2"
)

// SetupOrgRoutes sets up routes for organization-specific operations
func SetupOrgRoutes(app *fiber.App) {
	api := app.Group("/api/org")

	// Fetch all models for an organization
	// GET /api/org/:org/models?include_prs=true&include_discussion=true
	api.Get("/:org/models", controller.FetchOrgModels)

	// Fetch all datasets for an organization
	// GET /api/org/:org/datasets?include_prs=true&include_discussion=true
	api.Get("/:org/datasets", controller.FetchOrgDatasets)

	// Fetch all spaces for an organization
	// GET /api/org/:org/spaces?include_prs=true&include_discussion=true
	api.Get("/:org/spaces", controller.FetchOrgSpaces)
}

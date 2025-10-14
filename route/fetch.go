package route

import (
	"github.com/MishraShardendu22/Scanner/controller"
	"github.com/gofiber/fiber/v2"
)

// SetupFetchRoutes sets up routes for fetching individual models, datasets, and spaces
func SetupFetchRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Fetch single model
	// GET /api/model/:modelId?include_prs=true&include_discussion=true
	api.Get("/model/:modelId", controller.FetchModel)

	// Fetch single dataset
	// GET /api/dataset/:datasetId?include_prs=true&include_discussion=true
	api.Get("/dataset/:datasetId", controller.FetchDataset)

	// Fetch single space
	// GET /api/space/:spaceId?include_prs=true&include_discussion=true
	api.Get("/space/:spaceId", controller.FetchSpace)

	// Fetch PRs for a resource
	// GET /api/:type/:id/prs (type can be: models, datasets, spaces)
	api.Get("/:type/:id/prs", controller.FetchPRs)

	// Fetch discussions for a resource
	// GET /api/:type/:id/discussions (type can be: models, datasets, spaces)
	api.Get("/:type/:id/discussions", controller.FetchDiscussions)
}

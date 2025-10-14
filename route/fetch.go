package route

import (
	"github.com/MishraShardendu22/Scanner/controller"
	"github.com/gofiber/fiber/v2"
)

func SetupFetchRoutes(app *fiber.App) {

	api := app.Group("/api")
	api.Get("/model/:modelId", controller.FetchModel)
	api.Get("/dataset/:datasetId", controller.FetchDataset)
	api.Get("/space/:spaceId", controller.FetchSpace)
	api.Get("/:type/:id/prs", controller.FetchPRs)
	api.Get("/:type/:id/discussions", controller.FetchDiscussions)
}

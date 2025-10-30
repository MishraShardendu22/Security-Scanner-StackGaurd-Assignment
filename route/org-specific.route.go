package route

import (
	"github.com/MishraShardendu22/Scanner/controller"
	"github.com/gofiber/fiber/v2"
)

func SetupOrgRoutes(app *fiber.App) {

	api := app.Group("/api/org")
	api.Get("/:org/models", controller.FetchOrgModels)
	api.Get("/:org/datasets", controller.FetchOrgDatasets)
	api.Get("/:org/spaces", controller.FetchOrgSpaces)
}

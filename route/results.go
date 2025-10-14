package route

import (
	"github.com/MishraShardendu22/Scanner/controller"
	"github.com/gofiber/fiber/v2"
)

func SetupResultRoutes(app *fiber.App) {

	api := app.Group("/api")

	api.Get("/results/:scan_id", controller.GetScanResult)
	api.Get("/results", controller.GetAllResults)
	api.Get("/dashboard", controller.GetDashboard)
}

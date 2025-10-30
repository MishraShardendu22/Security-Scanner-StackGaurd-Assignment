package route

import (
	"github.com/MishraShardendu22/Scanner/controller"
	"github.com/gofiber/fiber/v2"
)

func SetupScanRoutes(app *fiber.App) {

	api := app.Group("/api")
	api.Post("/scan", controller.UnifiedScan)
	api.Post("/store", controller.StoreScanResult)

	scanAPI := app.Group("/api/scan")
	scanAPI.Post("/:request_id", controller.ScanRequest)
	scanAPI.Post("/by-id/:id", controller.ScanByID)
	scanAPI.Post("/org/:org/models", controller.ScanOrgModels)
	scanAPI.Post("/org/:org/datasets", controller.ScanOrgDatasets)
	scanAPI.Post("/org/:org/spaces", controller.ScanOrgSpaces)
}

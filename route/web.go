package route

import (
	"github.com/MishraShardendu22/Scanner/controller"
	"github.com/MishraShardendu22/Scanner/template"
	"github.com/gofiber/fiber/v2"
)

func RegisterWebRoutes(app *fiber.App) {
	// SEO files
	app.Get("/robots.txt", func(c *fiber.Ctx) error {
		return c.SendFile("./public/robots.txt")
	})

	app.Get("/sitemap.xml", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/xml")
		return c.SendFile("./public/sitemap.xml")
	})

	// Home page
	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return template.IndexNew().Render(c.Context(), c.Response().BodyWriter())
	})

	// Dashboard page
	app.Get("/dashboard", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return template.Dashboard().Render(c.Context(), c.Response().BodyWriter())
	})

	// Dashboard API endpoint
	app.Get("/api/dashboard", controller.GetDashboardStats)

	// Scan page
	app.Get("/scan", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return template.ScanForm().Render(c.Context(), c.Response().BodyWriter())
	})

	// Results pages
	app.Get("/results", controller.GetResultsPage)
	app.Get("/results/:request_id", controller.GetResultDetailPage)

	// API Testing page
	app.Get("/api-tester", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return template.APITester().Render(c.Context(), c.Response().BodyWriter())
	})
}

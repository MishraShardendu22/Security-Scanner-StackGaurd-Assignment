package route

import (
	"github.com/MishraShardendu22/Scanner/controller"
	"github.com/MishraShardendu22/Scanner/template"
	"github.com/gofiber/fiber/v2"
)

func RegisterWebRoutes(app *fiber.App) {
	// SEO files - cache for 1 day
	app.Get("/robots.txt", func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "public, max-age=86400") // 24 hours
		return c.SendFile("./public/robots.txt")
	})

	app.Get("/sitemap.xml", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/xml")
		c.Set("Cache-Control", "public, max-age=86400") // 24 hours
		return c.SendFile("./public/sitemap.xml")
	})

	// Static files route (if any)
	app.Static("/public", "./public", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        false,
		CacheDuration: 24 * 60 * 60, // 24 hours
		MaxAge:        86400,        // 1 day
	})

	// Home page - cache for 5 minutes
	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		c.Set("Cache-Control", "public, max-age=300") // 5 minutes
		return template.IndexNew().Render(c.Context(), c.Response().BodyWriter())
	})

	// Dashboard page - short cache since data changes
	app.Get("/dashboard", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		c.Set("Cache-Control", "public, max-age=60") // 1 minute
		return template.Dashboard().Render(c.Context(), c.Response().BodyWriter())
	})

	// Dashboard API endpoint - short cache
	app.Get("/api/dashboard", controller.GetDashboardStats)

	// Scan page - cache for 10 minutes
	app.Get("/scan", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		c.Set("Cache-Control", "public, max-age=600") // 10 minutes
		return template.ScanForm().Render(c.Context(), c.Response().BodyWriter())
	})

	// Results pages - short cache
	app.Get("/results", controller.GetResultsPage)
	app.Get("/results/:request_id", controller.GetResultDetailPage)

	// API Testing page - cache for 10 minutes
	app.Get("/api-tester", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		c.Set("Cache-Control", "public, max-age=600") // 10 minutes
		return template.APITester().Render(c.Context(), c.Response().BodyWriter())
	})
}

package controller

import (
	"github.com/MishraShardendu22/Scanner/models"
	"github.com/MishraShardendu22/Scanner/template"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

// GetResultsPage renders the results list page
func GetResultsPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	var results []models.SCAN_RESULT

	err := mgm.Coll(&models.SCAN_RESULT{}).SimpleFind(&results, bson.M{})
	if err != nil {
		results = []models.SCAN_RESULT{}
	}

	return template.ResultsListNew(results).Render(c.Context(), c.Response().BodyWriter())
}

// GetResultDetailPage renders a single result detail page
func GetResultDetailPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	requestID := c.Params("request_id")

	var result models.SCAN_RESULT
	err := mgm.Coll(&models.SCAN_RESULT{}).First(bson.M{"request_id": requestID}, &result)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Result not found")
	}

	return template.ResultDetailNew(result).Render(c.Context(), c.Response().BodyWriter())
}

package controller

import (
	"github.com/MishraShardendu22/Scanner/models"
	"github.com/MishraShardendu22/Scanner/template"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetResultsPage renders the results list page
func GetResultsPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")

	// Get page number from query, default to 1
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	perPage := 10 // Results per page
	skip := (page - 1) * perPage

	var results []models.SCAN_RESULT
	var totalCount int64

	// Get total count
	coll := mgm.Coll(&models.SCAN_RESULT{})
	totalCount, err := coll.CountDocuments(c.Context(), bson.M{})
	if err != nil {
		totalCount = 0
	}

	// Get paginated results - sort by created_at descending (newest first)
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(perPage))
	findOptions.SetSort(bson.M{"created_at": -1})

	cursor, err := coll.Find(c.Context(), bson.M{}, findOptions)
	if err != nil {
		results = []models.SCAN_RESULT{}
	} else {
		defer cursor.Close(c.Context())
		if err = cursor.All(c.Context(), &results); err != nil {
			results = []models.SCAN_RESULT{}
		}
	}

	totalPages := int((totalCount + int64(perPage) - 1) / int64(perPage))

	return template.ResultsListNew(results, page, totalPages).Render(c.Context(), c.Response().BodyWriter())
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

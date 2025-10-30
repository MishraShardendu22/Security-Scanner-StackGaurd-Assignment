package controller

import (
	"github.com/MishraShardendu22/Scanner/models"
	"github.com/MishraShardendu22/Scanner/templ_ms22"

	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetResultsPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")

	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	perPage := 10
	skip := (page - 1) * perPage

	var results []models.SCAN_RESULT
	var totalCount int64

	coll := mgm.Coll(&models.SCAN_RESULT{})
	totalCount, err := coll.CountDocuments(c.Context(), bson.M{})
	if err != nil {
		totalCount = 0
	}

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

	return templ_ms22.ResultsListNew(results, page, totalPages).Render(c.Context(), c.Response().BodyWriter())
}

func GetResultDetailPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	requestID := c.Params("request_id")

	var result models.SCAN_RESULT
	err := mgm.Coll(&models.SCAN_RESULT{}).First(bson.M{"request_id": requestID}, &result)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Result not found")
	}

	return templ_ms22.ResultDetailNew(result).Render(c.Context(), c.Response().BodyWriter())
}

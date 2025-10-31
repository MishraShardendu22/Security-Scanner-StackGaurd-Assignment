package controller

import (
	"time"

	"github.com/MishraShardendu22/Scanner/models"
	"github.com/MishraShardendu22/Scanner/templ_ms22"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DashboardStats struct {
	TotalScans            int              `json:"total_scans"`
	TotalFindings         int              `json:"total_findings"`
	TotalResourcesScanned int              `json:"total_resources_scanned"`
	HighSeverityFindings  int              `json:"high_severity_findings"`
	FindingsByType        map[string]int   `json:"findings_by_type"`
	RecentScans           []RecentScanInfo `json:"recent_scans"`
}

type RecentScanInfo struct {
	RequestID      string    `json:"request_id"`
	ResourcesCount int       `json:"resources_count"`
	FindingsCount  int       `json:"findings_count"`
	CreatedAt      time.Time `json:"created_at"`
}

func GetDashboardStats(c *fiber.Ctx) error {
	stats := DashboardStats{
		FindingsByType: make(map[string]int),
		RecentScans:    []RecentScanInfo{},
	}

	totalScans, err := mgm.Coll(&models.SCAN_RESULT{}).CountDocuments(mgm.Ctx(), bson.M{})
	if err == nil {
		stats.TotalScans = int(totalScans)
	}

	var scanResults []models.SCAN_RESULT
	err = mgm.Coll(&models.SCAN_RESULT{}).SimpleFind(&scanResults, bson.M{})
	if err == nil {
		for _, scan := range scanResults {
			stats.TotalResourcesScanned += len(scan.ScannedResources)

			for _, resource := range scan.ScannedResources {
				findingsCount := len(resource.Findings)
				stats.TotalFindings += findingsCount

				for _, finding := range resource.Findings {
					if finding.SecretType != "" {
						stats.FindingsByType[finding.SecretType]++
						stats.HighSeverityFindings++
					}
				}
			}
		}
	}

	var recentResults []models.SCAN_RESULT
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(10)
	cursor, err := mgm.Coll(&models.SCAN_RESULT{}).Find(mgm.Ctx(), bson.M{}, opts)
	if err == nil {
		defer cursor.Close(mgm.Ctx())
		if err = cursor.All(mgm.Ctx(), &recentResults); err == nil {
			for _, result := range recentResults {
				totalFindings := 0
				for _, resource := range result.ScannedResources {
					totalFindings += len(resource.Findings)
				}

				stats.RecentScans = append(stats.RecentScans, RecentScanInfo{
					RequestID:      result.RequestID,
					ResourcesCount: len(result.ScannedResources),
					FindingsCount:  totalFindings,
					CreatedAt:      result.CreatedAt,
				})
			}
		}
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   stats,
	})
}

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

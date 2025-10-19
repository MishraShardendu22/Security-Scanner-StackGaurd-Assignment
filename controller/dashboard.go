package controller

import (
	"time"

	"github.com/MishraShardendu22/Scanner/models"
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

	// Get total scans
	totalScans, err := mgm.Coll(&models.SCAN_RESULT{}).CountDocuments(mgm.Ctx(), bson.M{})
	if err == nil {
		stats.TotalScans = int(totalScans)
	}

	// Get all scan results to calculate stats
	var scanResults []models.SCAN_RESULT
	err = mgm.Coll(&models.SCAN_RESULT{}).SimpleFind(&scanResults, bson.M{})
	if err == nil {
		// Calculate total findings and resources
		for _, scan := range scanResults {
			stats.TotalResourcesScanned += len(scan.ScannedResources)

			for _, resource := range scan.ScannedResources {
				findingsCount := len(resource.Findings)
				stats.TotalFindings += findingsCount

				// Count findings by type
				for _, finding := range resource.Findings {
					if finding.SecretType != "" {
						stats.FindingsByType[finding.SecretType]++
						// All findings are considered critical by default since they are all security secrets
						stats.HighSeverityFindings++
					}
				}
			}
		}
	}

	// Get recent scans (last 10)
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

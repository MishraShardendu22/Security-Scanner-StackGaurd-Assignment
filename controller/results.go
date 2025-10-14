package controller

import (
	"github.com/MishraShardendu22/Scanner/models"
	"github.com/MishraShardendu22/Scanner/util"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetScanResult retrieves stored scan details and contextual metadata by scan_id
// GET /api/results/:scan_id
func GetScanResult(c *fiber.Ctx) error {
	scanID := c.Params("scan_id")
	if scanID == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Scan ID is required", nil, "")
	}

	// Convert string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(scanID)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid scan ID format", nil, "")
	}

	// Find the scan result
	scanResult := &models.SCAN_RESULT{}
	err = mgm.Coll(scanResult).FindByID(objectID, scanResult)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusNotFound, "Scan result not found", nil, "")
	}

	// Calculate statistics
	totalFindings := 0
	findingsByType := make(map[string]int)
	findingsByResource := make(map[string]int)

	for _, resource := range scanResult.ScannedResources {
		findingCount := len(resource.Findings)
		totalFindings += findingCount
		findingsByResource[resource.Type+"_"+resource.ID] = findingCount

		for _, finding := range resource.Findings {
			findingsByType[finding.SecretType]++
		}
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Scan result retrieved successfully", map[string]interface{}{
		"scan_id":              scanResult.ID.Hex(),
		"request_id":           scanResult.RequestID,
		"scanned_resources":    scanResult.ScannedResources,
		"total_findings":       totalFindings,
		"findings_by_type":     findingsByType,
		"findings_by_resource": findingsByResource,
		"created_at":           scanResult.CreatedAt,
		"updated_at":           scanResult.UpdatedAt,
	}, "")
}

// GetDashboard provides a lightweight dashboard showing all stored results grouped by resource type and severity
// GET /api/dashboard
func GetDashboard(c *fiber.Ctx) error {
	// Fetch all scan results
	scanResults := []models.SCAN_RESULT{}
	err := mgm.Coll(&models.SCAN_RESULT{}).SimpleFind(&scanResults, bson.M{})
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch scan results", nil, "")
	}

	// Initialize dashboard data structures
	dashboardData := map[string]interface{}{
		"total_scans":           len(scanResults),
		"total_findings":        0,
		"by_resource_type":      make(map[string]int),
		"by_secret_type":        make(map[string]int),
		"by_source_type":        make(map[string]int),
		"recent_scans":          []map[string]interface{}{},
		"high_risk_findings":    []models.Finding{},
		"resources_with_issues": 0,
	}

	totalFindings := 0
	byResourceType := make(map[string]int)
	bySecretType := make(map[string]int)
	bySourceType := make(map[string]int)
	resourcesWithIssues := 0

	// High-risk patterns (these are particularly sensitive)
	highRiskPatterns := map[string]bool{
		"AWS Access Key ID":       true,
		"GitHub PAT":              true,
		"OpenAI / LLM API Key":    true,
		"Stripe Secret Key":       true,
		"Database URI with creds": true,
		"PostgreSQL URI":          true,
		"MySQL URI":               true,
		"MongoDB URI":             true,
		"Google API Key":          true,
		"Kubernetes Bearer Token": true,
		"GitHub Actions Token":    true,
	}

	highRiskFindings := []models.Finding{}
	recentScans := []map[string]interface{}{}

	// Process each scan result
	for _, scan := range scanResults {
		scanFindings := 0

		for _, resource := range scan.ScannedResources {
			if len(resource.Findings) > 0 {
				resourcesWithIssues++
			}

			// Count by resource type
			byResourceType[resource.Type] += len(resource.Findings)

			for _, finding := range resource.Findings {
				scanFindings++
				totalFindings++

				// Count by secret type
				bySecretType[finding.SecretType]++

				// Count by source type
				bySourceType[finding.SourceType]++

				// Identify high-risk findings
				if highRiskPatterns[finding.SecretType] {
					highRiskFindings = append(highRiskFindings, finding)
				}
			}
		}

		// Add to recent scans (limit to top entries for performance)
		if len(recentScans) < 10 {
			recentScans = append(recentScans, map[string]interface{}{
				"scan_id":    scan.ID.Hex(),
				"request_id": scan.RequestID,
				"findings":   scanFindings,
				"resources":  len(scan.ScannedResources),
				"created_at": scan.CreatedAt,
			})
		}
	}

	// Update dashboard data
	dashboardData["total_findings"] = totalFindings
	dashboardData["by_resource_type"] = byResourceType
	dashboardData["by_secret_type"] = bySecretType
	dashboardData["by_source_type"] = bySourceType
	dashboardData["resources_with_issues"] = resourcesWithIssues
	dashboardData["recent_scans"] = recentScans
	dashboardData["high_risk_findings"] = highRiskFindings
	dashboardData["high_risk_count"] = len(highRiskFindings)

	// Calculate severity breakdown
	severityBreakdown := map[string]int{
		"high":   len(highRiskFindings),
		"medium": 0,
		"low":    0,
	}

	// Classify remaining findings as medium/low
	mediumRiskCount := 0
	lowRiskCount := 0

	for secretType, count := range bySecretType {
		if !highRiskPatterns[secretType] {
			// Simple heuristic: API keys are medium, others are low
			if secretType == "API Key" || secretType == "Access Token" {
				mediumRiskCount += count
			} else {
				lowRiskCount += count
			}
		}
	}

	severityBreakdown["medium"] = mediumRiskCount
	severityBreakdown["low"] = lowRiskCount
	dashboardData["severity_breakdown"] = severityBreakdown

	return util.ResponseAPI(c, fiber.StatusOK, "Dashboard data retrieved successfully", dashboardData, "")
}

// GetAllResults retrieves all scan results with pagination
// GET /api/results?page=1&limit=10
func GetAllResults(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	skip := (page - 1) * limit

	// Fetch scan results with pagination
	scanResults := []models.SCAN_RESULT{}
	err := mgm.Coll(&models.SCAN_RESULT{}).SimpleFind(&scanResults, bson.M{})
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch scan results", nil, "")
	}

	// Calculate total pages
	total := len(scanResults)
	totalPages := (total + limit - 1) / limit

	// Apply pagination
	start := skip
	end := skip + limit
	if start >= total {
		scanResults = []models.SCAN_RESULT{}
	} else {
		if end > total {
			end = total
		}
		scanResults = scanResults[start:end]
	}

	// Format results for response
	results := []map[string]interface{}{}
	for _, scan := range scanResults {
		findingCount := 0
		for _, resource := range scan.ScannedResources {
			findingCount += len(resource.Findings)
		}

		results = append(results, map[string]interface{}{
			"scan_id":    scan.ID.Hex(),
			"request_id": scan.RequestID,
			"findings":   findingCount,
			"resources":  len(scan.ScannedResources),
			"created_at": scan.CreatedAt,
		})
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Scan results retrieved successfully", map[string]interface{}{
		"results":     results,
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": totalPages,
	}, "")
}

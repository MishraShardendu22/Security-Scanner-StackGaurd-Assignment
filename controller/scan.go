package controller

import (
	"fmt"
	"log"

	"github.com/MishraShardendu22/Scanner/models"
	"github.com/MishraShardendu22/Scanner/util"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ScanRequest(c *fiber.Ctx) error {

	requestID := c.Params("request_id")

	if requestID == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Request ID is required", nil, "")
	}

	log.Printf("🔍 Starting scan for request ID: %s\n", requestID)

	aiRequest := &models.AI_REQUEST{}

	err := mgm.Coll(aiRequest).First(bson.M{"request_id": requestID}, aiRequest)

	if err != nil {
		return util.ResponseAPI(c, fiber.StatusNotFound, "Request not found", nil, "")
	}

	log.Printf("✅ Found request with %d files and %d discussions\n", len(aiRequest.Siblings), len(aiRequest.Discussions))

	log.Println("🔍 Scanning for secrets...")

	findings := util.ScanAIRequest(*aiRequest, util.SecretConfig)

	log.Printf("✅ Scan complete! Found %d potential secrets\n", len(findings))

	scannedResources := util.GroupFindingsByResource(findings)

	scanResult := &models.SCAN_RESULT{

		RequestID:        requestID,
		ScannedResources: scannedResources,
	}

	if err := mgm.Coll(scanResult).Create(scanResult); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save scan results", nil, "")
	}

	totalFindings := len(findings)

	findingsByType := util.CountFindingsByType(findings)
	findingsBySource := util.CountFindingsBySource(findings)

	return util.ResponseAPI(c, fiber.StatusOK, "Scan completed successfully", map[string]interface{}{

		"scan_id":            scanResult.ID.Hex(),
		"request_id":         requestID,
		"total_findings":     totalFindings,
		"findings_by_type":   findingsByType,
		"findings_by_source": findingsBySource,
		"scanned_resources":  scannedResources,
	}, "")
}

func ScanOrgModels(c *fiber.Ctx) error {
	org := c.Params("org")
	if org == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organization name is required", nil, "")
	}

	allFindings, scannedCount, err := util.ScanOrgResources(org, util.ResourceTypeModel)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == fmt.Sprintf("no %s found for this organization", util.ResourceTypeModel) {
			statusCode = fiber.StatusNotFound
		}
		return util.ResponseAPI(c, statusCode, err.Error(), nil, "")
	}

	scannedResources := util.GroupFindingsByResource(allFindings)
	scanResult, err := util.SaveScanResults("org-scan-"+org+"-models", scannedResources)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, err.Error(), nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Organization models scanned successfully", map[string]interface{}{
		"scan_id":           scanResult.ID.Hex(),
		"organization":      org,
		"models_scanned":    scannedCount,
		"total_findings":    len(allFindings),
		"scanned_resources": scannedResources,
	}, "")
}

func ScanOrgDatasets(c *fiber.Ctx) error {
	org := c.Params("org")
	if org == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organization name is required", nil, "")
	}

	allFindings, scannedCount, err := util.ScanOrgResources(org, util.ResourceTypeDataset)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == fmt.Sprintf("no %s found for this organization", util.ResourceTypeDataset) {
			statusCode = fiber.StatusNotFound
		}
		return util.ResponseAPI(c, statusCode, err.Error(), nil, "")
	}

	scannedResources := util.GroupFindingsByResource(allFindings)
	scanResult, err := util.SaveScanResults("org-scan-"+org+"-datasets", scannedResources)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, err.Error(), nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Organization datasets scanned successfully", map[string]interface{}{
		"scan_id":           scanResult.ID.Hex(),
		"organization":      org,
		"datasets_scanned":  scannedCount,
		"total_findings":    len(allFindings),
		"scanned_resources": scannedResources,
	}, "")
}

func ScanOrgSpaces(c *fiber.Ctx) error {
	org := c.Params("org")
	if org == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organization name is required", nil, "")
	}

	allFindings, scannedCount, err := util.ScanOrgResources(org, util.ResourceTypeSpace)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == fmt.Sprintf("no %s found for this organization", util.ResourceTypeSpace) {
			statusCode = fiber.StatusNotFound
		}
		return util.ResponseAPI(c, statusCode, err.Error(), nil, "")
	}

	scannedResources := util.GroupFindingsByResource(allFindings)
	scanResult, err := util.SaveScanResults("org-scan-"+org+"-spaces", scannedResources)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, err.Error(), nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Organization spaces scanned successfully", map[string]interface{}{
		"scan_id":           scanResult.ID.Hex(),
		"organization":      org,
		"spaces_scanned":    scannedCount,
		"total_findings":    len(allFindings),
		"scanned_resources": scannedResources,
	}, "")
}

// Helper functions moved to util/findings.go and util/org.scanner.go

func ScanByID(c *fiber.Ctx) error {

	id := c.Params("id")

	if id == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "ID is required", nil, "")
	}
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid ID format", nil, "")
	}

	aiRequest := &models.AI_REQUEST{}
	err = mgm.Coll(aiRequest).FindByID(objectID, aiRequest)

	if err != nil {
		return util.ResponseAPI(c, fiber.StatusNotFound, "Request not found", nil, "")
	}

	findings := util.ScanAIRequest(*aiRequest, util.SecretConfig)

	scannedResources := util.GroupFindingsByResource(findings)

	scanResult := &models.SCAN_RESULT{

		RequestID:        aiRequest.RequestID,
		ScannedResources: scannedResources,
	}
	if err := mgm.Coll(scanResult).Create(scanResult); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save scan results", nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Scan completed successfully", map[string]interface{}{

		"scan_id":           scanResult.ID.Hex(),
		"request_id":        aiRequest.RequestID,
		"total_findings":    len(findings),
		"scanned_resources": scannedResources,
	}, "")
}

package controller

import (
	"github.com/MishraShardendu22/Scanner/models"
	"github.com/MishraShardendu22/Scanner/util"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ScanRequest triggers a security scan on a stored AI_REQUEST
// POST /api/scan/:request_id
func ScanRequest(c *fiber.Ctx) error {
	requestID := c.Params("request_id")
	if requestID == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Request ID is required", nil, "")
	}

	// Find the AI_REQUEST by request_id field (not MongoDB _id)
	aiRequest := &models.AI_REQUEST{}
	err := mgm.Coll(aiRequest).First(bson.M{"request_id": requestID}, aiRequest)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusNotFound, "Request not found", nil, "")
	}

	// Perform the scan using the scanner utility
	findings := util.ScanAIRequest(*aiRequest, util.SecretConfig)

	// Group findings by resource
	scannedResources := groupFindingsByResource(findings, aiRequest)

	// Create and save scan result
	scanResult := &models.SCAN_RESULT{
		RequestID:        requestID,
		ScannedResources: scannedResources,
	}

	if err := mgm.Coll(scanResult).Create(scanResult); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save scan results", nil, "")
	}

	// Calculate statistics
	totalFindings := len(findings)
	findingsByType := make(map[string]int)
	findingsBySource := make(map[string]int)

	for _, finding := range findings {
		findingsByType[finding.SecretType]++
		findingsBySource[finding.SourceType]++
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Scan completed successfully", map[string]interface{}{
		"scan_id":            scanResult.ID.Hex(),
		"request_id":         requestID,
		"total_findings":     totalFindings,
		"findings_by_type":   findingsByType,
		"findings_by_source": findingsBySource,
		"scanned_resources":  scannedResources,
	}, "")
}

// ScanOrgModels scans all models for an organization
// POST /api/scan/org/:org/models
func ScanOrgModels(c *fiber.Ctx) error {
	org := c.Params("org")
	if org == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organization name is required", nil, "")
	}

	// Find all AI_Models for this organization
	aiModels := []models.AI_Models{}
	err := mgm.Coll(&models.AI_Models{}).SimpleFind(&aiModels, bson.M{"org": org})
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch organization models", nil, "")
	}

	if len(aiModels) == 0 {
		return util.ResponseAPI(c, fiber.StatusNotFound, "No models found for this organization", nil, "")
	}

	// Find all AI_REQUESTs related to these models
	// Note: This assumes you're storing the relationship. If not, you'll need to modify this logic
	var allFindings []models.Finding
	var scannedCount int

	// Get all AI_REQUESTs and scan them
	aiRequests := []models.AI_REQUEST{}
	mgm.Coll(&models.AI_REQUEST{}).SimpleFind(&aiRequests, bson.M{})

	for _, req := range aiRequests {
		// Scan this request
		findings := util.ScanAIRequest(req, util.SecretConfig)
		allFindings = append(allFindings, findings...)
		scannedCount++
	}

	// Group findings by resource
	scannedResources := groupAllFindings(allFindings)

	// Create and save scan result
	scanResult := &models.SCAN_RESULT{
		RequestID:        "org-scan-" + org + "-models",
		ScannedResources: scannedResources,
	}

	if err := mgm.Coll(scanResult).Create(scanResult); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save scan results", nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Organization models scanned successfully", map[string]interface{}{
		"scan_id":           scanResult.ID.Hex(),
		"organization":      org,
		"models_scanned":    scannedCount,
		"total_findings":    len(allFindings),
		"scanned_resources": scannedResources,
	}, "")
}

// ScanOrgDatasets scans all datasets for an organization
// POST /api/scan/org/:org/datasets
func ScanOrgDatasets(c *fiber.Ctx) error {
	org := c.Params("org")
	if org == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organization name is required", nil, "")
	}

	// Find all AI_DATASETS for this organization
	aiDatasets := []models.AI_DATASETS{}
	err := mgm.Coll(&models.AI_DATASETS{}).SimpleFind(&aiDatasets, bson.M{"org": org})
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch organization datasets", nil, "")
	}

	if len(aiDatasets) == 0 {
		return util.ResponseAPI(c, fiber.StatusNotFound, "No datasets found for this organization", nil, "")
	}

	var allFindings []models.Finding
	var scannedCount int

	// Find and scan all related AI_REQUESTs
	aiRequests := []models.AI_REQUEST{}
	mgm.Coll(&models.AI_REQUEST{}).SimpleFind(&aiRequests, bson.M{})

	for _, req := range aiRequests {
		findings := util.ScanAIRequest(req, util.SecretConfig)
		allFindings = append(allFindings, findings...)
		scannedCount++
	}

	scannedResources := groupAllFindings(allFindings)

	scanResult := &models.SCAN_RESULT{
		RequestID:        "org-scan-" + org + "-datasets",
		ScannedResources: scannedResources,
	}

	if err := mgm.Coll(scanResult).Create(scanResult); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save scan results", nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Organization datasets scanned successfully", map[string]interface{}{
		"scan_id":           scanResult.ID.Hex(),
		"organization":      org,
		"datasets_scanned":  scannedCount,
		"total_findings":    len(allFindings),
		"scanned_resources": scannedResources,
	}, "")
}

// ScanOrgSpaces scans all spaces for an organization
// POST /api/scan/org/:org/spaces
func ScanOrgSpaces(c *fiber.Ctx) error {
	org := c.Params("org")
	if org == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organization name is required", nil, "")
	}

	// Find all AI_SPACES for this organization
	aiSpaces := []models.AI_SPACES{}
	err := mgm.Coll(&models.AI_SPACES{}).SimpleFind(&aiSpaces, bson.M{"org": org})
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch organization spaces", nil, "")
	}

	if len(aiSpaces) == 0 {
		return util.ResponseAPI(c, fiber.StatusNotFound, "No spaces found for this organization", nil, "")
	}

	var allFindings []models.Finding
	var scannedCount int

	aiRequests := []models.AI_REQUEST{}
	mgm.Coll(&models.AI_REQUEST{}).SimpleFind(&aiRequests, bson.M{})

	for _, req := range aiRequests {
		findings := util.ScanAIRequest(req, util.SecretConfig)
		allFindings = append(allFindings, findings...)
		scannedCount++
	}

	scannedResources := groupAllFindings(allFindings)

	scanResult := &models.SCAN_RESULT{
		RequestID:        "org-scan-" + org + "-spaces",
		ScannedResources: scannedResources,
	}

	if err := mgm.Coll(scanResult).Create(scanResult); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save scan results", nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Organization spaces scanned successfully", map[string]interface{}{
		"scan_id":           scanResult.ID.Hex(),
		"organization":      org,
		"spaces_scanned":    scannedCount,
		"total_findings":    len(allFindings),
		"scanned_resources": scannedResources,
	}, "")
}

// Helper function to group findings by resource
func groupFindingsByResource(findings []models.Finding, aiRequest *models.AI_REQUEST) []models.SCANNED_RESOURCE {
	resourceMap := make(map[string]*models.SCANNED_RESOURCE)

	for _, finding := range findings {
		var resourceKey string
		var resourceType string
		var resourceID string

		if finding.SourceType == "file" {
			resourceType = "file"
			resourceID = finding.FileName
			resourceKey = "file:" + finding.FileName
		} else if finding.SourceType == "discussion" {
			resourceType = "discussion"
			resourceID = finding.DiscussionTitle
			resourceKey = "discussion:" + finding.DiscussionTitle
		}

		if _, exists := resourceMap[resourceKey]; !exists {
			resourceMap[resourceKey] = &models.SCANNED_RESOURCE{
				Type:     resourceType,
				ID:       resourceID,
				Findings: []models.Finding{},
			}
		}

		resourceMap[resourceKey].Findings = append(resourceMap[resourceKey].Findings, finding)
	}

	var scannedResources []models.SCANNED_RESOURCE
	for _, resource := range resourceMap {
		scannedResources = append(scannedResources, *resource)
	}

	return scannedResources
}

// Helper function to group all findings
func groupAllFindings(findings []models.Finding) []models.SCANNED_RESOURCE {
	resourceMap := make(map[string]*models.SCANNED_RESOURCE)

	for _, finding := range findings {
		var resourceKey string
		var resourceType string
		var resourceID string

		if finding.SourceType == "file" {
			resourceType = "file"
			resourceID = finding.FileName
			resourceKey = "file:" + finding.FileName
		} else if finding.SourceType == "discussion" {
			resourceType = "discussion"
			resourceID = finding.DiscussionTitle
			resourceKey = "discussion:" + finding.DiscussionTitle
		}

		if _, exists := resourceMap[resourceKey]; !exists {
			resourceMap[resourceKey] = &models.SCANNED_RESOURCE{
				Type:     resourceType,
				ID:       resourceID,
				Findings: []models.Finding{},
			}
		}

		resourceMap[resourceKey].Findings = append(resourceMap[resourceKey].Findings, finding)
	}

	var scannedResources []models.SCANNED_RESOURCE
	for _, resource := range resourceMap {
		scannedResources = append(scannedResources, *resource)
	}

	return scannedResources
}

// ScanByID scans a resource by its MongoDB ID
// POST /api/scan/by-id/:id
func ScanByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "ID is required", nil, "")
	}

	// Convert string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid ID format", nil, "")
	}

	// Find the AI_REQUEST by MongoDB _id
	aiRequest := &models.AI_REQUEST{}
	err = mgm.Coll(aiRequest).FindByID(objectID, aiRequest)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusNotFound, "Request not found", nil, "")
	}

	// Perform the scan
	findings := util.ScanAIRequest(*aiRequest, util.SecretConfig)
	scannedResources := groupFindingsByResource(findings, aiRequest)

	// Create and save scan result
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

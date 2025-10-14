package controller

import (
	"log"
	"sync"

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

	scannedResources := groupFindingsByResource(findings, aiRequest)

	scanResult := &models.SCAN_RESULT{

		RequestID:        requestID,
		ScannedResources: scannedResources,
	}

	if err := mgm.Coll(scanResult).Create(scanResult); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save scan results", nil, "")
	}

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

func ScanOrgModels(c *fiber.Ctx) error {

	org := c.Params("org")

	if org == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organization name is required", nil, "")
	}

	log.Printf("🏢 Starting organization model scan: %s\n", org)

	aiModels := []models.AI_Models{}

	err := mgm.Coll(&models.AI_Models{}).SimpleFind(&aiModels, bson.M{"org": org})

	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch organization models", nil, "")
	}
	if len(aiModels) == 0 {
		return util.ResponseAPI(c, fiber.StatusNotFound, "No models found for this organization", nil, "")
	}

	log.Printf("✅ Found %d models for organization\n", len(aiModels))
	var allFindings []models.Finding
	var scannedCount int

	aiRequests := []models.AI_REQUEST{}
	mgm.Coll(&models.AI_REQUEST{}).SimpleFind(&aiRequests, bson.M{})

	log.Printf("🔍 Scanning %d requests...\n", len(aiRequests))
	var wg sync.WaitGroup
	var mu sync.Mutex

	semaphore := make(chan struct{}, 10)

	for idx, req := range aiRequests {
		wg.Add(1)
		go func(r models.AI_REQUEST, index int) {

			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			log.Printf("  [%d/%d] Scanning request: %s\n", index+1, len(aiRequests), r.RequestID)
			findings := util.ScanAIRequest(r, util.SecretConfig)
			mu.Lock()
			allFindings = append(allFindings, findings...)
			scannedCount++
			mu.Unlock()
			if len(findings) > 0 {
				log.Printf("    ⚠️  Found %d secrets\n", len(findings))
			}
		}(req, idx)
	}
	wg.Wait()

	log.Printf("✅ Scan complete! Total findings: %d\n", len(allFindings))

	scannedResources := groupAllFindings(allFindings)

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

func ScanOrgDatasets(c *fiber.Ctx) error {

	org := c.Params("org")

	if org == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organization name is required", nil, "")
	}

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

	aiRequests := []models.AI_REQUEST{}
	mgm.Coll(&models.AI_REQUEST{}).SimpleFind(&aiRequests, bson.M{})

	log.Printf("🔍 Scanning %d requests...\n", len(aiRequests))
	var wg sync.WaitGroup
	var mu sync.Mutex

	semaphore := make(chan struct{}, 10)

	for idx, req := range aiRequests {
		wg.Add(1)
		go func(r models.AI_REQUEST, index int) {

			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			log.Printf("  [%d/%d] Scanning request: %s\n", index+1, len(aiRequests), r.RequestID)
			findings := util.ScanAIRequest(r, util.SecretConfig)
			mu.Lock()
			allFindings = append(allFindings, findings...)
			scannedCount++
			mu.Unlock()
			if len(findings) > 0 {
				log.Printf("    ⚠️  Found %d secrets\n", len(findings))
			}
		}(req, idx)
	}
	wg.Wait()

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

func ScanOrgSpaces(c *fiber.Ctx) error {

	org := c.Params("org")

	if org == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organization name is required", nil, "")
	}

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

	log.Printf("🔍 Scanning %d requests...\n", len(aiRequests))
	var wg sync.WaitGroup
	var mu sync.Mutex

	semaphore := make(chan struct{}, 10)

	for idx, req := range aiRequests {
		wg.Add(1)
		go func(r models.AI_REQUEST, index int) {

			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			log.Printf("  [%d/%d] Scanning request: %s\n", index+1, len(aiRequests), r.RequestID)
			findings := util.ScanAIRequest(r, util.SecretConfig)
			mu.Lock()
			allFindings = append(allFindings, findings...)
			scannedCount++
			mu.Unlock()
			if len(findings) > 0 {
				log.Printf("    ⚠️  Found %d secrets\n", len(findings))
			}
		}(req, idx)
	}
	wg.Wait()

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

func groupFindingsByResource(findings []models.Finding, _ *models.AI_REQUEST) []models.SCANNED_RESOURCE {

	resourceMap := make(map[string]*models.SCANNED_RESOURCE)

	for _, finding := range findings {
		var resourceKey string
		var resourceType string
		var resourceID string
		switch finding.SourceType {
		case "file":
			resourceType = "file"
			resourceID = finding.FileName
			resourceKey = "file:" + finding.FileName
		case "discussion":
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

func groupAllFindings(findings []models.Finding) []models.SCANNED_RESOURCE {

	resourceMap := make(map[string]*models.SCANNED_RESOURCE)

	for _, finding := range findings {
		var resourceKey string
		var resourceType string
		var resourceID string
		switch finding.SourceType {
		case "file":
			resourceType = "file"
			resourceID = finding.FileName
			resourceKey = "file:" + finding.FileName
		case "discussion":
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

	scannedResources := groupFindingsByResource(findings, aiRequest)

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

package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/MishraShardendu22/Scanner/models"
	"github.com/MishraShardendu22/Scanner/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kamva/mgm/v3"
)

// ScanRequestBody represents the unified scan request body
type ScanRequestBody struct {
	ModelID            string `json:"model_id"`
	DatasetID          string `json:"dataset_id"`
	SpaceID            string `json:"space_id"`
	Org                string `json:"org"`
	User               string `json:"user"`
	IncludeDiscussions bool   `json:"include_discussions"`
	IncludePRs         bool   `json:"include_prs"`
}

// UnifiedScan is the main scan endpoint as per assignment requirements
// POST /api/scan
func UnifiedScan(c *fiber.Ctx) error {
	var req ScanRequestBody
	if err := c.BodyParser(&req); err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid request body", nil, "")
	}

	// Generate scan ID
	scanID := fmt.Sprintf("SG-%s-%s", time.Now().Format("2006-0102"), uuid.New().String()[:8])
	requestID := uuid.New().String()

	// Create AI_REQUEST to store fetched data
	aiRequest := &models.AI_REQUEST{
		RequestID:   requestID,
		Siblings:    []models.SIBLING{},
		Discussions: []models.DISCUSSION{},
	}

	var scannedResources []models.SCANNED_RESOURCE
	var resourceType, resourceID string

	// Fetch and scan based on provided IDs
	if req.ModelID != "" {
		resourceType = "model"
		resourceID = req.ModelID
		if err := fetchAndAddToRequest(aiRequest, req.ModelID, "models", req.IncludePRs, req.IncludeDiscussions); err != nil {
			return util.ResponseAPI(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to fetch model: %v", err), nil, "")
		}
	} else if req.DatasetID != "" {
		resourceType = "dataset"
		resourceID = req.DatasetID
		if err := fetchAndAddToRequest(aiRequest, req.DatasetID, "datasets", req.IncludePRs, req.IncludeDiscussions); err != nil {
			return util.ResponseAPI(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to fetch dataset: %v", err), nil, "")
		}
	} else if req.SpaceID != "" {
		resourceType = "space"
		resourceID = req.SpaceID
		if err := fetchAndAddToRequest(aiRequest, req.SpaceID, "spaces", req.IncludePRs, req.IncludeDiscussions); err != nil {
			return util.ResponseAPI(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to fetch space: %v", err), nil, "")
		}
	} else if req.Org != "" {
		// Organization-level scan
		return scanOrganization(c, req.Org, req.IncludePRs, req.IncludeDiscussions, scanID)
	} else if req.User != "" {
		// User-level scan (treat as org)
		return scanOrganization(c, req.User, req.IncludePRs, req.IncludeDiscussions, scanID)
	} else {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "At least one of model_id, dataset_id, space_id, org, or user is required", nil, "")
	}

	// Save the AI_REQUEST
	if err := mgm.Coll(aiRequest).Create(aiRequest); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save request", nil, "")
	}

	// Perform the scan
	findings := util.ScanAIRequest(*aiRequest, util.SecretConfig)

	// Group findings by file/discussion
	findingsMap := make(map[string][]models.Finding)
	for _, finding := range findings {
		var key string
		if finding.SourceType == "file" {
			key = "file:" + finding.FileName
		} else {
			key = "discussion:" + finding.DiscussionTitle
		}
		findingsMap[key] = append(findingsMap[key], finding)
	}

	// Create scanned resource
	resourceFindings := []models.Finding{}
	for _, findingsList := range findingsMap {
		resourceFindings = append(resourceFindings, findingsList...)
	}

	scannedResource := models.SCANNED_RESOURCE{
		Type:     resourceType,
		ID:       resourceID,
		Findings: resourceFindings,
	}
	scannedResources = append(scannedResources, scannedResource)

	// Create and save scan result
	scanResult := &models.SCAN_RESULT{
		RequestID:        requestID,
		ScannedResources: scannedResources,
	}

	if err := mgm.Coll(scanResult).Create(scanResult); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save scan results", nil, "")
	}

	// Format response as per assignment requirements
	response := map[string]interface{}{
		"scan_id": scanID,
		"scanned_resources": []map[string]interface{}{
			{
				"type":     resourceType,
				"id":       resourceID,
				"findings": formatFindings(resourceFindings),
			},
		},
		"timestamp":      time.Now().Format(time.RFC3339),
		"total_findings": len(findings),
		"storage_id":     scanResult.ID.Hex(),
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Scan completed successfully", response, "")
}

// StoreScanResult stores scan results (POST /api/store)
func StoreScanResult(c *fiber.Ctx) error {
	var scanData map[string]interface{}
	if err := c.BodyParser(&scanData); err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid request body", nil, "")
	}

	// Extract scan_id
	scanIDValue, ok := scanData["scan_id"].(string)
	if !ok {
		scanIDValue = fmt.Sprintf("SG-%s-%s", time.Now().Format("2006-0102"), uuid.New().String()[:8])
	}

	// Convert to SCAN_RESULT format
	scanResult := &models.SCAN_RESULT{
		RequestID:        scanIDValue,
		ScannedResources: []models.SCANNED_RESOURCE{},
	}

	// Parse scanned_resources if present
	if resources, ok := scanData["scanned_resources"].([]interface{}); ok {
		for _, res := range resources {
			if resMap, ok := res.(map[string]interface{}); ok {
				resource := models.SCANNED_RESOURCE{}

				if resType, ok := resMap["type"].(string); ok {
					resource.Type = resType
				}
				if resID, ok := resMap["id"].(string); ok {
					resource.ID = resID
				}

				// Parse findings
				if findings, ok := resMap["findings"].([]interface{}); ok {
					for _, f := range findings {
						if fMap, ok := f.(map[string]interface{}); ok {
							finding := models.Finding{}
							if secretType, ok := fMap["secret_type"].(string); ok {
								finding.SecretType = secretType
							}
							if pattern, ok := fMap["pattern"].(string); ok {
								finding.Pattern = pattern
							}
							if secret, ok := fMap["secret"].(string); ok {
								finding.Secret = secret
							}
							if file, ok := fMap["file"].(string); ok {
								finding.FileName = file
								finding.SourceType = "file"
							}
							if line, ok := fMap["line"].(float64); ok {
								finding.Line = int(line)
							}
							resource.Findings = append(resource.Findings, finding)
						}
					}
				}

				scanResult.ScannedResources = append(scanResult.ScannedResources, resource)
			}
		}
	}

	// Save to database
	if err := mgm.Coll(scanResult).Create(scanResult); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to store scan results", nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Scan results stored successfully", map[string]interface{}{
		"status":     "stored",
		"scan_id":    scanIDValue,
		"storage_id": scanResult.ID.Hex(),
	}, "")
}

// Helper function to fetch and add data to AI_REQUEST
func fetchAndAddToRequest(aiRequest *models.AI_REQUEST, resourceID, resourceType string, includePRs, includeDiscussions bool) error {
	// Fetch resource metadata
	url := fmt.Sprintf("https://huggingface.co/api/%s/%s", resourceType, resourceID)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resource not found: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var resourceData map[string]interface{}
	if err := json.Unmarshal(body, &resourceData); err != nil {
		return err
	}

	// Extract and fetch siblings (files)
	if siblings, ok := resourceData["siblings"].([]interface{}); ok {
		for _, sib := range siblings {
			if sibMap, ok := sib.(map[string]interface{}); ok {
				if filename, ok := sibMap["rfilename"].(string); ok {
					sibling := fetchFileContentHelper(resourceID, filename)
					aiRequest.Siblings = append(aiRequest.Siblings, sibling)
				}
			}
		}
	}

	// Fetch discussions/PRs if requested
	if includePRs || includeDiscussions {
		discussions, _ := fetchDiscussionsHelper(resourceID, resourceType, includePRs, includeDiscussions)
		aiRequest.Discussions = discussions
	}

	return nil
}

// Helper function to fetch discussions (uses existing fetchDiscussions and getDiscussionsFromURL from fetch.go)
func fetchDiscussionsHelper(id, resourceType string, includePRs, includeDiscussion bool) ([]models.DISCUSSION, error) {
	var discussions []models.DISCUSSION

	if includePRs {
		url := fmt.Sprintf("https://huggingface.co/api/%s/%s/discussions?types=pr&status=all", resourceType, id)
		prs, _ := fetchDiscussionsFromURL(url)
		discussions = append(discussions, prs...)
	}

	if includeDiscussion {
		url := fmt.Sprintf("https://huggingface.co/api/%s/%s/discussions?types=discussion&status=all", resourceType, id)
		discs, _ := fetchDiscussionsFromURL(url)
		discussions = append(discussions, discs...)
	}

	return discussions, nil
}

// fetchDiscussionsFromURL fetches discussions from a URL
func fetchDiscussionsFromURL(url string) ([]models.DISCUSSION, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rawDiscussions []map[string]interface{}
	if err := json.Unmarshal(body, &rawDiscussions); err != nil {
		return nil, err
	}

	var discussions []models.DISCUSSION
	for _, disc := range rawDiscussions {
		discussion := models.DISCUSSION{}

		if num, ok := disc["num"].(float64); ok {
			discussion.Num = int64(num)
		}
		if title, ok := disc["title"].(string); ok {
			discussion.Title = title
		}
		if status, ok := disc["status"].(string); ok {
			discussion.Status = status
		}
		if isPR, ok := disc["isPullRequest"].(bool); ok {
			discussion.IsPullRequest = isPR
		}
		if createdAt, ok := disc["createdAt"].(string); ok {
			discussion.CreatedAt = createdAt
		}
		if author, ok := disc["author"].(map[string]interface{}); ok {
			if name, ok := author["name"].(string); ok {
				discussion.AuthorName = name
			}
		}
		if repo, ok := disc["repo"].(map[string]interface{}); ok {
			if name, ok := repo["name"].(string); ok {
				discussion.RepoName = name
			}
		}
		if numComments, ok := disc["numComments"].(float64); ok {
			discussion.NumComments = int64(numComments)
		}
		if pinned, ok := disc["pinned"].(bool); ok {
			discussion.Pinned = pinned
		}

		discussions = append(discussions, discussion)
	}

	return discussions, nil
}

// fetchFileContentHelper fetches file content from HuggingFace
func fetchFileContentHelper(resourceID, filename string) models.SIBLING {
	sibling := models.SIBLING{
		RFilename:   filename,
		FileContent: "",
	}

	fileURL := fmt.Sprintf("https://huggingface.co/%s/resolve/main/%s", resourceID, filename)
	resp, err := http.Get(fileURL)
	if err != nil {
		return sibling
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return sibling
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return sibling
	}

	sibling.FileContent = string(body)
	return sibling
}

// Helper to scan organization
func scanOrganization(c *fiber.Ctx, org string, includePRs, includeDiscussions bool, scanID string) error {
	// Fetch all models for org
	modelsURL := fmt.Sprintf("https://huggingface.co/api/models?author=%s&full=true", org)
	resp, err := http.Get(modelsURL)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch organization models", nil, "")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to read response", nil, "")
	}

	var modelsData []map[string]interface{}
	if err := json.Unmarshal(body, &modelsData); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to parse response", nil, "")
	}

	var allScannedResources []models.SCANNED_RESOURCE
	var totalFindings int

	// Scan each model (limit to first 5 for performance)
	limit := 5
	if len(modelsData) < limit {
		limit = len(modelsData)
	}

	for i := 0; i < limit; i++ {
		modelData := modelsData[i]
		modelID, ok := modelData["id"].(string)
		if !ok {
			continue
		}

		aiRequest := &models.AI_REQUEST{
			RequestID:   uuid.New().String(),
			Siblings:    []models.SIBLING{},
			Discussions: []models.DISCUSSION{},
		}

		if err := fetchAndAddToRequest(aiRequest, modelID, "models", includePRs, includeDiscussions); err != nil {
			continue
		}

		findings := util.ScanAIRequest(*aiRequest, util.SecretConfig)
		totalFindings += len(findings)

		if len(findings) > 0 {
			scannedResource := models.SCANNED_RESOURCE{
				Type:     "model",
				ID:       modelID,
				Findings: findings,
			}
			allScannedResources = append(allScannedResources, scannedResource)
		}
	}

	// Save scan result
	scanResult := &models.SCAN_RESULT{
		RequestID:        "org-" + org,
		ScannedResources: allScannedResources,
	}

	if err := mgm.Coll(scanResult).Create(scanResult); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save scan results", nil, "")
	}

	// Format response
	formattedResources := []map[string]interface{}{}
	for _, resource := range allScannedResources {
		formattedResources = append(formattedResources, map[string]interface{}{
			"type":     resource.Type,
			"id":       resource.ID,
			"findings": formatFindings(resource.Findings),
		})
	}

	response := map[string]interface{}{
		"scan_id":           scanID,
		"scanned_resources": formattedResources,
		"timestamp":         time.Now().Format(time.RFC3339),
		"total_findings":    totalFindings,
		"models_scanned":    limit,
		"storage_id":        scanResult.ID.Hex(),
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Organization scan completed successfully", response, "")
}

// Helper to format findings as per assignment requirements
func formatFindings(findings []models.Finding) []map[string]interface{} {
	formatted := []map[string]interface{}{}
	for _, finding := range findings {
		item := map[string]interface{}{
			"secret_type": finding.SecretType,
			"pattern":     maskSecret(finding.Secret),
			"secret":      maskSecret(finding.Secret),
		}

		if finding.SourceType == "file" {
			item["file"] = finding.FileName
			item["line"] = finding.Line
		} else if finding.SourceType == "discussion" {
			item["discussion"] = finding.DiscussionTitle
			item["discussion_num"] = finding.DiscussionNum
		}

		formatted = append(formatted, item)
	}
	return formatted
}

// Helper to mask secrets in output
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "****"
	}
	return secret[:4] + "********" + secret[len(secret)-4:]
}

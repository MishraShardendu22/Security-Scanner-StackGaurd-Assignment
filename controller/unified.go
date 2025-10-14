package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/MishraShardendu22/Scanner/models"
	"github.com/MishraShardendu22/Scanner/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kamva/mgm/v3"
)

var httpClient = &http.Client{

	Timeout: 45 * time.Second,
	Transport: &http.Transport{

		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 50,
		MaxConnsPerHost:     100,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
		DisableCompression:  false,
	},
}

func UnifiedScan(c *fiber.Ctx) error {
	var req models.ScanRequestBody

	if err := c.BodyParser(&req); err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid request body", nil, "")
	}

	log.Println("🚀 Starting unified scan request...")

	scanID := fmt.Sprintf("SG-%s-%s", time.Now().Format("2006-0102"), uuid.New().String()[:8])

	requestID := uuid.New().String()

	aiRequest := &models.AI_REQUEST{

		RequestID:   requestID,
		Siblings:    []models.SIBLING{},
		Discussions: []models.DISCUSSION{},
	}
	var scannedResources []models.SCANNED_RESOURCE
	var resourceType, resourceID string

	if req.ModelID != "" {
		resourceType = "model"
		resourceID = req.ModelID

		log.Printf("📦 Fetching model: %s\n", req.ModelID)
		if err := fetchAndAddToRequest(aiRequest, req.ModelID, "models", req.IncludePRs, req.IncludeDiscussions); err != nil {
			return util.ResponseAPI(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to fetch model: %v", err), nil, "")
		}

		log.Printf("✅ Model fetched successfully with %d files\n", len(aiRequest.Siblings))
	} else if req.DatasetID != "" {
		resourceType = "dataset"
		resourceID = req.DatasetID

		log.Printf("📦 Fetching dataset: %s\n", req.DatasetID)
		if err := fetchAndAddToRequest(aiRequest, req.DatasetID, "datasets", req.IncludePRs, req.IncludeDiscussions); err != nil {
			return util.ResponseAPI(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to fetch dataset: %v", err), nil, "")
		}

		log.Printf("✅ Dataset fetched successfully with %d files\n", len(aiRequest.Siblings))
	} else if req.SpaceID != "" {
		resourceType = "space"
		resourceID = req.SpaceID

		log.Printf("📦 Fetching space: %s\n", req.SpaceID)
		if err := fetchAndAddToRequest(aiRequest, req.SpaceID, "spaces", req.IncludePRs, req.IncludeDiscussions); err != nil {
			return util.ResponseAPI(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to fetch space: %v", err), nil, "")
		}

		log.Printf("✅ Space fetched successfully with %d files\n", len(aiRequest.Siblings))
	} else if req.Org != "" {
		log.Printf("🏢 Starting organization scan: %s\n", req.Org)
		return scanOrganization(c, req.Org, req.IncludePRs, req.IncludeDiscussions, scanID)
	} else if req.User != "" {
		log.Printf("👤 Starting user scan: %s\n", req.User)
		return scanOrganization(c, req.User, req.IncludePRs, req.IncludeDiscussions, scanID)
	} else {

		return util.ResponseAPI(c, fiber.StatusBadRequest, "At least one of model_id, dataset_id, space_id, org, or user is required", nil, "")
	}

	log.Println("💾 Saving request to database...")

	if err := mgm.Coll(aiRequest).Create(aiRequest); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save request", nil, "")
	}

	log.Println("🔍 Starting security scan...")

	findings := util.ScanAIRequest(*aiRequest, util.SecretConfig)

	log.Printf("✅ Scan complete! Found %d potential secrets\n", len(findings))

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

	scanResult := &models.SCAN_RESULT{

		RequestID:        requestID,
		ScannedResources: scannedResources,
	}
	if err := mgm.Coll(scanResult).Create(scanResult); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save scan results", nil, "")
	}

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

func StoreScanResult(c *fiber.Ctx) error {

	var scanData map[string]interface{}

	if err := c.BodyParser(&scanData); err != nil {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Invalid request body", nil, "")
	}
	scanIDValue, ok := scanData["scan_id"].(string)

	if !ok {
		scanIDValue = fmt.Sprintf("SG-%s-%s", time.Now().Format("2006-0102"), uuid.New().String()[:8])
	}

	scanResult := &models.SCAN_RESULT{

		RequestID:        scanIDValue,
		ScannedResources: []models.SCANNED_RESOURCE{},
	}
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
	if err := mgm.Coll(scanResult).Create(scanResult); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to store scan results", nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Scan results stored successfully", map[string]interface{}{

		"status":     "stored",
		"scan_id":    scanIDValue,
		"storage_id": scanResult.ID.Hex(),
	}, "")
}

func fetchAndAddToRequest(aiRequest *models.AI_REQUEST, resourceID, resourceType string, includePRs, includeDiscussions bool) error {

	url := fmt.Sprintf("https://huggingface.co/api/%s/%s", resourceType, resourceID)

	log.Printf("📡 Fetching metadata from: %s\n", url)
	resp, err := httpClient.Get(url)

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
	if siblings, ok := resourceData["siblings"].([]interface{}); ok {
		log.Printf("📂 Found %d files to fetch\n", len(siblings))
		var wg sync.WaitGroup
		var mu sync.Mutex
		semaphore := make(chan struct{}, 30)
		for idx, sib := range siblings {
			if sibMap, ok := sib.(map[string]interface{}); ok {
				if filename, ok := sibMap["rfilename"].(string); ok {
					wg.Add(1)
					go func(fname string, index int) {

						defer wg.Done()
						semaphore <- struct{}{}
						defer func() { <-semaphore }()

						log.Printf("  📄 [%d/%d] Fetching file: %s\n", index+1, len(siblings), fname)
						sibling := fetchFileContentHelper(resourceID, fname)
						mu.Lock()
						aiRequest.Siblings = append(aiRequest.Siblings, sibling)
						mu.Unlock()
					}(filename, idx)
				}
			}
		}
		wg.Wait()

		log.Printf("✅ All %d files fetched successfully\n", len(siblings))
	}
	if includePRs || includeDiscussions {
		log.Println("💬 Fetching discussions and PRs...")
		discussions, _ := fetchDiscussionsHelper(resourceID, resourceType, includePRs, includeDiscussions)
		aiRequest.Discussions = discussions

		log.Printf("✅ Fetched %d discussions/PRs\n", len(discussions))
	}

	return nil
}

func fetchDiscussionsHelper(id, resourceType string, includePRs, includeDiscussion bool) ([]models.DISCUSSION, error) {

	var discussions []models.DISCUSSION
	var wg sync.WaitGroup
	var mu sync.Mutex

	if includePRs {
		wg.Add(1)
		go func() {

			defer wg.Done()
			url := fmt.Sprintf("https://huggingface.co/api/%s/%s/discussions?types=pr&status=all", resourceType, id)

			log.Printf("  🔀 Fetching PRs from: %s\n", url)
			prs, _ := fetchDiscussionsFromURL(url)
			mu.Lock()
			discussions = append(discussions, prs...)
			mu.Unlock()

			log.Printf("  ✅ Fetched %d PRs\n", len(prs))
		}()
	}
	if includeDiscussion {
		wg.Add(1)
		go func() {

			defer wg.Done()
			url := fmt.Sprintf("https://huggingface.co/api/%s/%s/discussions?types=discussion&status=all", resourceType, id)

			log.Printf("  💬 Fetching discussions from: %s\n", url)
			discs, _ := fetchDiscussionsFromURL(url)
			mu.Lock()
			discussions = append(discussions, discs...)
			mu.Unlock()

			log.Printf("  ✅ Fetched %d discussions\n", len(discs))
		}()
	}
	wg.Wait()

	return discussions, nil
}

func fetchDiscussionsFromURL(url string) ([]models.DISCUSSION, error) {

	resp, err := httpClient.Get(url)

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

func fetchFileContentHelper(resourceID, filename string) models.SIBLING {

	sibling := models.SIBLING{

		RFilename:   filename,
		FileContent: "",
	}

	fileURL := fmt.Sprintf("https://huggingface.co/%s/resolve/main/%s", resourceID, filename)
	resp, err := httpClient.Get(fileURL)

	if err != nil {
		log.Printf("  ⚠️  Failed to fetch %s: %v\n", filename, err)
		return sibling
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("  ⚠️  File %s returned status %d\n", filename, resp.StatusCode)
		return sibling
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Printf("  ⚠️  Failed to read %s: %v\n", filename, err)
		return sibling
	}
	sibling.FileContent = string(body)

	return sibling
}

func scanOrganization(c *fiber.Ctx, org string, includePRs, includeDiscussions bool, scanID string) error {

	modelsURL := fmt.Sprintf("https://huggingface.co/api/models?author=%s&full=true", org)

	log.Printf("📡 Fetching organization models from: %s\n", modelsURL)
	resp, err := httpClient.Get(modelsURL)

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

	log.Printf("✅ Found %d models for organization\n", len(modelsData))
	var allScannedResources []models.SCANNED_RESOURCE
	var totalFindings int
	var mu sync.Mutex

	limit := 10

	if len(modelsData) < limit {
		limit = len(modelsData)
	}

	log.Printf("🔍 Scanning first %d models concurrently...\n", limit)
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, 10)

	for i := 0; i < limit; i++ {
		modelData := modelsData[i]
		modelID, ok := modelData["id"].(string)
		if !ok {
			continue
		}
		wg.Add(1)
		go func(id string, index int) {

			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			log.Printf("  🔍 [%d/%d] Scanning model: %s\n", index+1, limit, id)
			aiRequest := &models.AI_REQUEST{

				RequestID:   uuid.New().String(),
				Siblings:    []models.SIBLING{},
				Discussions: []models.DISCUSSION{},
			}
			if err := fetchAndAddToRequest(aiRequest, id, "models", includePRs, includeDiscussions); err != nil {
				log.Printf("  ⚠️  Failed to fetch model %s: %v\n", id, err)
				return
			}
			findings := util.ScanAIRequest(*aiRequest, util.SecretConfig)
			mu.Lock()
			totalFindings += len(findings)
			mu.Unlock()
			if len(findings) > 0 {
				log.Printf("  ⚠️  Found %d secrets in model: %s\n", len(findings), id)
				scannedResource := models.SCANNED_RESOURCE{

					Type:     "model",
					ID:       id,
					Findings: findings,
				}
				mu.Lock()
				allScannedResources = append(allScannedResources, scannedResource)
				mu.Unlock()
			} else {

				log.Printf("  ✅ No secrets found in model: %s\n", id)
			}
		}(modelID, i)
	}
	wg.Wait()

	log.Printf("✅ Organization scan complete! Total findings: %d\n", totalFindings)

	scanResult := &models.SCAN_RESULT{

		RequestID:        "org-" + org,
		ScannedResources: allScannedResources,
	}
	if err := mgm.Coll(scanResult).Create(scanResult); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save scan results", nil, "")
	}

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

func formatFindings(findings []models.Finding) []map[string]interface{} {
	formatted := []map[string]interface{}{}

	for _, finding := range findings {
		item := map[string]interface{}{

			"secret_type": finding.SecretType,
			"pattern":     util.MaskSecret(finding.Secret),
			"secret":      util.MaskSecret(finding.Secret),
		}

		switch finding.SourceType {
		case "file":
			item["file"] = finding.FileName
			item["line"] = finding.Line
		case "discussion":
			item["discussion"] = finding.DiscussionTitle
			item["discussion_num"] = finding.DiscussionNum
		}
		formatted = append(formatted, item)
	}

	return formatted
}

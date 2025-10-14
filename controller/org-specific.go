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

// HTTP client for org operations with increased limits
var orgHTTPClient = &http.Client{
	Timeout: 45 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        200, // Increased from 100
		MaxIdleConnsPerHost: 50,  // Increased from 10
		MaxConnsPerHost:     100, // New: limit max connections per host
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
		DisableCompression:  false,
	},
}

func FetchOrgModels(c *fiber.Ctx) error {
	org := c.Params("org")
	if org == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organization name is required", nil, "")
	}

	log.Printf("🏢 Fetching models for organization: %s\n", org)

	url := fmt.Sprintf("https://huggingface.co/api/models?author=%s&full=true", org)
	resp, err := orgHTTPClient.Get(url)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch models", nil, "")
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

	log.Printf("✅ Found %d models\n", len(modelsData))

	includePRs := c.Query("include_prs", "false") == "true"
	includeDiscussion := c.Query("include_discussion", "false") == "true"

	// Save each model to database concurrently
	var savedModels []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 20) // Increased from 10 to 20

	for idx, modelData := range modelsData {
		modelID, ok := modelData["id"].(string)
		if !ok {
			continue
		}

		wg.Add(1)
		go func(id string, index int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			log.Printf("  💾 [%d/%d] Saving model: %s\n", index+1, len(modelsData), id)

			aiModel := &models.AI_Models{
				BaseAI: models.BaseAI{
					Org:               org,
					IncludePRS:        includePRs,
					IncludeDiscussion: includeDiscussion,
				},
				Model_ID: id,
			}

			if err := mgm.Coll(aiModel).Create(aiModel); err == nil {
				mu.Lock()
				savedModels = append(savedModels, id)
				mu.Unlock()
			}
		}(modelID, idx)
	}

	wg.Wait()
	log.Printf("✅ Saved %d models to database\n", len(savedModels))

	return util.ResponseAPI(c, fiber.StatusOK, fmt.Sprintf("Fetched %d models for organization %s", len(modelsData), org), map[string]interface{}{
		"organization": org,
		"count":        len(modelsData),
		"saved_count":  len(savedModels),
		"models":       modelsData,
	}, "")
}

// FetchOrgDatasets fetches all datasets for an organization
func FetchOrgDatasets(c *fiber.Ctx) error {
	org := c.Params("org")
	if org == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organization name is required", nil, "")
	}

	url := fmt.Sprintf("https://huggingface.co/api/datasets?author=%s&full=true", org)
	resp, err := http.Get(url)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch datasets", nil, "")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to read response", nil, "")
	}

	var datasetsData []map[string]interface{}
	if err := json.Unmarshal(body, &datasetsData); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to parse response", nil, "")
	}

	log.Printf("✅ Found %d datasets\n", len(datasetsData))

	includePRs := c.Query("include_prs", "false") == "true"
	includeDiscussion := c.Query("include_discussion", "false") == "true"

	// Save each dataset to database concurrently
	var savedDatasets []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 20) // Increased concurrency

	for idx, datasetData := range datasetsData {
		datasetID, ok := datasetData["id"].(string)
		if !ok {
			continue
		}

		wg.Add(1)
		go func(id string, index int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			log.Printf("  💾 [%d/%d] Saving dataset: %s\n", index+1, len(datasetsData), id)

			aiDataset := &models.AI_DATASETS{
				BaseAI: models.BaseAI{
					Org:               org,
					IncludePRS:        includePRs,
					IncludeDiscussion: includeDiscussion,
				},
				Dataset_ID: id,
			}

			if err := mgm.Coll(aiDataset).Create(aiDataset); err == nil {
				mu.Lock()
				savedDatasets = append(savedDatasets, id)
				mu.Unlock()
			}
		}(datasetID, idx)
	}

	wg.Wait()
	log.Printf("✅ Saved %d datasets to database\n", len(savedDatasets))

	return util.ResponseAPI(c, fiber.StatusOK, fmt.Sprintf("Fetched %d datasets for organization %s", len(datasetsData), org), map[string]interface{}{
		"organization": org,
		"count":        len(datasetsData),
		"saved_count":  len(savedDatasets),
		"datasets":     datasetsData,
	}, "")
}

// FetchOrgSpaces fetches all spaces for an organization
func FetchOrgSpaces(c *fiber.Ctx) error {
	org := c.Params("org")
	if org == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organization name is required", nil, "")
	}

	url := fmt.Sprintf("https://huggingface.co/api/spaces?author=%s&full=true", org)
	resp, err := http.Get(url)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch spaces", nil, "")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to read response", nil, "")
	}

	var spacesData []map[string]interface{}
	if err := json.Unmarshal(body, &spacesData); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to parse response", nil, "")
	}

	log.Printf("✅ Found %d spaces\n", len(spacesData))

	includePRs := c.Query("include_prs", "false") == "true"
	includeDiscussion := c.Query("include_discussion", "false") == "true"

	// Save each space to database concurrently
	var savedSpaces []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 20) // Increased concurrency

	for idx, spaceData := range spacesData {
		spaceID, ok := spaceData["id"].(string)
		if !ok {
			continue
		}

		wg.Add(1)
		go func(id string, index int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			log.Printf("  💾 [%d/%d] Saving space: %s\n", index+1, len(spacesData), id)

			aiSpace := &models.AI_SPACES{
				BaseAI: models.BaseAI{
					Org:               org,
					IncludePRS:        includePRs,
					IncludeDiscussion: includeDiscussion,
				},
				Space_ID: id,
			}

			if err := mgm.Coll(aiSpace).Create(aiSpace); err == nil {
				mu.Lock()
				savedSpaces = append(savedSpaces, id)
				mu.Unlock()
			}
		}(spaceID, idx)
	}

	wg.Wait()
	log.Printf("✅ Saved %d spaces to database\n", len(savedSpaces))

	return util.ResponseAPI(c, fiber.StatusOK, fmt.Sprintf("Fetched %d spaces for organization %s", len(spacesData), org), map[string]interface{}{
		"organization": org,
		"count":        len(spacesData),
		"saved_count":  len(savedSpaces),
		"spaces":       spacesData,
	}, "")
}

// FetchPRs fetches PRs for a specific resource
func FetchPRs(c *fiber.Ctx) error {
	resourceType := c.Params("type") // models, datasets, or spaces
	resourceID := c.Params("id")

	if resourceType == "" || resourceID == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Resource type and ID are required", nil, "")
	}

	url := fmt.Sprintf("https://huggingface.co/api/%s/%s/discussions?types=pr&status=all", resourceType, resourceID)
	resp, err := http.Get(url)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch PRs", nil, "")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to read response", nil, "")
	}

	var rawDiscussions []map[string]interface{}
	if err := json.Unmarshal(body, &rawDiscussions); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to parse response", nil, "")
	}

	// Convert to DISCUSSION model and save
	aiRequest := &models.AI_REQUEST{
		RequestID:   uuid.New().String(),
		Siblings:    []models.SIBLING{},
		Discussions: []models.DISCUSSION{},
	}

	for _, disc := range rawDiscussions {
		discussion := models.DISCUSSION{IsPullRequest: true}

		if num, ok := disc["num"].(float64); ok {
			discussion.Num = int64(num)
		}
		if title, ok := disc["title"].(string); ok {
			discussion.Title = title
		}
		if status, ok := disc["status"].(string); ok {
			discussion.Status = status
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

		aiRequest.Discussions = append(aiRequest.Discussions, discussion)
	}

	if err := mgm.Coll(aiRequest).Create(aiRequest); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save to database", nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "PRs fetched successfully", map[string]interface{}{
		"request_id": aiRequest.RequestID,
		"count":      len(aiRequest.Discussions),
		"prs":        aiRequest.Discussions,
	}, "")
}

// FetchDiscussions fetches discussions for a specific resource
func FetchDiscussions(c *fiber.Ctx) error {
	resourceType := c.Params("type") // models, datasets, or spaces
	resourceID := c.Params("id")

	if resourceType == "" || resourceID == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Resource type and ID are required", nil, "")
	}

	url := fmt.Sprintf("https://huggingface.co/api/%s/%s/discussions?types=discussion&status=all", resourceType, resourceID)
	resp, err := http.Get(url)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch discussions", nil, "")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to read response", nil, "")
	}

	var rawDiscussions []map[string]interface{}
	if err := json.Unmarshal(body, &rawDiscussions); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to parse response", nil, "")
	}

	aiRequest := &models.AI_REQUEST{
		RequestID:   uuid.New().String(),
		Siblings:    []models.SIBLING{},
		Discussions: []models.DISCUSSION{},
	}

	for _, disc := range rawDiscussions {
		discussion := models.DISCUSSION{IsPullRequest: false}

		if num, ok := disc["num"].(float64); ok {
			discussion.Num = int64(num)
		}
		if title, ok := disc["title"].(string); ok {
			discussion.Title = title
		}
		if status, ok := disc["status"].(string); ok {
			discussion.Status = status
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

		aiRequest.Discussions = append(aiRequest.Discussions, discussion)
	}

	if err := mgm.Coll(aiRequest).Create(aiRequest); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save to database", nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Discussions fetched successfully", map[string]interface{}{
		"request_id":  aiRequest.RequestID,
		"count":       len(aiRequest.Discussions),
		"discussions": aiRequest.Discussions,
	}, "")
}

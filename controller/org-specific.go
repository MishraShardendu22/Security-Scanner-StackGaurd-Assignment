package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MishraShardendu22/Scanner/models"
	"github.com/MishraShardendu22/Scanner/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kamva/mgm/v3"
)

// FetchOrgModels fetches all models for an organization
func FetchOrgModels(c *fiber.Ctx) error {
	org := c.Params("org")
	if org == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organization name is required", nil, "")
	}

	url := fmt.Sprintf("https://huggingface.co/api/models?author=%s&full=true", org)
	resp, err := http.Get(url)
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

	includePRs := c.Query("include_prs", "false") == "true"
	includeDiscussion := c.Query("include_discussion", "false") == "true"

	// Save each model to database
	var savedModels []string
	for _, modelData := range modelsData {
		modelID, ok := modelData["id"].(string)
		if !ok {
			continue
		}

		aiModel := &models.AI_Models{
			BaseAI: models.BaseAI{
				Org:               org,
				IncludePRS:        includePRs,
				IncludeDiscussion: includeDiscussion,
			},
			Model_ID: modelID,
		}

		if err := mgm.Coll(aiModel).Create(aiModel); err == nil {
			savedModels = append(savedModels, modelID)
		}
	}

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

	includePRs := c.Query("include_prs", "false") == "true"
	includeDiscussion := c.Query("include_discussion", "false") == "true"

	var savedDatasets []string
	for _, datasetData := range datasetsData {
		datasetID, ok := datasetData["id"].(string)
		if !ok {
			continue
		}

		aiDataset := &models.AI_DATASETS{
			BaseAI: models.BaseAI{
				Org:               org,
				IncludePRS:        includePRs,
				IncludeDiscussion: includeDiscussion,
			},
			Dataset_ID: datasetID,
		}

		if err := mgm.Coll(aiDataset).Create(aiDataset); err == nil {
			savedDatasets = append(savedDatasets, datasetID)
		}
	}

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

	includePRs := c.Query("include_prs", "false") == "true"
	includeDiscussion := c.Query("include_discussion", "false") == "true"

	var savedSpaces []string
	for _, spaceData := range spacesData {
		spaceID, ok := spaceData["id"].(string)
		if !ok {
			continue
		}

		aiSpace := &models.AI_SPACES{
			BaseAI: models.BaseAI{
				Org:               org,
				IncludePRS:        includePRs,
				IncludeDiscussion: includeDiscussion,
			},
			Space_ID: spaceID,
		}

		if err := mgm.Coll(aiSpace).Create(aiSpace); err == nil {
			savedSpaces = append(savedSpaces, spaceID)
		}
	}

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

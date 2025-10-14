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

// FetchModel fetches a single model from HuggingFace API
func FetchModel(c *fiber.Ctx) error {
	modelID := c.Params("modelId")
	if modelID == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Model ID is required", nil, "")
	}

	url := fmt.Sprintf("https://huggingface.co/api/models/%s", modelID)
	resp, err := http.Get(url)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch model", nil, "")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return util.ResponseAPI(c, resp.StatusCode, "Model not found", nil, "")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to read response", nil, "")
	}

	var modelData map[string]interface{}
	if err := json.Unmarshal(body, &modelData); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to parse response", nil, "")
	}

	// Parse query params for PR and discussion
	includePRs := c.Query("include_prs", "false") == "true"
	includeDiscussion := c.Query("include_discussion", "false") == "true"

	// Create AI_REQUEST record
	aiRequest := &models.AI_REQUEST{
		RequestID:   uuid.New().String(),
		Siblings:    []models.SIBLING{},
		Discussions: []models.DISCUSSION{},
	}

	// Extract siblings
	if siblings, ok := modelData["siblings"].([]interface{}); ok {
		for _, sib := range siblings {
			if sibMap, ok := sib.(map[string]interface{}); ok {
				if filename, ok := sibMap["rfilename"].(string); ok {
					aiRequest.Siblings = append(aiRequest.Siblings, models.SIBLING{RFilename: filename})
				}
			}
		}
	}

	// Fetch discussions/PRs if requested
	if includePRs || includeDiscussion {
		discussions, _ := fetchDiscussions(modelID, "models", includePRs, includeDiscussion)
		aiRequest.Discussions = discussions
	}

	// Save to database
	if err := mgm.Coll(aiRequest).Create(aiRequest); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save to database", nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Model fetched successfully", map[string]interface{}{
		"request_id":  aiRequest.RequestID,
		"model":       modelData,
		"siblings":    aiRequest.Siblings,
		"discussions": aiRequest.Discussions,
	}, "")
}

// FetchDataset fetches a single dataset from HuggingFace API
func FetchDataset(c *fiber.Ctx) error {
	datasetID := c.Params("datasetId")
	if datasetID == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Dataset ID is required", nil, "")
	}

	url := fmt.Sprintf("https://huggingface.co/api/datasets/%s", datasetID)
	resp, err := http.Get(url)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch dataset", nil, "")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return util.ResponseAPI(c, resp.StatusCode, "Dataset not found", nil, "")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to read response", nil, "")
	}

	var datasetData map[string]interface{}
	if err := json.Unmarshal(body, &datasetData); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to parse response", nil, "")
	}

	includePRs := c.Query("include_prs", "false") == "true"
	includeDiscussion := c.Query("include_discussion", "false") == "true"

	aiRequest := &models.AI_REQUEST{
		RequestID:   uuid.New().String(),
		Siblings:    []models.SIBLING{},
		Discussions: []models.DISCUSSION{},
	}

	if siblings, ok := datasetData["siblings"].([]interface{}); ok {
		for _, sib := range siblings {
			if sibMap, ok := sib.(map[string]interface{}); ok {
				if filename, ok := sibMap["rfilename"].(string); ok {
					aiRequest.Siblings = append(aiRequest.Siblings, models.SIBLING{RFilename: filename})
				}
			}
		}
	}

	if includePRs || includeDiscussion {
		discussions, _ := fetchDiscussions(datasetID, "datasets", includePRs, includeDiscussion)
		aiRequest.Discussions = discussions
	}

	if err := mgm.Coll(aiRequest).Create(aiRequest); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save to database", nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Dataset fetched successfully", map[string]interface{}{
		"request_id":  aiRequest.RequestID,
		"dataset":     datasetData,
		"siblings":    aiRequest.Siblings,
		"discussions": aiRequest.Discussions,
	}, "")
}

// FetchSpace fetches a single space from HuggingFace API
func FetchSpace(c *fiber.Ctx) error {
	spaceID := c.Params("spaceId")
	if spaceID == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Space ID is required", nil, "")
	}

	url := fmt.Sprintf("https://huggingface.co/api/spaces/%s", spaceID)
	resp, err := http.Get(url)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to fetch space", nil, "")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return util.ResponseAPI(c, resp.StatusCode, "Space not found", nil, "")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to read response", nil, "")
	}

	var spaceData map[string]interface{}
	if err := json.Unmarshal(body, &spaceData); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to parse response", nil, "")
	}

	includePRs := c.Query("include_prs", "false") == "true"
	includeDiscussion := c.Query("include_discussion", "false") == "true"

	aiRequest := &models.AI_REQUEST{
		RequestID:   uuid.New().String(),
		Siblings:    []models.SIBLING{},
		Discussions: []models.DISCUSSION{},
	}

	if siblings, ok := spaceData["siblings"].([]interface{}); ok {
		for _, sib := range siblings {
			if sibMap, ok := sib.(map[string]interface{}); ok {
				if filename, ok := sibMap["rfilename"].(string); ok {
					aiRequest.Siblings = append(aiRequest.Siblings, models.SIBLING{RFilename: filename})
				}
			}
		}
	}

	if includePRs || includeDiscussion {
		discussions, _ := fetchDiscussions(spaceID, "spaces", includePRs, includeDiscussion)
		aiRequest.Discussions = discussions
	}

	if err := mgm.Coll(aiRequest).Create(aiRequest); err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, "Failed to save to database", nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Space fetched successfully", map[string]interface{}{
		"request_id":  aiRequest.RequestID,
		"space":       spaceData,
		"siblings":    aiRequest.Siblings,
		"discussions": aiRequest.Discussions,
	}, "")
}

// Helper function to fetch discussions and PRs
func fetchDiscussions(id, resourceType string, includePRs, includeDiscussion bool) ([]models.DISCUSSION, error) {
	var discussions []models.DISCUSSION

	if includePRs {
		url := fmt.Sprintf("https://huggingface.co/api/%s/%s/discussions?types=pr&status=all", resourceType, id)
		prs, _ := getDiscussionsFromURL(url)
		discussions = append(discussions, prs...)
	}

	if includeDiscussion {
		url := fmt.Sprintf("https://huggingface.co/api/%s/%s/discussions?types=discussion&status=all", resourceType, id)
		discs, _ := getDiscussionsFromURL(url)
		discussions = append(discussions, discs...)
	}

	return discussions, nil
}

func getDiscussionsFromURL(url string) ([]models.DISCUSSION, error) {
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

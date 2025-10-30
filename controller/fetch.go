package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/MishraShardendu22/Scanner/models"
	"github.com/MishraShardendu22/Scanner/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kamva/mgm/v3"
)

var fetchHTTPClient = util.SharedHTTPClient()

func FetchModel(c *fiber.Ctx) error {

	modelID := c.Params("modelId")

	if modelID == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Model ID is required", nil, "")
	}

	log.Printf("🚀 Fetching model: %s\n", modelID)

	url := fmt.Sprintf("https://huggingface.co/api/models/%s", modelID)
	resp, err := fetchHTTPClient.Get(url)

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

	includePRs, includeDiscussion := util.ParseIncludeFlags(c)

	aiRequest := &models.AI_REQUEST{

		RequestID:   uuid.New().String(),
		Siblings:    []models.SIBLING{},
		Discussions: []models.DISCUSSION{},
	}

	if siblings, ok := modelData["siblings"].([]interface{}); ok {
		aiRequest.Siblings = util.FetchFilesFromSiblings(modelID, siblings)
	}

	if includePRs || includeDiscussion {
		log.Println("💬 Fetching discussions/PRs...")
		discussions, _ := util.FetchDiscussions(modelID, "models", includePRs, includeDiscussion)
		aiRequest.Discussions = discussions
		log.Printf("✅ Fetched %d discussions/PRs\n", len(discussions))
	}

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

	includePRs, includeDiscussion := util.ParseIncludeFlags(c)

	log.Printf("🚀 Fetching dataset: %s\n", datasetID)

	aiRequest := &models.AI_REQUEST{

		RequestID:   uuid.New().String(),
		Siblings:    []models.SIBLING{},
		Discussions: []models.DISCUSSION{},
	}
	if siblings, ok := datasetData["siblings"].([]interface{}); ok {
		aiRequest.Siblings = util.FetchFilesFromSiblings(datasetID, siblings)
	}
	if includePRs || includeDiscussion {
		log.Println("💬 Fetching discussions/PRs...")
		discussions, _ := util.FetchDiscussions(datasetID, "datasets", includePRs, includeDiscussion)
		aiRequest.Discussions = discussions
		log.Printf("✅ Fetched %d discussions/PRs\n", len(discussions))
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

	includePRs, includeDiscussion := util.ParseIncludeFlags(c)

	log.Printf("🚀 Fetching space: %s\n", spaceID)

	aiRequest := &models.AI_REQUEST{

		RequestID:   uuid.New().String(),
		Siblings:    []models.SIBLING{},
		Discussions: []models.DISCUSSION{},
	}
	if siblings, ok := spaceData["siblings"].([]interface{}); ok {
		aiRequest.Siblings = util.FetchFilesFromSiblings(spaceID, siblings)
	}
	if includePRs || includeDiscussion {
		log.Println("💬 Fetching discussions/PRs...")
		discussions, _ := util.FetchDiscussions(spaceID, "spaces", includePRs, includeDiscussion)
		aiRequest.Discussions = discussions
		log.Printf("✅ Fetched %d discussions/PRs\n", len(discussions))
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
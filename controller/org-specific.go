package controller

import (
	"fmt"
	"log"

	"github.com/MishraShardendu22/Scanner/util"
	"github.com/gofiber/fiber/v2"
)

func FetchOrgModels(c *fiber.Ctx) error {
	org := c.Params("org")
	if org == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organization name is required", nil, "")
	}

	log.Printf("🏢 Fetching models for organization: %s\n", org)
	includePRs, includeDiscussion := util.ParseIncludeFlags(c)

	modelsData, savedModels, err := util.FetchOrgResources(c, org, util.ResourceTypeModel, includePRs, includeDiscussion)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, err.Error(), nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, fmt.Sprintf("Fetched %d models for organization %s", len(modelsData), org), map[string]interface{}{
		"organization": org,
		"count":        len(modelsData),
		"saved_count":  len(savedModels),
		"models":       modelsData,
	}, "")
}

func FetchOrgDatasets(c *fiber.Ctx) error {
	org := c.Params("org")
	if org == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organization name is required", nil, "")
	}

	includePRs, includeDiscussion := util.ParseIncludeFlags(c)

	datasetsData, savedDatasets, err := util.FetchOrgResources(c, org, util.ResourceTypeDataset, includePRs, includeDiscussion)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, err.Error(), nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, fmt.Sprintf("Fetched %d datasets for organization %s", len(datasetsData), org), map[string]interface{}{
		"organization": org,
		"count":        len(datasetsData),
		"saved_count":  len(savedDatasets),
		"datasets":     datasetsData,
	}, "")
}

func FetchOrgSpaces(c *fiber.Ctx) error {
	org := c.Params("org")
	if org == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Organization name is required", nil, "")
	}

	includePRs, includeDiscussion := util.ParseIncludeFlags(c)

	spacesData, savedSpaces, err := util.FetchOrgResources(c, org, util.ResourceTypeSpace, includePRs, includeDiscussion)
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, err.Error(), nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, fmt.Sprintf("Fetched %d spaces for organization %s", len(spacesData), org), map[string]interface{}{
		"organization": org,
		"count":        len(spacesData),
		"saved_count":  len(savedSpaces),
		"spaces":       spacesData,
	}, "")
}

func FetchPRs(c *fiber.Ctx) error {
	resourceType := c.Params("type")
	resourceID := c.Params("id")

	if resourceType == "" || resourceID == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Resource type and ID are required", nil, "")
	}

	aiRequest, err := util.FetchAndSaveDiscussionsByType(resourceType, resourceID, "pr")
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, err.Error(), nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "PRs fetched successfully", map[string]interface{}{
		"request_id": aiRequest.RequestID,
		"count":      len(aiRequest.Discussions),
		"prs":        aiRequest.Discussions,
	}, "")
}

func FetchDiscussions(c *fiber.Ctx) error {
	resourceType := c.Params("type")
	resourceID := c.Params("id")

	if resourceType == "" || resourceID == "" {
		return util.ResponseAPI(c, fiber.StatusBadRequest, "Resource type and ID are required", nil, "")
	}

	aiRequest, err := util.FetchAndSaveDiscussionsByType(resourceType, resourceID, "discussion")
	if err != nil {
		return util.ResponseAPI(c, fiber.StatusInternalServerError, err.Error(), nil, "")
	}

	return util.ResponseAPI(c, fiber.StatusOK, "Discussions fetched successfully", map[string]interface{}{
		"request_id":  aiRequest.RequestID,
		"count":       len(aiRequest.Discussions),
		"discussions": aiRequest.Discussions,
	}, "")
}

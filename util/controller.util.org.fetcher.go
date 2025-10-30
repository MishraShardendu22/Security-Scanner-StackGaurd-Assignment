package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/MishraShardendu22/Scanner/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kamva/mgm/v3"
)

type ResourceType string

const (
	ResourceTypeModel   ResourceType = "models"
	ResourceTypeSpace   ResourceType = "spaces"
	ResourceTypeDataset ResourceType = "datasets"
)

func FetchOrgResources(
	c *fiber.Ctx,
	org string,
	resourceType ResourceType,
	includePRs, includeDiscussion bool,
) ([]map[string]interface{}, []string, error) {
	url := fmt.Sprintf("https://huggingface.co/api/%s?author=%s&full=true", resourceType, org)
	httpClient := SharedHTTPClient()
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch %s", resourceType)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response")
	}

	var resourcesData []map[string]interface{}
	if err := json.Unmarshal(body, &resourcesData); err != nil {
		return nil, nil, fmt.Errorf("failed to parse response")
	}

	log.Printf("✅ Found %d %s\n", len(resourcesData), resourceType)

	var saved []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 20)

	for idx, resourceData := range resourcesData {
		resourceID, ok := resourceData["id"].(string)
		if !ok {
			continue
		}
		wg.Add(1)
		go func(id string, index int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			log.Printf("  💾 [%d/%d] Saving %s: %s\n", index+1, len(resourcesData), resourceType, id)

			switch resourceType {
			case ResourceTypeModel:
				model := &models.AI_Models{
					BaseAI: models.BaseAI{
						Org:               org,
						IncludePRS:        includePRs,
						IncludeDiscussion: includeDiscussion,
					},
					Model_ID: id,
				}
				if err := mgm.Coll(model).Create(model); err == nil {
					mu.Lock()
					saved = append(saved, id)
					mu.Unlock()
				}
			case ResourceTypeDataset:
				model := &models.AI_DATASETS{
					BaseAI: models.BaseAI{
						Org:               org,
						IncludePRS:        includePRs,
						IncludeDiscussion: includeDiscussion,
					},
					Dataset_ID: id,
				}
				if err := mgm.Coll(model).Create(model); err == nil {
					mu.Lock()
					saved = append(saved, id)
					mu.Unlock()
				}
			case ResourceTypeSpace:
				model := &models.AI_SPACES{
					BaseAI: models.BaseAI{
						Org:               org,
						IncludePRS:        includePRs,
						IncludeDiscussion: includeDiscussion,
					},
					Space_ID: id,
				}
				if err := mgm.Coll(model).Create(model); err == nil {
					mu.Lock()
					saved = append(saved, id)
					mu.Unlock()
				}
			}
		}(resourceID, idx)
	}
	wg.Wait()

	log.Printf("✅ Saved %d %s to database\n", len(saved), resourceType)
	return resourcesData, saved, nil
}

func FetchSingleResource(
	resourceID string,
	resourceType ResourceType,
	includePRs, includeDiscussion bool,
) (*models.AI_REQUEST, map[string]interface{}, error) {
	url := fmt.Sprintf("https://huggingface.co/api/%s/%s", resourceType, resourceID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch %s", resourceType)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("%s not found", resourceType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response")
	}

	var resourceData map[string]interface{}
	if err := json.Unmarshal(body, &resourceData); err != nil {
		return nil, nil, fmt.Errorf("failed to parse response")
	}

	log.Printf("🚀 Fetching %s: %s\n", resourceType, resourceID)

	aiRequest := &models.AI_REQUEST{
		RequestID:    uuid.New().String(),
		ResourceType: string(resourceType),
		ResourceID:   resourceID,
		Siblings:     []models.SIBLING{},
		Discussions:  []models.DISCUSSION{},
	}

	if siblings, ok := resourceData["siblings"].([]interface{}); ok {
		aiRequest.Siblings = FetchFilesFromSiblings(resourceID, siblings)
	}

	if includePRs || includeDiscussion {
		log.Println("💬 Fetching discussions/PRs...")
		discussions, _ := FetchDiscussions(resourceID, string(resourceType), includePRs, includeDiscussion)
		aiRequest.Discussions = discussions
		log.Printf("✅ Fetched %d discussions/PRs\n", len(discussions))
	}

	return aiRequest, resourceData, nil
}

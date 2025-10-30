package util

import (
	"fmt"

	"github.com/MishraShardendu22/Scanner/models"
	"github.com/google/uuid"
	"github.com/kamva/mgm/v3"
)

func FetchAndSaveDiscussionsByType(resourceType, resourceID, discussionType string) (*models.AI_REQUEST, error) {
	url := fmt.Sprintf("https://huggingface.co/api/%s/%s/discussions?types=%s&status=all", resourceType, resourceID, discussionType)

	discussions, err := GetDiscussionsFromURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s", discussionType)
	}

	aiRequest := &models.AI_REQUEST{
		RequestID:   uuid.New().String(),
		Siblings:    []models.SIBLING{},
		Discussions: discussions,
	}

	if err := mgm.Coll(aiRequest).Create(aiRequest); err != nil {
		return nil, fmt.Errorf("failed to save to database")
	}

	return aiRequest, nil
}

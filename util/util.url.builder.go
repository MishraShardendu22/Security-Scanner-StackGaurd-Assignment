package util

import "fmt"

func BuildHuggingFaceFileURL(resourceType, resourceID, fileName string, lineNumber int) string {
	baseURL := fmt.Sprintf("https://huggingface.co/%s/blob/main/%s", resourceID, fileName)
	if lineNumber > 0 {
		return fmt.Sprintf("%s?line=%d", baseURL, lineNumber)
	}
	return baseURL
}

func BuildHuggingFaceDiscussionURL(resourceType, resourceID string, discussionNum int64) string {
	return fmt.Sprintf("https://huggingface.co/%s/%s/discussions/%d", resourceType, resourceID, discussionNum)
}

func ExtractOrgFromResourceID(resourceID string) string {
	for i := 0; i < len(resourceID); i++ {
		if resourceID[i] == '/' {
			return resourceID[:i]
		}
	}
	return resourceID
}

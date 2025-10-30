package util

import (
	"github.com/MishraShardendu22/Scanner/models"
)

func GroupFindingsByResource(findings []models.Finding) []models.SCANNED_RESOURCE {
	resourceMap := make(map[string]*models.SCANNED_RESOURCE)

	for _, finding := range findings {
		var resourceKey string
		var resourceType string
		var resourceID string

		switch finding.SourceType {
		case "file":
			resourceType = "file"
			resourceID = finding.FileName
			resourceKey = "file:" + finding.FileName
		case "discussion":
			resourceType = "discussion"
			resourceID = finding.DiscussionTitle
			resourceKey = "discussion:" + finding.DiscussionTitle
		}

		if _, exists := resourceMap[resourceKey]; !exists {
			resourceMap[resourceKey] = &models.SCANNED_RESOURCE{
				Type:     resourceType,
				ID:       resourceID,
				Findings: []models.Finding{},
			}
		}

		resourceMap[resourceKey].Findings = append(resourceMap[resourceKey].Findings, finding)
	}

	var scannedResources []models.SCANNED_RESOURCE
	for _, resource := range resourceMap {
		scannedResources = append(scannedResources, *resource)
	}

	return scannedResources
}

func CountFindingsByType(findings []models.Finding) map[string]int {
	findingsByType := make(map[string]int)
	for _, finding := range findings {
		if finding.SecretType != "" {
			findingsByType[finding.SecretType]++
		}
	}
	return findingsByType
}

func CountFindingsBySource(findings []models.Finding) map[string]int {
	findingsBySource := make(map[string]int)
	for _, finding := range findings {
		findingsBySource[finding.SourceType]++
	}
	return findingsBySource
}

func FormatFindings(findings []models.Finding) []map[string]interface{} {
	formatted := []map[string]interface{}{}

	for _, finding := range findings {
		item := map[string]interface{}{
			"secret_type": finding.SecretType,
			"pattern":     finding.Secret,
			"secret":      finding.Secret,
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

func CountTotalFindings(scannedResources []models.SCANNED_RESOURCE) int {
	total := 0
	for _, resource := range scannedResources {
		total += len(resource.Findings)
	}
	return total
}

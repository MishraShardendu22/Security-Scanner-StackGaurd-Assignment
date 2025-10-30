package util

import (
	"fmt"
	"log"
	"sync"

	"github.com/MishraShardendu22/Scanner/models"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

// ScanOrgResources scans all resources for an organization
func ScanOrgResources(
	org string,
	resourceType ResourceType,
) ([]models.Finding, int, error) {
	log.Printf("🏢 Starting organization %s scan: %s\n", resourceType, org)

	var count int64
	var err error

	switch resourceType {
	case ResourceTypeModel:
		count, err = mgm.Coll(&models.AI_Models{}).CountDocuments(mgm.Ctx(), bson.M{"org": org})
	case ResourceTypeDataset:
		count, err = mgm.Coll(&models.AI_DATASETS{}).CountDocuments(mgm.Ctx(), bson.M{"org": org})
	case ResourceTypeSpace:
		count, err = mgm.Coll(&models.AI_SPACES{}).CountDocuments(mgm.Ctx(), bson.M{"org": org})
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch organization %s", resourceType)
	}
	if count == 0 {
		return nil, 0, fmt.Errorf("no %s found for this organization", resourceType)
	}

	log.Printf("✅ Found %d %s for organization\n", count, resourceType)

	var allFindings []models.Finding
	var scannedCount int

	aiRequests := []models.AI_REQUEST{}
	mgm.Coll(&models.AI_REQUEST{}).SimpleFind(&aiRequests, bson.M{})

	log.Printf("🔍 Scanning %d requests...\n", len(aiRequests))
	var wg sync.WaitGroup
	var mu sync.Mutex
	semaphore := make(chan struct{}, 10)

	for idx, req := range aiRequests {
		wg.Add(1)
		go func(r models.AI_REQUEST, index int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			log.Printf("  [%d/%d] Scanning request: %s\n", index+1, len(aiRequests), r.RequestID)
			findings := ScanAIRequest(r, SecretConfig)
			mu.Lock()
			allFindings = append(allFindings, findings...)
			scannedCount++
			mu.Unlock()
			if len(findings) > 0 {
				log.Printf("    ⚠️  Found %d secrets\n", len(findings))
			}
		}(req, idx)
	}
	wg.Wait()

	log.Printf("✅ Scan complete! Total findings: %d\n", len(allFindings))
	return allFindings, scannedCount, nil
}

// SaveScanResults saves scan results to database
func SaveScanResults(requestID string, scannedResources []models.SCANNED_RESOURCE) (*models.SCAN_RESULT, error) {
	scanResult := &models.SCAN_RESULT{
		RequestID:        requestID,
		ScannedResources: scannedResources,
	}
	if err := mgm.Coll(scanResult).Create(scanResult); err != nil {
		return nil, fmt.Errorf("failed to save scan results")
	}
	return scanResult, nil
}

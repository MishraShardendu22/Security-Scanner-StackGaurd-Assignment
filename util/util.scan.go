package util

import (
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/MishraShardendu22/Scanner/models"
	util_model "github.com/MishraShardendu22/Scanner/util/model"
)

func ScanFile(file models.SIBLING, patterns []util_model.SecretPattern, resourceType, resourceID string) []models.Finding {

	ext := strings.ToLower(filepath.Ext(file.RFilename))

	if !TextExtensions[ext] {
		return nil
	}
	var findings []models.Finding

	lines := strings.Split(file.FileContent, "\n")

	organization := ExtractOrgFromResourceID(resourceID)

	for i, line := range lines {
		for _, pattern := range patterns {
			re := regexp.MustCompile(pattern.Regex)
			matches := re.FindAllString(line, -1)
			for _, match := range matches {
				lineNum := i + 1
				findings = append(findings, models.Finding{
					SecretType:   pattern.Name,
					Pattern:      pattern.Regex,
					Secret:       match,
					SourceType:   "file",
					Organization: organization,
					ResourceID:   resourceID,
					ResourceType: resourceType,
					FileName:     file.RFilename,
					Line:         lineNum,
					URL:          BuildHuggingFaceFileURL(resourceType, resourceID, file.RFilename, lineNum),
				})
			}
		}
	}

	return findings
}

func ScanDiscussion(disc models.DISCUSSION, patterns []util_model.SecretPattern, resourceType, resourceID string) []models.Finding {

	var findings []models.Finding

	text := disc.Title + " " + disc.RepoName

	organization := ExtractOrgFromResourceID(resourceID)

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern.Regex)
		matches := re.FindAllString(text, -1)
		for _, match := range matches {
			findings = append(findings, models.Finding{
				SecretType:      pattern.Name,
				Pattern:         pattern.Regex,
				Secret:          match,
				SourceType:      "discussion",
				Organization:    organization,
				ResourceID:      resourceID,
				ResourceType:    resourceType,
				DiscussionNum:   disc.Num,
				DiscussionTitle: disc.Title,
				DiscussionRepo:  disc.RepoName,
				URL:             BuildHuggingFaceDiscussionURL(resourceType, resourceID, disc.Num),
			})
		}
	}

	return findings
}

func ScanAIRequest(req models.AI_REQUEST, patterns []util_model.SecretPattern, resourceType, resourceID string) []models.Finding {

	var wg sync.WaitGroup

	ch := make(chan []models.Finding, 10)

	results := []models.Finding{}

	totalItems := len(req.Siblings) + len(req.Discussions)
	var scannedCount int32

	log.Printf("  🔍 Scanning %d files and %d discussions...\n", len(req.Siblings), len(req.Discussions))

	semaphore := make(chan struct{}, 50)

	for _, f := range req.Siblings {
		wg.Add(1)
		go func(file models.SIBLING) {

			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			findings := ScanFile(file, patterns, resourceType, resourceID)
			ch <- findings
			count := atomic.AddInt32(&scannedCount, 1)
			if len(findings) > 0 {
				log.Printf("    [%d/%d] ⚠️  %s: Found %d secrets\n", count, totalItems, file.RFilename, len(findings))
			} else if count%10 == 0 {
				log.Printf("    [%d/%d] Scanned...\n", count, totalItems)
			}
		}(f)
	}

	for _, d := range req.Discussions {
		wg.Add(1)
		go func(disc models.DISCUSSION) {

			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			findings := ScanDiscussion(disc, patterns, resourceType, resourceID)
			ch <- findings
			count := atomic.AddInt32(&scannedCount, 1)
			if len(findings) > 0 {
				log.Printf("    [%d/%d] ⚠️  Discussion '%s': Found %d secrets\n", count, totalItems, disc.Title, len(findings))
			}
		}(d)
	}
	go func() {

		wg.Wait()
		close(ch)
	}()

	for f := range ch {
		results = append(results, f...)
	}

	return results
}

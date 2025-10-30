package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/MishraShardendu22/Scanner/models"
)

var httpClient = SharedHTTPClient()

func FetchDiscussions(id, resourceType string, includePRs, includeDiscussion bool) ([]models.DISCUSSION, error) {
	var discussions []models.DISCUSSION
	var wg sync.WaitGroup
	var mu sync.Mutex

	if includePRs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url := fmt.Sprintf("https://huggingface.co/api/%s/%s/discussions?types=pr&status=all", resourceType, id)
			log.Printf("  🔀 Fetching PRs from: %s\n", url)
			prs, _ := GetDiscussionsFromURL(url)
			mu.Lock()
			discussions = append(discussions, prs...)
			mu.Unlock()
			log.Printf("  ✅ Fetched %d PRs\n", len(prs))
		}()
	}

	if includeDiscussion {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url := fmt.Sprintf("https://huggingface.co/api/%s/%s/discussions?types=discussion&status=all", resourceType, id)
			log.Printf("  💬 Fetching discussions from: %s\n", url)
			discs, _ := GetDiscussionsFromURL(url)
			mu.Lock()
			discussions = append(discussions, discs...)
			mu.Unlock()
			log.Printf("  ✅ Fetched %d discussions\n", len(discs))
		}()
	}

	wg.Wait()
	return discussions, nil
}

func GetDiscussionsFromURL(url string) ([]models.DISCUSSION, error) {
	resp, err := httpClient.Get(url)
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

func FetchFileContent(resourceID, filename string) models.SIBLING {
	sibling := models.SIBLING{
		RFilename:   filename,
		FileContent: "",
	}

	fileURL := fmt.Sprintf("https://huggingface.co/%s/resolve/main/%s", resourceID, filename)
	resp, err := httpClient.Get(fileURL)
	if err != nil {
		log.Printf("  ⚠️  Failed to fetch %s: %v\n", filename, err)
		return sibling
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("  ⚠️  File %s returned status %d\n", filename, resp.StatusCode)
		return sibling
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("  ⚠️  Failed to read %s: %v\n", filename, err)
		return sibling
	}

	sibling.FileContent = string(body)
	return sibling
}

func FetchFilesFromSiblings(resourceID string, siblings []interface{}) []models.SIBLING {
	var result []models.SIBLING
	var wg sync.WaitGroup
	var mu sync.Mutex
	semaphore := make(chan struct{}, 20)

	log.Printf("📂 Found %d files to fetch\n", len(siblings))

	for idx, sib := range siblings {
		if sibMap, ok := sib.(map[string]interface{}); ok {
			if filename, ok := sibMap["rfilename"].(string); ok {

				ext := strings.ToLower(filepath.Ext(filename))
				if !TextExtensions[ext] {
					log.Printf("  ⏭️  [SKIP] Non-readable: %s", filename)
					continue
				}
				wg.Add(1)
				go func(fname string, index int) {
					defer wg.Done()
					semaphore <- struct{}{}
					defer func() { <-semaphore }()

					log.Printf("  📄 [%d/%d] Fetching: %s\n", index+1, len(siblings), fname)
					sibling := FetchFileContent(resourceID, fname)
					mu.Lock()
					result = append(result, sibling)
					mu.Unlock()
				}(filename, idx)
			}
		}
	}

	wg.Wait()
	log.Printf("✅ All %d files fetched\n", len(siblings))
	return result
}

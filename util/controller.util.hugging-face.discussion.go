package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/MishraShardendu22/Scanner/models"
	"github.com/google/uuid"
)

func FetchDiscussions(id, resourceType string, includePRs, includeDiscussion bool) ([]models.DISCUSSION, error) {
	var discussions []models.DISCUSSION
	var wg sync.WaitGroup
	var mu sync.Mutex

	start := time.Now()
	requestID := uuid.New().String()

	log.Printf(
		"op=FetchDiscussions stage=start request_id=%s resource_type=%s id=%s include_prs=%t include_discussion=%t",
		requestID, resourceType, id, includePRs, includeDiscussion,
	)

	if includePRs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url := fmt.Sprintf("https://huggingface.co/api/%s/%s/discussions?types=pr&status=all", resourceType, id)
			localStart := time.Now()
			log.Printf(
				"op=FetchDiscussions stage=fetch_prs_start request_id=%s resource_type=%s id=%s url=%s",
				requestID, resourceType, id, url,
			)
			prs, err := GetDiscussionsFromURL(url)
			if err != nil {
				log.Printf(
					"op=FetchDiscussions stage=fetch_prs_error request_id=%s resource_type=%s id=%s url=%s error=%v elapsed=%s",
					requestID, resourceType, id, url, err, time.Since(localStart),
				)
				return
			}
			mu.Lock()
			discussions = append(discussions, prs...)
			mu.Unlock()
			log.Printf(
				"op=FetchDiscussions stage=fetch_prs_done request_id=%s resource_type=%s id=%s count=%d elapsed=%s",
				requestID, resourceType, id, len(prs), time.Since(localStart),
			)
		}()
	}

	if includeDiscussion {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url := fmt.Sprintf("https://huggingface.co/api/%s/%s/discussions?types=discussion&status=all", resourceType, id)
			localStart := time.Now()
			log.Printf(
				"op=FetchDiscussions stage=fetch_discussions_start request_id=%s resource_type=%s id=%s url=%s",
				requestID, resourceType, id, url,
			)
			discs, err := GetDiscussionsFromURL(url)
			if err != nil {
				log.Printf(
					"op=FetchDiscussions stage=fetch_discussions_error request_id=%s resource_type=%s id=%s url=%s error=%v elapsed=%s",
					requestID, resourceType, id, url, err, time.Since(localStart),
				)
				return
			}
			mu.Lock()
			discussions = append(discussions, discs...)
			mu.Unlock()
			log.Printf(
				"op=FetchDiscussions stage=fetch_discussions_done request_id=%s resource_type=%s id=%s count=%d elapsed=%s",
				requestID, resourceType, id, len(discs), time.Since(localStart),
			)
		}()
	}

	wg.Wait()

	log.Printf(
		"op=FetchDiscussions stage=success request_id=%s resource_type=%s id=%s total_count=%d total_elapsed=%s",
		requestID, resourceType, id, len(discussions), time.Since(start),
	)

	return discussions, nil
}

func GetDiscussionsFromURL(url string) ([]models.DISCUSSION, error) {
	start := time.Now()
	requestID := uuid.New().String()

	log.Printf(
		"op=GetDiscussionsFromURL stage=start request_id=%s url=%s",
		requestID, url,
	)

	resp, err := httpClient.Get(url)
	if err != nil {
		log.Printf(
			"op=GetDiscussionsFromURL stage=http_get_error request_id=%s url=%s error=%v elapsed=%s",
			requestID, url, err, time.Since(start),
		)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf(
		"op=GetDiscussionsFromURL stage=got_response request_id=%s url=%s status=%d content_length=%q",
		requestID, url, resp.StatusCode, resp.Header.Get("Content-Length"),
	)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf(
			"op=GetDiscussionsFromURL stage=read_error request_id=%s url=%s error=%v elapsed=%s",
			requestID, url, err, time.Since(start),
		)
		return nil, err
	}

	log.Printf(
		"op=GetDiscussionsFromURL stage=read_ok request_id=%s url=%s bytes=%d",
		requestID, url, len(body),
	)

	var rawDiscussions []map[string]interface{}
	if err := json.Unmarshal(body, &rawDiscussions); err != nil {
		log.Printf(
			"op=GetDiscussionsFromURL stage=json_unmarshal_error request_id=%s url=%s error=%v elapsed=%s",
			requestID, url, err, time.Since(start),
		)
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

	log.Printf(
		"op=GetDiscussionsFromURL stage=success request_id=%s url=%s count=%d total_elapsed=%s",
		requestID, url, len(discussions), time.Since(start),
	)

	return discussions, nil
}

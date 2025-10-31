package util

import (
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/MishraShardendu22/Scanner/models"

	"github.com/google/uuid"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

var httpClient = SharedHTTPClient()

// fetches only readable file content, (sirf padhne layak siblings)
func FetchFileContent(resourceID, filename string) models.SIBLING {
	start := time.Now()
	requestID := uuid.New().String()

	// ye ek sample siling hai jisme fetched content of a file content save karte hai ham loog
	sibling := models.SIBLING{
		RFilename:   filename,
		FileContent: "",
	}

	fileURL := fmt.Sprintf("https://huggingface.co/%s/resolve/main/%s", resourceID, filename)

	log.Printf(
		"op=FetchFileContent stage=start request_id=%s resource_id=%s filename=%s url=%s",
		requestID, resourceID, filename, fileURL,
	)

	resp, err := httpClient.Get(fileURL)
	if err != nil {
		log.Printf(
			"op=FetchFileContent stage=http_get_error request_id=%s resource_id=%s filename=%s url=%s error=%v elapsed=%s",
			requestID, resourceID, filename, fileURL, err, time.Since(start),
		)
		return sibling
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf(
			"op=FetchFileContent stage=not_ok request_id=%s resource_id=%s filename=%s status=%d elapsed=%s",
			requestID, resourceID, filename, resp.StatusCode, time.Since(start),
		)
		return sibling
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf(
			"op=FetchFileContent stage=read_error request_id=%s resource_id=%s filename=%s error=%v elapsed=%s",
			requestID, resourceID, filename, err, time.Since(start),
		)
		return sibling
	}

	log.Printf(
		"op=FetchFileContent stage=success request_id=%s resource_id=%s filename=%s bytes=%d total_elapsed=%s",
		requestID, resourceID, filename, len(body), time.Since(start),
	)

	sibling.FileContent = string(body)
	return sibling
}

// the siblings have the files names and we need to fetch their content,
// we use above helper function to fetch content of readable files concurrently
func FetchFilesFromSiblings(resourceID string, siblings []interface{}) []models.SIBLING {
	start := time.Now()
	requestID := uuid.New().String()

	var result []models.SIBLING
	var wg sync.WaitGroup
	var mu sync.Mutex
	semaphore := make(chan struct{}, 20)

	log.Printf(
		"op=FetchFilesFromSiblings stage=start request_id=%s resource_id=%s total_candidates=%d",
		requestID, resourceID, len(siblings),
	)

	total := len(siblings)

	for idx, sib := range siblings {
		if sibMap, ok := sib.(map[string]interface{}); ok {
			if filename, ok := sibMap["rfilename"].(string); ok {
				ext := strings.ToLower(filepath.Ext(filename))

				// agar wo files readable nahi hai toh skip kar denge
				if !TextExtensions[ext] {
					log.Printf(
						"op=FetchFilesFromSiblings stage=skip_non_readable request_id=%s resource_id=%s filename=%s ext=%s",
						requestID, resourceID, filename, ext,
					)
					continue
				}
				wg.Add(1)
				go func(fname string, index int, totalN int) {
					defer wg.Done()

					// semaphore basically is amde using buffered channel
					// adds  some thing to a semaphore so it's size becomes one less (limits concurrent goroutines)
					semaphore <- struct{}{}

					// removes from semaphore to free up space (makes size one more)
					defer func() { <-semaphore }()

					// defer is a treated like a stack

					localStart := time.Now()
					log.Printf(
						"op=FetchFilesFromSiblings stage=fetch_start request_id=%s resource_id=%s index=%d total=%d filename=%s",
						requestID, resourceID, index+1, totalN, fname,
					)

					sibling := FetchFileContent(resourceID, fname)

					// sibling jo extract hua hai usko add karenge atomic tareh se
					mu.Lock()
					result = append(result, sibling)
					mu.Unlock()
					log.Printf(
						"op=FetchFilesFromSiblings stage=fetch_done request_id=%s resource_id=%s index=%d total=%d filename=%s elapsed=%s",
						requestID, resourceID, index+1, totalN, fname, time.Since(localStart),
					)
				}(filename, idx, total)
			}
		}
	}

	wg.Wait()

	log.Printf(
		"op=FetchFilesFromSiblings stage=success request_id=%s resource_id=%s fetched=%d total_candidates=%d total_elapsed=%s",
		requestID, resourceID, len(result), len(siblings), time.Since(start),
	)

	return result
}

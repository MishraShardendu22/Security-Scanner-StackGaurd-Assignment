package util

import (
	"github.com/MishraShardendu22/Scanner/models"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
)

var SecretConfig = []models.SecretPattern{

	{Name: "Box API Key", Regex: `[0-9a-z]{40}`},
	{Name: "Algolia API Key", Regex: `[a-z0-9]{32}`},
	{Name: "GitHub PAT", Regex: `ghp_[A-Za-z0-9]{36}`},
	{Name: "Zoom JWT Key", Regex: `z0m[A-Za-z0-9]{32,}`},
	{Name: "GitLab PAT", Regex: `glpat-[A-Za-z0-9]{20,}`},
	{Name: "AWS Access Key ID", Regex: `AKIA[0-9A-Z]{16}`},
	{Name: "Oracle JDBC", Regex: `jdbc:oracle:thin:@[^\s]+`},
	{Name: "Netlify API Key", Regex: `ntl_[A-Za-z0-9]{32,}`},
	{Name: "Vercel Token", Regex: `vercel_[A-Za-z0-9]{32,}`},
	{Name: "Shopify API Key", Regex: `shp_[A-Za-z0-9]{32,}`},
	{Name: "Dropbox API Key", Regex: `sl.A[A-Za-z0-9]{32,}`},
	{Name: "Google API Key", Regex: `AIza[0-9A-Za-z\-_]{35}`},
	{Name: "PayPal Client ID", Regex: `Abcd[A-Za-z0-9]{20,}`},
	{Name: "Slack App Token", Regex: `xapp-[A-Za-z0-9]{24,}`},
	{Name: "Twilio API Key / SID", Regex: `AC[0-9a-fA-F]{32}`},
	{Name: "Anthropic API Key", Regex: `api-[A-Za-z0-9]{32,}`},
	{Name: "Firebase API Key", Regex: `AAAA[A-Za-z0-9\-_]{7,}`},
	{Name: "Datadog API Key", Regex: `ddapikey[A-Za-z0-9]{32,}`},
	{Name: "GitHub Actions Token", Regex: `gho_[A-Za-z0-9]{36}`},
	{Name: "OpenAI / LLM API Key", Regex: `sk-[A-Za-z0-9]{32,}`},
	{Name: "Hugging Face API Key", Regex: `hf_[A-Za-z0-9]{20,}`},
	{Name: "Facebook Access Token", Regex: `EAAG[A-Za-z0-9]{30,}`},
	{Name: "Stripe Secret Key", Regex: `sk_live_[A-Za-z0-9]{24,}`},
	{Name: "Instagram Access Token", Regex: `IGQV[A-Za-z0-9]{30,}`},
	{Name: "CircleCI API Token", Regex: `circleci-[A-Za-z0-9]{32,}`},
	{Name: "Twitter Bearer Token", Regex: `AAAAAAAA[A-Za-z0-9]{30,}`},
	{Name: "Google Cloud Service Account Key ID", Regex: `[A-Z0-9]{32}`},
	{Name: "Azure Access Key", Regex: `AZURE_STORAGE_KEY_[A-Za-z0-9]{40,}`},
	{Name: "Kubernetes Bearer Token", Regex: `eyJhbGciOiJSUzI1Ni[A-Za-z0-9\-_]+`},
	{Name: "Slack Bot Token", Regex: `xoxb-[0-9]{10,}-[0-9]{10,}-[A-Za-z0-9]{24,}`},
	{Name: "MySQL URI", Regex: `mysql:\/\/(?:[^@\s]+@)?[^\s\/:]+(?::\d+)?\/[^\s]*`},
	{Name: "Redis URI", Regex: `rediss?:\/\/(?:[^@\s]+@)?[^\s\/:]+(?::\d+)?(?:\/\d+)?\b`},
	{Name: "JDBC generic", Regex: `jdbc:[a-z0-9]+:\/\/(?:[^@\s]+@)?[^\s\/:]+(?::\d+)?\/[^\s]*`},
	{Name: "AMQP / RabbitMQ URI", Regex: `amqps?:\/\/(?:[^@\s]+@)?[^\s\/:]+(?::\d+)?\/?[^\s]*`},
	{Name: "Elasticsearch Basic Auth", Regex: `https?:\/\/[^:\s]+:[^@\s]+@[^\/\s:]+(?::\d+)?\/?`},
	{Name: "MongoDB URI", Regex: `mongodb(?:\+srv)?:\/\/(?:[^@\s]+@)?[^\s\/:]+(?::\d+)?\/[^\s]*`},
	{Name: "PostgreSQL URI", Regex: `postgres(?:ql)?:\/\/(?:[^@\s]+@)?[^\s\/:]+(?::\d+)?\/[^\s]*`},
	{Name: "MSSQL URI", Regex: `(?:mssql|sqlserver):\/\/(?:[^@\s]+@)?[^\s\/:]+(?::\d+)?\/?[^\s]*`},
	{Name: "CockroachDB URI", Regex: `cockroach(?:db)?:\/\/(?:[^@\s]+@)?[^\s\/:]+(?::\d+)?\/[^\s]*`},
	{Name: "Generic DB URI with creds", Regex: `[a-z0-9+\-]+:\/\/[^:\s\/]+:[^@\s\/]+@[^\/\s:]+(?::\d+)?\/[^\s]*`},
}
var TextExtensions = map[string]bool{

	".sh": true, ".bash": true, ".zsh": true,
	".csv": true, ".tsv": true, ".ipynb": true,
	".docker-compose": true, ".docker-compose.yml": true,
	".jks": true, ".kdb": true, ".pub": true, ".asc": true,
	".html": true, ".htm": true, ".ejs": true, ".vue": true,
	".ini": true, ".cfg": true, ".conf": true, ".env": true,
	".env.example": true, ".lock": true, ".properties": true,
	".css": true, ".scss": true, ".sass": true, ".less": true,
	".json": true, ".yaml": true, ".yml": true, ".toml": true,
	".ini.local": true, ".yml.local": true, ".yaml.local": true,
	".dockerfile": true, ".gitignore": true, ".gitattributes": true,
	".java": true, ".cpp": true, ".c": true, ".h": true, ".cs": true,
	".py": true, ".js": true, ".jsx": true, ".ts": true, ".tsx": true,
	".go": true, ".rs": true, ".rb": true, ".php": true, ".mjs": true,
	".md": true, ".mdx": true, ".rst": true, ".txt": true, ".log": true,
	".secrets": true, ".key": true, ".pem": true, ".crt": true, ".p12": true,
	".env.local": true, ".env.dev": true, ".env.prod": true, ".env.test": true,
}

func ScanFile(file models.SIBLING, patterns []models.SecretPattern) []models.Finding {

	ext := strings.ToLower(filepath.Ext(file.RFilename))

	if !TextExtensions[ext] {
		return nil
	}
	var findings []models.Finding

	lines := strings.Split(file.FileContent, "\n")

	for i, line := range lines {
		for _, pattern := range patterns {
			re := regexp.MustCompile(pattern.Regex)
			matches := re.FindAllString(line, -1)
			for _, match := range matches {
				findings = append(findings, models.Finding{

					SecretType: pattern.Name,
					Pattern:    pattern.Regex,
					Secret:     match,
					SourceType: "file",
					FileName:   file.RFilename,
					Line:       i + 1,
				})
			}
		}
	}

	return findings
}

func ScanDiscussion(disc models.DISCUSSION, patterns []models.SecretPattern) []models.Finding {

	var findings []models.Finding

	text := disc.Title + " " + disc.RepoName

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern.Regex)
		matches := re.FindAllString(text, -1)
		for _, match := range matches {
			findings = append(findings, models.Finding{

				SecretType:      pattern.Name,
				Pattern:         pattern.Regex,
				Secret:          match,
				SourceType:      "discussion",
				DiscussionNum:   disc.Num,
				DiscussionTitle: disc.Title,
				DiscussionRepo:  disc.RepoName,
			})
		}
	}

	return findings
}

func ScanAIRequest(req models.AI_REQUEST, patterns []models.SecretPattern) []models.Finding {

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
			findings := ScanFile(file, patterns)
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
			findings := ScanDiscussion(disc, patterns)
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

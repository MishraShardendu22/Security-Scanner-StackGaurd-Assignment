package util

import "github.com/MishraShardendu22/Scanner/models"

var SecretConfig = []models.SecretPattern{
	// API / Access Tokens
	{Name: "GitHub PAT", Regex: `ghp_[A-Za-z0-9]{36}`},
	{Name: "GitHub Actions Token", Regex: `gho_[A-Za-z0-9]{36}`},
	{Name: "GitLab PAT", Regex: `glpat-[A-Za-z0-9]{20,}`},
	{Name: "OpenAI / LLM API Key", Regex: `sk-[A-Za-z0-9]{32,}`},
	{Name: "Anthropic API Key", Regex: `api-[A-Za-z0-9]{32,}`},
	{Name: "Hugging Face API Key", Regex: `hf_[A-Za-z0-9]{20,}`},
	{Name: "Slack Bot Token", Regex: `xoxb-[0-9]{10,}-[0-9]{10,}-[A-Za-z0-9]{24,}`},
	{Name: "Slack App Token", Regex: `xapp-[A-Za-z0-9]{24,}`},
	{Name: "Stripe Secret Key", Regex: `sk_live_[A-Za-z0-9]{24,}`},
	{Name: "Azure Access Key", Regex: `AZURE_STORAGE_KEY_[A-Za-z0-9]{40,}`},
	{Name: "Twilio API Key / SID", Regex: `AC[0-9a-fA-F]{32}`},
	{Name: "Twitter Bearer Token", Regex: `AAAAAAAA[A-Za-z0-9]{30,}`},
	{Name: "Dropbox API Key", Regex: `sl.A[A-Za-z0-9]{32,}`},
	{Name: "Zoom JWT Key", Regex: `z0m[A-Za-z0-9]{32,}`},
	{Name: "PayPal Client ID", Regex: `Abcd[A-Za-z0-9]{20,}`},
	{Name: "Algolia API Key", Regex: `[a-z0-9]{32}`},
	{Name: "Box API Key", Regex: `[0-9a-z]{40}`},
	{Name: "Datadog API Key", Regex: `ddapikey[A-Za-z0-9]{32,}`},
	{Name: "CircleCI API Token", Regex: `circleci-[A-Za-z0-9]{32,}`},
	{Name: "Netlify API Key", Regex: `ntl_[A-Za-z0-9]{32,}`},
	{Name: "Vercel Token", Regex: `vercel_[A-Za-z0-9]{32,}`},
	{Name: "Shopify API Key", Regex: `shp_[A-Za-z0-9]{32,}`},
	{Name: "Instagram Access Token", Regex: `IGQV[A-Za-z0-9]{30,}`},
	{Name: "Facebook Access Token", Regex: `EAAG[A-Za-z0-9]{30,}`},

	// Cloud / Service Keys
	{Name: "AWS Access Key ID", Regex: `AKIA[0-9A-Z]{16}`},
	{Name: "Google API Key", Regex: `AIza[0-9A-Za-z\-_]{35}`},
	{Name: "Firebase API Key", Regex: `AAAA[A-Za-z0-9\-_]{7,}`},
	{Name: "Google Cloud Service Account Key ID", Regex: `[A-Z0-9]{32}`},

	// Database / Connection Strings
	{Name: "MongoDB URI", Regex: `mongodb(?:\+srv)?:\/\/(?:[^@\s]+@)?[^\s\/:]+(?::\d+)?\/[^\s]*`},
	{Name: "PostgreSQL URI", Regex: `postgres(?:ql)?:\/\/(?:[^@\s]+@)?[^\s\/:]+(?::\d+)?\/[^\s]*`},
	{Name: "MySQL URI", Regex: `mysql:\/\/(?:[^@\s]+@)?[^\s\/:]+(?::\d+)?\/[^\s]*`},
	{Name: "Redis URI", Regex: `rediss?:\/\/(?:[^@\s]+@)?[^\s\/:]+(?::\d+)?(?:\/\d+)?\b`},
	{Name: "MSSQL URI", Regex: `(?:mssql|sqlserver):\/\/(?:[^@\s]+@)?[^\s\/:]+(?::\d+)?\/?[^\s]*`},
	{Name: "Oracle JDBC", Regex: `jdbc:oracle:thin:@[^\s]+`},
	{Name: "JDBC generic", Regex: `jdbc:[a-z0-9]+:\/\/(?:[^@\s]+@)?[^\s\/:]+(?::\d+)?\/[^\s]*`},
	{Name: "AMQP / RabbitMQ URI", Regex: `amqps?:\/\/(?:[^@\s]+@)?[^\s\/:]+(?::\d+)?\/?[^\s]*`},
	{Name: "Elasticsearch Basic Auth", Regex: `https?:\/\/[^:\s]+:[^@\s]+@[^\/\s:]+(?::\d+)?\/?`},
	{Name: "CockroachDB URI", Regex: `cockroach(?:db)?:\/\/(?:[^@\s]+@)?[^\s\/:]+(?::\d+)?\/[^\s]*`},
	{Name: "Generic DB URI with creds", Regex: `[a-z0-9+\-]+:\/\/[^:\s\/]+:[^@\s\/]+@[^\/\s:]+(?::\d+)?\/[^\s]*`},
	{Name: "Kubernetes Bearer Token", Regex: `eyJhbGciOiJSUzI1Ni[A-Za-z0-9\-_]+`},
}

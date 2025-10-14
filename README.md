# StackGuard Assignment: Hugging Face Secret Scanner

A comprehensive API-based system that scans Hugging Face-hosted assets (models, datasets, and spaces) for potential secrets and tokens using regex-based pattern matching.

## Overview

This system extends StackGuard's secret scanning capabilities to AI model ecosystems by:
- Fetching files from Hugging Face resources (Model / Dataset / Space)
- Scanning content using 52+ regex-based secret detectors
- Returning structured findings via API responses
- Persisting and contextualizing results for visualization

## Features

✅ **Unified Scan Endpoint** - Single endpoint to scan models, datasets, spaces, or entire organizations
✅ **52+ Secret Patterns** - Detects AWS keys, GitHub tokens, API keys, database URIs, and more
✅ **Organization-Level Scanning** - Scan all resources under a given org/user
✅ **Results Storage** - Persist scan results in MongoDB with full metadata
✅ **Dashboard API** - View aggregated statistics grouped by resource type and severity
✅ **Concurrent Scanning** - Fast parallel scanning using goroutines
✅ **RESTful API** - Clean, well-documented API endpoints

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Fiber (Fast HTTP framework)
- **Database**: MongoDB (via MGM - MongoDB Go Manager)
- **External API**: Hugging Face Hub API

## Prerequisites

- Go 1.21 or higher
- MongoDB instance (local or cloud)
- Internet connection for Hugging Face API access

## Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/MishraShardendu22/Security-Scanner-StackGaurd-Assignment.git
cd Security-Scanner-StackGaurd-Assignment
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Configure Environment

Create a `.env` file in the root directory:

```env
PORT=8080
DB_NAME=security_scanner
MONGODB_URI=mongodb://localhost:27017
LOG_LEVEL=info
ENVIRONMENT=development
CORS_ALLOW_ORIGINS=*
```

### 4. Build the Project

```bash
go build -o scanner
```

### 5. Run the Server

```bash
./scanner
# or
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### 1. Main Scan Endpoint (Assignment Required)

**POST /api/scan**

Unified endpoint to scan any Hugging Face resource.

**Request Body:**
```json
{
  "model_id": "microsoft/phi-3",
  "dataset_id": null,
  "space_id": null,
  "org": "microsoft",
  "user": null,
  "include_discussions": true,
  "include_prs": false
}
```

**Response:**
```json
{
  "status": 200,
  "message": "Scan completed successfully",
  "data": {
    "scan_id": "SG-2025-1015-a1b2c3d4",
    "scanned_resources": [
      {
        "type": "model",
        "id": "microsoft/phi-3",
        "findings": [
          {
            "secret_type": "AWS Access Key",
            "pattern": "AKIA********EXAMPLE",
            "file": "config.json",
            "line": 24
          }
        ]
      }
    ],
    "timestamp": "2025-10-15T12:30:00Z",
    "total_findings": 5,
    "storage_id": "64a7f8b9c123456789abcdef"
  }
}
```

### 2. Store Scan Results (Assignment Required)

**POST /api/store**

Store scan results with metadata.

**Request Body:**
```json
{
  "scan_id": "SG-2025-1015-001",
  "scanned_resources": [
    {
      "type": "model",
      "id": "microsoft/phi-3",
      "findings": [...]
    }
  ],
  "timestamp": "2025-10-15T12:30:00Z"
}
```

**Response:**
```json
{
  "status": 200,
  "message": "Scan results stored successfully",
  "data": {
    "status": "stored",
    "scan_id": "SG-2025-1015-001",
    "storage_id": "64a7f8b9c123456789abcdef"
  }
}
```

### 3. Fetch Results (Assignment Required)

**GET /api/results/{scan_id}**

Retrieve stored scan details and contextual metadata.

**Response:**
```json
{
  "status": 200,
  "message": "Scan result retrieved successfully",
  "data": {
    "scan_id": "64a7f8b9c123456789abcdef",
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "scanned_resources": [...],
    "total_findings": 15,
    "findings_by_type": {
      "AWS Access Key ID": 5,
      "GitHub PAT": 3
    },
    "created_at": "2025-10-15T10:30:00Z"
  }
}
```

### 4. Dashboard (Showcase Feature)

**GET /api/dashboard**

Lightweight dashboard showing all stored results grouped by resource type and severity.

**Response:**
```json
{
  "status": 200,
  "message": "Dashboard data retrieved successfully",
  "data": {
    "total_scans": 25,
    "total_findings": 340,
    "severity_breakdown": {
      "high": 108,
      "medium": 150,
      "low": 82
    },
    "by_resource_type": {
      "file": 280,
      "discussion": 60
    },
    "recent_scans": [...]
  }
}
```

## Secret Detection Patterns

The scanner detects 52+ types of secrets including:

### API Keys & Tokens
- GitHub PAT (`ghp_*`, `gho_*`)
- GitLab PAT (`glpat-*`)
- OpenAI/LLM API Keys (`sk-*`)
- Hugging Face API Keys (`hf_*`)
- AWS Access Key ID (`AKIA*`)
- Google API Key (`AIza*`)
- Stripe Secret Key (`sk_live_*`)
- Slack Bot/App Tokens
- And 30+ more...

### Database Connection Strings
- PostgreSQL, MySQL, MongoDB URIs
- Redis, MSSQL, CockroachDB
- JDBC, AMQP/RabbitMQ
- Generic DB URIs with credentials

### Cloud & Service Keys
- AWS Access Keys
- Azure Storage Keys
- Google Cloud Service Account Keys
- Kubernetes Bearer Tokens

## Usage Examples

### Example 1: Scan a Single Model

```bash
curl -X POST http://localhost:8080/api/scan \
  -H "Content-Type: application/json" \
  -d '{
    "model_id": "meta-llama/Meta-Llama-3-70B-Instruct",
    "include_discussions": true,
    "include_prs": false
  }'
```

### Example 2: Scan an Entire Organization

```bash
curl -X POST http://localhost:8080/api/scan \
  -H "Content-Type: application/json" \
  -d '{
    "org": "microsoft",
    "include_discussions": true,
    "include_prs": true
  }'
```

### Example 3: Get Dashboard Statistics

```bash
curl http://localhost:8080/api/dashboard
```

### Example 4: Retrieve Scan Results

```bash
curl http://localhost:8080/api/results/64a7f8b9c123456789abcdef
```

## Project Structure

```
.
├── main.go                 # Application entry point
├── controller/
│   ├── fetch.go           # Individual resource fetching
│   ├── org-specific.go    # Organization-level fetching
│   ├── unified.go         # Main scan endpoint (POST /api/scan)
│   ├── scan.go            # Legacy scan endpoints
│   └── results.go         # Results & dashboard endpoints
├── route/
│   ├── fetch.go           # Fetch routes
│   ├── org-specific.go    # Organization routes
│   ├── scan.go            # Scan routes
│   └── results.go         # Results routes
├── models/
│   ├── model.go           # Data models
│   ├── secret.go          # Secret pattern definitions
│   └── config.go          # Configuration model
├── util/
│   ├── secret.go          # Scanner logic
│   ├── apiResponse.go     # Response helper
│   └── getEnv.go          # Environment helper
└── database/
    └── database.go        # MongoDB connection
```

## Approach & Design Decisions

### 1. Architecture
- **MVC Pattern**: Clean separation of concerns with controllers, models, and routes
- **Modular Design**: Each component has a single responsibility
- **RESTful API**: Standard HTTP methods and status codes

### 2. Scanning Strategy
- **Regex-Based Detection**: Fast, efficient pattern matching
- **Concurrent Scanning**: Goroutines for parallel file processing
- **File Type Filtering**: Only scans text-based files (70+ extensions)
- **Line-by-Line Analysis**: Precise location tracking

### 3. Data Storage
- **MongoDB**: Flexible schema for varied scan results
- **MGM Library**: Automatic timestamps and ID management
- **Structured Results**: Consistent JSON format for easy integration

### 4. Performance Optimizations
- **Concurrent Scanning**: WaitGroups and channels
- **Compiled Regex**: Pre-compiled patterns for speed
- **Selective Fetching**: Fetch only necessary files
- **Organization Limits**: Limit scans to first 5 models for org scans

### 5. Error Handling
- **Graceful Degradation**: Continue scanning even if individual files fail
- **Comprehensive Logging**: Structured logging with slog
- **HTTP Status Codes**: Proper error responses

## Error Handling

The system includes:
- Request validation
- MongoDB connection error handling
- Hugging Face API error handling
- Rate limiting considerations
- Structured error responses

## Rate Limiting & Best Practices

- Uses Hugging Face public API (no auth required for public resources)
- Implements sensible limits for organization scans (5 models max)
- Includes request timeouts
- Supports graceful shutdown

## Testing

Run tests:
```bash
go test ./...
```

Test the API:
```bash
# Health check
curl http://localhost:8080/api/test

# Scan a model
curl -X POST http://localhost:8080/api/scan \
  -H "Content-Type: application/json" \
  -d '{"model_id": "microsoft/phi-3-mini-4k-instruct"}'
```

## Future Enhancements

- [ ] Webhook notifications for high-risk findings
- [ ] Export results to CSV/JSON
- [ ] Custom regex pattern support
- [ ] Integration with secret management tools
- [ ] Real-time scanning via WebSocket
- [ ] API authentication/authorization
- [ ] Rate limiting middleware

## License

MIT License

## Author

Shardendu Mishra

## Acknowledgments

- Hugging Face for their excellent Hub API
- StackGuard for the assignment opportunity
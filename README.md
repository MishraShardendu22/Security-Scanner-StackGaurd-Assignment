# StackGuard Security Scanner

<div align="center">

**A comprehensive API-based security scanner for detecting secrets and sensitive information in Hugging Face repositories**

[Features](#features) • [Quick Start](#installation--setup) • [API Documentation](#api-endpoints) • [Architecture](#architecture--design)

</div>

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Installation & Setup](#installation--setup)
- [API Endpoints](#api-endpoints)
- [Secret Detection Patterns](#secret-detection-patterns)
- [Usage Examples](#usage-examples)
- [Project Structure](#project-structure)
- [Architecture & Design](#architecture--design)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

**StackGuard Security Scanner** is a production-ready security tool that extends secret scanning capabilities to AI model ecosystems hosted on Hugging Face. It provides comprehensive scanning of models, datasets, and spaces to detect potential security vulnerabilities such as exposed API keys, tokens, and credentials.

### Why This Matters

With the rapid growth of AI and machine learning, developers often inadvertently commit sensitive information to model repositories. This scanner helps identify these security risks before they become vulnerabilities.

### Key Capabilities

- **Comprehensive Scanning**: Analyze models, datasets, spaces, and entire organizations
- **52+ Secret Patterns**: Advanced regex-based detection for various secret types
- **High Performance**: Concurrent scanning using Go's goroutines
- **Persistent Storage**: MongoDB-backed result storage with full metadata
- **Dashboard Analytics**: Aggregated statistics and insights
- **RESTful API**: Well-documented, easy-to-integrate endpoints
- **Production Ready**: Error handling, logging, graceful shutdown

---

## Features

| Feature | Description |
|---------|-------------|
| **Unified Scan Endpoint** | Single endpoint to scan any Hugging Face resource type |
| **Multi-Pattern Detection** | 52+ regex patterns for detecting secrets, keys, and tokens |
| **Organization Scanning** | Scan all resources under a specific organization or user |
| **Concurrent Processing** | Fast parallel file scanning using goroutines and channels |
| **Result Persistence** | Store scan results with full metadata in MongoDB |
| **Dashboard Analytics** | View aggregated statistics by resource type and severity |
| **File Type Filtering** | Smart filtering of 70+ text-based file extensions |
| **Line-Level Precision** | Exact line number and context for each finding |
| **Graceful Degradation** | Continue scanning even if individual files fail |
| **Structured Logging** | Comprehensive logging using `slog` for debugging |

---

## Tech Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Language** | Go 1.24.8 | High-performance backend |
| **Web Framework** | Fiber v2.52.9 | Fast HTTP server and routing |
| **Database** | MongoDB | Flexible document storage |
| **ODM** | MGM v3.5.0 | MongoDB Go Manager for models |
| **External API** | Hugging Face Hub API | Repository data fetching |
| **Configuration** | godotenv v1.5.1 | Environment variable management |
| **UUID Generation** | google/uuid v1.6.0 | Unique identifier generation |

---

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go**: Version 1.21 or higher ([Download](https://go.dev/dl/))
- **MongoDB**: Local instance or cloud connection ([MongoDB Atlas](https://www.mongodb.com/cloud/atlas))
- **Git**: For cloning the repository
- **Internet Connection**: Required for Hugging Face API access

### System Requirements

- **OS**: Linux, macOS, or Windows
- **RAM**: Minimum 2GB (4GB recommended for large scans)
- **Disk Space**: 500MB for application and dependencies

---

## Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/MishraShardendu22/Security-Scanner-StackGaurd-Assignment.git
cd Security-Scanner-StackGaurd-Assignment
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

### 3. Configure Environment Variables

Create a `.env` file in the root directory with the following configuration:

```env
# Server Configuration
PORT=8080
ENVIRONMENT=development

# Database Configuration
DB_NAME=security_scanner
MONGODB_URI=mongodb://localhost:27017

# Logging
LOG_LEVEL=info

# CORS Configuration
CORS_ALLOW_ORIGINS=*
```

**Configuration Options Explained:**

- `PORT`: Server port (default: 8080)
- `DB_NAME`: MongoDB database name
- `MONGODB_URI`: MongoDB connection string
- `LOG_LEVEL`: Logging verbosity (debug, info, warn, error)
- `ENVIRONMENT`: Runtime environment (development, production)

### 4. Start MongoDB

**Option A: Local MongoDB**

```bash
# Start MongoDB service
sudo systemctl start mongod

# Or using Docker
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

**Option B: MongoDB Atlas (Cloud)**

1. Create a free cluster at [MongoDB Atlas](https://www.mongodb.com/cloud/atlas)
2. Get your connection string
3. Update `MONGODB_URI` in `.env` file

### 5. Build the Project

```bash
# Build executable
go build -o scanner

# Or build with optimizations
go build -ldflags="-s -w" -o scanner
```

### 6. Run the Server

```bash
# Using compiled binary
./scanner

# Or run directly
go run main.go
```

The server will start on `http://localhost:8080`

**Expected Output:**

```
Stack Guard Assignment
INFO Starting Security Scanner environment=development port=8080 log_level=info
INFO Database connected successfully database=security_scanner
INFO Server listening on http://localhost:8080
```

### 7. Verify Installation

Test the server health endpoint:

```bash
curl http://localhost:8080/api/test
```

Expected response:

```json
{
  "status": 200,
  "message": "API is working!",
  "data": null
}
```

---

## API Endpoints

### 1. Main Scan Endpoint (Unified)

#### `POST /api/scan`

Unified endpoint to scan any Hugging Face resource type.

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

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `model_id` | string | No* | Hugging Face model ID (e.g., "microsoft/phi-3") |
| `dataset_id` | string | No* | Hugging Face dataset ID |
| `space_id` | string | No* | Hugging Face space ID |
| `org` | string | No* | Organization name to scan all resources |
| `user` | string | No* | User name to scan all resources |
| `include_discussions` | boolean | No | Include discussion scanning (default: false) |
| `include_prs` | boolean | No | Include PR/commit scanning (default: false) |

*At least one of: `model_id`, `dataset_id`, `space_id`, `org`, or `user` must be provided.

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
            "line": 24,
            "context": "AWS_ACCESS_KEY=AKIA..."
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

---

### 2. Store Scan Results

#### `POST /api/store`

Store scan results with metadata for later retrieval.

**Request Body:**

```json
{
  "scan_id": "SG-2025-1015-001",
  "scanned_resources": [
    {
      "type": "model",
      "id": "microsoft/phi-3",
      "findings": []
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

---

### 3. Fetch Scan Results

#### `GET /api/results/{scan_id}`

Retrieve stored scan details and contextual metadata.

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `scan_id` | string | MongoDB ObjectID of the scan |

**Response:**

```json
{
  "status": 200,
  "message": "Scan result retrieved successfully",
  "data": {
    "scan_id": "64a7f8b9c123456789abcdef",
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "scanned_resources": [],
    "total_findings": 15,
    "findings_by_type": {
      "AWS Access Key ID": 5,
      "GitHub PAT": 3,
      "OpenAI API Key": 7
    },
    "created_at": "2025-10-15T10:30:00Z",
    "updated_at": "2025-10-15T10:30:00Z"
  }
}
```

---

### 4. Dashboard Analytics

#### `GET /api/dashboard`

View aggregated statistics and insights from all scans.

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
    "recent_scans": [
      {
        "scan_id": "64a7f8b9c123456789abcdef",
        "timestamp": "2025-10-15T12:30:00Z",
        "findings_count": 15
      }
    ]
  }
}
```

---

### 5. Health Check

#### `GET /api/test`

Verify API availability.

**Response:**

```json
{
  "status": 200,
  "message": "API is working!",
  "data": null
}
```

---

## Secret Detection Patterns

The scanner uses advanced regex patterns to detect **52+ types of secrets** including:

### API Keys & Tokens

- GitHub PAT (`ghp_*`, `gho_*`, `ghs_*`, `ghr_*`)
- GitLab PAT (`glpat-*`)
- OpenAI/LLM API Keys (`sk-*`, `sk-proj-*`)
- Hugging Face API Keys (`hf_*`)
- AWS Access Key ID (`AKIA*`, `ASIA*`)
- Google API Key (`AIza*`)
- Stripe Secret Key (`sk_live_*`, `rk_live_*`)
- Slack Bot/App Tokens (`xoxb-*`, `xoxa-*`, `xoxp-*`)
- Anthropic API Keys
- Twilio API Keys
- SendGrid API Keys
- Discord Bot Tokens
- Heroku API Keys
- Mailgun API Keys
- And 30+ more...

### Database Connection Strings

- PostgreSQL (`postgresql://`, `postgres://`)
- MySQL (`mysql://`)
- MongoDB (`mongodb://`, `mongodb+srv://`)
- Redis (`redis://`)
- Microsoft SQL Server (`mssql://`)
- CockroachDB (`cockroachdb://`)
- JDBC Connection Strings
- AMQP/RabbitMQ (`amqp://`)
- Generic DB URIs with embedded credentials

### Cloud & Service Keys

- AWS Access Keys & Secret Keys
- Azure Storage Account Keys
- Google Cloud Service Account Keys
- Kubernetes Bearer Tokens
- Docker Registry Credentials
- SSH Private Keys (RSA, DSA, EC, OPENSSH)

### Authentication Tokens

- JWT Tokens
- Basic Auth credentials in URLs
- Bearer tokens
- OAuth tokens
- Session tokens

### Sensitive Patterns

- Generic passwords in configuration
- Email addresses
- Credit card numbers (basic detection)
- Social Security Numbers (SSN)
- Private keys and certificates

---

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

### Example 2: Scan a Dataset

```bash
curl -X POST http://localhost:8080/api/scan \
  -H "Content-Type: application/json" \
  -d '{
    "dataset_id": "squad",
    "include_discussions": false
  }'
```

### Example 3: Scan an Entire Organization

```bash
curl -X POST http://localhost:8080/api/scan \
  -H "Content-Type: application/json" \
  -d '{
    "org": "microsoft",
    "include_discussions": true,
    "include_prs": true
  }'
```

### Example 4: Scan a Space

```bash
curl -X POST http://localhost:8080/api/scan \
  -H "Content-Type: application/json" \
  -d '{
    "space_id": "stabilityai/stable-diffusion",
    "include_discussions": false
  }'
```

### Example 5: Get Dashboard Statistics

```bash
curl http://localhost:8080/api/dashboard
```

### Example 6: Retrieve Scan Results

```bash
curl http://localhost:8080/api/results/64a7f8b9c123456789abcdef
```

### Example 7: Using with jq for Pretty Output

```bash
curl -s http://localhost:8080/api/dashboard | jq '.'
```

---

## Project Structure

```plaintext
Security-Scanner-StackGaurd-Assignment/
├── main.go                              # Application entry point & server setup
├── go.mod                               # Go module dependencies
├── go.sum                               # Dependency checksums
├── .env                                 # Environment configuration
├── README.md                            # Project documentation
├── TODO.md                              # Task tracking
│
├── controller/                          # Business logic layer
│   ├── unified.go                       # Main unified scan endpoint handler
│   ├── fetch.go                         # Individual resource fetching logic
│   ├── org-specific.go                  # Organization-level scanning
│   ├── scan.go                          # Legacy scan endpoints
│   └── results.go                       # Results retrieval & dashboard
│
├── route/                               # API route definitions
│   ├── scan.go                          # Scan endpoint routes
│   ├── fetch.go                         # Fetch endpoint routes
│   ├── org-specific.go                  # Organization endpoint routes
│   └── results.go                       # Results endpoint routes
│
├── models/                              # Data models & schemas
│   ├── model.go                         # Core data structures
│   ├── secret.go                        # Secret pattern definitions (52+)
│   └── config.go                        # Configuration model
│
├── util/                                # Utility functions
│   ├── secret.go                        # Secret scanning core logic
│   ├── apiResponse.go                   # Standardized API response helper
│   └── getEnv.go                        # Environment variable utilities
│
└── database/                            # Database layer
    └── database.go                      # MongoDB connection & setup
```

---

## Architecture & Design

### Design Principles

#### 1. **MVC Architecture**

- **Model**: Data structures and business entities (`models/`)
- **View**: JSON API responses (`util/apiResponse.go`)
- **Controller**: Request handling and business logic (`controller/`)
- **Router**: Route definitions and middleware (`route/`)

#### 2. **Scanning Strategy**

- **Regex-Based Detection**: 52+ pre-compiled regex patterns for optimal performance
- **Concurrent Processing**: Goroutines with WaitGroups for parallel file scanning
- **File Type Filtering**: Smart filtering for 70+ text-based file extensions
- **Line-by-Line Analysis**: Precise location tracking with context
- **Fail-Safe Design**: Continue scanning even if individual files fail

#### 3. **Data Flow**

```plaintext
Request → Router → Controller → Scanner → Hugging Face API
                                    ↓
                              Pattern Matching
                                    ↓
                              Store in MongoDB
                                    ↓
                              Return Results
```

#### 4. **Performance Optimizations**

- **Pre-compiled Regex**: All patterns compiled at startup
- **Concurrent Scanning**: Parallel processing of files
- **Selective Fetching**: Only fetch scannable file types
- **Connection Pooling**: MongoDB connection reuse
- **Organization Limits**: Default limit of 5 models per org scan

#### 5. **Error Handling**

- **Graceful Degradation**: Partial results on non-critical failures
- **Structured Logging**: `slog` for production-grade logging
- **Context Propagation**: Request context throughout the stack
- **HTTP Status Codes**: Proper REST status codes

---

## Security & Best Practices

### Rate Limiting

- Respects Hugging Face API rate limits
- Configurable request timeouts
- Organization scan limits to prevent abuse

### Error Handling

- Input validation and sanitization
- MongoDB connection error recovery
- External API failure handling
- Comprehensive error messages

### Logging

- Structured logging with severity levels
- Request/response tracking
- Performance metrics
- Error stack traces

### Graceful Shutdown

- Signal handling (SIGINT, SIGTERM)
- Connection cleanup
- In-flight request completion

---

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Manual API Testing

#### Health Check

```bash
curl http://localhost:8080/api/test
```

#### Scan a Model

```bash
curl -X POST http://localhost:8080/api/scan \
  -H "Content-Type: application/json" \
  -d '{
    "model_id": "microsoft/phi-3-mini-4k-instruct",
    "include_discussions": false
  }'
```

#### View Dashboard

```bash
curl http://localhost:8080/api/dashboard | jq '.'
```

### Using Postman

Import the provided Postman collection:

- `StackGuard-Security-Scanner.postman_collection.json`

---

## Deployment

### Docker Deployment (Recommended)

#### 1. Create Dockerfile

```dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o scanner

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/scanner .
COPY .env .

EXPOSE 8080
CMD ["./scanner"]
```

#### 2. Build and Run

```bash
# Build image
docker build -t security-scanner .

# Run container
docker run -d -p 8080:8080 \
  -e MONGODB_URI=mongodb://host.docker.internal:27017 \
  --name scanner security-scanner
```

### Cloud Deployment

#### Heroku

```bash
# Login to Heroku
heroku login

# Create app
heroku create your-scanner-app

# Add MongoDB addon
heroku addons:create mongolab

# Deploy
git push heroku main
```

#### AWS EC2

1. Launch EC2 instance (Ubuntu 22.04)
2. Install Go and MongoDB
3. Clone repository
4. Configure environment variables
5. Set up systemd service
6. Configure security groups (port 8080)

#### DigitalOcean

Similar to AWS EC2, using droplet with MongoDB.

---

## Environment Variables Reference

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `8080` | No |
| `DB_NAME` | MongoDB database name | `security_scanner` | Yes |
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017` | Yes |
| `LOG_LEVEL` | Logging level (debug/info/warn/error) | `info` | No |
| `ENVIRONMENT` | Runtime environment | `development` | No |
| `CORS_ALLOW_ORIGINS` | CORS allowed origins | `*` | No |

---

## Troubleshooting

### Common Issues

#### MongoDB Connection Failed

```bash
# Check MongoDB is running
sudo systemctl status mongod

# Test connection
mongosh --host localhost --port 27017
```

#### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

#### Module Not Found

```bash
# Clean and reinstall dependencies
go clean -modcache
go mod download
go mod tidy
```

---

## Performance Metrics

- **Average Scan Time**: 2-5 seconds per model
- **Concurrent Files**: Up to 50 files simultaneously
- **Memory Usage**: ~100-200MB average
- **Throughput**: ~10-15 scans per minute

---

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow Go conventions and best practices
- Use `gofmt` for formatting
- Add comments for exported functions
- Write tests for new features

---

## Future Enhancements

- [ ] Webhook notifications for critical findings
- [ ] Export results to CSV/JSON/PDF
- [ ] Custom regex pattern support via UI
- [ ] Integration with Vault/AWS Secrets Manager
- [ ] Real-time scanning via WebSocket
- [ ] API authentication with JWT
- [ ] Rate limiting middleware
- [ ] Scheduled automated scans
- [ ] Slack/Discord integration
- [ ] Multi-language support
- [ ] False positive reduction using ML
- [ ] Browser extension for quick scans

---

## License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

## Author

**Shardendu Mishra**

- GitHub: [@MishraShardendu22](https://github.com/MishraShardendu22)
- LinkedIn: [Shardendu Mishra](https://linkedin.com/in/shardendu-mishra)

---

## Acknowledgments

- **Hugging Face** for their excellent Hub API and ecosystem
- **StackGuard** for the assignment opportunity and inspiration
- **Go Fiber** team for the amazing web framework
- **MongoDB** for flexible document storage
- Open source community for various libraries and tools

---

## Support

For issues, questions, or feature requests, please:

1. Check existing [Issues](https://github.com/MishraShardendu22/Security-Scanner-StackGaurd-Assignment/issues)
2. Create a new issue with detailed information
3. Include logs, error messages, and reproduction steps

---

<div align="center">

**Star this repository if you find it helpful!**

Made by Shardendu Mishra

</div>

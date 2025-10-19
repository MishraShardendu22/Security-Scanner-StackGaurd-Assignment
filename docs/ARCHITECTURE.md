# StackGuard Security Scanner - Architecture Documentation

## Table of Contents
1. [Overview](#overview)
2. [Technology Stack](#technology-stack)
3. [Project Structure](#project-structure)
4. [Core Components](#core-components)
5. [Design Patterns](#design-patterns)
6. [Data Flow](#data-flow)

---

## Overview

StackGuard is a comprehensive security scanner designed to detect leaked secrets, API keys, and sensitive information in AI/ML resources hosted on platforms like Hugging Face. The application scans models, datasets, and spaces for potential security vulnerabilities.

### Key Features
- **Multi-Resource Scanning**: Scan AI models, datasets, and spaces
- **Pattern-Based Detection**: 15+ secret patterns (AWS keys, GitHub PATs, API keys, etc.)
- **Real-Time Dashboard**: Live updates using HTMX
- **Concurrent Processing**: Efficient parallel scanning using Go routines
- **Web Interface**: Modern, responsive UI with Templ templates

---

## Technology Stack

### Backend Technologies

#### 1. **Go (Golang) 1.21+**
- **Purpose**: Primary backend language
- **Why Go**: 
  - Excellent concurrency support (goroutines and channels)
  - Fast compilation and execution
  - Strong standard library
  - Memory efficient
  - Built-in HTTP server

#### 2. **Fiber Web Framework**
- **Version**: v2.x
- **Purpose**: HTTP web framework
- **Why Fiber**:
  - Express.js-inspired API (easy to learn)
  - Fastest Go web framework
  - Built on fasthttp
  - Rich middleware ecosystem
  - Low memory footprint
  
**Key Features Used**:
- Routing and middleware
- Request/response handling
- JSON parsing and rendering
- Error handling
- Static file serving
- CORS support
- Request logging
- Panic recovery

#### 3. **MongoDB**
- **Purpose**: Primary database
- **Why MongoDB**:
  - Flexible schema for varying scan results
  - Good for nested document structures
  - Fast writes for scan results
  - Easy to scale horizontally
  
**ODM**: MGM (Mongo Go Models)
- Simplifies MongoDB operations
- Model-based approach
- Built on official MongoDB Go driver

#### 4. **Templ Template Engine**
- **Purpose**: Type-safe HTML templating
- **Why Templ**:
  - Compile-time type safety
  - Go syntax (no new template language)
  - Better IDE support
  - Performance (compiles to Go code)
  - Component-based architecture

### Frontend Technologies

#### 1. **HTMX**
- **Version**: 1.9.10
- **Purpose**: Dynamic HTML interactions
- **Why HTMX**:
  - No JavaScript frameworks needed
  - Declarative AJAX calls
  - Automatic UI updates
  - Small footprint (~14kb)
  - Progressive enhancement
  
**Features Used**:
- `hx-get`: GET requests
- `hx-post`: Form submissions
- `hx-trigger`: Event handling
- `hx-swap`: Content replacement
- `hx-target`: Element targeting
- Polling for live updates

#### 2. **Tailwind CSS**
- **Purpose**: Utility-first CSS framework
- **Why Tailwind**:
  - Rapid development
  - Consistent design system
  - Responsive by default
  - Small production builds
  - Customizable

#### 3. **Font Awesome**
- **Version**: 6.4.0
- **Purpose**: Icon library
- **Usage**: UI icons throughout the application

---

## Project Structure

```
Security-Scanner-StackGaurd-Assignment/
├── controller/           # Business logic and handlers
│   ├── dashboard.go     # Dashboard statistics
│   ├── fetch.go         # Fetch resources from HuggingFace
│   ├── org-specific.go  # Organization-wide scanning
│   ├── results.go       # Results retrieval
│   ├── scan.go          # Main scanning logic
│   ├── unified.go       # Unified scan operations
│   └── web.go           # Web page handlers
├── database/            # Database connection and setup
│   └── database.go
├── models/              # Data models and schemas
│   ├── config.go        # Configuration model
│   ├── model.go         # Core data models
│   └── secret.go        # Secret pattern definitions
├── route/               # Route definitions
│   ├── fetch.go         # Fetch route setup
│   ├── org-specific.go  # Org-specific routes
│   ├── results.go       # Results routes
│   ├── scan.go          # Scan routes
│   └── web.go           # Web page routes
├── template/            # Templ templates
│   ├── layout.templ     # Base layout
│   ├── index.templ      # Home page
│   ├── dashboard.templ  # Dashboard page
│   ├── scan.templ       # Scan form page
│   ├── results.templ    # Results pages
│   ├── api-tester.templ # API testing page
│   └── *_templ.go       # Generated Go files
├── util/                # Utility functions
│   ├── apiResponse.go   # API response helpers
│   ├── getEnv.go        # Environment variables
│   ├── mask.go          # Secret masking
│   └── secret.go        # Secret detection logic
├── logs/                # Application logs
├── public/              # Static assets
├── main.go              # Application entry point
├── go.mod               # Go module definition
└── .env                 # Environment configuration
```

---

## Core Components

### 1. Main Application (`main.go`)

**Responsibilities**:
- Application bootstrapping
- Configuration loading
- Database initialization
- Server setup
- Graceful shutdown
- Logging configuration

**Key Functions**:
```go
func main()                                    // Entry point
func loadConfig() *models.Config              // Load configuration
func setupLogger(config *models.Config)       // Configure logging
func setupMiddleware(app *fiber.App, ...)     // Setup middleware
func SetUpRoutes(app *fiber.App, ...)         // Register routes
func gracefulShutdown(app *fiber.App, ...)    // Handle shutdown
```

### 2. Controllers

#### Scan Controller (`controller/scan.go`)
- Orchestrates scanning operations
- Manages concurrent scanning
- Coordinates with fetch operations
- Saves scan results

**Key Functions**:
```go
func ScanModelHandler(c *fiber.Ctx)           // Scan AI models
func ScanDatasetHandler(c *fiber.Ctx)         // Scan datasets
func ScanSpaceHandler(c *fiber.Ctx)           // Scan spaces
func scanResources(aiRequest *models.AI_REQUEST, ...) // Core scanning
```

#### Dashboard Controller (`controller/dashboard.go`)
- Aggregates statistics
- Provides dashboard data
- Calculates metrics

**Key Functions**:
```go
func GetDashboardStats(c *fiber.Ctx)          // Get dashboard data
```

#### Results Controller (`controller/results.go`)
- Retrieves scan results
- Formats result data
- Renders result pages

**Key Functions**:
```go
func GetResultsPage(c *fiber.Ctx)             // List all results
func GetResultDetailPage(c *fiber.Ctx)        // Single result detail
func GetAllResults(c *fiber.Ctx)              // API endpoint
```

### 3. Models

#### Core Models (`models/model.go`)

```go
type AI_REQUEST struct {
    RequestID    string      // Unique scan identifier
    Org          string      // Organization/user name
    URLs         []string    // Resource URLs to scan
    ResourceType string      // Type: model/dataset/space
    CreatedAt    time.Time   // Request timestamp
}

type SCAN_RESULT struct {
    RequestID        string            // Links to AI_REQUEST
    ScannedResources []ScannedResource // All scanned items
    CreatedAt        time.Time         // Scan completion time
}

type ScannedResource struct {
    ID       string           // Resource identifier
    Type     string           // Resource type
    Findings []SecretFinding  // Detected secrets
}

type SecretFinding struct {
    SecretType       string   // Type of secret (AWS Key, etc.)
    Pattern          string   // Regex pattern used
    Secret           string   // Masked secret value
    SourceType       string   // Where found (file/PR/discussion)
    FileName         string   // File name if applicable
    Line             int      // Line number if applicable
    DiscussionTitle  string   // Discussion title if applicable
    DiscussionNum    int      // Discussion number if applicable
}
```

#### Secret Patterns (`models/secret.go`)

Defines 15+ secret detection patterns:
- AWS Access Key ID
- AWS Secret Access Key
- GitHub Personal Access Token
- GitHub OAuth Token
- GitLab PAT
- Slack Token
- Slack Webhook
- Google API Key
- Google OAuth
- Algolia API Key
- Firebase URL
- Heroku API Key
- MailChimp API Key
- MailGun API Key
- PayPal Braintree
- Picatic API Key

### 4. Database Layer

#### Connection Manager (`database/database.go`)

```go
func ConnectDatabase(dbName, uri string) error
```

**Features**:
- Connection pooling
- Error handling
- Context management
- MGM initialization

### 5. Utilities

#### Secret Detection (`util/secret.go`)

```go
func ScanSecret(text string) []models.SecretFinding
```

- Regex-based pattern matching
- Multi-pattern scanning
- Efficient text processing

#### Secret Masking (`util/mask.go`)

```go
func MaskSecret(secret string) string
```

- Preserves first/last 4 characters
- Masks middle with asterisks
- Handles short secrets

#### API Response (`util/apiResponse.go`)

```go
func ResponseAPI(c *fiber.Ctx, code int, msg string, data interface{}, requestID string)
```

- Standardized response format
- Consistent error handling
- Request ID tracking

---

## Design Patterns

### 1. **MVC Architecture**
- **Models**: Data structures and database schemas
- **Views**: Templ templates for rendering HTML
- **Controllers**: Business logic and request handling

### 2. **Repository Pattern**
- Abstraction over database operations
- Centralized data access
- Easy to test and mock

### 3. **Middleware Pattern**
- Request/response interceptors
- Cross-cutting concerns (logging, CORS, recovery)
- Composable and reusable

### 4. **Concurrent Worker Pool**
- Semaphore-based concurrency control
- Prevents resource exhaustion
- Efficient parallel processing

### 5. **Graceful Shutdown Pattern**
- Signal handling
- Connection draining
- Clean resource cleanup

---

## Data Flow

### Scan Request Flow

```
User Request (Web/API)
    ↓
Fiber Router
    ↓
Scan Controller
    ↓
Create AI_REQUEST → MongoDB
    ↓
Fetch Resource Data (concurrent)
    ↓
Scan Each Resource (concurrent)
    ↓
Pattern Matching (Regex)
    ↓
Collect Findings
    ↓
Create SCAN_RESULT → MongoDB
    ↓
Return Response (Request ID)
    ↓
User Views Results
```

### Dashboard Update Flow

```
Browser (HTMX every 5s)
    ↓
GET /api/dashboard
    ↓
Dashboard Controller
    ↓
Query MongoDB (Aggregations)
    ↓
Calculate Statistics
    ↓
Return JSON
    ↓
HTMX Updates DOM
```

### Results Viewing Flow

```
User Requests Results List
    ↓
GET /results
    ↓
Results Controller
    ↓
Query MongoDB
    ↓
Render Templ Template
    ↓
Return HTML
    ↓
HTMX Auto-refresh (every 10s)
```

---

## Performance Considerations

### 1. **Concurrency**
- Multiple resources scanned in parallel
- Semaphore limits prevent overload
- Efficient goroutine management

### 2. **Database Optimization**
- Indexed fields (RequestID, CreatedAt)
- Efficient queries
- Connection pooling

### 3. **Caching Strategy**
- Template compilation at build time
- Static asset serving
- Browser caching headers

### 4. **Memory Management**
- Stream processing for large files
- Bounded worker pools
- Efficient string operations

---

## Security Features

1. **Secret Masking**: Detected secrets are masked before storage
2. **Input Validation**: All inputs validated
3. **Error Handling**: Proper error messages without leaking info
4. **Logging**: Security events logged
5. **CORS Configuration**: Controlled API access

---

## Scalability

### Horizontal Scaling
- Stateless application design
- MongoDB sharding support
- Load balancer compatible

### Vertical Scaling
- Configurable worker pool sizes
- Adjustable connection limits
- Memory-efficient operations

---

## Monitoring & Observability

1. **Structured Logging**: JSON logs with context
2. **Error Tracking**: Comprehensive error logging
3. **Metrics**: Request counts, timing, success rates
4. **Health Checks**: Application status endpoint

---

## Future Enhancements

1. **Caching Layer**: Redis for dashboard data
2. **Queue System**: Background job processing
3. **WebSocket Support**: Real-time scan progress
4. **API Rate Limiting**: Protect against abuse
5. **Authentication**: User management system
6. **Audit Trail**: Detailed activity logs

---

*Last Updated: October 19, 2025*

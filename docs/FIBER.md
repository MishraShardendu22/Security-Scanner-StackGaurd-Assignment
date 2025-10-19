# Fiber Framework Implementation Guide

## Table of Contents
1. [Overview](#overview)
2. [Why Fiber?](#why-fiber)
3. [Fiber Setup](#fiber-setup)
4. [Routing](#routing)
5. [Middleware](#middleware)
6. [Request Handling](#request-handling)
7. [Response Formatting](#response-formatting)
8. [Error Handling](#error-handling)
9. [Best Practices](#best-practices)

---

## Overview

Fiber is the HTTP web framework powering StackGuard. It's built on top of Fasthttp, making it one of the fastest Go web frameworks available.

### Fiber vs Other Frameworks

| Feature | Fiber | Gin | Echo | net/http |
|---------|-------|-----|------|----------|
| Performance | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| Express-like API | ✅ | ❌ | ❌ | ❌ |
| Built-in middleware | ✅ | ✅ | ✅ | ❌ |
| Low memory | ✅ | ✅ | ✅ | ❌ |
| Zero allocation | ✅ | ❌ | ❌ | ❌ |

---

## Why Fiber?

### 1. **Performance**
- Built on Fasthttp (fastest HTTP engine for Go)
- Zero memory allocation router
- Optimized for speed and low memory footprint

### 2. **Developer Experience**
- Express.js-inspired API (familiar to JavaScript developers)
- Extensive documentation
- Active community
- 100+ middleware packages

### 3. **Features**
- Route parameters and wildcards
- Static file serving
- Middleware support
- Template engines
- WebSocket support
- Server-Sent Events (SSE)

### 4. **Benchmarks**
```
Framework       Requests/sec    Latency     Memory
Fiber           634,000         0.28ms      1.2 MB
Gin             462,000         0.38ms      1.8 MB
Echo            453,000         0.39ms      1.7 MB
net/http        298,000         0.59ms      2.4 MB
```

---

## Fiber Setup

### Installation

```bash
go get -u github.com/gofiber/fiber/v2
```

### Basic Configuration

**File**: `main.go`, Lines: 51-69

```go
app := fiber.New(fiber.Config{
    // Application name
    AppName: "Security Scanner",
    
    // Custom server header
    ServerHeader: "Security-Scanner",
    
    // Timeouts
    ReadTimeout:  30 * time.Second,  // Max time to read request
    WriteTimeout: 30 * time.Second,  // Max time to write response
    IdleTimeout:  120 * time.Second, // Max keep-alive time
    
    // Custom error handler
    ErrorHandler: func(c *fiber.Ctx, err error) error {
        // Log the error
        logger.Error("request error", slog.Group("req",
            slog.String("method", c.Method()),
            slog.String("path", c.Path()),
            slog.String("error", err.Error()),
        ))
        
        // Determine status code
        code := fiber.StatusInternalServerError
        if e, ok := err.(*fiber.Error); ok {
            code = e.Code
        }
        
        // Return standardized error response
        return util.ResponseAPI(c, code, err.Error(), nil, "")
    },
})
```

**Key Configuration Options**:

1. **AppName**: Identifies app in headers
2. **ServerHeader**: Custom server identification
3. **ReadTimeout**: Prevents slow-read attacks
4. **WriteTimeout**: Prevents slow-write attacks
5. **IdleTimeout**: Connection keep-alive duration
6. **ErrorHandler**: Centralized error handling

---

## Routing

### Route Registration

**File**: `main.go`, Lines: 89-98

```go
func SetUpRoutes(app *fiber.App, logger *slog.Logger, config *models.Config) {
    // Web routes (HTML pages)
    route.RegisterWebRoutes(app)

    // API routes
    app.Get("/api/test", func(c *fiber.Ctx) error {
        return util.ResponseAPI(c, fiber.StatusOK, "API is working fine", nil, "")
    })
    
    route.SetupFetchRoutes(app)
    route.SetupOrgRoutes(app)
    route.SetupScanRoutes(app)
    route.SetupResultsRoutes(app)
}
```

### Web Routes

**File**: `route/web.go`

```go
func RegisterWebRoutes(app *fiber.App) {
    // Home page
    app.Get("/", func(c *fiber.Ctx) error {
        c.Set("Content-Type", "text/html; charset=utf-8")
        return template.IndexNew().Render(c.Context(), c.Response().BodyWriter())
    })

    // Dashboard page
    app.Get("/dashboard", func(c *fiber.Ctx) error {
        c.Set("Content-Type", "text/html; charset=utf-8")
        return template.Dashboard().Render(c.Context(), c.Response().BodyWriter())
    })

    // Dashboard API endpoint
    app.Get("/api/dashboard", controller.GetDashboardStats)

    // Scan page
    app.Get("/scan", func(c *fiber.Ctx) error {
        c.Set("Content-Type", "text/html; charset=utf-8")
        return template.ScanForm().Render(c.Context(), c.Response().BodyWriter())
    })

    // Results pages
    app.Get("/results", controller.GetResultsPage)
    app.Get("/results/:request_id", controller.GetResultDetailPage)

    // API Testing page
    app.Get("/api-tester", func(c *fiber.Ctx) error {
        c.Set("Content-Type", "text/html; charset=utf-8")
        return template.APITester().Render(c.Context(), c.Response().BodyWriter())
    })
}
```

### API Routes

**File**: `route/scan.go`

```go
func SetupScanRoutes(app *fiber.App) {
    // Unified scan endpoint
    app.Post("/scan", controller.UnifiedScan)
    
    // Specific resource type endpoints
    app.Post("/scan/model", controller.ScanModelHandler)
    app.Post("/scan/dataset", controller.ScanDatasetHandler)
    app.Post("/scan/space", controller.ScanSpaceHandler)
}
```

### Route Groups

```go
// Group related routes
api := app.Group("/api")
{
    api.Get("/dashboard", controller.GetDashboardStats)
    api.Get("/results", controller.GetAllResults)
    
    // Version 1 API
    v1 := api.Group("/v1")
    {
        v1.Post("/scan", controller.UnifiedScan)
        v1.Get("/results/:id", controller.GetResultByID)
    }
}
```

### Route Parameters

```go
// Dynamic route parameter
app.Get("/results/:request_id", func(c *fiber.Ctx) error {
    requestID := c.Params("request_id")
    // Use requestID to fetch data
    return c.JSON(fiber.Map{
        "request_id": requestID,
    })
})

// Multiple parameters
app.Get("/users/:id/posts/:postId", func(c *fiber.Ctx) error {
    userID := c.Params("id")
    postID := c.Params("postId")
    return c.SendString(fmt.Sprintf("User: %s, Post: %s", userID, postID))
})

// Optional parameter
app.Get("/user/:name?", func(c *fiber.Ctx) error {
    name := c.Params("name", "Guest") // Default to "Guest"
    return c.SendString("Hello, " + name)
})
```

### Wildcard Routes

```go
// Catch-all
app.Get("/static/*", func(c *fiber.Ctx) error {
    path := c.Params("*")
    return c.SendFile("./public/" + path)
})
```

---

## Middleware

### Built-in Middleware Used

**File**: `main.go`, Lines: 116-134

```go
func setupMiddleware(app *fiber.App, config *models.Config) {
    // CORS middleware
    app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    }))

    // Request logging middleware
    app.Use(logger.New(logger.Config{
        Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
        TimeFormat: "2006-01-02 15:04:05",
        Output: os.Stdout,
    }))

    // Panic recovery middleware
    app.Use(recover.New(recover.Config{
        EnableStackTrace: true,
    }))

    // Static file serving
    app.Static("/public", "./public")
}
```

### 1. CORS Middleware

**Purpose**: Handle Cross-Origin Resource Sharing

```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "*",                              // Allow all origins
    AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",    // Allowed HTTP methods
    AllowHeaders: "Origin, Content-Type, Accept",   // Allowed headers
    AllowCredentials: true,                         // Allow cookies
    ExposeHeaders: "Content-Length",                // Exposed headers
    MaxAge: 3600,                                   // Preflight cache duration
}))
```

**When It Runs**: Before every request

**What It Does**:
- Adds CORS headers to responses
- Handles preflight OPTIONS requests
- Allows/blocks cross-origin requests

### 2. Logger Middleware

**Purpose**: Log all HTTP requests

```go
app.Use(logger.New(logger.Config{
    Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
    TimeFormat: "2006-01-02 15:04:05",
    Output: os.Stdout,
}))
```

**Output Example**:
```
[2025-10-19 15:30:42] 200 - 45ms GET /dashboard
[2025-10-19 15:30:43] 201 - 1.2s POST /scan
[2025-10-19 15:30:44] 404 - 2ms GET /notfound
```

**Available Tags**:
- `${time}`: Request timestamp
- `${status}`: HTTP status code
- `${latency}`: Request duration
- `${method}`: HTTP method
- `${path}`: Request path
- `${ip}`: Client IP
- `${ua}`: User agent
- `${body}`: Request body
- `${error}`: Error message

### 3. Recover Middleware

**Purpose**: Recover from panics

```go
app.Use(recover.New(recover.Config{
    EnableStackTrace: true,
    StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
        log.Printf("Panic: %v\n", e)
    },
}))
```

**What It Does**:
- Catches panics in handlers
- Logs stack trace
- Returns 500 error instead of crashing
- Keeps server running

**Example**:
```go
app.Get("/panic", func(c *fiber.Ctx) error {
    panic("Something went wrong!") // Server won't crash
})
```

### Custom Middleware

```go
// Authentication middleware
func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := c.Get("Authorization")
        
        if token == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Unauthorized",
            })
        }
        
        // Validate token
        // ...
        
        // Continue to next handler
        return c.Next()
    }
}

// Usage
api := app.Group("/api", AuthMiddleware())
```

---

## Request Handling

### Reading Request Data

#### 1. Path Parameters

```go
app.Get("/results/:request_id", func(c *fiber.Ctx) error {
    requestID := c.Params("request_id")
    return c.SendString("Request ID: " + requestID)
})
```

#### 2. Query Parameters

```go
app.Get("/search", func(c *fiber.Ctx) error {
    query := c.Query("q")           // Get query parameter
    page := c.Query("page", "1")    // With default value
    
    return c.JSON(fiber.Map{
        "query": query,
        "page": page,
    })
})
// URL: /search?q=test&page=2
```

#### 3. Request Body (JSON)

**File**: `controller/scan.go`

```go
func UnifiedScan(c *fiber.Ctx) error {
    var req models.ScanRequest
    
    // Parse JSON body
    if err := c.BodyParser(&req); err != nil {
        return util.ResponseAPI(c, fiber.StatusBadRequest, 
            "Invalid request body", nil, "")
    }
    
    // Validate required fields
    if req.Org == "" {
        return util.ResponseAPI(c, fiber.StatusBadRequest, 
            "Organization is required", nil, "")
    }
    
    // Process request
    // ...
}
```

**Request Example**:
```json
{
  "org": "huggingface",
  "model_id": "bert-base-uncased",
  "include_prs": false,
  "include_discussions": false
}
```

#### 4. Form Data

```go
app.Post("/upload", func(c *fiber.Ctx) error {
    name := c.FormValue("name")
    file, err := c.FormFile("file")
    
    if err != nil {
        return err
    }
    
    return c.SaveFile(file, "./uploads/" + file.Filename)
})
```

#### 5. Headers

```go
app.Get("/headers", func(c *fiber.Ctx) error {
    contentType := c.Get("Content-Type")
    userAgent := c.Get("User-Agent")
    
    // Set response header
    c.Set("X-Custom-Header", "value")
    
    return c.SendString("Headers received")
})
```

### Context Methods

```go
// Request info
c.Method()        // GET, POST, etc.
c.Path()          // /api/scan
c.IP()            // Client IP
c.Hostname()      // example.com
c.Protocol()      // http or https
c.BaseURL()       // http://example.com

// Headers
c.Get("key")      // Get header
c.Set("key", "val") // Set header
c.GetReqHeaders() // All request headers

// Cookies
c.Cookies("name") // Get cookie
c.Cookie(&fiber.Cookie{...}) // Set cookie

// Request body
c.Body()          // Raw body bytes
c.BodyParser(&v)  // Parse to struct

// Query & Params
c.Query("key")    // Query parameter
c.Params("key")   // Path parameter

// Response
c.Status(200)     // Set status
c.SendString("hi") // Send string
c.JSON(data)      // Send JSON
c.SendFile(path)  // Send file
```

---

## Response Formatting

### Standard JSON Response

**File**: `util/apiResponse.go`

```go
func ResponseAPI(c *fiber.Ctx, code int, msg string, data interface{}, requestID string) error {
    response := fiber.Map{
        "status": func() string {
            if code >= 200 && code < 300 {
                return "success"
            }
            return "error"
        }(),
        "message": msg,
        "data":    data,
    }
    
    if requestID != "" {
        response["request_id"] = requestID
    }
    
    return c.Status(code).JSON(response)
}
```

**Usage**:
```go
// Success response
return util.ResponseAPI(c, fiber.StatusOK, 
    "Scan completed", scanResult, requestID)

// Error response
return util.ResponseAPI(c, fiber.StatusBadRequest, 
    "Invalid input", nil, "")
```

**Response Format**:
```json
{
  "status": "success",
  "message": "Scan completed",
  "data": {
    "request_id": "abc123",
    "resources": [...]
  },
  "request_id": "abc123"
}
```

### Different Response Types

```go
// JSON
c.JSON(fiber.Map{"key": "value"})

// String
c.SendString("Hello, World!")

// HTML
c.Type("html").SendString("<h1>Hello</h1>")

// File
c.SendFile("./file.pdf")

// Download
c.Download("./file.pdf", "custom-name.pdf")

// Redirect
c.Redirect("/new-url", 301)

// Stream
c.SendStream(reader)

// Status only
c.SendStatus(204)
```

---

## Error Handling

### Global Error Handler

**File**: `main.go`, Lines: 59-69

```go
ErrorHandler: func(c *fiber.Ctx, err error) error {
    // Log error with context
    logger.Error("request error", slog.Group("req",
        slog.String("method", c.Method()),
        slog.String("path", c.Path()),
        slog.String("error", err.Error()),
    ))
    
    // Determine status code
    code := fiber.StatusInternalServerError
    if e, ok := err.(*fiber.Error); ok {
        code = e.Code
    }
    
    // Return standardized response
    return util.ResponseAPI(c, code, err.Error(), nil, "")
}
```

### Throwing Errors

```go
// Using fiber.NewError
if invalid {
    return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
}

// Using fiber.ErrXXX constants
if notFound {
    return fiber.ErrNotFound
}

// Custom error
return errors.New("something went wrong")
```

### Built-in Error Constants

```go
fiber.ErrBadRequest           // 400
fiber.ErrUnauthorized         // 401
fiber.ErrForbidden            // 403
fiber.ErrNotFound             // 404
fiber.ErrMethodNotAllowed     // 405
fiber.ErrInternalServerError  // 500
fiber.ErrServiceUnavailable   // 503
```

---

## Best Practices

### 1. Use Route Groups

```go
// Group related routes
api := app.Group("/api")
v1 := api.Group("/v1")
v1.Post("/scan", scanHandler)
v1.Get("/results", resultsHandler)
```

### 2. Reuse Context

```go
// DON'T: Create new context
return c.Status(200).JSON(data)

// DO: Chain methods
return c.Status(200).JSON(data)
```

### 3. Validate Input

```go
if req.Org == "" {
    return fiber.NewError(400, "Organization required")
}
```

### 4. Use Constants

```go
const (
    ErrInvalidOrg = "Invalid organization"
    ErrNoResources = "No resources found"
)
```

### 5. Centralize Responses

```go
// Use helper function
return util.ResponseAPI(c, status, message, data, reqID)
```

---

*Last Updated: October 19, 2025*

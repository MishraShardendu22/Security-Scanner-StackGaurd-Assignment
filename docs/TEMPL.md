# Templ Template Engine Documentation

## Table of Contents
1. [Overview](#overview)
2. [Why Templ?](#why-templ)
3. [Installation & Setup](#installation--setup)
4. [Syntax Guide](#syntax-guide)
5. [Components](#components)
6. [Integration with Fiber](#integration-with-fiber)
7. [HTMX Integration](#htmx-integration)
8. [Best Practices](#best-practices)

---

## Overview

Templ is a type-safe templating language for Go that compiles to Go code. Unlike traditional template engines that use string-based templates, Templ provides compile-time type safety and better IDE support.

### Key Features
- **Type-safe**: Compile-time checks catch errors early
- **Go syntax**: No new template language to learn
- **Fast**: Compiles to native Go code
- **IDE support**: Full autocomplete and type checking
- **Component-based**: Reusable template components
- **Props validation**: Function parameters for props

---

## Why Templ?

### Comparison with Alternatives

| Feature | Templ | html/template | Pongo2 | Jet |
|---------|-------|---------------|--------|-----|
| Type Safety | ✅ | ❌ | ❌ | ❌ |
| Compile Time | ✅ | ❌ | ❌ | ❌ |
| IDE Support | ✅ | Limited | Limited | Limited |
| Performance | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| Go Syntax | ✅ | ❌ | ❌ | ❌ |

### Problems with Traditional Templates

```go
// html/template - No type safety!
tmpl.Execute(w, map[string]interface{}{
    "name": user.Name,      // Typos not caught
    "age": user.Age,        // Wrong type not caught
    "missing": "value",     // Extra fields not caught
})
```

### Templ Solution

```templ
// Templ - Type-safe!
templ UserProfile(user User) {
    <h1>{user.Name}</h1>    // Autocomplete works
    <p>{user.Age}</p>       // Type checked
}
// Missing fields = compile error ✅
```

---

## Installation & Setup

### 1. Install Templ CLI

```bash
go install github.com/a-h/templ/cmd/templ@latest
```

### 2. Add to Go Module

```bash
go get github.com/a-h/templ
```

### 3. Create Templ File

**File**: `template/layout.templ`

```templ
package template

templ Layout(title string) {
    <!DOCTYPE html>
    <html>
        <head>
            <title>{title}</title>
        </head>
        <body>
            {children...}
        </body>
    </html>
}
```

### 4. Generate Go Code

```bash
templ generate
```

This creates `layout_templ.go` with compiled Go code.

### 5. Use in Go Code

```go
func HomeHandler(c *fiber.Ctx) error {
    c.Set("Content-Type", "text/html; charset=utf-8")
    return template.Layout("Home").Render(c.Context(), c.Response().BodyWriter())
}
```

---

## Syntax Guide

### Basic Interpolation

```templ
templ Hello(name string) {
    <h1>Hello, {name}!</h1>
}
```

### Conditionals

```templ
templ UserBadge(user User) {
    if user.IsAdmin {
        <span class="badge admin">Admin</span>
    } else if user.IsModerator {
        <span class="badge mod">Moderator</span>
    } else {
        <span class="badge user">User</span>
    }
}
```

### Loops

```templ
templ UserList(users []User) {
    <ul>
        for _, user := range users {
            <li>{user.Name}</li>
        }
    </ul>
}
```

### Switch Statements

```templ
templ StatusBadge(status string) {
    switch status {
        case "active":
            <span class="badge green">Active</span>
        case "pending":
            <span class="badge yellow">Pending</span>
        case "inactive":
            <span class="badge red">Inactive</span>
        default:
            <span class="badge gray">Unknown</span>
    }
}
```

### Component Composition

```templ
templ Layout(title string) {
    <!DOCTYPE html>
    <html>
        <head><title>{title}</title></head>
        <body>
            @Header()
            {children...}
            @Footer()
        </body>
    </html>
}

templ Header() {
    <header>
        <nav>Navigation</nav>
    </header>
}

templ Footer() {
    <footer>© 2025</footer>
}
```

### Using Components

```templ
templ HomePage() {
    @Layout("Home") {
        <main>
            <h1>Welcome</h1>
            <p>This is the home page</p>
        </main>
    }
}
```

---

## Components

### Layout Component

**File**: `template/layout.templ`

```templ
package template

templ Layout(title string) {
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <meta charset="UTF-8"/>
            <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
            <title>{ title } - Security Scanner</title>
            <script src="https://cdn.tailwindcss.com"></script>
            <script src="https://unpkg.com/htmx.org@1.9.10"></script>
            <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css"/>
            <style>
                :root {
                    --primary-yellow: #FFC107;
                    --dark-yellow: #FFA000;
                    --light-yellow: #FFECB3;
                    --primary-black: #1A1A1A;
                    --secondary-black: #2D2D2D;
                    --accent-black: #404040;
                }
                
                .htmx-swapping {
                    opacity: 0;
                    transition: opacity 200ms ease-out;
                }
                
                .htmx-settling {
                    opacity: 1;
                    transition: opacity 200ms ease-in;
                }
                
                .htmx-indicator {
                    display: none;
                }
                
                .htmx-request .htmx-indicator {
                    display: block;
                }
            </style>
        </head>
        <body class="bg-black min-h-screen">
            @Navigation()
            <main class="w-full px-8 py-8">
                { children... }
            </main>
        </body>
    </html>
}

templ Navigation() {
    <nav class="bg-black border-b-4 border-yellow-400 shadow-lg">
        <div class="w-full px-8 py-4">
            <div class="flex items-center justify-between">
                <div class="flex items-center space-x-3">
                    <i class="fas fa-shield-alt text-2xl text-yellow-400"></i>
                    <h1 class="text-2xl font-bold text-white">StackGuard Security Scanner</h1>
                </div>
                <div class="flex space-x-6">
                    <a href="/" class="text-white hover:text-yellow-400 transition font-medium">
                        <i class="fas fa-home mr-2"></i>Home
                    </a>
                    <a href="/dashboard" class="text-white hover:text-yellow-400 transition font-medium">
                        <i class="fas fa-chart-line mr-2"></i>Dashboard
                    </a>
                    <a href="/scan" class="text-white hover:text-yellow-400 transition font-medium">
                        <i class="fas fa-search mr-2"></i>New Scan
                    </a>
                    <a href="/results" class="text-white hover:text-yellow-400 transition font-medium">
                        <i class="fas fa-list mr-2"></i>Results
                    </a>
                    <a href="/api-tester" class="text-white hover:text-yellow-400 transition font-medium">
                        <i class="fas fa-flask mr-2"></i>API Tester
                    </a>
                </div>
            </div>
        </div>
    </nav>
}
```

**Usage**:
```templ
templ Dashboard() {
    @Layout("Dashboard") {
        <h1>Dashboard Content</h1>
    }
}
```

### Results List Component

**File**: `template/results.templ`

```templ
import (
    "github.com/MishraShardendu22/Scanner/models"
    "fmt"
)

templ ResultsList(results []models.SCAN_RESULT) {
    @Layout("Scan Results") {
        <div class="max-w-6xl mx-auto">
            <div class="bg-black border-4 border-yellow-400 rounded-lg shadow-lg p-8">
                <h2 class="text-3xl font-bold mb-6 text-yellow-400">
                    <i class="fas fa-list-alt mr-3"></i>All Scan Results
                </h2>
                if len(results) == 0 {
                    <div class="text-center py-12">
                        <i class="fas fa-inbox text-6xl text-gray-600 mb-4"></i>
                        <p class="text-xl text-gray-400 mb-4">No scan results found</p>
                        <a href="/scan" class="inline-block bg-yellow-400 hover:bg-yellow-500 text-black font-bold py-3 px-6 rounded-lg transition">
                            <i class="fas fa-play mr-2"></i>Start New Scan
                        </a>
                    </div>
                } else {
                    <div class="grid gap-6">
                        for _, result := range results {
                            @ResultItem(result)
                        }
                    </div>
                }
            </div>
        </div>
    }
}

templ ResultItem(result models.SCAN_RESULT) {
    <div class="border-l-4 border-yellow-400 bg-gray-900 rounded-r-lg p-6">
        <div class="flex items-center justify-between mb-4">
            <div>
                <h3 class="text-xl font-bold text-yellow-400">Request ID: { result.RequestID }</h3>
                <p class="text-sm text-gray-400 mt-1">Created: { result.CreatedAt.Format("2006-01-02 15:04:05") }</p>
            </div>
            <a href={ templ.URL(fmt.Sprintf("/results/%s", result.RequestID)) } 
               class="bg-yellow-400 hover:bg-yellow-500 text-black font-bold py-2 px-4 rounded-lg transition">
                <i class="fas fa-eye mr-2"></i>View Details
            </a>
        </div>
        <div class="flex items-center space-x-4 text-sm">
            <span class="text-gray-300 font-medium">
                <i class="fas fa-folder mr-2 text-yellow-400"></i>
                { fmt.Sprintf("%d", len(result.ScannedResources)) } Resources
            </span>
            <span class="text-gray-300 font-medium">
                <i class="fas fa-exclamation-triangle mr-2 text-yellow-400"></i>
                { fmt.Sprintf("%d", countTotalFindings(result.ScannedResources)) } Findings
            </span>
        </div>
    </div>
}

// Helper function
func countTotalFindings(resources []models.ScannedResource) int {
    total := 0
    for _, r := range resources {
        total += len(r.Findings)
    }
    return total
}
```

---

## Integration with Fiber

### Rendering in Handlers

**Pattern**:
```go
func Handler(c *fiber.Ctx) error {
    // Set content type
    c.Set("Content-Type", "text/html; charset=utf-8")
    
    // Render templ component
    return template.ComponentName(params).Render(
        c.Context(),                  // Context
        c.Response().BodyWriter(),    // Writer
    )
}
```

**Example - Home Page**:
```go
// File: route/web.go
app.Get("/", func(c *fiber.Ctx) error {
    c.Set("Content-Type", "text/html; charset=utf-8")
    return template.IndexNew().Render(c.Context(), c.Response().BodyWriter())
})
```

**Example - Dashboard**:
```go
app.Get("/dashboard", func(c *fiber.Ctx) error {
    c.Set("Content-Type", "text/html; charset=utf-8")
    return template.Dashboard().Render(c.Context(), c.Response().BodyWriter())
})
```

**Example - Results with Data**:
```go
// File: controller/results.go
func GetResultsPage(c *fiber.Ctx) error {
    var results []models.SCAN_RESULT
    
    // Fetch from database
    err := mgm.Coll(&models.SCAN_RESULT{}).SimpleFind(&results, bson.M{})
    if err != nil {
        return err
    }
    
    // Render with data
    c.Set("Content-Type", "text/html; charset=utf-8")
    return template.ResultsList(results).Render(c.Context(), c.Response().BodyWriter())
}
```

---

## HTMX Integration

Templ works seamlessly with HTMX for dynamic updates.

### Auto-Refreshing Dashboard

```templ
templ Dashboard() {
    @Layout("Dashboard") {
        <div class="w-full" 
            hx-get="/api/dashboard" 
            hx-trigger="load, every 5s" 
            hx-target="#dashboardContent"
            hx-swap="innerHTML">
            <div id="dashboardContent">
                <!-- Content here -->
            </div>
        </div>
    }
}
```

### Form with HTMX

```templ
templ ScanForm() {
    @Layout("New Scan") {
        <form id="scanForm" 
            hx-post="/scan" 
            hx-trigger="submit" 
            hx-target="#scanResults" 
            hx-swap="innerHTML"
            hx-indicator="#loadingSpinner">
            
            <input type="text" name="org" required/>
            <button type="submit">Scan</button>
        </form>
        
        <div id="scanResults"></div>
        <div id="loadingSpinner" class="htmx-indicator">Loading...</div>
    }
}
```

### Partial Templates for HTMX

```templ
// Partial that can be returned by HTMX endpoint
templ DashboardStats(stats DashboardData) {
    <div id="statsOverview">
        <div class="stat-card">
            <h3>Total Scans</h3>
            <p>{fmt.Sprintf("%d", stats.TotalScans)}</p>
        </div>
        <div class="stat-card">
            <h3>Total Findings</h3>
            <p>{fmt.Sprintf("%d", stats.TotalFindings)}</p>
        </div>
    </div>
}
```

**API Endpoint**:
```go
app.Get("/api/dashboard", func(c *fiber.Ctx) error {
    stats := getDashboardStats()
    
    // Return HTML fragment (not full page)
    c.Set("Content-Type", "text/html; charset=utf-8")
    return template.DashboardStats(stats).Render(c.Context(), c.Response().BodyWriter())
})
```

---

## Best Practices

### 1. Component Organization

```
template/
├── layout.templ          # Base layout
├── components/
│   ├── header.templ      # Reusable header
│   ├── footer.templ      # Reusable footer
│   └── card.templ        # Reusable card
├── pages/
│   ├── home.templ        # Home page
│   ├── dashboard.templ   # Dashboard page
│   └── results.templ     # Results page
└── partials/
    ├── stats.templ       # HTMX partials
    └── scan-item.templ   # Small components
```

### 2. Type-Safe Props

```templ
// Define clear types for props
type UserCardProps struct {
    Name  string
    Email string
    Role  string
}

templ UserCard(props UserCardProps) {
    <div class="card">
        <h3>{props.Name}</h3>
        <p>{props.Email}</p>
        <span>{props.Role}</span>
    </div>
}
```

### 3. Extract Reusable Components

```templ
// Reusable button component
templ Button(text string, class string, onclick string) {
    <button class={class} onclick={onclick}>
        {text}
    </button>
}

// Usage
@Button("Submit", "btn-primary", "handleSubmit()")
```

### 4. Use Helper Functions

```templ
// In .templ file
func formatDate(t time.Time) string {
    return t.Format("Jan 2, 2006")
}

templ DateDisplay(date time.Time) {
    <span>{formatDate(date)}</span>
}
```

### 5. Conditional Classes

```templ
templ Alert(message string, isError bool) {
    <div class={templ.Classes(
        "alert",
        templ.If(isError, "alert-error"),
        templ.If(!isError, "alert-success"),
    )}>
        {message}
    </div>
}
```

### 6. Safe URLs

```templ
import "fmt"

templ Link(id string) {
    <a href={templ.URL(fmt.Sprintf("/results/%s", id))}>
        View Result
    </a>
}
```

### 7. Script Elements

```templ
templ PageWithScript() {
    @Layout("Page") {
        <div id="content">Content</div>
        
        <script type="text/javascript">
            document.getElementById('content').addEventListener('click', function() {
                alert('Clicked!');
            });
        </script>
    }
}
```

---

## Development Workflow

### 1. Watch Mode

```bash
# Auto-regenerate on file changes
templ generate --watch
```

### 2. Format Templates

```bash
templ fmt template/
```

### 3. LSP Support

Install Templ LSP for your editor:

**VS Code**: Install "templ" extension
**Neovim**: Configure templ LSP

### 4. Build Process

```bash
# 1. Generate Go code from templ files
templ generate

# 2. Build Go application
go build -o app

# 3. Run
./app
```

---

## Common Patterns

### Loading State

```templ
templ DataList(items []Item, loading bool) {
    if loading {
        <div class="spinner">Loading...</div>
    } else if len(items) == 0 {
        <p>No items found</p>
    } else {
        for _, item := range items {
            @ItemCard(item)
        }
    }
}
```

### Error State

```templ
templ DataDisplay(data Data, err error) {
    if err != nil {
        <div class="error">
            <i class="fas fa-exclamation-triangle"></i>
            <p>Error: {err.Error()}</p>
        </div>
    } else {
        @DataContent(data)
    }
}
```

### Pagination

```templ
templ Paginated(items []Item, page int, totalPages int) {
    <div class="items">
        for _, item := range items {
            @ItemCard(item)
        }
    </div>
    
    <div class="pagination">
        if page > 1 {
            <a href={templ.URL(fmt.Sprintf("?page=%d", page-1))}>Previous</a>
        }
        <span>Page {fmt.Sprintf("%d", page)} of {fmt.Sprintf("%d", totalPages)}</span>
        if page < totalPages {
            <a href={templ.URL(fmt.Sprintf("?page=%d", page+1))}>Next</a>
        }
    </div>
}
```

---

## Advantages Over Traditional Templates

### 1. Compile-Time Safety

```go
// html/template - Runtime error
tmpl.Execute(w, map[string]interface{}{
    "nam": user.Name, // Typo! Only fails at runtime
})

// Templ - Compile error
template.UserProfile(user) // Wrong type? Compile error!
```

### 2. IDE Support

- Autocomplete for struct fields
- Go to definition
- Refactoring support
- Type hints

### 3. Performance

- Compiles to Go code
- No runtime parsing
- Zero reflection overhead
- Inline optimizations

### 4. Type Safety

```templ
// Props are type-checked
templ UserCard(user User) {
    <h1>{user.Name}</h1>        // ✅ Type checked
    <p>{user.MissingField}</p>  // ❌ Compile error!
}
```

---

*Last Updated: October 19, 2025*

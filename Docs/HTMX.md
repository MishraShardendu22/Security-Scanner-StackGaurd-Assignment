# HTMX Complete Documentation

> High-performance web applications without JavaScript

## Table of Contents

- [Introduction](#introduction)
- [Quick Start](#quick-start)
- [Performance Optimization](#performance-optimization)
- [Core Concepts](#core-concepts)
- [Complete Attribute Reference](#complete-attribute-reference)
- [Real-World Examples](#real-world-examples)
- [Server-Side Integration (Go Fiber)](#server-side-integration-go-fiber)
- [Advanced Patterns](#advanced-patterns)
- [CSS & Styling](#css--styling)
- [Troubleshooting](#troubleshooting)
- [Best Practices](#best-practices)
- [Migration from JavaScript](#migration-from-javascript)

***

## Introduction

### What is HTMX?

HTMX extends HTML with attributes that enable modern browser features like AJAX, WebSockets, and Server-Sent Events without writing JavaScript. It's a 14KB library that can replace 200KB+ JavaScript frameworks.

### Why Use HTMX?

- **Simplicity**: No build step, no npm, no complex tooling
- **Performance**: Smaller bundle size, less JavaScript parsing
- **Productivity**: Write HTML, not JavaScript
- **Progressive Enhancement**: Works without JavaScript
- **Server-Side Rendering**: Natural fit for Go, Python, Ruby, PHP

### Browser Support

- Chrome/Edge: ✅
- Firefox: ✅
- Safari: ✅
- Mobile browsers: ✅
- IE11: ❌ (use HTMX 1.x)

---

## Quick Start

### Installation

#### CDN (Recommended for Production)

```html
<!DOCTYPE html>
<html>
<head>
    <!-- Performance optimization -->
    <link rel="preconnect" href="https://unpkg.com"/>
    <link rel="dns-prefetch" href="https://unpkg.com"/>
    
    <!-- HTMX -->
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
    <!-- Your content -->
</body>
</html>
```

#### NPM

```bash
npm install htmx.org
```

```javascript
import 'htmx.org';
```

### Hello World

```html
<button hx-get="/api/hello" hx-target="#result">
    Say Hello
</button>
<div id="result"></div>
```

When clicked, this button:
1. Makes GET request to `/api/hello`
2. Takes response HTML
3. Inserts it into `#result`

**No JavaScript required!**

***

## Performance Optimization

### The 5 Resource Optimization Techniques

#### 1. dns-prefetch

Resolve DNS for external domains early.

```html
<link rel="dns-prefetch" href="https://api.example.com"/>
```

**Saves**: 20-120ms  
**Use for**: Multiple third-party domains, fallback

#### 2. preconnect

Full connection setup: DNS + TCP + TLS.

```html
<link rel="preconnect" href="https://unpkg.com"/>
```

**Saves**: 100-500ms  
**Use for**: Critical CDNs (limit to 4-6 domains)  
**Cost**: High resource usage

#### 3. preload

Download specific resource for current page.

```html
<link rel="preload" href="critical.css" as="style"/>
<link rel="preload" href="hero.jpg" as="image"/>
<link rel="preload" href="font.woff2" as="font" crossorigin/>
```

**Required**: `as` attribute  
**Use for**: Critical CSS, fonts, hero images

#### 4. prefetch

Download resource for next page.

```html
<link rel="prefetch" href="/next-page.html"/>
```

**Priority**: Low (loads when idle)  
**Use for**: Predictable navigation flows

#### 5. fetchpriority

Control download priority.

```html
<!-- High priority for LCP image -->
<img src="hero.jpg" fetchpriority="high"/>

<!-- Low priority for below-fold -->
<img src="footer.jpg" fetchpriority="low"/>

<!-- With preload -->
<link rel="preload" href="critical.js" as="script" fetchpriority="high"/>
```

**Values**: `high`, `low`, `auto`

### Complete Performance Setup

```html
<!DOCTYPE html>
<html>
<head>
    <!-- 1. DNS Prefetch (fallback) -->
    <link rel="dns-prefetch" href="https://unpkg.com"/>
    <link rel="dns-prefetch" href="https://cdn.tailwindcss.com"/>
    <link rel="dns-prefetch" href="https://cdnjs.cloudflare.com"/>
    
    <!-- 2. Preconnect (critical domains) -->
    <link rel="preconnect" href="https://unpkg.com"/>
    <link rel="preconnect" href="https://cdn.tailwindcss.com"/>
    <link rel="preconnect" href="https://cdnjs.cloudflare.com"/>
    
    <!-- 3. Preload (critical resources) -->
    <link rel="preload" href="https://unpkg.com/htmx.org@1.9.10" as="script"/>
    <link rel="preload" href="/hero.jpg" as="image" fetchpriority="high"/>
    
    <!-- 4. Actual resources -->
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    
    <!-- 5. High priority LCP image -->
    <img src="/hero.jpg" fetchpriority="high" alt="Hero"/>
</head>
<body>
    <!-- Content -->
</body>
</html>
```

### Performance Checklist

- [ ] Preconnect to 2-4 critical CDNs
- [ ] DNS-prefetch for additional domains
- [ ] Preload critical CSS/fonts/images
- [ ] Use `fetchpriority="high"` for LCP elements
- [ ] Use `fetchpriority="low"` for below-fold content
- [ ] Limit preconnect to 6 domains maximum

***

## Core Concepts

### How HTMX Works

```
Traditional JavaScript:
User Action → JavaScript Handler → fetch() → Parse JSON → Update DOM

HTMX:
User Action → HTMX Sends Request → Server Returns HTML → HTMX Updates DOM
```

### Basic Pattern

```html
<button 
    hx-METHOD="/url"           <!-- HTTP method + endpoint -->
    hx-trigger="EVENT"         <!-- What triggers request -->
    hx-target="#element"       <!-- Where response goes -->
    hx-swap="METHOD">          <!-- How to insert response -->
    Click Me
</button>
```

### Request Lifecycle

1. **Trigger** - User clicks, page loads, timer fires
2. **Request** - HTMX sends HTTP request
3. **Response** - Server returns HTML
4. **Swap** - HTMX updates target element

***

## Complete Attribute Reference

### HTTP Requests

| Attribute | Description | Example |
|-----------|-------------|---------|
| `hx-get` | GET request | `<button hx-get="/api/users">` |
| `hx-post` | POST request | `<form hx-post="/submit">` |
| `hx-put` | PUT request | `<button hx-put="/api/user/1">` |
| `hx-delete` | DELETE request | `<button hx-delete="/api/user/1">` |
| `hx-patch` | PATCH request | `<button hx-patch="/api/user/1">` |

### Targeting

```html
<!-- Specific element -->
<button hx-get="/data" hx-target="#output">Load</button>

<!-- Self -->
<button hx-get="/data" hx-target="this">Load</button>

<!-- Closest parent with class -->
<button hx-get="/data" hx-target="closest .card">Load</button>

<!-- Next sibling -->
<button hx-get="/data" hx-target="next .result">Load</button>

<!-- Previous sibling -->
<button hx-get="/data" hx-target="previous .result">Load</button>

<!-- Find anywhere -->
<button hx-get="/data" hx-target="find .result">Load</button>
```

### Swapping

| Value | Description | Use Case |
|-------|-------------|----------|
| `innerHTML` | Replace content (default) | Most common |
| `outerHTML` | Replace entire element | Replace component |
| `beforebegin` | Insert before element | Add item above |
| `afterbegin` | Insert at start | Prepend to list |
| `beforeend` | Insert at end | Append to list |
| `afterend` | Insert after element | Add item below |
| `delete` | Remove element | Delete operation |
| `none` | Don't swap | Side effects only |

**Modifiers:**

```html
<!-- Scroll to top after swap -->
<div hx-swap="innerHTML scroll:top">...</div>

<!-- Scroll to bottom -->
<div hx-swap="innerHTML scroll:bottom">...</div>

<!-- Show swap after delay -->
<div hx-swap="innerHTML swap:200ms">...</div>

<!-- Settle animation duration -->
<div hx-swap="innerHTML settle:500ms">...</div>

<!-- Focus after swap -->
<div hx-swap="innerHTML focus-scroll:true">...</div>

<!-- Combined -->
<div hx-swap="innerHTML swap:200ms settle:500ms scroll:top">...</div>
```

### Triggers

```html
<!-- Click (default for buttons) -->
<button hx-get="/data">Click</button>

<!-- Page load -->
<div hx-get="/data" hx-trigger="load"></div>

<!-- Form submit (default for forms) -->
<form hx-post="/submit"></form>

<!-- Input change -->
<input hx-get="/search" hx-trigger="change"/>

<!-- Keyup -->
<input hx-get="/search" hx-trigger="keyup"/>

<!-- Keyup with delay (debounce) -->
<input hx-get="/search" hx-trigger="keyup changed delay:500ms"/>

<!-- Multiple triggers -->
<div hx-get="/data" hx-trigger="load, click"></div>

<!-- Polling -->
<div hx-get="/data" hx-trigger="every 5s"></div>

<!-- Polling only when visible -->
<div hx-get="/data" hx-trigger="every 5s[document.visibilityState === 'visible']"></div>

<!-- Intersection (infinite scroll) -->
<div hx-get="/page2" hx-trigger="intersect once"></div>

<!-- Revealed (lazy load) -->
<div hx-get="/content" hx-trigger="revealed"></div>

<!-- Custom event -->
<div hx-get="/data" hx-trigger="customEvent from:body"></div>

<!-- Mouse events -->
<div hx-get="/data" hx-trigger="mouseenter once">Hover</div>
```

### Including Data

```html
<!-- Include specific elements -->
<textarea id="notes"></textarea>
<button hx-post="/save" hx-include="#notes">Save</button>

<!-- Include multiple -->
<input id="name"/>
<input id="email"/>
<button hx-post="/submit" hx-include="#name, #email">Submit</button>

<!-- Include closest parent -->
<div>
    <input name="field1"/>
    <button hx-post="/submit" hx-include="closest div">Submit</button>
</div>

<!-- Add static values -->
<button hx-post="/api" hx-vals='{"type": "premium"}'>Upgrade</button>

<!-- Add dynamic values -->
<button hx-post="/api" hx-vals='js:{timestamp: Date.now()}'>Submit</button>
```

### Loading States

```html
<!-- Loading indicator -->
<button hx-get="/data" hx-indicator="#spinner">Load</button>
<div id="spinner" class="htmx-indicator">Loading...</div>

<!-- Disable during request -->
<button hx-post="/submit" hx-disabled-elt="this">Submit</button>

<!-- Disable multiple elements -->
<form hx-post="/submit" hx-disabled-elt="button, input">
    <input name="name"/>
    <button type="submit">Submit</button>
</form>

<!-- Button text change -->
<button hx-post="/submit">
    <span class="htmx-indicator:hidden">Submit</span>
    <span class="hidden htmx-request:inline">Submitting...</span>
</button>
```

### Navigation

```html
<!-- Push URL to history -->
<a hx-get="/page2" hx-push-url="true">Next Page</a>

<!-- Push custom URL -->
<a hx-get="/api/page2" hx-push-url="/page2">Next Page</a>

<!-- Replace current URL (no history entry) -->
<a hx-get="/page2" hx-replace-url="true">Next Page</a>
```

### Confirmation

```html
<button hx-delete="/item/1" hx-confirm="Are you sure?">
    Delete
</button>
```

### Headers

```html
<button hx-post="/api" 
        hx-headers='{"X-API-Key": "secret", "X-Custom": "value"}'>
    Submit
</button>
```

### Selecting Response Content

```html
<!-- Server returns full page, only use #content -->
<div hx-get="/page2" hx-select="#content"></div>

<!-- Multiple selections -->
<div hx-get="/page2" hx-select="#header, #content, #footer"></div>
```

### Out-of-Band Swaps

Update multiple elements from one response.

```html
<!-- Page HTML -->
<div id="stats">Old stats</div>
<div id="notifications">Old notifications</div>
<div id="main">Main content</div>
```

**Server response:**

```html
<!-- Main content (normal swap) -->
<div id="main">New main content</div>

<!-- Out-of-band swaps -->
<div id="stats" hx-swap-oob="true">New stats</div>
<div id="notifications" hx-swap-oob="true">New notifications</div>
```

### Request Synchronization

```html
<!-- Drop in-flight requests (for search) -->
<input hx-get="/search" hx-sync="this:drop"/>

<!-- Queue requests -->
<button hx-post="/process" hx-sync="this:queue">Process</button>

<!-- Replace in-flight request -->
<button hx-get="/data" hx-sync="this:replace">Load</button>

<!-- Abort existing requests -->
<button hx-get="/data" hx-sync="this:abort">Load</button>
```

***

## Real-World Examples

### 1. API Endpoint Tester

Complete implementation for testing REST APIs.

```html
<div class="api-tester">
    <h2>POST /scan</h2>
    
    <!-- Request Body -->
    <textarea id="scanBody" name="body" class="request-input">
{
  "org": "huggingface",
  "model_id": "bert-base-uncased",
  "include_prs": false
}
    </textarea>
    
    <!-- Test Button -->
    <button 
        hx-post="/scan"
        hx-include="#scanBody"
        hx-target="#scanResponse"
        hx-swap="innerHTML"
        hx-indicator="#scanSpinner"
        hx-disabled-elt="this"
        class="btn-primary">
        <span class="htmx-indicator:hidden">Test Endpoint</span>
        <span class="hidden htmx-request:inline">Testing...</span>
    </button>
    
    <!-- Loading Spinner -->
    <div id="scanSpinner" class="htmx-indicator">
        <div class="spinner"></div>
        <p>Running scan...</p>
    </div>
    
    <!-- Response Display -->
    <div id="scanResponse"></div>
</div>

<style>
.htmx-indicator { display: none; }
.htmx-request .htmx-indicator { display: block; }
.spinner {
    width: 40px;
    height: 40px;
    border: 4px solid #f3f3f3;
    border-top: 4px solid #3498db;
    border-radius: 50%;
    animation: spin 1s linear infinite;
}
@keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
}
</style>
```

**Server (Go):**

```go
func ScanHandler(c *fiber.Ctx) error {
    body := c.Body()
    
    // Process request
    result, err := performScan(body)
    if err != nil {
        return c.SendString(fmt.Sprintf(`
            <div class="error">
                <i class="icon-error"></i>
                <span>%s</span>
            </div>
        `, err.Error()))
    }
    
    // Return formatted response
    html := fmt.Sprintf(`
        <div class="response success">
            <label>Response (200 OK):</label>
            <pre>%s</pre>
        </div>
    `, formatJSON(result))
    
    return c.SendString(html)
}
```

### 2. Auto-Refreshing Dashboard

Live dashboard with polling and manual refresh.

```html
<div id="dashboard" 
     hx-get="/api/dashboard" 
     hx-trigger="load, every 30s"
     hx-swap="innerHTML">
    
    <!-- Stats Grid -->
    <div class="stats-grid">
        <div class="stat-card">
            <h3>Total Scans</h3>
            <p id="totalScans">0</p>
        </div>
        <div class="stat-card">
            <h3>Total Findings</h3>
            <p id="totalFindings">0</p>
        </div>
        <div class="stat-card">
            <h3>Resources Scanned</h3>
            <p id="totalResources">0</p>
        </div>
    </div>
    
    <!-- Manual Refresh -->
    <button hx-get="/api/dashboard" 
            hx-target="#dashboard"
            hx-indicator="#refreshing">
        <i class="icon-refresh"></i> Refresh
    </button>
    <span id="refreshing" class="htmx-indicator">Refreshing...</span>
    
    <!-- Recent Activity -->
    <div id="recentActivity">
        Loading activity...
    </div>
</div>
```

**Server with OOB:**

```go
func DashboardHandler(c *fiber.Ctx) error {
    data := getDashboardData()
    
    html := fmt.Sprintf(`
        <!-- Update stats (OOB) -->
        <p id="totalScans" hx-swap-oob="true">%d</p>
        <p id="totalFindings" hx-swap-oob="true">%d</p>
        <p id="totalResources" hx-swap-oob="true">%d</p>
        
        <!-- Main content -->
        <div class="stats-grid">
            <div class="stat-card">
                <h3>Total Scans</h3>
                <p>%d</p>
            </div>
            <div class="stat-card">
                <h3>Total Findings</h3>
                <p>%d</p>
            </div>
            <div class="stat-card">
                <h3>Resources Scanned</h3>
                <p>%d</p>
            </div>
        </div>
        
        <!-- Recent Activity (OOB) -->
        <div id="recentActivity" hx-swap-oob="true">
            %s
        </div>
    `, data.TotalScans, data.TotalFindings, data.TotalResources,
       data.TotalScans, data.TotalFindings, data.TotalResources,
       renderActivity(data.RecentActivity))
    
    return c.SendString(html)
}
```

### 3. Pagination Without Full Reload

Fast, SPA-like pagination.

```html
<div id="results-container">
    <!-- Results List -->
    <div class="results">
        <div class="result-item">Result 1</div>
        <div class="result-item">Result 2</div>
        <div class="result-item">Result 3</div>
    </div>
    
    <!-- Pagination -->
    <div class="pagination">
        <button 
            hx-get="/results?page=1"
            hx-target="#results-container"
            hx-swap="innerHTML"
            hx-push-url="/results?page=1"
            disabled>
            Previous
        </button>
        
        <button class="active">1</button>
        
        <button 
            hx-get="/results?page=2"
            hx-target="#results-container"
            hx-swap="innerHTML"
            hx-push-url="/results?page=2">
            2
        </button>
        
        <button 
            hx-get="/results?page=3"
            hx-target="#results-container"
            hx-swap="innerHTML"
            hx-push-url="/results?page=3">
            3
        </button>
        
        <button 
            hx-get="/results?page=2"
            hx-target="#results-container"
            hx-swap="innerHTML"
            hx-push-url="/results?page=2">
            Next
        </button>
    </div>
</div>
```

**Server:**

```go
func ResultsHandler(c *fiber.Ctx) error {
    page := c.QueryInt("page", 1)
    results := getResults(page)
    totalPages := getTotalPages()
    
    return Render(c, ResultsList(results, page, totalPages))
}
```

### 4. Form with Validation & Redirect

Complete form handling with server-side validation.

```html
<form hx-post="/scan" 
      hx-indicator="#loading"
      hx-disabled-elt="button[type='submit']">
    
    <!-- Resource Type -->
    <div class="radio-group">
        <label>
            <input type="radio" name="type" value="model" checked/>
            <span>AI Model</span>
        </label>
        <label>
            <input type="radio" name="type" value="dataset"/>
            <span>Dataset</span>
        </label>
        <label>
            <input type="radio" name="type" value="space"/>
            <span>Space</span>
        </label>
    </div>
    
    <!-- Text Inputs -->
    <input 
        name="org" 
        placeholder="Organization (e.g., huggingface)" 
        required
        minlength="2"/>
    
    <input 
        name="resourceId" 
        placeholder="Resource ID (e.g., bert-base-uncased)" 
        required/>
    
    <!-- Checkboxes -->
    <label>
        <input type="checkbox" name="includePRs" value="true"/>
        Include Pull Requests
    </label>
    
    <label>
        <input type="checkbox" name="includeDiscussions" value="true"/>
        Include Discussions
    </label>
    
    <!-- Submit -->
    <button type="submit">
        <span class="htmx-indicator:hidden">Start Scan</span>
        <span class="hidden htmx-request:inline">Starting...</span>
    </button>
    
    <!-- Loading -->
    <div id="loading" class="htmx-indicator">
        <div class="spinner"></div>
        Processing your scan...
    </div>
</form>

<!-- Errors -->
<div id="errors"></div>
```

**Server:**

```go
func ScanFormHandler(c *fiber.Ctx) error {
    var req struct {
        Type               string `form:"type"`
        Org                string `form:"org"`
        ResourceID         string `form:"resourceId"`
        IncludePRs         string `form:"includePRs"`
        IncludeDiscussions string `form:"includeDiscussions"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        c.Set("HX-Retarget", "#errors")
        return c.SendString(`
            <div class="error-message">
                Invalid form data
            </div>
        `)
    }
    
    // Validate
    if len(req.Org) < 2 {
        c.Set("HX-Retarget", "#errors")
        return c.SendString(`
            <div class="error-message">
                Organization must be at least 2 characters
            </div>
        `)
    }
    
    if req.ResourceID == "" {
        c.Set("HX-Retarget", "#errors")
        return c.SendString(`
            <div class="error-message">
                Resource ID is required
            </div>
        `)
    }
    
    // Process
    result := performScan(ScanRequest{
        Type:               req.Type,
        Org:                req.Org,
        ResourceID:         req.ResourceID,
        IncludePRs:         req.IncludePRs == "true",
        IncludeDiscussions: req.IncludeDiscussions == "true",
    })
    
    // Redirect on success
    c.Set("HX-Redirect", fmt.Sprintf("/results/%s", result.RequestID))
    return c.SendStatus(200)
}
```

### 5. Live Search with Debounce

Search as you type with delay.

```html
<input 
    type="search"
    name="q"
    placeholder="Search..."
    hx-get="/search"
    hx-trigger="keyup changed delay:500ms"
    hx-target="#search-results"
    hx-indicator="#search-loading"
    hx-sync="this:drop"/>

<div id="search-loading" class="htmx-indicator">
    Searching...
</div>

<div id="search-results">
    <!-- Results appear here -->
</div>
```

### 6. Infinite Scroll

Load more content as user scrolls.

```html
<div id="content">
    <div class="item">Item 1</div>
    <div class="item">Item 2</div>
    <div class="item">Item 3</div>
</div>

<!-- Load more trigger -->
<div hx-get="/items?page=2" 
     hx-trigger="intersect once"
     hx-target="#content"
     hx-swap="beforeend">
    <div class="loading">Loading more...</div>
</div>
```

**Server returns:**

```html
<div class="item">Item 4</div>
<div class="item">Item 5</div>
<div class="item">Item 6</div>

<!-- Next page trigger -->
<div hx-get="/items?page=3" 
     hx-trigger="intersect once"
     hx-target="this"
     hx-swap="outerHTML">
    <div class="loading">Loading more...</div>
</div>
```

### 7. Shopping Cart

Real-time cart updates.

```html
<!-- Product Card -->
<div class="product">
    <h3>Product Name</h3>
    <p>$29.99</p>
    <button 
        hx-post="/cart/add"
        hx-vals='{"product_id": 123, "quantity": 1}'
        hx-target="#cart-count"
        hx-swap="innerHTML"
        hx-on::after-request="showToast('Added to cart')">
        Add to Cart
    </button>
</div>

<!-- Cart Badge -->
<div id="cart-badge">
    <i class="icon-cart"></i>
    <span id="cart-count">0</span>
</div>

<!-- Cart Drawer (updates on trigger) -->
<div id="cart-drawer" 
     hx-get="/cart" 
     hx-trigger="cartUpdated from:body">
    <!-- Cart contents -->
</div>
```

**Server:**

```go
func AddToCartHandler(c *fiber.Ctx) error {
    var req struct {
        ProductID int `json:"product_id"`
        Quantity  int `json:"quantity"`
    }
    c.BodyParser(&req)
    
    cart := addToCart(req.ProductID, req.Quantity)
    
    // Trigger cart drawer update
    c.Set("HX-Trigger", "cartUpdated")
    
    return c.SendString(fmt.Sprintf("%d", cart.ItemCount))
}
```

### 8. Todo List

CRUD operations with animations.

```html
<!-- Add Form -->
<form hx-post="/todos" 
      hx-target="#todo-list"
      hx-swap="afterbegin"
      hx-on::after-request="this.reset()">
    <input name="text" placeholder="New todo" required/>
    <button type="submit">Add</button>
</form>

<!-- Todo List -->
<div id="todo-list">
    <div class="todo" id="todo-1">
        <input 
            type="checkbox" 
            hx-patch="/todos/1/toggle"
            hx-target="closest .todo"
            hx-swap="outerHTML"/>
        <span>Buy milk</span>
        <button 
            hx-delete="/todos/1"
            hx-target="closest .todo"
            hx-swap="outerHTML swap:1s"
            hx-confirm="Delete this todo?">
            Delete
        </button>
    </div>
</div>

<style>
.todo.htmx-swapping {
    opacity: 0;
    transition: opacity 1s;
}
</style>
```

### 9. Modal Dialog

Load modal content dynamically.

```html
<!-- Trigger -->
<button hx-get="/modal/user/123" 
        hx-target="#modal-container"
        hx-swap="innerHTML">
    View Profile
</button>

<!-- Modal Container -->
<div id="modal-container"></div>
```

**Server returns:**

```html
<div class="modal-backdrop" onclick="this.parentElement.innerHTML=''">
    <div class="modal-content" onclick="event.stopPropagation()">
        <div class="modal-header">
            <h2>User Profile</h2>
            <button onclick="document.getElementById('modal-container').innerHTML=''">
                ×
            </button>
        </div>
        <div class="modal-body">
            <p>Name: John Doe</p>
            <p>Email: john@example.com</p>
        </div>
    </div>
</div>
```

### 10. Real-Time Notifications

Server-Sent Events for live updates.

```html
<div hx-sse="connect:/notifications/stream">
    <div id="notifications" 
         hx-sse="swap:message"
         hx-swap="afterbegin">
        <!-- Notifications appear here -->
    </div>
</div>
```

**Server (Go):**

```go
func NotificationsStream(c *fiber.Ctx) error {
    c.Set("Content-Type", "text/event-stream")
    c.Set("Cache-Control", "no-cache")
    c.Set("Connection", "keep-alive")
    
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            notification := getLatestNotification()
            html := fmt.Sprintf(`
                <div class="notification">
                    <i class="icon-bell"></i>
                    <span>%s</span>
                </div>
            `, notification.Message)
            
            fmt.Fprintf(c, "event: message\n")
            fmt.Fprintf(c, "data: %s\n\n", html)
            c.Context().Response.Flush()
        case <-c.Context().Done():
            return nil
        }
    }
}
```

***

## Server-Side Integration (Go Fiber)

### Basic Setup

```go
package main

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()
    
    // HTMX routes
    app.Get("/", IndexHandler)
    app.Get("/api/data", DataHandler)
    app.Post("/api/submit", SubmitHandler)
    
    app.Listen(":3000")
}
```

### Rendering HTML

```go
// Simple string
func Handler(c *fiber.Ctx) error {
    return c.SendString(`<div>Hello World</div>`)
}

// With Templ
func Handler(c *fiber.Ctx) error {
    data := getData()
    return Render(c, Component(data))
}

func Render(c *fiber.Ctx, component templ.Component) error {
    c.Set("Content-Type", "text/html")
    return component.Render(c.Context(), c.Response().BodyWriter())
}
```

### Form Handling

```go
func FormHandler(c *fiber.Ctx) error {
    var req struct {
        Name  string `form:"name"`
        Email string `form:"email"`
        Active string `form:"active"` // checkbox: "true" or empty
    }
    
    if err := c.BodyParser(&req); err != nil {
        return sendError(c, "Invalid form data")
    }
    
    // Validate
    if req.Name == "" {
        return sendError(c, "Name is required")
    }
    
    // Process
    user := createUser(req.Name, req.Email, req.Active == "true")
    
    return c.SendString(fmt.Sprintf(`
        <div class="success">
            User %s created successfully!
        </div>
    `, user.Name))
}
```

### Error Handling

```go
func sendError(c *fiber.Ctx, message string) error {
    c.Set("HX-Retarget", "#errors")
    c.Set("HX-Reswap", "innerHTML")
    return c.SendString(fmt.Sprintf(`
        <div class="error">
            <i class="icon-error"></i>
            <span>%s</span>
        </div>
    `, message))
}
```

### Redirects

```go
func Handler(c *fiber.Ctx) error {
    // Process request
    result := process()
    
    // Redirect via HTMX
    c.Set("HX-Redirect", fmt.Sprintf("/results/%s", result.ID))
    return c.SendStatus(200)
}
```

### Out-of-Band Swaps

```go
func Handler(c *fiber.Ctx) error {
    stats := getStats()
    content := getContent()
    
    html := fmt.Sprintf(`
        <!-- Main content (normal swap) -->
        <div id="main-content">
            %s
        </div>
        
        <!-- OOB: Update stats -->
        <div id="stats" hx-swap-oob="true">
            <p>Total: %d</p>
        </div>
        
        <!-- OOB: Update notification badge -->
        <span id="notification-count" hx-swap-oob="true">
            %d
        </span>
    `, content, stats.Total, stats.Unread)
    
    return c.SendString(html)
}
```

### Triggering Client Events

```go
func Handler(c *fiber.Ctx) error {
    // Process
    deleteItem()
    
    // Trigger event on client
    c.Set("HX-Trigger", "itemDeleted")
    
    // With event data
    c.Set("HX-Trigger", `{"itemDeleted": {"id": 123, "name": "Item"}}`)
    
    return c.SendString(`<div>Deleted</div>`)
}
```

**Client listens:**

```html
<div id="list" 
     hx-get="/items" 
     hx-trigger="itemDeleted from:body">
    <!-- List reloads when itemDeleted fires -->
</div>
```

### Response Headers Reference

```go
// Redirect
c.Set("HX-Redirect", "/new-page")

// Refresh page
c.Set("HX-Refresh", "true")

// Push URL (update history)
c.Set("HX-Push-Url", "/new-url")

// Replace URL (no history)
c.Set("HX-Replace-Url", "/current-url")

// Change target
c.Set("HX-Retarget", "#different-element")

// Change swap method
c.Set("HX-Reswap", "outerHTML")

// Trigger client event
c.Set("HX-Trigger", "myEvent")

// Trigger after settle
c.Set("HX-Trigger-After-Settle", "myEvent")

// Trigger after swap
c.Set("HX-Trigger-After-Swap", "myEvent")
```

### Detecting HTMX Requests

```go
func Handler(c *fiber.Ctx) error {
    if c.Get("HX-Request") == "true" {
        // Return partial HTML for HTMX
        return c.SendString(`<div>Partial content</div>`)
    }
    
    // Return full page for normal request
    return Render(c, FullPage())
}
```

### JSON Responses (with json-enc extension)

```go
func Handler(c *fiber.Ctx) error {
    var data map[string]interface{}
    
    if err := c.BodyParser(&data); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid JSON",
        })
    }
    
    result := process(data)
    
    return c.JSON(fiber.Map{
        "status": "success",
        "data": result,
    })
}
```

***

## Advanced Patterns

### Request Debouncing

```html
<!-- Wait 500ms after typing stops -->
<input 
    hx-get="/search"
    hx-trigger="keyup changed delay:500ms"
    hx-sync="this:drop"/>
```

### Request Cancellation

```html
<!-- Cancel in-flight requests when new one starts -->
<input 
    hx-get="/search"
    hx-trigger="keyup"
    hx-sync="this:drop"/>
```

### Conditional Requests

```html
<!-- Only send if input has 3+ characters -->
<input 
    hx-get="/search"
    hx-trigger="keyup[target.value.length > 3]"/>

<!-- Only send if shift key held -->
<button hx-get="/data" hx-trigger="click[shiftKey]">
    Shift+Click
</button>

<!-- Only poll when tab visible -->
<div 
    hx-get="/data" 
    hx-trigger="every 5s[document.visibilityState === 'visible']">
</div>
```

### Optimistic Updates

```html
<button 
    hx-post="/like"
    hx-swap="outerHTML"
    onclick="this.textContent='Liked!'; this.disabled=true;">
    Like
</button>
```

### Progressive Enhancement

```html
<!-- Works with or without JavaScript -->
<form action="/submit" method="POST" 
      hx-post="/submit" 
      hx-target="#result">
    <!-- If HTMX loads: AJAX submission -->
    <!-- If HTMX fails: Normal form submission -->
</form>
```

### Polling with Backoff

```html
<div 
    hx-get="/status"
    hx-trigger="load, every 2s">
    Checking status...
</div>

<script>
// Stop polling after 10 attempts
let attempts = 0;
htmx.on('htmx:afterRequest', function(evt) {
    attempts++;
    if (attempts >= 10) {
        evt.target.removeAttribute('hx-trigger');
    }
});
</script>
```

### File Upload

```html
<form 
    hx-post="/upload"
    hx-encoding="multipart/form-data"
    hx-target="#result">
    
    <input type="file" name="file" accept="image/*"/>
    <button type="submit">Upload</button>
    
    <progress id="progress" value="0" max="100"></progress>
</form>

<div id="result"></div>

<script>
htmx.on('htmx:xhr:progress', function(evt) {
    const progress = document.getElementById('progress');
    progress.value = (evt.detail.loaded / evt.detail.total) * 100;
});
</script>
```

### WebSocket

```html
<div hx-ws="connect:/ws">
    <form hx-ws="send">
        <input name="message"/>
        <button type="submit">Send</button>
    </form>
    
    <div id="messages"></div>
</div>
```

### Custom Events

```javascript
// Trigger from JavaScript
htmx.trigger('#element', 'customEvent', {detail: {data: 'value'}});
```

```html
<!-- Listen in HTMX -->
<div hx-get="/data" hx-trigger="customEvent"></div>
```

***

## CSS & Styling

### HTMX Classes

```css
/* Request in progress */
.htmx-request {
    opacity: 0.7;
    pointer-events: none;
}

/* Swapping (removing old content) */
.htmx-swapping {
    opacity: 0;
    transition: opacity 1s;
}

/* Settling (showing new content) */
.htmx-settling {
    opacity: 0;
}

.htmx-settling:not(.htmx-swapping) {
    opacity: 1;
    transition: opacity 1s;
}

/* Added content (for beforeend/afterbegin) */
.htmx-added {
    animation: fadeIn 0.5s;
}

@keyframes fadeIn {
    from { opacity: 0; transform: translateY(-10px); }
    to { opacity: 1; transform: translateY(0); }
}

/* Loading indicators */
.htmx-indicator {
    display: none;
}

.htmx-request .htmx-indicator,
.htmx-request.htmx-indicator {
    display: block;
}
```

### Loading Spinner

```html
<style>
.spinner {
    display: inline-block;
    width: 40px;
    height: 40px;
    border: 4px solid rgba(0,0,0,.1);
    border-left-color: #000;
    border-radius: 50%;
    animation: spin 1s linear infinite;
}

@keyframes spin {
    to { transform: rotate(360deg); }
}
</style>

<button hx-get="/data" hx-indicator=".spinner">
    Load
    <div class="spinner htmx-indicator"></div>
</button>
```

### Skeleton Loader

```html
<style>
.skeleton {
    background: linear-gradient(
        90deg,
        #f0f0f0 25%,
        #e0e0e0 50%,
        #f0f0f0 75%
    );
    background-size: 200% 100%;
    animation: loading 1.5s infinite;
    border-radius: 4px;
}

@keyframes loading {
    0% { background-position: 200% 0; }
    100% { background-position: -200% 0; }
}
</style>

<div hx-get="/content" hx-trigger="load">
    <div class="skeleton h-8 w-full mb-4"></div>
    <div class="skeleton h-8 w-3/4 mb-4"></div>
    <div class="skeleton h-8 w-1/2"></div>
</div>
```

### Progress Bar

```html
<style>
.progress {
    position: fixed;
    top: 0;
    left: 0;
    height: 3px;
    background: #3498db;
    width: 0;
    transition: width 0.3s;
    z-index: 9999;
}

.htmx-request .progress {
    width: 70%;
}
</style>

<div class="progress htmx-indicator"></div>
```

***

## Troubleshooting

### Common Issues

#### Request Not Firing

**Symptoms**: Button does nothing when clicked

**Check**:
1. Is HTMX loaded? Console: `typeof htmx`
2. Correct trigger? (default: `click` for buttons, `submit` for forms)
3. Is element disabled?
4. Any JavaScript errors?

**Debug**:
```html
<button hx-get="/data"
        hx-on::before-request="console.log('Request starting')"
        hx-on::after-request="console.log('Request complete')">
    Test
</button>
```

#### Response Not Showing

**Symptoms**: Request succeeds but page doesn't update

**Check**:
1. Does target element exist? `document.getElementById('target')`
2. Is `hx-target` correct?
3. Is server returning HTML (not JSON)?
4. Check swap method

**Debug**:
```html
<div hx-get="/data"
     hx-on::after-request="console.log(event.detail.xhr.response)">
</div>
```

#### Form Not Submitting

**Symptoms**: Form submit button doesn't work

**Check**:
1. Is `hx-post` on `<form>`, not button?
2. Do inputs have `name` attributes?
3. No conflicting `onsubmit` handlers?

**Debug**:
```html
<form hx-post="/submit"
      hx-on::before-request="console.log('Form data:', new FormData(event.target))">
</form>
```

#### OOB Swaps Not Working

**Symptoms**: Only main content updates, OOB targets don't

**Check**:
1. Target elements exist with correct IDs
2. Response elements have `hx-swap-oob="true"`
3. IDs match exactly

**Example**:
```html
<!-- Page -->
<div id="stats">Old</div>

<!-- Server must return -->
<div id="stats" hx-swap-oob="true">New</div>
```

#### Polling Stopped

**Symptoms**: Auto-refresh stops working

**Check**:
1. Is tab visible? Use visibility condition
2. Any errors in console?
3. Was element removed from DOM?

**Fix**:
```html
<div hx-get="/data" 
     hx-trigger="every 5s[document.visibilityState === 'visible']">
</div>
```

### Enable Debug Logging

```html
<script>
// Log all HTMX events
htmx.logAll();
</script>
```

### Event Listeners

```javascript
// Before request
htmx.on('htmx:beforeRequest', function(evt) {
    console.log('Sending request to:', evt.detail.requestConfig.path);
});

// After request
htmx.on('htmx:afterRequest', function(evt) {
    console.log('Response status:', evt.detail.xhr.status);
    console.log('Response body:', evt.detail.xhr.response);
});

// On error
htmx.on('htmx:responseError', function(evt) {
    console.error('Request failed:', evt.detail.error);
});

// Before swap
htmx.on('htmx:beforeSwap', function(evt) {
    console.log('About to swap into:', evt.detail.target);
});

// After swap
htmx.on('htmx:afterSwap', function(evt) {
    console.log('Swap complete');
});
```

### Network Inspection

Open DevTools → Network tab

**Request headers to look for**:
- `HX-Request: true`
- `HX-Target: element-id`
- `HX-Trigger: trigger-name`

**Response headers to look for**:
- `HX-Redirect`
- `HX-Refresh`
- `HX-Trigger`
- `HX-Retarget`

***

## Best Practices

### 1. Keep Responses Small

```go
// Good: Minimal HTML
func Handler(c *fiber.Ctx) error {
    return c.SendString(`<div>Updated</div>`)
}

// Bad: Full page
func Handler(c *fiber.Ctx) error {
    return Render(c, FullPageWithLayout())
}
```

### 2. Use Debouncing for Search

```html
<!-- Good: Wait for user to stop typing -->
<input hx-get="/search" hx-trigger="keyup changed delay:500ms"/>

<!-- Bad: Request on every keystroke -->
<input hx-get="/search" hx-trigger="keyup"/>
```

### 3. Cancel In-Flight Requests

```html
<input hx-get="/search" hx-sync="this:drop"/>
```

### 4. Reasonable Polling Intervals

```html
<!-- Good: 30 seconds -->
<div hx-get="/data" hx-trigger="every 30s"></div>

<!-- Bad: Every second (wastes resources) -->
<div hx-get="/data" hx-trigger="every 1s"></div>
```

### 5. Disable Buttons During Requests

```html
<button hx-post="/submit" hx-disabled-elt="this">
    Submit
</button>
```

### 6. Show Loading Indicators

```html
<button hx-get="/data" hx-indicator="#loading">
    Load
</button>
<div id="loading" class="htmx-indicator">Loading...</div>
```

### 7. Validate Server-Side

```go
func Handler(c *fiber.Ctx) error {
    var req Request
    if err := c.BodyParser(&req); err != nil {
        return sendError(c, "Invalid input")
    }
    
    // Always validate
    if req.Email == "" || !isValidEmail(req.Email) {
        return sendError(c, "Invalid email")
    }
    
    // Sanitize
    req.Email = sanitize(req.Email)
    
    // Process
    return c.SendString(`<div>Success</div>`)
}
```

### 8. Use Progressive Enhancement

```html
<!-- Still works without HTMX -->
<form action="/submit" method="POST" hx-post="/submit">
    <button type="submit">Submit</button>
</form>
```

### 9. Batch Updates with OOB

Instead of multiple requests, use one with OOB swaps:

```html
<!-- One request updates multiple elements -->
<div id="stats" hx-swap-oob="true">Updated stats</div>
<div id="notifications" hx-swap-oob="true">Updated notifications</div>
<div id="main">Main content</div>
```

### 10. Cache Appropriately

```go
func Handler(c *fiber.Ctx) error {
    // Static content
    c.Set("Cache-Control", "public, max-age=3600")
    
    // Dynamic content
    c.Set("Cache-Control", "no-cache")
    
    return c.SendString(`<div>Content</div>`)
}
```

***

## Migration from JavaScript

### Fetch API → HTMX

**Before (JavaScript)**:
```javascript
async function loadData() {
    const response = await fetch('/api/data');
    const html = await response.text();
    document.getElementById('result').innerHTML = html;
}
button.addEventListener('click', loadData);
```

**After (HTMX)**:
```html
<button hx-get="/api/data" hx-target="#result">Load</button>
<div id="result"></div>
```

**Eliminated**: 6 lines of JavaScript

### Event Listeners → hx-trigger

**Before**:
```javascript
input.addEventListener('keyup', function() {
    if (input.value.length > 3) {
        search(input.value);
    }
});
```

**After**:
```html
<input hx-get="/search" 
       hx-trigger="keyup[target.value.length > 3]"/>
```

### Polling → hx-trigger every

**Before**:
```javascript
setInterval(() => {
    fetch('/api/data')
        .then(r => r.text())
        .then(html => element.innerHTML = html);
}, 5000);
```

**After**:
```html
<div hx-get="/api/data" hx-trigger="every 5s"></div>
```

### Form Submission → hx-post

**Before**:
```javascript
form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const formData = new FormData(form);
    const response = await fetch('/submit', {
        method: 'POST',
        body: formData
    });
    const result = await response.text();
    resultDiv.innerHTML = result;
});
```

**After**:
```html
<form hx-post="/submit" hx-target="#result"></form>
<div id="result"></div>
```

**Eliminated**: 10 lines of JavaScript

### Complete Migration Example

**Before (Traditional SPA)**:
```html
<div id="app"></div>

<script>
let currentPage = 1;

async function loadPage(page) {
    showLoading();
    try {
        const response = await fetch(`/api/page/${page}`);
        const data = await response.json();
        renderPage(data);
        history.pushState({page}, '', `/page/${page}`);
    } catch (error) {
        showError(error.message);
    } finally {
        hideLoading();
    }
}

function renderPage(data) {
    const html = `
        <h1>${data.title}</h1>
        <p>${data.content}</p>
        <button onclick="loadPage(${data.next})">Next</button>
    `;
    document.getElementById('app').innerHTML = html;
}

function showLoading() {
    document.getElementById('loading').style.display = 'block';
}

function hideLoading() {
    document.getElementById('loading').style.display = 'none';
}

function showError(message) {
    document.getElementById('error').textContent = message;
}

window.addEventListener('load', () => loadPage(1));
</script>
```

**After (HTMX)**:
```html
<div id="app" 
     hx-get="/page/1" 
     hx-trigger="load"
     hx-push-url="true"
     hx-indicator="#loading">
</div>

<div id="loading" class="htmx-indicator">Loading...</div>
<div id="error"></div>
```

**Server returns:**
```html
<h1>Page Title</h1>
<p>Page content</p>
<button hx-get="/page/2" 
        hx-target="#app"
        hx-push-url="true">
    Next
</button>
```

**Result**: ~40 lines of JavaScript → 0 lines

***

## Quick Reference Card

```
╔═══════════════════════════════════════════════════════════════╗
║                    HTMX QUICK REFERENCE                       ║
╠═══════════════════════════════════════════════════════════════╣
║ HTTP Requests                                                 ║
║   hx-get="/url"              GET request                      ║
║   hx-post="/url"             POST request                     ║
║   hx-put="/url"              PUT request                      ║
║   hx-delete="/url"           DELETE request                   ║
║   hx-patch="/url"            PATCH request                    ║
╠═══════════════════════════════════════════════════════════════╣
║ Targeting                                                     ║
║   hx-target="#id"            Specific element                 ║
║   hx-target="this"           Self                             ║
║   hx-target="closest .cls"   Nearest parent                   ║
║   hx-target="body"           Full page                        ║
╠═══════════════════════════════════════════════════════════════╣
║ Swapping                                                      ║
║   hx-swap="innerHTML"        Replace content (default)        ║
║   hx-swap="outerHTML"        Replace element                  ║
║   hx-swap="beforeend"        Append                           ║
║   hx-swap="afterbegin"       Prepend                          ║
║   hx-swap="none"             No swap                          ║
╠═══════════════════════════════════════════════════════════════╣
║ Triggers                                                      ║
║   hx-trigger="click"         On click (default)               ║
║   hx-trigger="load"          On page load                     ║
║   hx-trigger="every 5s"      Polling                          ║
║   hx-trigger="keyup delay:500ms"  Debounced                   ║
║   hx-trigger="intersect"     Scroll into view                 ║
╠═══════════════════════════════════════════════════════════════╣
║ Loading                                                       ║
║   hx-indicator="#id"         Show loading indicator           ║
║   hx-disabled-elt="this"     Disable during request           ║
╠═══════════════════════════════════════════════════════════════╣
║ Navigation                                                    ║
║   hx-push-url="true"         Update browser URL               ║
║   hx-swap-oob="true"         Out-of-band swap                 ║
╠═══════════════════════════════════════════════════════════════╣
║ Extras                                                        ║
║   hx-confirm="Sure?"         Confirmation dialog              ║
║   hx-include="#id"           Include other elements           ║
║   hx-vals='{"k":"v"}'        Add values to request            ║
║   hx-sync="this:drop"        Cancel in-flight requests        ║
╚═══════════════════════════════════════════════════════════════╝
```

***

## Performance Optimization Summary

### Resource Hints Checklist

```html
<!-- 1. Preconnect to critical CDNs (2-4 max) -->
<link rel="preconnect" href="https://unpkg.com"/>
<link rel="preconnect" href="https://cdn.tailwindcss.com"/>

<!-- 2. DNS prefetch as fallback -->
<link rel="dns-prefetch" href="https://unpkg.com"/>
<link rel="dns-prefetch" href="https://cdn.tailwindcss.com"/>

<!-- 3. Preload critical resources -->
<link rel="preload" href="critical.css" as="style"/>
<link rel="preload" href="hero.jpg" as="image"/>

<!-- 4. Prefetch next page -->
<link rel="prefetch" href="/next-page.html"/>

<!-- 5. Set fetch priorities -->
<img src="hero.jpg" fetchpriority="high"/>
<img src="footer.jpg" fetchpriority="low"/>
```

### HTMX Performance Tips

- ✅ Use reasonable polling intervals (30s+)
- ✅ Debounce search inputs (500ms)
- ✅ Cancel in-flight requests (`hx-sync="this:drop"`)
- ✅ Lazy load below-fold content (`hx-trigger="revealed"`)
- ✅ Disable buttons during requests
- ✅ Return minimal HTML from server
- ✅ Use HTTP caching
- ✅ Batch updates with OOB swaps

***

## Resources

- **Official Docs**: https://htmx.org/docs/
- **Examples**: https://htmx.org/examples/
- **Discord**: https://htmx.org/discord
- **GitHub**: https://github.com/bigskysoftware/htmx

## License

HTMX is licensed under the BSD 2-Clause License.

***

**Last Updated**: October 2025  
**HTMX Version**: 1.9.10  
**Author**: Comprehensive documentation based on real-world implementation

[1](https://htmx.org/docs/)
[2](https://www.reddit.com/r/htmx/comments/1c304wh/doceaser_interactive_documentation_with_markdown/)
[3](https://www.markdownguide.org/basic-syntax/)
[4](https://docs.github.com/github/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax)
[5](https://kabartolo.github.io/chicago-docs-demo/docs/mdx-guide/writing/)
[6](https://experienceleague.adobe.com/en/docs/contributor/contributor-guide/writing-essentials/markdown)
[7](https://docs.astro.build/en/guides/markdown-content/)
[8](https://developer.mozilla.org/en-US/docs/MDN/Writing_guidelines/Howto/Markdown_in_MDN)
[9](https://www.sphinx-doc.org/en/master/usage/markdown.html)
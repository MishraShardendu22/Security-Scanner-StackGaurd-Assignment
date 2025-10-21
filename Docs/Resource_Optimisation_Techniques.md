# The 5 Resource Optimization Techniques

## 1. dns-prefetch

**What it does:** DNS lookup only[1] [2]

```html
<link rel="dns-prefetch" href="https://example.com"/>
```

**Use case:** Many third-party domains, fallback for older browsers

**Time saved:** 20-120ms

**Cost:** Very low

## 2. preconnect

**What it does:** DNS + TCP handshake + TLS negotiation[3] [1]

```html
<link rel="preconnect" href="https://unpkg.com"/>
```

**Use case:** Critical domains you WILL use immediately

**Time saved:** 100-500ms

**Cost:** High (resource intensive)

**Limit:** 4-6 domains maximum[4]

## 3. preload

**What it does:** Download specific resource for **current page**[5] [2]

```html
<link rel="preload" href="main.js" as="script"/>
<link rel="preload" href="font.woff2" as="font" crossorigin/>
<link rel="preload" href="hero.jpg" as="image"/>
```

**Use case:** Critical resources needed for initial render

**Required attribute:** `as` (script, style, font, image, etc.)

**Time saved:** Full download time

## 4. prefetch

**What it does:** Download resource for **future navigation** [next page](5) [1]

```html
<link rel="prefetch" href="/page2.html"/>
<link rel="prefetch" href="next-image.jpg"/>
```

**Use case:** Resources likely needed on next page

**Priority:** Low (browsers load it when idle)

**Best for:** Multi-page flows where next step is predictable

## 5. fetchpriority

**What it does:** Controls download priority relative to other resources[6] [7]

```html
<!-- High priority for LCP image -->
<img src="hero.jpg" fetchpriority="high"/>

<!-- Low priority for below-fold image -->
<img src="footer-logo.jpg" fetchpriority="low"/>

<!-- With preload -->
<link rel="preload" href="critical.css" as="style" fetchpriority="high"/>
```

**Values:**

- `high` - Load before other similar resources
- `low` - Load after other similar resources  
- `auto` - Browser decides (default)

**Use case:** Fine-tune loading order of similar resources[8] [9]

## Comparison Table

| Technique | Target | Timing | Priority | Browser Support |
|-----------|--------|--------|----------|-----------------|
| **dns-prefetch** | Domain | Before request | Low cost | Excellent[2] |
| **preconnect** | Domain | Before request | High cost | Excellent[3] |
| **preload** | Specific file | Current page | High | Excellent[5] |
| **prefetch** | Specific file | Future page | Low | Good[1] |
| **fetchpriority** | Any resource | Priority hint | Varies | Modern browsers[6] |

## Real-World Example

```html
<head>
    <!-- 1. DNS-PREFETCH: Fallback for older browsers -->
    <link rel="dns-prefetch" href="https://unpkg.com"/>
    
    <!-- 2. PRECONNECT: Critical CDN (full connection) -->
    <link rel="preconnect" href="https://unpkg.com"/>
    
    <!-- 3. PRELOAD: Critical resources for current page -->
    <link rel="preload" href="https://unpkg.com/htmx.org@1.9.10" as="script"/>
    <link rel="preload" href="/hero.jpg" as="image" fetchpriority="high"/>
    
    <!-- 4. PREFETCH: Next page resources -->
    <link rel="prefetch" href="/dashboard.html"/>
    
    <!-- 5. FETCHPRIORITY: LCP image needs high priority -->
    <img src="/hero.jpg" fetchpriority="high" alt="Hero"/>
    
    <!-- Lower priority for footer -->
    <img src="/footer.jpg" fetchpriority="low" alt="Footer"/>
</head>
```

## When to Use Each

### dns-prefetch

✅ Many external domains  
✅ Fallback with preconnect  
✅ Conditional resources  
❌ Same-origin (already connected)

### preconnect

✅ 2-4 critical external domains  
✅ CDNs you know you'll use  
✅ API endpoints called immediately  
❌ More than 6 domains [wastes resources](4)

### preload

✅ Critical CSS/JS  
✅ Web fonts  
✅ Hero images [LCP elements](7)
❌ Resources not used on current page

### prefetch

✅ Next page in a flow  
✅ Predictable navigation  
✅ Idle time optimization  
❌ Critical resources (use preload instead)

### fetchpriority

✅ LCP images (set to `high`) [9]
✅ Below-fold images (set to `low`)  
✅ Non-critical scripts (set to `low`)  
❌ Overuse (defeats browser heuristics)

### BONUS: prerender (Deprecated)

There was a 6th technique called **prerender** that would render entire next pages in background, but it's been deprecated and replaced by Speculation Rules API.[3] [1]

```html
<!-- OLD (deprecated) -->
<link rel="prerender" href="/next-page.html"/>
```

### Quick Decision Tree

```md
Need to optimize a resource?
│
├─ Is it on another domain?
│  ├─ Yes → Will you definitely use it?
│  │  ├─ Yes → preconnect
│  │  └─ Maybe → dns-prefetch
│  │
│  └─ No → Continue below
│
├─ Is it for the current page?
│  ├─ Yes → Is it critical?
│     ├─ Yes → preload + fetchpriority="high"
│     └─ No → fetchpriority="low"
│
└─ Is it for the next page?
   └─ Yes → prefetch
```

The **5 main techniques** are: dns-prefetch, preconnect, preload, prefetch, and fetchpriority.[2] [6] [5]

[1](https://www.w3.org/TR/2023/DISC-resource-hints-20230314/)
[2](https://developer.mozilla.org/en-US/docs/Web/Performance/Guides/Speculative_loading)
[3](https://nitropack.io/blog/post/resource-hints-performance-optimization)
[4](https://stackoverflow.com/questions/47273743/preconnect-vs-dns-prefetch-resource-hints)
[5](https://web.dev/learn/performance/resource-hints)
[6](https://web.dev/articles/fetch-priority)
[7](https://developer.mozilla.org/en-US/docs/Web/API/HTMLImageElement/fetchPriority)
[8](https://www.debugbear.com/blog/fetchpriority-attribute)
[9](https://rabbitloader.com/articles/fetchpriority/)
[10](https://www.debugbear.com/blog/resource-hints-rel-preload-prefetch-preconnect)
[11](https://www.keycdn.com/blog/resource-hints)
[12](https://www.oxyplug.com/optimization/how-resource-hints-can-help-website-performance/)
[13](https://developer.mozilla.org/en-US/docs/Web/API/HTMLLinkElement/fetchPriority)
[14](https://almanac.httparchive.org/en/2020/resource-hints)
[15](https://nitropack.io/blog/post/priority-hints)
[16](https://leapcell.io/blog/optimizing-resource-loading-with-fetchpriority)
[17](https://www.smashingmagazine.com/2022/04/boost-resource-loading-new-priority-hint-fetchpriority/)
[18](https://perfmatters.io/docs/fetch-priority/)
[19](https://carlos.sanchezdonate.com/en/article/fetchpriority/)

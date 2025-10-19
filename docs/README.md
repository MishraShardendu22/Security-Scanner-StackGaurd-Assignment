# StackGuard Security Scanner - Complete Documentation Index

Welcome to the comprehensive documentation for the StackGuard Security Scanner project!

## 📚 Documentation Overview

This project includes detailed documentation covering all aspects of the application, from architecture to implementation details.

---

## 📖 Main Documentation Files

### 1. [Architecture Documentation](./ARCHITECTURE.md)
**Comprehensive system design and architecture overview**

**Topics Covered**:
- Complete technology stack breakdown
- Project structure and organization
- Core components explained
- Design patterns used
- Data flow diagrams
- Performance considerations
- Security features
- Scalability approach
- Monitoring and observability
- Future enhancements

**Perfect for**: Understanding the big picture and system design

---

### 2. [Concurrency Guide](./CONCURRENCY.md)
**Deep dive into goroutines, channels, and concurrent programming**

**Topics Covered**:
- Concurrency patterns (Worker Pool, Semaphore, WaitGroup)
- Goroutine implementation examples
- Channel usage (buffered, unbuffered, signal channels)
- Semaphore pattern for rate limiting
- Real-world code examples with line numbers
- Performance impact measurements
- Best practices and common pitfalls
- Debugging concurrent code
- Race condition prevention

**Perfect for**: Understanding parallel processing and Go concurrency

---

### 3. [Fiber Framework Guide](./FIBER.md)
**Complete guide to Fiber web framework usage**

**Topics Covered**:
- Why Fiber over other frameworks
- Configuration and setup
- Routing (parameters, groups, wildcards)
- Middleware implementation
- Request handling (body, query, params, headers)
- Response formatting
- Error handling strategies
- Best practices
- Performance benchmarks

**Perfect for**: Understanding HTTP handling and API design

---

### 4. [Templ Template Engine](./TEMPL.md)
**Type-safe templating with Templ**

**Topics Covered**:
- Why Templ over html/template
- Installation and setup
- Syntax guide (conditionals, loops, components)
- Component composition
- Integration with Fiber
- HTMX integration patterns
- Best practices and patterns
- Development workflow

**Perfect for**: Understanding frontend rendering and type-safe templates

---

### 5. [HTMX Integration](../HTMX_INTEGRATION.md)
**Dynamic HTML without JavaScript frameworks**

**Topics Covered**:
- HTMX setup and configuration
- Dashboard auto-refresh implementation
- Form submission without page reload
- Results page live updates
- CSS transitions
- HTMX attributes used
- Performance benefits
- Future enhancement ideas

**Perfect for**: Understanding dynamic UI updates

---

## 🚀 Quick Start Guides

### For Developers New to the Project

**Read in this order**:
1. [Architecture Documentation](./ARCHITECTURE.md) - Get the big picture
2. [Fiber Framework Guide](./FIBER.md) - Understand HTTP handling
3. [Concurrency Guide](./CONCURRENCY.md) - Learn about performance
4. [Templ Template Engine](./TEMPL.md) - Master the frontend
5. [HTMX Integration](../HTMX_INTEGRATION.md) - Add interactivity

### For Backend Developers

**Focus on**:
- [Architecture Documentation](./ARCHITECTURE.md) - System design
- [Fiber Framework Guide](./FIBER.md) - API development
- [Concurrency Guide](./CONCURRENCY.md) - Performance optimization

### For Frontend Developers

**Focus on**:
- [Templ Template Engine](./TEMPL.md) - Template development
- [HTMX Integration](../HTMX_INTEGRATION.md) - Dynamic updates
- [Architecture Documentation](./ARCHITECTURE.md) - API endpoints

### For DevOps/Infrastructure

**Focus on**:
- [Architecture Documentation](./ARCHITECTURE.md) - Deployment requirements
- [Concurrency Guide](./CONCURRENCY.md) - Resource usage
- README.md - Setup and configuration

---

## 🔍 Topic-Specific Navigation

### Performance & Optimization
- [Concurrency Guide](./CONCURRENCY.md) - Goroutines and channels
- [Fiber Framework Guide](./FIBER.md) - Fast HTTP processing
- [Architecture Documentation](./ARCHITECTURE.md) - Performance considerations

### Security
- [Architecture Documentation](./ARCHITECTURE.md) - Security features
- README.md - Secret detection patterns
- [Fiber Framework Guide](./FIBER.md) - Input validation

### Database
- [Architecture Documentation](./ARCHITECTURE.md) - MongoDB usage
- README.md - Database setup

### API Development
- [Fiber Framework Guide](./FIBER.md) - Complete API guide
- [Architecture Documentation](./ARCHITECTURE.md) - API design
- README.md - API endpoints

### Frontend
- [Templ Template Engine](./TEMPL.md) - Template development
- [HTMX Integration](../HTMX_INTEGRATION.md) - Dynamic UI
- [Architecture Documentation](./ARCHITECTURE.md) - Frontend stack

---

## 📊 Code Examples by Topic

### Goroutines & Channels
See: [Concurrency Guide - Real-World Examples](./CONCURRENCY.md#real-world-examples)

### Fiber Routing
See: [Fiber Guide - Routing Section](./FIBER.md#routing)

### Templ Components
See: [Templ Guide - Components](./TEMPL.md#components)

### HTMX Patterns
See: [HTMX Integration](../HTMX_INTEGRATION.md)

---

## 🛠️ Development Workflow

### 1. Setup Environment
```bash
# Install dependencies
go mod download

# Install Templ CLI
go install github.com/a-h/templ/cmd/templ@latest

# Setup MongoDB
# See README.md for details
```

### 2. Generate Templates
```bash
# Generate Go code from templ files
templ generate

# Watch mode for development
templ generate --watch
```

### 3. Run Application
```bash
# Development
go run main.go

# Production
go build -o scanner
./scanner
```

### 4. Testing
```bash
# Run tests
go test ./...

# With race detection
go test -race ./...

# With coverage
go test -cover ./...
```

---

## 🎯 Learning Paths

### Path 1: Backend Engineer
1. Understand [Architecture](./ARCHITECTURE.md)
2. Master [Fiber](./FIBER.md)
3. Learn [Concurrency](./CONCURRENCY.md)
4. Implement new API endpoints

### Path 2: Full-Stack Developer
1. Read [Architecture](./ARCHITECTURE.md)
2. Learn [Fiber](./FIBER.md) for backend
3. Master [Templ](./TEMPL.md) for frontend
4. Add [HTMX](../HTMX_INTEGRATION.md) features

### Path 3: Performance Engineer
1. Study [Concurrency Guide](./CONCURRENCY.md)
2. Review [Architecture - Performance](./ARCHITECTURE.md#performance-considerations)
3. Analyze [Fiber benchmarks](./FIBER.md#benchmarks)
4. Optimize scanning algorithms

### Path 4: Frontend Developer
1. Learn [Templ basics](./TEMPL.md)
2. Understand [HTMX patterns](../HTMX_INTEGRATION.md)
3. Review [Architecture - Frontend](./ARCHITECTURE.md#frontend-technologies)
4. Build new UI components

---

## 📝 Additional Resources

### External Documentation
- [Go Documentation](https://go.dev/doc/)
- [Fiber Documentation](https://docs.gofiber.io/)
- [Templ Documentation](https://templ.guide/)
- [HTMX Documentation](https://htmx.org/docs/)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)

### Useful Commands
```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run

# Generate docs
godoc -http=:6060

# Format templ files
templ fmt template/

# Build for production
go build -ldflags="-s -w" -o scanner
```

---

## 🤝 Contributing

When contributing to this project:

1. **Read the relevant documentation first**
2. **Follow the established patterns**
3. **Write tests for new features**
4. **Update documentation if needed**
5. **Use the project's code style**

### Code Style Guide
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use meaningful variable names
- Add comments for exported functions
- Keep functions small and focused
- Write idiomatic Go code

---

## 🐛 Troubleshooting

### Common Issues

**Templates not updating?**
```bash
# Regenerate templates
templ generate
```

**MongoDB connection failed?**
- Check `.env` file
- Ensure MongoDB is running
- Verify connection string

**Goroutine leaks?**
```bash
# Run with race detector
go run -race main.go
```

**Build errors?**
```bash
# Clean and rebuild
go clean
go build
```

---

## 📈 Project Statistics

- **Lines of Code**: ~5,000+
- **Go Packages**: 8 (controller, database, models, route, template, util, main)
- **Templates**: 10+ Templ components
- **API Endpoints**: 15+
- **Concurrent Workers**: Up to 30 goroutines
- **Secret Patterns**: 15+ detection patterns
- **Performance**: 10x speedup with concurrency

---

## 🎓 Learning Objectives

After reading this documentation, you should understand:

✅ How the application is structured and organized  
✅ Why specific technologies were chosen  
✅ How goroutines and channels work together  
✅ How to use Fiber for web development  
✅ How Templ provides type-safe templates  
✅ How HTMX enables dynamic updates  
✅ How to scan for security vulnerabilities  
✅ How to optimize for performance  
✅ How to extend and maintain the system  

---

## 📞 Support

For questions or issues:
1. Check the relevant documentation file
2. Review code examples in the docs
3. Check the troubleshooting section
4. Review GitHub issues

---

## 📅 Documentation Updates

This documentation is actively maintained. Last major update: **October 19, 2025**

**Recent Changes**:
- Added HTMX integration documentation
- Enhanced concurrency guide with examples
- Added Fiber framework detailed guide
- Created Templ template engine guide
- Reorganized architecture documentation

---

## 🎉 Conclusion

This documentation suite provides everything you need to understand, develop, and maintain the StackGuard Security Scanner. Each document focuses on a specific aspect of the system, making it easy to find the information you need.

**Happy Learning and Coding!** 🚀

---

*For the main project README, see [README.md](../README.md)*

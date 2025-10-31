package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MishraShardendu22/Scanner/util"

	"github.com/MishraShardendu22/Scanner/database"
	"github.com/MishraShardendu22/Scanner/models"
	"github.com/MishraShardendu22/Scanner/route"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

var (
	logToFile = flag.Bool("log-to-file", false, "Save all logs to files in the logs/ directory")
	logFile   *os.File
)

func main() {

	flag.Parse()
	fmt.Println("Stack Guard Assignment")

	config := loadConfig()

	if err := database.ConnectDatabase(config.DbName, config.MongoURI); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	setupLogger(config)

	logger := slog.Default()
	logger.Info("Starting Security Scanner",
		"environment", config.Environment,
		"port", config.Port,
		"log_level", config.LogLevel,
	)

	app := fiber.New(fiber.Config{
		AppName:               "Security Scanner",
		ServerHeader:          "Security-Scanner",
		ReadTimeout:           30 * time.Second,
		WriteTimeout:          30 * time.Second,
		IdleTimeout:           120 * time.Second,
		DisableStartupMessage: false,

		// Performance optimizations
		Prefork:              false, // Set to true in production for multi-core usage
		StrictRouting:        false,
		CaseSensitive:        false,
		UnescapePath:         true,
		ETag:                 true,            // Enable ETag for caching
		BodyLimit:            4 * 1024 * 1024, // 4MB body limit
		Concurrency:          256 * 1024,      // Max concurrent connections
		ReadBufferSize:       4096,
		WriteBufferSize:      4096,
		CompressedFileSuffix: ".fiber.gz",
		ProxyHeader:          fiber.HeaderXForwardedFor,

		// Disable unnecessary features for performance
		DisableKeepalive:          false,
		DisableDefaultDate:        false,
		DisableDefaultContentType: false,
		DisableHeaderNormalizing:  false,

		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Error("request error", slog.Group("req",
				slog.String("method", c.Method()),
				slog.String("path", c.Path()),
				slog.String("error", err.Error()),
			))
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return util.ResponseAPI(c, code, err.Error(), nil, "")
		},
	})

	setupMiddleware(app, config)
	SetUpRoutes(app, logger, config)

	go func() {

		logger.Info("Server starting", "port", config.Port)
		if err := app.Listen(":" + config.Port); err != nil {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	gracefulShutdown(app, logger)
}

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
	route.SetupResultRoutes(app)
}

func setupLogger(config *models.Config) {

	var level slog.Level
	switch config.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{

		Level:     level,
		AddSource: true,
	}
	var writer io.Writer = os.Stdout

	if *logToFile {
		if err := os.MkdirAll("logs", 0755); err != nil {
			log.Printf("Failed to create logs directory: %v", err)
		} else {

			timestamp := time.Now().Format("2006-01-02_15-04-05")
			logFileName := fmt.Sprintf("logs/server_%s.log", timestamp)
			var err error
			logFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Printf("Failed to open log file: %v", err)
			} else {

				writer = io.MultiWriter(os.Stdout, logFile)
				log.SetOutput(writer)
				fmt.Printf("📝 Logs will be saved to: %s\n", logFileName)
			}
		}
	}

	handler := util.NewPrettyTextHandler(writer, opts)

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func loadConfig() *models.Config {

	config := &models.Config{

		Port:             util.GetEnv("PORT", "8080"),
		DbName:           util.GetEnv("DB_NAME", "main"),
		LogLevel:         util.GetEnv("LOG_LEVEL", "debug"),
		CorsAllowOrigins: util.GetEnv("CORS_ALLOW_ORIGINS", ""),
		Environment:      util.GetEnv("ENVIRONMENT", "development"),
		MongoURI:         util.GetEnv("MONGODB_URI", "mongodb+srv://shardendu:some_password@cluster0.0uz8vjv.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"),
	}

	return config
}

func setupMiddleware(app *fiber.App, config *models.Config) {
	// 1. Recover middleware - must be first
	app.Use(recover.New(recover.Config{
		EnableStackTrace: config.Environment == "development",
	}))

	// 2. Compression middleware - compress responses (gzip)
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // Balance between speed and compression
	}))

	// 3. ETag middleware - enable caching with ETags
	app.Use(etag.New(etag.Config{
		Weak: true, // Use weak ETags for better performance
	}))

	// 4. Cache middleware for static content and API responses
	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			// Skip cache for POST, PUT, DELETE, PATCH requests
			return c.Method() != "GET" && c.Method() != "HEAD"
		},
		Expiration:           5 * time.Minute, // Cache for 5 minutes
		CacheControl:         true,            // Add Cache-Control headers
		CacheHeader:          "X-Cache",       // Custom header to show cache status
		Methods:              []string{"GET", "HEAD"},
		StoreResponseHeaders: true,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Path() + "?" + c.Request().URI().QueryArgs().String()
		},
	}))

	// 5. Rate limiting to prevent abuse
	app.Use(limiter.New(limiter.Config{
		Max:        100,             // 100 requests
		Expiration: 1 * time.Minute, // per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests, please try again later",
			})
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		LimiterMiddleware:      limiter.SlidingWindow{},
	}))

	// 6. Favicon middleware - cache favicon
	app.Use(favicon.New(favicon.Config{
		File:         "./public/favicon.ico",
		CacheControl: "public, max-age=31536000", // Cache for 1 year
	}))

	// 7. CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:  config.CorsAllowOrigins,
		AllowMethods:  "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders: "Content-Length",
		MaxAge:        86400,
	}))

	// 8. Logger middleware - should be last
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))
}

func gracefulShutdown(app *fiber.App, logger *slog.Logger) {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}
	logger.Info("Server exited")

	if logFile != nil {
		logFile.Close()
	}
}

func init() {

	currEnv := "development"

	if currEnv == "development" {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: error loading .env file: %v", err)
		}
	}
}

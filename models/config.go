package models

type Config struct {
	Port             string
	DbName           string
	LogLevel         string
	MongoURI         string
	Environment      string
	CorsAllowOrigins string
}

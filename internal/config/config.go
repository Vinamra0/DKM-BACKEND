package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port           int
	MongoURI       string
	MongoDB        string
	JWTSecret      string
	AdminEmail     string
	AdminPassword  string
	AdminName      string
	AllowedOrigins []string
}

func Load() Config {
	port := 3001
	if p := os.Getenv("PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			port = v
		} else {
			log.Printf("invalid PORT, defaulting to %d", port)
		}
	}

	allowed := []string{"*"}
	if s := os.Getenv("ALLOWED_ORIGINS"); s != "" {
		parts := strings.Split(s, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		allowed = parts
	}

	return Config{
		Port:           port,
		MongoURI:       getenv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:        getenv("MONGO_DB", "dkm_pharma"),
		JWTSecret:      getenv("JWT_SECRET", "dev-secret"),
		AdminEmail:     getenv("ADMIN_EMAIL", "admin@example.com"),
		AdminPassword:  getenv("ADMIN_PASSWORD", "changeme"),
		AdminName:      getenv("ADMIN_NAME", "Admin"),
		AllowedOrigins: allowed,
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

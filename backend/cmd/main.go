// @title           JiraFlow API
// @version         1.0
// @description     Project management backend API (Jira-like).
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Value formatı: **Bearer &lt;token&gt;**
package main

import (
	"log"

	"github.com/joho/godotenv"
	_ "github.com/jira-backend/jiraflow-backend/docs"
	"github.com/jira-backend/jiraflow-backend/internal/app"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/config"
)

func main() {
	// Load .env if present; ignored silently in production where vars are injected externally.
	_ = godotenv.Load()

	cfg := config.Load()
	if err := app.Run(cfg); err != nil {
		log.Fatalf("app run: %v", err)
	}
}

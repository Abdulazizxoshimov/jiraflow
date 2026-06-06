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

	_ "github.com/jira-backend/jiraflow-backend/docs"
	"github.com/jira-backend/jiraflow-backend/internal/app"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/config"
)

func main() {
	cfg := config.Load()
	if err := app.Run(cfg); err != nil {
		log.Fatalf("app run: %v", err)
	}
}

# JiraFlow

Full-stack project management platform — Jira + Confluence in one system.

## Structure

```
jiraflow/
├── backend/     # Go REST API (Gin, PostgreSQL, Redis, MinIO)
├── frontend/    # UI prototype (JSX/HTML)
└── docs/        # Architecture & API documentation
```

## Backend

```bash
cd backend
go run ./cmd/main.go
```

**Requirements:** PostgreSQL, Redis, MinIO

## Frontend

Open `frontend/index.html` in a browser — runs fully against mock API, no server needed.  
To connect to the real backend: remove the `api.jsx` script tag from `index.html`.

## API

Base URL: `http://localhost:8080/api/v1`  
Docs: `http://localhost:8080/swagger/index.html`

package page_export

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"time"

	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/tiptap"
)

type useCase struct {
	pageRepo repository.PageRepository
}

func New(pageRepo repository.PageRepository) UseCase {
	return &useCase{pageRepo: pageRepo}
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="uz">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}}</title>
<style>
  body{font-family:system-ui,sans-serif;max-width:860px;margin:40px auto;padding:0 20px;color:#1a1a1a;line-height:1.6}
  h1{font-size:2rem;border-bottom:2px solid #e5e7eb;padding-bottom:12px}
  h2{font-size:1.5rem;margin-top:2em}h3{font-size:1.2rem}
  pre{background:#f3f4f6;padding:16px;border-radius:6px;overflow-x:auto}
  code{background:#f3f4f6;padding:2px 4px;border-radius:3px;font-size:0.9em}
  blockquote{border-left:4px solid #6b7280;margin:0;padding-left:1em;color:#6b7280}
  table{border-collapse:collapse;width:100%}
  td,th{border:1px solid #e5e7eb;padding:8px 12px}
  th{background:#f9fafb;font-weight:600}
  .meta{color:#6b7280;font-size:0.9rem;margin-bottom:2em}
  hr{border:none;border-top:1px solid #e5e7eb;margin:2em 0}
</style>
</head>
<body>
<h1>{{.Title}}</h1>
<div class="meta">Exported: {{.ExportedAt}} · Version {{.Version}}</div>
<hr>
{{.Body}}
</body>
</html>`

func (uc *useCase) ExportHTML(ctx context.Context, pageID string) ([]byte, string, error) {
	page, err := uc.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportHTML: %w", err)
	}

	var bodyHTML string
	if len(page.Content) > 0 {
		contentJSON, _ := json.Marshal(page.Content)
		doc, parseErr := tiptap.Parse(contentJSON)
		if parseErr == nil {
			bodyHTML = doc.RenderHTML()
		} else {
			bodyHTML = "<p>" + template.HTMLEscapeString(page.ContentText) + "</p>"
		}
	}

	tmpl, err := template.New("export").Parse(htmlTemplate)
	if err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportHTML template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]any{
		"Title":      page.Title,
		"ExportedAt": time.Now().UTC().Format("2006-01-02 15:04 UTC"),
		"Version":    page.CurrentVersion,
		"Body":       template.HTML(bodyHTML),
	}); err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportHTML render: %w", err)
	}

	filename := sanitizeFilename(page.Title) + ".html"
	return buf.Bytes(), filename, nil
}

func sanitizeFilename(title string) string {
	var b bytes.Buffer
	for _, r := range title {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			b.WriteRune(r)
		} else if r == ' ' {
			b.WriteRune('_')
		}
	}
	result := b.String()
	if result == "" {
		result = "page"
	}
	return result
}

// ExportMarkdown converts the page content to GitHub-flavored Markdown.
func (uc *useCase) ExportMarkdown(ctx context.Context, pageID string) ([]byte, string, error) {
	page, err := uc.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportMarkdown: %w", err)
	}

	var md string
	if len(page.Content) > 0 {
		contentJSON, _ := json.Marshal(page.Content)
		doc, parseErr := tiptap.Parse(contentJSON)
		if parseErr == nil {
			md = "# " + page.Title + "\n\n" + doc.RenderMarkdown()
		}
	}
	if md == "" {
		md = "# " + page.Title + "\n\n" + page.ContentText
	}

	filename := sanitizeFilename(page.Title) + ".md"
	return []byte(md), filename, nil
}

// ExportDOCX converts the page to .docx using pandoc (must be installed on the system).
func (uc *useCase) ExportDOCX(ctx context.Context, pageID string) ([]byte, string, error) {
	mdBytes, _, err := uc.ExportMarkdown(ctx, pageID)
	if err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportDOCX markdown: %w", err)
	}

	if _, err := exec.LookPath("pandoc"); err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportDOCX: pandoc not installed on this server")
	}

	tmpMD, err := os.CreateTemp("", "page-*.md")
	if err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportDOCX tmp md: %w", err)
	}
	defer os.Remove(tmpMD.Name())
	if _, err := tmpMD.Write(mdBytes); err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportDOCX write md: %w", err)
	}
	tmpMD.Close()

	tmpDOCX, err := os.CreateTemp("", "page-*.docx")
	if err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportDOCX tmp docx: %w", err)
	}
	tmpDOCX.Close()
	defer os.Remove(tmpDOCX.Name())

	cmd := exec.CommandContext(ctx, "pandoc", "-f", "markdown", "-t", "docx", "-o", tmpDOCX.Name(), tmpMD.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportDOCX pandoc: %w — %s", err, string(out))
	}

	docxBytes, err := os.ReadFile(tmpDOCX.Name())
	if err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportDOCX read docx: %w", err)
	}

	page, err := uc.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, "", err
	}
	filename := sanitizeFilename(page.Title) + ".docx"
	return docxBytes, filename, nil
}

// ExportPDF converts the page to PDF using wkhtmltopdf (must be installed on the system).
func (uc *useCase) ExportPDF(ctx context.Context, pageID string) ([]byte, string, error) {
	htmlBytes, _, err := uc.ExportHTML(ctx, pageID)
	if err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportPDF html: %w", err)
	}

	if _, err := exec.LookPath("wkhtmltopdf"); err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportPDF: wkhtmltopdf not installed on this server")
	}

	tmpHTML, err := os.CreateTemp("", "page-*.html")
	if err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportPDF tmp html: %w", err)
	}
	defer os.Remove(tmpHTML.Name())
	if _, err := tmpHTML.Write(htmlBytes); err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportPDF write html: %w", err)
	}
	tmpHTML.Close()

	tmpPDF, err := os.CreateTemp("", "page-*.pdf")
	if err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportPDF tmp pdf: %w", err)
	}
	tmpPDF.Close()
	defer os.Remove(tmpPDF.Name())

	cmd := exec.CommandContext(ctx, "wkhtmltopdf", "--quiet", tmpHTML.Name(), tmpPDF.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportPDF wkhtmltopdf: %w — %s", err, string(out))
	}

	pdfBytes, err := os.ReadFile(tmpPDF.Name())
	if err != nil {
		return nil, "", fmt.Errorf("pageExport.ExportPDF read pdf: %w", err)
	}

	page, err := uc.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, "", err
	}
	filename := sanitizeFilename(page.Title) + ".pdf"
	return pdfBytes, filename, nil
}

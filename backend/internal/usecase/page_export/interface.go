package page_export

import "context"

type UseCase interface {
	ExportHTML(ctx context.Context, pageID string) ([]byte, string, error)     // content, filename, error
	ExportPDF(ctx context.Context, pageID string) ([]byte, string, error)      // content, filename, error
	ExportMarkdown(ctx context.Context, pageID string) ([]byte, string, error) // content, filename, error
	ExportDOCX(ctx context.Context, pageID string) ([]byte, string, error)     // content, filename, error
}

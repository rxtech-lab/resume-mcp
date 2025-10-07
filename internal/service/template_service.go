package service

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/rxtech-lab/resume-mcp/internal/models"
)

type stringBuilder struct {
	content string
}

func (sb *stringBuilder) Write(p []byte) (n int, err error) {
	sb.content += string(p)
	return len(p), nil
}

func (sb *stringBuilder) String() string {
	return sb.content
}

type TemplateService struct {
}

func NewTemplateService() *TemplateService {
	return &TemplateService{}
}

func (s *TemplateService) GeneratePreview(templateStr, css string, resume models.Resume) (string, error) {
	return s.GeneratePreviewWithOptions(templateStr, css, resume, false, "")
}

func (s *TemplateService) GeneratePreviewWithOptions(templateStr, css string, resume models.Resume, includeDownloadButton bool, downloadURL string) (string, error) {

	tmpl, err := template.New("resume").Parse(templateStr)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.SetFlags(0)
		log.Printf("Template parse error: %v", err)
		log.SetOutput(io.Discard)
		return "", fmt.Errorf("Template parse error: %v", err)
	}

	var html string
	builder := &stringBuilder{}
	if err := tmpl.Execute(builder, resume); err != nil {
		log.SetOutput(os.Stderr)
		log.SetFlags(0)
		log.Printf("Template execution error: %v", err)
		log.SetOutput(io.Discard)
		return "", fmt.Errorf("Template execution error: %v", err)
	}
	html = builder.String()

	var cssStyle string
	if css != "" {
		cssStyle = "<style>" + css + "</style>"
	}

	var downloadButton string
	if includeDownloadButton && downloadURL != "" {
		downloadButton = `
    <a href="` + downloadURL + `"
       class="fixed top-4 right-4 bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-lg shadow-lg transition-colors duration-200 no-print z-50"
       download="resume.pdf">
        Download PDF
    </a>
    <style>
        @media print {
            .no-print {
                display: none !important;
            }
        }
    </style>`
	}

	fullHTML := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Resume Preview</title>
    <script src="https://cdn.tailwindcss.com"></script>
    ` + cssStyle + `
</head>
<body>
    ` + downloadButton + `
    ` + html + `
</body>
</html>`

	return fullHTML, nil
}

func (s *TemplateService) GeneratePDF(templateStr, css string, resume models.Resume) ([]byte, error) {
	// Generate HTML without download button
	html, err := s.GeneratePreviewWithOptions(templateStr, css, resume, false, "")
	if err != nil {
		return nil, err
	}

	// Get remote Chrome URL from environment
	remoteURL := os.Getenv("CHROMEDP_REMOTE_URL")

	var ctx context.Context
	var cancel context.CancelFunc

	if remoteURL != "" {
		// Use remote Chrome instance
		allocCtx, allocCancel := chromedp.NewRemoteAllocator(context.Background(), remoteURL)
		defer allocCancel()
		ctx, cancel = chromedp.NewContext(allocCtx)
	} else {
		// Use local Chrome instance
		ctx, cancel = chromedp.NewContext(context.Background())
	}
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var pdfBuffer []byte

	// Navigate to data URL and generate PDF
	err = chromedp.Run(ctx,
		chromedp.Navigate("data:text/html,"+html),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Use page.PrintToPDF with print background enabled
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}
			pdfBuffer = buf
			return nil
		}),
	)

	if err != nil {
		log.SetOutput(os.Stderr)
		log.SetFlags(0)
		log.Printf("PDF generation error: %v", err)
		log.SetOutput(io.Discard)
		return nil, fmt.Errorf("PDF generation error: %v", err)
	}

	return pdfBuffer, nil
}

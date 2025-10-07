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

	var appBar string
	if includeDownloadButton && downloadURL != "" {
		appBar = `
    <div class="app-bar no-print" style="position: fixed; top: 0; left: 0; right: 0; height: 56px; background: white; border-bottom: 1px solid hsl(214.3 31.8% 91.4%); display: flex; align-items: center; justify-content: space-between; padding: 0 24px; z-index: 50;">
        <h1 style="font-size: 18px; font-weight: 600; color: hsl(222.2 47.4% 11.2%); margin: 0;">Resume Preview</h1>
        <button id="download-btn"
                class="inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 h-9 w-9"
                style="background: transparent; border: 1px solid hsl(214.3 31.8% 91.4%); color: hsl(222.2 47.4% 11.2%);"
                onclick="downloadPDF()"
                title="Download PDF">
            <svg id="download-icon" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                <polyline points="7 10 12 15 17 10"></polyline>
                <line x1="12" y1="15" x2="12" y2="3"></line>
            </svg>
            <svg id="download-spinner" class="hidden animate-spin" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 12a9 9 0 1 1-6.219-8.56"></path>
            </svg>
        </button>
    </div>
    <div class="no-print" style="height: 56px;"></div>
    <style>
        @media print {
            .no-print {
                display: none !important;
            }
        }
        @keyframes spin {
            from {
                transform: rotate(0deg);
            }
            to {
                transform: rotate(360deg);
            }
        }
        .animate-spin {
            animation: spin 1s linear infinite;
        }
        .hidden {
            display: none;
        }
        #download-btn:hover {
            background: hsl(214.3 31.8% 91.4%) !important;
        }
        #download-btn:disabled {
            opacity: 0.5;
            cursor: not-allowed;
        }
    </style>
    <script>
        async function downloadPDF() {
            const btn = document.getElementById('download-btn');
            const icon = document.getElementById('download-icon');
            const spinner = document.getElementById('download-spinner');

            // Disable button and show loading state
            btn.disabled = true;
            icon.classList.add('hidden');
            spinner.classList.remove('hidden');

            try {
                const response = await fetch('` + downloadURL + `');
                if (!response.ok) throw new Error('Download failed');

                const blob = await response.blob();
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.style.display = 'none';
                a.href = url;
                a.download = 'resume.pdf';
                document.body.appendChild(a);
                a.click();
                window.URL.revokeObjectURL(url);
                document.body.removeChild(a);

                // Reset state after download
                setTimeout(() => {
                    icon.classList.remove('hidden');
                    spinner.classList.add('hidden');
                    btn.disabled = false;
                }, 500);
            } catch (error) {
                console.error('Download error:', error);
                icon.classList.remove('hidden');
                spinner.classList.add('hidden');
                btn.disabled = false;
                alert('Failed to download PDF. Please try again.');
            }
        }
    </script>`
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
    ` + appBar + `
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
	var allocCancel context.CancelFunc

	if remoteURL != "" {
		// Use remote Chrome instance
		var allocCtx context.Context
		allocCtx, allocCancel = chromedp.NewRemoteAllocator(context.Background(), remoteURL)
		// Ensure allocator is always cancelled
		defer func() {
			if allocCancel != nil {
				allocCancel()
			}
		}()
		ctx, cancel = chromedp.NewContext(allocCtx)
	} else {
		// Use local Chrome instance
		ctx, cancel = chromedp.NewContext(context.Background())
	}
	// Ensure context is always cancelled
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	// Set timeout - use a separate variable to avoid shadowing cancel
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 30*time.Second)
	defer timeoutCancel()

	var pdfBuffer []byte

	// Navigate to data URL and generate PDF
	err = chromedp.Run(timeoutCtx,
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

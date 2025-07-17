package service

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"

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
    ` + html + `
</body>
</html>`

	return fullHTML, nil
}

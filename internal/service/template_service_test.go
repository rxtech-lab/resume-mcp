package service

import (
	"strings"
	"testing"
	"time"

	"github.com/rxtech-lab/resume-mcp/internal/models"
)

func TestStringBuilder_Write(t *testing.T) {
	tests := []struct {
		name     string
		inputs   [][]byte
		expected string
	}{
		{
			name:     "single write",
			inputs:   [][]byte{[]byte("hello")},
			expected: "hello",
		},
		{
			name:     "multiple writes",
			inputs:   [][]byte{[]byte("hello"), []byte(" "), []byte("world")},
			expected: "hello world",
		},
		{
			name:     "empty write",
			inputs:   [][]byte{[]byte("")},
			expected: "",
		},
		{
			name:     "unicode characters",
			inputs:   [][]byte{[]byte("Hello 世界")},
			expected: "Hello 世界",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &stringBuilder{}

			for _, input := range tt.inputs {
				n, err := sb.Write(input)
				if err != nil {
					t.Errorf("Write() error = %v, wantErr false", err)
					return
				}
				if n != len(input) {
					t.Errorf("Write() returned %d, want %d", n, len(input))
				}
			}

			if got := sb.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStringBuilder_String(t *testing.T) {
	sb := &stringBuilder{}

	// Test initial empty state
	if got := sb.String(); got != "" {
		t.Errorf("String() on empty builder = %v, want empty string", got)
	}

	// Test after writing content
	sb.Write([]byte("test content"))
	if got := sb.String(); got != "test content" {
		t.Errorf("String() = %v, want %v", got, "test content")
	}
}

func TestNewTemplateService(t *testing.T) {
	service := NewTemplateService()

	if service == nil {
		t.Error("NewTemplateService() returned nil")
	}

	if _, ok := interface{}(service).(*TemplateService); !ok {
		t.Error("NewTemplateService() did not return *TemplateService")
	}
}

func TestTemplateService_GeneratePreview(t *testing.T) {
	service := NewTemplateService()

	// Create sample resume data for testing
	startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	sampleResume := models.Resume{
		Name:        "John Doe",
		Description: "Software Engineer",
		Contacts: []models.Contact{
			{Key: "email", Value: "john@example.com"},
			{Key: "phone", Value: "123-456-7890"},
		},
		WorkExperiences: []models.WorkExperience{
			{
				Company:   "Tech Corp",
				JobTitle:  "Software Engineer",
				StartDate: startDate,
				EndDate:   &endDate,
				FeatureMaps: []models.FeatureMap{
					{Key: "description", Value: "Developed web applications"},
				},
			},
		},
	}

	tests := []struct {
		name        string
		templateStr string
		css         string
		resume      models.Resume
		wantErr     bool
		errContains string
		contains    []string
	}{
		{
			name:        "valid template without CSS",
			templateStr: "<h1>{{.Name}}</h1><p>{{.Description}}</p>",
			css:         "",
			resume:      sampleResume,
			wantErr:     false,
			contains:    []string{"John Doe", "Software Engineer", "<!DOCTYPE html>", "<script src=\"https://cdn.tailwindcss.com\"></script>"},
		},
		{
			name:        "valid template with CSS",
			templateStr: "<h1>{{.Name}}</h1>",
			css:         "h1 { color: red; }",
			resume:      sampleResume,
			wantErr:     false,
			contains:    []string{"John Doe", "<style>h1 { color: red; }</style>", "<!DOCTYPE html>"},
		},
		{
			name:        "template with work experience",
			templateStr: "{{range .WorkExperiences}}<div>{{.Company}} - {{.JobTitle}}</div>{{end}}",
			css:         "",
			resume:      sampleResume,
			wantErr:     false,
			contains:    []string{"Tech Corp - Software Engineer"},
		},
		{
			name:        "invalid template syntax",
			templateStr: "{{.InvalidSyntax",
			css:         "",
			resume:      sampleResume,
			wantErr:     true,
			errContains: "Template parse error",
		},
		{
			name:        "template execution error - invalid field",
			templateStr: "{{.NonExistentField.Test}}",
			css:         "",
			resume:      sampleResume,
			wantErr:     true,
			errContains: "Template execution error",
		},
		{
			name:        "empty template",
			templateStr: "",
			css:         "",
			resume:      sampleResume,
			wantErr:     false,
			contains:    []string{"<!DOCTYPE html>", "<body>", "</body>"},
		},
		{
			name:        "template with special characters",
			templateStr: "<p>Special chars: &lt; &gt; &amp;</p>",
			css:         "",
			resume:      sampleResume,
			wantErr:     false,
			contains:    []string{"Special chars: &lt; &gt; &amp;"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GeneratePreview(tt.templateStr, tt.css, tt.resume)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GeneratePreview() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("GeneratePreview() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("GeneratePreview() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check that all expected content is present
			for _, content := range tt.contains {
				if !strings.Contains(got, content) {
					t.Errorf("GeneratePreview() result does not contain expected content: %v", content)
				}
			}

			// Verify basic HTML structure
			if !strings.Contains(got, "<!DOCTYPE html>") {
				t.Error("GeneratePreview() result missing DOCTYPE declaration")
			}
			if !strings.Contains(got, "<html lang=\"en\">") {
				t.Error("GeneratePreview() result missing html tag")
			}
			if !strings.Contains(got, "</html>") {
				t.Error("GeneratePreview() result missing closing html tag")
			}
		})
	}
}

func TestTemplateService_GeneratePreview_ComplexResume(t *testing.T) {
	service := NewTemplateService()

	// Create a more complex resume for testing
	startDate1 := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate1 := time.Date(2018, 6, 1, 0, 0, 0, 0, time.UTC)
	endDate2 := time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC)
	eduStartDate := time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC)
	eduEndDate := time.Date(2018, 5, 31, 0, 0, 0, 0, time.UTC)

	complexResume := models.Resume{
		Name:        "Jane Smith",
		Description: "Experienced Software Engineer",
		Contacts: []models.Contact{
			{Key: "email", Value: "jane.smith@example.com"},
			{Key: "phone", Value: "555-123-4567"},
			{Key: "address", Value: "123 Main St, City, State 12345"},
		},
		WorkExperiences: []models.WorkExperience{
			{
				Company:   "TechCorp Inc",
				JobTitle:  "Senior Software Engineer",
				StartDate: startDate1,
				FeatureMaps: []models.FeatureMap{
					{Key: "description", Value: "Lead development of microservices architecture"},
				},
			},
			{
				Company:   "StartupXYZ",
				JobTitle:  "Full Stack Developer",
				StartDate: endDate1,
				EndDate:   &endDate2,
				FeatureMaps: []models.FeatureMap{
					{Key: "description", Value: "Built web applications using React and Node.js"},
				},
			},
		},
		Educations: []models.Education{
			{
				SchoolName: "University of Technology",
				StartDate:  eduStartDate,
				EndDate:    &eduEndDate,
				FeatureMaps: []models.FeatureMap{
					{Key: "degree", Value: "Bachelor of Computer Science"},
					{Key: "gpa", Value: "3.8"},
				},
			},
		},
	}

	templateStr := `
<div class="resume">
	<header>
		<h1>{{.Name}}</h1>
		<p>{{.Description}}</p>
		{{range .Contacts}}
			{{if eq .Key "email"}}<span>{{.Value}}</span>{{end}}
			{{if eq .Key "phone"}} | <span>{{.Value}}</span>{{end}}
		{{end}}
		{{range .Contacts}}{{if eq .Key "address"}}<p>{{.Value}}</p>{{end}}{{end}}
	</header>
	
	{{if .WorkExperiences}}
	<section>
		<h2>Work Experience</h2>
		{{range .WorkExperiences}}
		<div class="job">
			<h3>{{.JobTitle}} at {{.Company}}</h3>
			<p>{{.StartDate.Format "2006-01"}}{{if .EndDate}} - {{.EndDate.Format "2006-01"}}{{else}} - Present{{end}}</p>
			{{range .FeatureMaps}}{{if eq .Key "description"}}<p>{{.Value}}</p>{{end}}{{end}}
		</div>
		{{end}}
	</section>
	{{end}}
	
	{{if .Educations}}
	<section>
		<h2>Education</h2>
		{{range .Educations}}
		<div class="education">
			{{range .FeatureMaps}}{{if eq .Key "degree"}}<h3>{{.Value}}</h3>{{end}}{{end}}
			<p>{{.SchoolName}}</p>
			<p>{{.StartDate.Format "2006-01"}}{{if .EndDate}} - {{.EndDate.Format "2006-01"}}{{end}}</p>
			{{range .FeatureMaps}}{{if eq .Key "gpa"}}<p>GPA: {{.Value}}</p>{{end}}{{end}}
		</div>
		{{end}}
	</section>
	{{end}}
</div>`

	css := `
.resume { font-family: Arial, sans-serif; }
.job, .education { margin-bottom: 1rem; }
h1 { color: #333; }
h2 { color: #666; border-bottom: 1px solid #ccc; }
`

	result, err := service.GeneratePreview(templateStr, css, complexResume)
	if err != nil {
		t.Fatalf("GeneratePreview() failed: %v", err)
	}

	expectedContent := []string{
		"Jane Smith",
		"Experienced Software Engineer",
		"jane.smith@example.com",
		"555-123-4567",
		"123 Main St, City, State 12345",
		"Senior Software Engineer at TechCorp Inc",
		"Full Stack Developer at StartupXYZ",
		"University of Technology",
		"Bachelor of Computer Science",
		"GPA: 3.8",
		"Work Experience",
		"Education",
		".resume { font-family: Arial, sans-serif; }",
	}

	for _, content := range expectedContent {
		if !strings.Contains(result, content) {
			t.Errorf("Result does not contain expected content: %v", content)
		}
	}
}

// Benchmark tests
func BenchmarkStringBuilder_Write(b *testing.B) {
	sb := &stringBuilder{}
	data := []byte("benchmark test data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sb.Write(data)
	}
}

func BenchmarkTemplateService_GeneratePreview(b *testing.B) {
	service := NewTemplateService()
	resume := models.Resume{
		Name:        "Benchmark User",
		Description: "Test User",
		Contacts: []models.Contact{
			{Key: "email", Value: "bench@example.com"},
		},
	}
	templateStr := "<h1>{{.Name}}</h1><p>{{.Description}}</p>"
	css := "h1 { color: blue; }"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GeneratePreview(templateStr, css, resume)
		if err != nil {
			b.Fatal(err)
		}
	}
}

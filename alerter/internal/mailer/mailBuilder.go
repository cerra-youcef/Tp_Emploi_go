package mailer

import (
	"bytes"
	"embed"
	"errors"
	"html/template"
	"strings"

	"github.com/adrg/frontmatter"
)

//go:embed templates
var embeddedTemplates embed.FS

type TemplateData struct {
	EventName   string
	NewDate     string
	NewLocation string
}

type Matter struct {
	Subject string `yaml:"subject"`
}

func GetEmailContent(templatePath string, data TemplateData) (subject string, content string, err error) {
	// Charger le template
	tplContent, err := embeddedTemplates.ReadFile(templatePath)
	if err != nil {
		return "", "", errors.New("failed to read template: " + err.Error())
	}

	// Parser le frontmatter
	var matter Matter
	remainingContent, err := frontmatter.Parse(strings.NewReader(string(tplContent)), &matter)
	if err != nil {
		return "", "", errors.New("failed to parse frontmatter: " + err.Error())
	}

	// Exécuter le template avec les données
	tmpl, err := template.New("email").Parse(string(remainingContent))
	if err != nil {
		return "", "", errors.New("failed to parse template: " + err.Error())
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", "", errors.New("failed to execute template: " + err.Error())
	}

	return matter.Subject, buf.String(), nil
}

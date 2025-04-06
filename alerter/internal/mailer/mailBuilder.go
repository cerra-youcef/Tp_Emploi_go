package mailer

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"strings"

	"github.com/adrg/frontmatter"
)

//go:embed templates
var embeddedTemplates embed.FS

type TemplateData struct {
	EventName string
	Date      string
	Location  string
	Changes   map[string]string
}

type Matter struct {
	Subject string `yaml:"subject"`
}

func GetEmailContent(templatePath string, data TemplateData, alertType string) (subject string, content string, err error) {
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

	if alertType == "event.updated" && len(data.Changes) > 0 {
		for key, value := range data.Changes {
			buf.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		}
	}

	return matter.Subject, buf.String(), nil
}

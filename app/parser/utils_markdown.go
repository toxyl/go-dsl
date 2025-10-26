package main

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func (dsl *dslCollection) renderMarkdownToTerminal(markdown string) string {
	// Create a custom renderer with our theme colors
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		return fmt.Sprintf("Error creating renderer: %v", err)
	}

	out, err := renderer.Render(markdown)
	if err != nil {
		return fmt.Sprintf("Error rendering markdown: %v", err)
	}

	return out
}

func (dsl *dslCollection) renderMarkdownToHTML(markdown string) string {
	// Convert markdown to HTML using goldmark
	var buf bytes.Buffer
	parser := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			extension.Typographer,
		),
	)
	if err := parser.Convert([]byte(markdown), &buf); err != nil {
		return fmt.Sprintf("<div class=\"error\">Error converting markdown to HTML: %v</div>", err)
	}
	convertedHTML := buf.String()

	// Add language class to all code blocks
	convertedHTML = strings.ReplaceAll(convertedHTML, "<code>", fmt.Sprintf("<code class=\"language-%s\">", dsl.id))
	convertedHTML = strings.ReplaceAll(convertedHTML, "<pre><code", fmt.Sprintf("<pre><code class=\"language-%s\"", dsl.id))

	type TemplateData struct {
		Name    string
		ID      string
		Content template.HTML
		Theme   *dslColorTheme
	}

	data := TemplateData{
		Name:    dsl.name,
		ID:      dsl.id,
		Content: template.HTML(convertedHTML),
		Theme:   dsl.theme,
	}

	tmpl, err := template.ParseFS(dslTemplates, "template_html.tmpl")
	if err != nil {
		return fmt.Sprintf("<div class=\"error\">Error parsing template: %v</div>", err)
	}

	var htmlBuf bytes.Buffer
	if err := tmpl.Execute(&htmlBuf, data); err != nil {
		return fmt.Sprintf("<div class=\"error\">Error executing template: %v</div>", err)
	}

	return htmlBuf.String()
}

package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var placeholderPattern = regexp.MustCompile(`\{\{[A-Z0-9_]+\}\}`)

// RenderDocxTemplate reads a .docx file (a template containing {{PLACEHOLDER}}
// tokens in word/document.xml), replaces every token with the matching value
// from `values`, and returns the resulting .docx as bytes. Every other part
// of the .docx (styles, headers, media, etc.) is copied through unchanged,
// so the template's formatting/logo/images are preserved untouched.
//
// A placeholder with no matching key in `values` is left as-is in the
// output rather than silently blanked - that makes a template/code mismatch
// obvious when someone opens the generated document, instead of a quiet
// missing field.
func RenderDocxTemplate(templatePath string, values map[string]string) ([]byte, error) {
	r, err := zip.OpenReader(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open docx template at %s: %w", templatePath, err)
	}
	defer r.Close()

	var out bytes.Buffer
	w := zip.NewWriter(&out)

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open %s in template: %w", f.Name, err)
		}
		content, readErr := io.ReadAll(rc)
		rc.Close()
		if readErr != nil {
			return nil, fmt.Errorf("failed to read %s in template: %w", f.Name, readErr)
		}

		if f.Name == "word/document.xml" {
			content = replacePlaceholders(content, values)
		}

		fw, err := w.CreateHeader(&zip.FileHeader{
			Name:   f.Name,
			Method: f.Method,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to write %s to output docx: %w", f.Name, err)
		}
		if _, err := fw.Write(content); err != nil {
			return nil, fmt.Errorf("failed to write %s content: %w", f.Name, err)
		}
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to finalize generated docx: %w", err)
	}
	return out.Bytes(), nil
}

func replacePlaceholders(xmlContent []byte, values map[string]string) []byte {
	return placeholderPattern.ReplaceAllFunc(xmlContent, func(match []byte) []byte {
		key := strings.TrimSuffix(strings.TrimPrefix(string(match), "{{"), "}}")
		val, ok := values[key]
		if !ok {
			return match
		}
		return []byte(escapeXMLText(val))
	})
}

// escapeXMLText escapes text for safe insertion into a Word document.xml
// text node, and turns "\n" into a real Word line break (<w:br/>) so
// multi-line values (like the document checklist) render as separate
// lines in Word instead of a literal backslash-n.
func escapeXMLText(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '\n':
			b.WriteString(`</w:t></w:r><w:r><w:br/><w:t xml:space="preserve">`)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// InsertImageIntoDocx takes docx bytes (already text-rendered via
// RenderDocxTemplate, still containing a literal {{SIGNATURE_IMAGE}} token
// in word/document.xml) and an image, and returns new docx bytes with that
// token replaced by an inline image. widthEMU/heightEMU control the
// rendered size (914400 EMU = 1 inch).
func InsertImageIntoDocx(docxBytes []byte, imageBytes []byte, imageExt string, widthEMU, heightEMU int64) ([]byte, error) {
	r, err := zip.NewReader(bytes.NewReader(docxBytes), int64(len(docxBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to open docx: %w", err)
	}

	var documentXML, relsXML, contentTypesXML []byte
	otherFiles := map[string][]byte{}

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		content, readErr := io.ReadAll(rc)
		rc.Close()
		if readErr != nil {
			return nil, readErr
		}
		switch f.Name {
		case "word/document.xml":
			documentXML = content
		case "word/_rels/document.xml.rels":
			relsXML = content
		case "[Content_Types].xml":
			contentTypesXML = content
		default:
			otherFiles[f.Name] = content
		}
	}
	if documentXML == nil || relsXML == nil || contentTypesXML == nil {
		return nil, fmt.Errorf("docx is missing document.xml, its rels, or content types")
	}

	const relID = "rIdSignatureImage"
	relTarget := "media/signature." + imageExt
	mediaZipPath := "word/" + relTarget

	relEntry := fmt.Sprintf(
		`<Relationship Id="%s" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="%s"/>`,
		relID, relTarget,
	)
	relsStr := strings.Replace(string(relsXML), "</Relationships>", relEntry+"</Relationships>", 1)

	ctStr := string(contentTypesXML)
	if !strings.Contains(ctStr, `Extension="`+imageExt+`"`) {
		contentType := "image/png"
		if imageExt == "jpg" || imageExt == "jpeg" {
			contentType = "image/jpeg"
		}
		def := fmt.Sprintf(`<Default Extension="%s" ContentType="%s"/>`, imageExt, contentType)
		ctStr = strings.Replace(ctStr, "</Types>", def+"</Types>", 1)
	}

	drawingXML := fmt.Sprintf(
		`<w:r><w:drawing><wp:inline xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" distT="0" distB="0" distL="0" distR="0">`+
			`<wp:extent cx="%d" cy="%d"/><wp:effectExtent l="0" t="0" r="0" b="0"/>`+
			`<wp:docPr id="1" name="Signature"/>`+
			`<wp:cNvGraphicFramePr><a:graphicFrameLocks xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" noChangeAspect="1"/></wp:cNvGraphicFramePr>`+
			`<a:graphic xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">`+
			`<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture">`+
			`<pic:pic xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">`+
			`<pic:nvPicPr><pic:cNvPr id="0" name="Signature"/><pic:cNvPicPr/></pic:nvPicPr>`+
			`<pic:blipFill><a:blip r:embed="%s" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/><a:stretch><a:fillRect/></a:stretch></pic:blipFill>`+
			`<pic:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="%d" cy="%d"/></a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></pic:spPr>`+
			`</pic:pic></a:graphicData></a:graphic></wp:inline></w:drawing></w:r>`,
		widthEMU, heightEMU, relID, widthEMU, heightEMU,
	)

	docStr := string(documentXML)
	if !strings.Contains(docStr, "{{SIGNATURE_IMAGE}}") {
		return nil, fmt.Errorf("template does not contain the {{SIGNATURE_IMAGE}} token")
	}
	// IMPORTANT: match EVERY run one at a time (lazy .*? stops at the
	// nearest </w:r>), then only replace the ONE specific run whose
	// content actually contains the token. Using ReplaceAllString (no
	// per-match check) would replace EVERY run in the document with the
	// drawing, since runPattern legitimately matches every <w:r> in the
	// file - not just the one we want.
	runPattern := regexp.MustCompile(`<w:r(?:\s[^>]*)?>(?:(?s).)*?</w:r>`)
	found := false
	docStr = runPattern.ReplaceAllStringFunc(docStr, func(runXML string) string {
		if strings.Contains(runXML, "{{SIGNATURE_IMAGE}}") {
			found = true
			return drawingXML
		}
		return runXML
	})
	if !found {
		return nil, fmt.Errorf("could not locate the signature run in the template")
	}

	var out bytes.Buffer
	w := zip.NewWriter(&out)
	write := func(name string, content []byte) error {
		fw, err := w.Create(name)
		if err != nil {
			return err
		}
		_, err = fw.Write(content)
		return err
	}
	if err := write("word/document.xml", []byte(docStr)); err != nil {
		return nil, err
	}
	if err := write("word/_rels/document.xml.rels", []byte(relsStr)); err != nil {
		return nil, err
	}
	if err := write("[Content_Types].xml", []byte(ctStr)); err != nil {
		return nil, err
	}
	if err := write(mediaZipPath, imageBytes); err != nil {
		return nil, err
	}
	for name, content := range otherFiles {
		if err := write(name, content); err != nil {
			return nil, err
		}
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

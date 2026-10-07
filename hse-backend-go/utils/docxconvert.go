package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ConvertDocxToPDF shells out to LibreOffice (headless) to convert a .docx
// to .pdf. This is the most reliable way to get an accurate, Word-faithful
// PDF - there's no pure-Go library that reproduces .docx layout correctly,
// so we lean on a real office suite installed in the Docker image
// (see Dockerfile: libreoffice-writer).
func ConvertDocxToPDF(docxBytes []byte) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "docx-convert-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.docx")
	if err := os.WriteFile(inputPath, docxBytes, 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp docx: %w", err)
	}

	// A dedicated -env:UserInstallation profile per conversion avoids two
	// approvals happening close together from fighting over the same
	// LibreOffice user-profile lock.
	userProfileDir := "file://" + filepath.Join(tmpDir, "lo-profile")
	cmd := exec.Command("soffice",
		"--headless", "--norestore", "--nologo", "--nofirststartwizard",
		"-env:UserInstallation="+userProfileDir,
		"--convert-to", "pdf", "--outdir", tmpDir, inputPath,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("libreoffice conversion failed: %w (output: %s)", err, string(output))
	}

	pdfPath := filepath.Join(tmpDir, "input.pdf")
	pdfBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("converted pdf not found (libreoffice output: %s): %w", string(output), err)
	}
	return pdfBytes, nil
}

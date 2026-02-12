package fileparser

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

// Supported file types
const (
	FileTypePDF  = "pdf"
	FileTypeDOCX = "docx"
	FileTypeDOC  = "doc"
	FileTypeTXT  = "txt"
)

// MaxFileSize is 10MB
const MaxFileSize = 10 * 1024 * 1024

var (
	ErrUnsupportedFileType = errors.New("unsupported file type")
	ErrFileTooLarge        = errors.New("file too large")
	ErrEmptyFile           = errors.New("file is empty")
	ErrInvalidFile         = errors.New("invalid file format")
)

// FileParserService handles parsing of various file formats
type FileParserService struct{}

// NewFileParserService creates a new file parser service
func NewFileParserService() *FileParserService {
	return &FileParserService{}
}

// ParseFile extracts text content from a file based on its type
func (s *FileParserService) ParseFile(data []byte, filename string) (string, error) {
	if len(data) == 0 {
		return "", ErrEmptyFile
	}

	if len(data) > MaxFileSize {
		return "", ErrFileTooLarge
	}

	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		return s.parsePDF(data)
	case ".docx":
		return s.parseDOCX(data)
	case ".doc":
		return s.parseDOC(data)
	case ".txt":
		return s.parseTXT(data)
	default:
		return "", fmt.Errorf("%w: %s", ErrUnsupportedFileType, ext)
	}
}

// GetSupportedExtensions returns a list of supported file extensions
func (s *FileParserService) GetSupportedExtensions() []string {
	return []string{".pdf", ".docx", ".doc", ".txt"}
}

// parsePDF extracts text from a PDF file
func (s *FileParserService) parsePDF(data []byte) (string, error) {
	reader := bytes.NewReader(data)
	pdfReader, err := pdf.NewReader(reader, int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("failed to read PDF: %w", err)
	}

	var textBuilder strings.Builder
	numPages := pdfReader.NumPage()

	for i := 1; i <= numPages; i++ {
		page := pdfReader.Page(i)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			// Try to continue with other pages
			continue
		}

		textBuilder.WriteString(text)
		if i < numPages {
			textBuilder.WriteString("\n\n")
		}
	}

	result := textBuilder.String()
	if strings.TrimSpace(result) == "" {
		return "", fmt.Errorf("could not extract text from PDF - the file may be image-based or protected")
	}

	return s.cleanText(result), nil
}

// parseDOCX extracts text from a DOCX file
func (s *FileParserService) parseDOCX(data []byte) (string, error) {
	reader := bytes.NewReader(data)
	zipReader, err := zip.NewReader(reader, int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("failed to read DOCX: %w", err)
	}

	var documentXML *zip.File
	for _, file := range zipReader.File {
		if file.Name == "word/document.xml" {
			documentXML = file
			break
		}
	}

	if documentXML == nil {
		return "", ErrInvalidFile
	}

	rc, err := documentXML.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open document.xml: %w", err)
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("failed to read document.xml: %w", err)
	}

	text := s.extractTextFromDOCXML(content)
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("could not extract text from DOCX")
	}

	return s.cleanText(text), nil
}

// extractTextFromDOCXML extracts text content from DOCX XML
func (s *FileParserService) extractTextFromDOCXML(xmlContent []byte) string {
	type Text struct {
		Content string `xml:",chardata"`
	}

	type Run struct {
		Text []Text `xml:"t"`
	}

	type Paragraph struct {
		Runs []Run `xml:"r"`
	}

	type Body struct {
		Paragraphs []Paragraph `xml:"p"`
	}

	type Document struct {
		Body Body `xml:"body"`
	}

	var doc Document
	if err := xml.Unmarshal(xmlContent, &doc); err != nil {
		// Fallback: use regex to extract text
		return s.extractTextFromDOCXMLRegex(xmlContent)
	}

	var textBuilder strings.Builder
	for i, para := range doc.Body.Paragraphs {
		var paraText strings.Builder
		for _, run := range para.Runs {
			for _, text := range run.Text {
				paraText.WriteString(text.Content)
			}
		}
		if paraText.Len() > 0 {
			if i > 0 {
				textBuilder.WriteString("\n")
			}
			textBuilder.WriteString(paraText.String())
		}
	}

	return textBuilder.String()
}

// extractTextFromDOCXMLRegex is a fallback method using regex
func (s *FileParserService) extractTextFromDOCXMLRegex(xmlContent []byte) string {
	// Remove all XML tags except text content
	re := regexp.MustCompile(`<w:t[^>]*>([^<]*)</w:t>`)
	matches := re.FindAllSubmatch(xmlContent, -1)

	var textBuilder strings.Builder
	for _, match := range matches {
		if len(match) > 1 {
			textBuilder.Write(match[1])
		}
	}

	// Also try to get paragraph breaks
	content := textBuilder.String()
	// Replace paragraph markers with newlines
	content = regexp.MustCompile(`</w:p>`).ReplaceAllString(string(xmlContent), "\n")

	// Re-extract text after paragraph processing
	re = regexp.MustCompile(`<w:t[^>]*>([^<]*)</w:t>`)
	matches = re.FindAllSubmatch([]byte(content), -1)

	textBuilder.Reset()
	for _, match := range matches {
		if len(match) > 1 {
			textBuilder.Write(match[1])
		}
	}

	return textBuilder.String()
}

// parseDOC handles legacy .doc files
// Note: Full DOC parsing requires complex OLE2 handling
// For basic support, we attempt to extract readable text
func (s *FileParserService) parseDOC(data []byte) (string, error) {
	// DOC files are OLE2 compound documents
	// For simple cases, we can try to extract ASCII text
	// This is a simplified approach that works for many DOC files

	var textBuilder strings.Builder
	inText := false
	textStart := 0

	for i := 0; i < len(data); i++ {
		b := data[i]
		// Look for printable ASCII characters
		if b >= 32 && b <= 126 {
			if !inText {
				inText = true
				textStart = i
			}
		} else if b == '\n' || b == '\r' || b == '\t' {
			// Allow whitespace in text
			continue
		} else {
			if inText {
				// End of text run
				textRun := string(data[textStart:i])
				// Only keep runs that look like real text (at least 3 chars)
				if len(textRun) >= 3 {
					textBuilder.WriteString(textRun)
					textBuilder.WriteString(" ")
				}
				inText = false
			}
		}
	}

	// Handle trailing text
	if inText && len(data)-textStart >= 3 {
		textBuilder.WriteString(string(data[textStart:]))
	}

	result := textBuilder.String()

	// Clean up the result
	result = s.cleanDOCText(result)

	if strings.TrimSpace(result) == "" {
		return "", fmt.Errorf("could not extract text from DOC file - please convert to DOCX or PDF format")
	}

	return result, nil
}

// cleanDOCText cleans extracted DOC text
func (s *FileParserService) cleanDOCText(text string) string {
	// Remove common binary artifacts
	text = regexp.MustCompile(`[\x00-\x1f]+`).ReplaceAllString(text, " ")

	// Remove excessive whitespace
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Try to identify and extract the main document content
	// DOC files often have headers/footers/metadata mixed in
	lines := strings.Split(text, " ")
	var cleanLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip very short fragments that are likely artifacts
		if len(line) > 2 {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, " ")
}

// parseTXT handles plain text files
func (s *FileParserService) parseTXT(data []byte) (string, error) {
	text := string(data)
	if strings.TrimSpace(text) == "" {
		return "", ErrEmptyFile
	}
	return s.cleanText(text), nil
}

// cleanText normalizes and cleans extracted text
func (s *FileParserService) cleanText(text string) string {
	// Normalize line endings
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// Remove null characters
	text = strings.ReplaceAll(text, "\x00", "")

	// Collapse multiple blank lines into at most two
	text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")

	// Remove leading/trailing whitespace from each line
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	text = strings.Join(lines, "\n")

	// Trim overall
	text = strings.TrimSpace(text)

	return text
}

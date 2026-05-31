package fileparser

import (
	"archive/zip"
	"bytes"
	"errors"
	"strings"
	"testing"
)

func newDOCX(t *testing.T, documentXML string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create("word/document.xml")
	if err != nil {
		t.Fatalf("failed to create zip entry: %v", err)
	}
	if _, err := w.Write([]byte(documentXML)); err != nil {
		t.Fatalf("failed to write zip entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("failed to close zip: %v", err)
	}
	return buf.Bytes()
}

func TestParseFileEmpty(t *testing.T) {
	s := NewFileParserService()
	if _, err := s.ParseFile(nil, "resume.txt"); !errors.Is(err, ErrEmptyFile) {
		t.Errorf("expected ErrEmptyFile, got %v", err)
	}
}

func TestParseFileTooLarge(t *testing.T) {
	s := NewFileParserService()
	big := make([]byte, MaxFileSize+1)
	if _, err := s.ParseFile(big, "resume.txt"); !errors.Is(err, ErrFileTooLarge) {
		t.Errorf("expected ErrFileTooLarge, got %v", err)
	}
}

func TestParseFileUnsupportedType(t *testing.T) {
	s := NewFileParserService()
	if _, err := s.ParseFile([]byte("hello"), "resume.rtf"); !errors.Is(err, ErrUnsupportedFileType) {
		t.Errorf("expected ErrUnsupportedFileType, got %v", err)
	}
}

func TestParseFileExtensionIsCaseInsensitive(t *testing.T) {
	s := NewFileParserService()
	out, err := s.ParseFile([]byte("Hello World"), "RESUME.TXT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "Hello World" {
		t.Errorf("got %q, want %q", out, "Hello World")
	}
}

func TestParseTXTCleansText(t *testing.T) {
	s := NewFileParserService()
	// CRLF line endings, trailing spaces, and 4 blank lines should be normalized.
	in := "  John Doe  \r\n\r\n\r\n\r\nSenior Engineer\r\n"
	out, err := s.ParseFile([]byte(in), "resume.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "John Doe\n\nSenior Engineer"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestParseTXTWhitespaceOnly(t *testing.T) {
	s := NewFileParserService()
	if _, err := s.ParseFile([]byte("   \n\t  "), "resume.txt"); !errors.Is(err, ErrEmptyFile) {
		t.Errorf("expected ErrEmptyFile, got %v", err)
	}
}

func TestParseDOCX(t *testing.T) {
	s := NewFileParserService()
	docXML := `<?xml version="1.0"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p><w:r><w:t>Jane Smith</w:t></w:r></w:p>
    <w:p><w:r><w:t>Product </w:t></w:r><w:r><w:t>Manager</w:t></w:r></w:p>
  </w:body>
</w:document>`
	data := newDOCX(t, docXML)

	out, err := s.ParseFile(data, "resume.docx")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Jane Smith") {
		t.Errorf("expected output to contain name, got %q", out)
	}
	if !strings.Contains(out, "Product Manager") {
		t.Errorf("expected runs to be joined into 'Product Manager', got %q", out)
	}
}

func TestParseDOCXMissingDocumentXML(t *testing.T) {
	s := NewFileParserService()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("word/other.xml")
	_, _ = w.Write([]byte("nope"))
	_ = zw.Close()

	if _, err := s.ParseFile(buf.Bytes(), "resume.docx"); !errors.Is(err, ErrInvalidFile) {
		t.Errorf("expected ErrInvalidFile, got %v", err)
	}
}

func TestGetSupportedExtensions(t *testing.T) {
	s := NewFileParserService()
	got := s.GetSupportedExtensions()
	want := map[string]bool{".pdf": true, ".docx": true, ".doc": true, ".txt": true}
	if len(got) != len(want) {
		t.Fatalf("got %d extensions, want %d: %v", len(got), len(want), got)
	}
	for _, ext := range got {
		if !want[ext] {
			t.Errorf("unexpected extension %q", ext)
		}
	}
}

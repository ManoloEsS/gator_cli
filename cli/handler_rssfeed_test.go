package cli

import (
	"testing"
	"time"
)

func TestStripHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain text",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "simple HTML tags",
			input:    "<p>Hello World</p>",
			expected: "Hello World",
		},
		{
			name:     "nested HTML tags",
			input:    "<div><p>Hello <strong>World</strong></p></div>",
			expected: "Hello World",
		},
		{
			name:     "HTML with attributes",
			input:    `<a href="https://example.com">Click here</a>`,
			expected: "Click here",
		},
		{
			name:     "multiple paragraphs",
			input:    "<p>First paragraph</p><p>Second paragraph</p>",
			expected: "First paragraphSecond paragraph",
		},
		{
			name:     "HTML with line breaks",
			input:    "Hello<br/>World",
			expected: "HelloWorld",
		},
		{
			name:     "HTML entities get unescaped",
			input:    "<p>Test &amp; Example</p>",
			expected: "Test & Example",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "HTML with scripts (script content is extracted)",
			input:    "<script>alert('test')</script>Content",
			expected: "alert('test')Content",
		},
		{
			name:     "HTML with styles (style content is extracted)",
			input:    "<style>body { color: red; }</style>Content",
			expected: "body { color: red; }Content",
		},
		{
			name:     "complex real-world example",
			input:    `<div class="content"><h1>Title</h1><p>This is a <strong>test</strong> paragraph with <a href="#">a link</a>.</p></div>`,
			expected: "TitleThis is a test paragraph with a link.",
		},
		{
			name:     "self-closing tags",
			input:    "<img src='test.jpg' />Text<br />More text",
			expected: "TextMore text",
		},
		{
			name:     "HTML comments",
			input:    "Before<!-- comment -->After",
			expected: "BeforeAfter",
		},
		{
			name:     "whitespace handling",
			input:    "  <p>  Text with spaces  </p>  ",
			expected: "Text with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripHTML(tt.input)
			if result != tt.expected {
				t.Errorf("StripHTML(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParsePubDate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectOk  bool
	}{
		{
			name:     "RFC1123 format",
			input:    "Mon, 02 Jan 2006 15:04:05 MST",
			expectOk: true,
		},
		{
			name:     "RFC1123Z format",
			input:    "Mon, 02 Jan 2006 15:04:05 -0700",
			expectOk: true,
		},
		{
			name:     "RFC822 format",
			input:    "02 Jan 06 15:04 MST",
			expectOk: true,
		},
		{
			name:     "RFC822Z format",
			input:    "02 Jan 06 15:04 -0700",
			expectOk: true,
		},
		{
			name:     "RFC3339 format",
			input:    "2006-01-02T15:04:05Z",
			expectOk: true,
		},
		{
			name:     "custom format without seconds",
			input:    "Mon, 02 Jan 2006 15:04 MST",
			expectOk: true,
		},
		{
			name:     "custom format with explicit offset",
			input:    "Mon, 02 Jan 2006 15:04:05 -0700",
			expectOk: true,
		},
		{
			name:     "invalid date format",
			input:    "not a date",
			expectOk: false,
		},
		{
			name:     "empty string",
			input:    "",
			expectOk: false,
		},
		{
			name:     "string with whitespace",
			input:    "  Mon, 02 Jan 2006 15:04:05 MST  ",
			expectOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := parsePubDate(tt.input)
			if ok != tt.expectOk {
				t.Errorf("parsePubDate(%q) ok = %v, want %v", tt.input, ok, tt.expectOk)
			}
			if ok && result.IsZero() {
				t.Errorf("parsePubDate(%q) returned zero time but ok=true", tt.input)
			}
		})
	}
}

func TestNewNullTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		isZero   bool
		wantNull bool
	}{
		{
			name:     "valid time",
			input:    "2006-01-02T15:04:05Z",
			isZero:   false,
			wantNull: false,
		},
		{
			name:     "zero time",
			input:    "",
			isZero:   true,
			wantNull: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputTime time.Time
			if !tt.isZero {
				parsed, err := time.Parse(time.RFC3339, tt.input)
				if err != nil {
					t.Fatalf("failed to parse test input time: %v", err)
				}
				inputTime = parsed
			}

			result := newNullTime(inputTime)

			if tt.wantNull {
				if result.Valid {
					t.Error("expected Valid to be false for zero time")
				}
			} else {
				if !result.Valid {
					t.Error("expected Valid to be true for non-zero time")
				}
				if result.Time.IsZero() {
					t.Error("expected Time to be non-zero")
				}
			}
		})
	}
}

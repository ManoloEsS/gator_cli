package cli

import (
	"testing"
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

package color

import "testing"

func TestColorString(t *testing.T) {
	tests := []struct {
		name     string
		color    Color
		input    string
		expected string
	}{
		{"Blue color", Blue, "test", "\033[34mtest\033[0m"},
		{"Red color", Red, "test", "\033[31mtest\033[0m"},
		{"Empty string", Blue, "", "\033[34m\033[0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ColorString(tt.color, tt.input)
			if result != tt.expected {
				t.Errorf("ColorString(%v, %v) = %v; want %v", tt.color, tt.input, result, tt.expected)
			}
		})
	}
}

func TestAnsiLink(t *testing.T) {
	tests := []struct {
		name     string
		href     string
		input    string
		expected string
	}{
		{"Valid link", "http://example.com", "example", "\033]8;;http://example.com\033\\example\033]8;;\033\\"},
		{"Empty href", "", "example", "\033]8;;\033\\example\033]8;;\033\\"},
		{"Empty string", "http://example.com", "", "\033]8;;http://example.com\033\\\033]8;;\033\\"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AnsiLink(tt.href, tt.input)
			if result != tt.expected {
				t.Errorf("AnsiLink(%v, %v) = %v; want %v", tt.href, tt.input, result, tt.expected)
			}
		})
	}
}

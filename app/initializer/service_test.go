package initializer

import (
	"errors"
	"testing"
)

func TestValidateProjectName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantError error
	}{
		{
			name:      "Empty name",
			input:     "",
			wantError: errors.New("project name cannot be empty"),
		},
		{
			name:      "Too short",
			input:     "ab",
			wantError: errors.New("project name must be at least 3 characters"),
		},
		{
			name:      "Too long",
			input:     string(make([]byte, 51)),
			wantError: errors.New("project name must be less than 50 characters"),
		},
		{
			name:      "Invalid characters",
			input:     "bad$name",
			wantError: errors.New("project name can only contain letters, numbers, hyphens, and underscores"),
		},
		{
			name:      "Valid name - simple",
			input:     "my_project",
			wantError: nil,
		},
		{
			name:      "Valid name - hyphen",
			input:     "my-project",
			wantError: nil,
		},
		{
			name:      "Valid name - alphanumeric",
			input:     "proj123",
			wantError: nil,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateProjectName(tt.input)
			if tt.wantError == nil && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
			if tt.wantError != nil && (err == nil || err.Error() != tt.wantError.Error()) {
				t.Errorf("expected error %q, got: %v", tt.wantError, err)
			}
		})
	}
}


func TestSanitizeProjectName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Normal name", "MyProject", "myproject"},
		{"Trim whitespace", "  My Project  ", "my-project"},
		{"Multiple spaces", "My   Project", "my-project"},
		{"Underscores collapse", "My__Project", "my-project"},
		{"Special characters", "Hello!@#$%^&*()World", "hello-world"},
		{"Only symbols", "!!!@@@", ""},
		{"Only underscores", "____", ""},
		{"Only dashes", "----", ""},
		{"Dashes and underscores", "a--__--b", "a-b"},
		{"Ends with dash/underscore", "end--_", "end"},
		{"Mixed uppercase", "Go_Initializr!", "go-initializr"},
		{"Numbers in name", "Proj123", "proj123"},
		{"Empty string", "", ""},
		{"Already clean", "clean-name", "clean-name"},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := sanitizeProjectName(tt.input)
			if output != tt.expected {
				t.Errorf("sanitizeProjectName(%q) = %q; want %q", tt.input, output, tt.expected)
			}
		})
	}
}

func TestValidateModuleName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantError error
	}{
		{"Empty name", "", errors.New("module name cannot be empty")},
		{"Contains space", "github.com user/project", errors.New("module name must not contain spaces")},
		{"Starts with dash", "-github.com/user", errors.New("module name must start with a lowercase letter or number")},
		{"Starts with uppercase", "Github.com/user", errors.New("module name must start with a lowercase letter or number")},
		{"Contains invalid chars", "github.com/user@project", errors.New("module name contains invalid characters")},
		{"Valid with slashes", "github.com/user/project", nil},
		{"Valid with dash", "example.com/user-name", nil},
		{"Valid with dots", "gopkg.in/yaml.v2", nil},
		{"Only lowercase", "myproject", nil},
		{"Numbers okay", "123project", nil},
		{"Mixed valid chars", "domain.com/user-1/pkg.v2", nil},
	}

	for _, tt := range tests {
		tt := tt // capture
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateModuleName(tt.input)

			if tt.wantError == nil && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.wantError != nil && (err == nil || err.Error() != tt.wantError.Error()) {
				t.Errorf("expected error: %q, got: %v", tt.wantError, err)
			}
		})
	}
}


func TestSanitizeModuleName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Whitespace", "   github.com/user/project   ", "github.com/user/project"},
		{"Uppercase", "GITHUB.Com/USER/Project", "github.com/user/project"},
		{"Spaces to dash", "github com user project", "github-com-user-project"},
		{"Invalid chars removed", "github.com/user@project!", "github.com/user-project"},
		{"Multiple dashes", "user---project", "user-project"},
		{"Multiple slashes", "github.com//user///project", "github.com/user/project"},
		{"Trim leading/trailing", "--/github.com/user/project/--", "github.com/user/project"},
		{"Dots preserved", "gopkg.in/yaml.v2", "gopkg.in/yaml.v2"},
		{"Dash + slash combo", "-/user--name///pkg..v2", "user-name/pkg..v2"},
		{"Only invalid chars", "@@@@!!!", ""},
		{"Complex mix", "  GitHub.COM//My_Project!@# v2 ", "github.com/my-project-v2"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := sanitizeModuleName(tt.input)
			if got != tt.expected {
				t.Errorf("sanitizeModuleName(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
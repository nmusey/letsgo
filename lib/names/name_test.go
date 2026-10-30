package names_test

import (
	"testing"

	"github.com/nmusey/letsgo/lib/names"
)

func TestStructName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"single lowercase word", "user", "User"},
		{"already capitalised", "User", "User"},
		{"space separated", "user profile", "UserProfile"},
		{"snake case", "user_profile", "UserProfile"},
		{"kebab case", "user-profile", "UserProfile"},
		{"camel case is preserved", "userProfile", "UserProfile"},
		{"pascal case is preserved", "UserProfile", "UserProfile"},
		{"acronyms are not split", "HTTPServer", "HTTPServer"},
		{"digits are kept", "oauth2 token", "Oauth2Token"},
		{"repeated and leading separators", "__user -- profile__", "UserProfile"},
		{"only separators", "---", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := names.StructName(tt.input); got != tt.want {
				t.Errorf("StructName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestPackageName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"single word", "user", "user"},
		{"uppercase is lowered", "User", "user"},
		{"space separated", "User Profile", "userprofile"},
		{"snake case", "user_profile", "userprofile"},
		{"kebab case", "user-profile", "userprofile"},
		{"pascal case", "UserProfile", "userprofile"},
		{"digits are kept", "OAuth2", "oauth2"},
		{"only separators", "-_ ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := names.PackageName(tt.input); got != tt.want {
				t.Errorf("PackageName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

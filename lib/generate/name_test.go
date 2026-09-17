package generate_test

import (
	"testing"

	"github.com/nmusey/letsgo/lib/generate"
)

func TestNewName(t *testing.T) {
	cases := []struct {
		raw         string
		wantStruct  string
		wantPackage string
	}{
		{"BlogPost", "BlogPost", "blogpost"},
		{"blog_post", "BlogPost", "blogpost"},
		{"blog-post", "BlogPost", "blogpost"},
		{"blog post", "BlogPost", "blogpost"},
		{"blogPost", "BlogPost", "blogpost"},
		{"User", "User", "user"},
		{"user", "User", "user"},
	}

	for _, c := range cases {
		t.Run(c.raw, func(t *testing.T) {
			name := generate.NewName(c.raw)

			if name.Struct != c.wantStruct {
				t.Errorf("NewName(%q).Struct = %q, want %q", c.raw, name.Struct, c.wantStruct)
			}
			if name.Package != c.wantPackage {
				t.Errorf("NewName(%q).Package = %q, want %q", c.raw, name.Package, c.wantPackage)
			}
		})
	}
}

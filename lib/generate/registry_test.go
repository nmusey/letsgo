package generate_test

import (
	"testing"

	"github.com/nmusey/letsgo/lib/generate"
)

func TestKindByName_Found(t *testing.T) {
	kind, ok := generate.KindByName("repository")
	if !ok {
		t.Fatal(`KindByName("repository") ok = false, want true`)
	}

	if kind.OutputFile != "repository.go" {
		t.Errorf("kind.OutputFile = %q, want %q", kind.OutputFile, "repository.go")
	}
	if kind.FlagShort != "r" {
		t.Errorf("kind.FlagShort = %q, want %q", kind.FlagShort, "r")
	}
}

func TestKindByName_NotFound(t *testing.T) {
	if _, ok := generate.KindByName("controller"); ok {
		t.Fatal(`KindByName("controller") ok = true, want false`)
	}
}

func TestKinds_ContainsModelRepositoryServiceInOrder(t *testing.T) {
	var names []string
	for _, kind := range generate.Kinds {
		names = append(names, kind.Name)
	}

	want := []string{"model", "repository", "service"}
	if len(names) != len(want) {
		t.Fatalf("generate.Kinds has %d entries, want %d", len(names), len(want))
	}
	for i, name := range names {
		if name != want[i] {
			t.Errorf("generate.Kinds[%d].Name = %q, want %q", i, name, want[i])
		}
	}
}

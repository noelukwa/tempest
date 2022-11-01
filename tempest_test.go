package tempest_test

import (
	"strings"
	"testing"

	"github.com/noelukwa/tempest"
	"github.com/noelukwa/tempest/test"
)

func TestLoadFs(t *testing.T) {
	temps := tempest.LoadFs(test.Files, tempest.Config{})

	if len(temps) != 2 {
		t.Errorf("expected 2 templates, got %d", len(temps))
	}

	if _, ok := temps["index"]; !ok {
		t.Error("expected index template to be loaded")
	}

	if _, ok := temps["about"]; !ok {
		t.Error("expected about template to be loaded")
	}

	{
		// Test index template
		tmpl := temps["index"]

		// initialise content to be type strings Builder and load content of temp into it, then compare
		content := strings.Builder{}
		tmpl.Execute(&content, nil)

		if content.String() != "main-header index main-footer" {
			t.Errorf("expected index template to be loaded, got %s", content.String())
		}

	}
}

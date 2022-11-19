package tempest_test

import (
	"strings"
	"testing"

	"github.com/noelukwa/tempest"
	"github.com/noelukwa/tempest/test"
)

func TestLoadFs(t *testing.T) {
	tempst := tempest.New()

	temps, err := tempst.LoadFS(test.Files)
	if err != nil {
		t.Fatal(err)
	}

	if len(temps) != 3 {
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
			t.Errorf("expected : %s got : %s", "main-header index main-footer", content.String())
		}
	}
	{
		// Test about template
		tmpl := temps["about"]

		// initialise content to be type strings Builder and load content of temp into it, then compare
		content := strings.Builder{}
		tmpl.Execute(&content, nil)

		if content.String() != "main-header about main-footer" {
			t.Errorf("expected : %s got : %s", "main-header about main-footer", content.String())
		}

	}
	{
		// Test nested template
		tmpl := temps["admin/dash"]

		// initialise content to be type strings Builder and load content of temp into it, then compare
		content := strings.Builder{}
		tmpl.Execute(&content, nil)

		if content.String() != "main-header admin-layout admin-dash main-footer" {
			t.Errorf("expected : %s got : %s", "main-header admin-layout admin-dash main-footer", content.String())
		}
	}
}

func TestLoadWithConfig(t *testing.T) {
	tempest := tempest.WithConfig(&tempest.Config{
		Layout:      "base",
		IncludesDir: "partials",
	})
	temps, err := tempest.LoadFS(test.SpecialFiles)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := temps["home"]; !ok {
		t.Error("expected home template to be loaded")
	}

	{
		// Test index template
		tmpl := temps["home"]

		// initialise content to be type strings Builder and load content of temp into it, then compare
		content := strings.Builder{}
		tmpl.Execute(&content, nil)

		if content.String() != "header home-page footer" {
			t.Errorf("expected : %s got : %s", "header home-page footer", content.String())
		}
	}

}

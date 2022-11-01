package tempest

import (
	"html/template"
	"io/fs"
)

type Config struct {
	PartialDir string

	BaseLayout string

	FuncMap template.FuncMap
}

type templateInfo struct {
	baseLayouts string

	partials []string

	main string
}

func getConfig(conf Config) (partials, baseLayout string, funcMap template.FuncMap) {

	if conf.PartialDir == "" {
		conf.PartialDir = "partials"
	}

	if conf.BaseLayout == "" {
		conf.BaseLayout = "base"
	}
	if conf.FuncMap == nil {
		conf.FuncMap = template.FuncMap{}
	}

	return conf.PartialDir, conf.BaseLayout, conf.FuncMap
}

func LoadFs(files fs.FS, conf Config) (templates map[string]*template.Template) {
	// partials, baseLayout, funcMap := getConfig(conf)

	templates = map[string]*template.Template{}
	// Load partials
	walkFn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		return nil
	}

	fs.WalkDir(files, ".", walkFn)

	return
}

// Copyright (c) 2022 Noel Ukwa. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

// Package tempest is a simple template grouping helper for GO.
//
// It is intended to be used in conjuction with go html/template package,
// and not a replacement for the html/template package.
//
// A simple example:
//
// 	package main
//
// 	import (
// 		"fmt"
// 		"html/template"
// 		"log"
// 		"net/http"
// 		"os"
//
// 		"github.com/noelukwa/tempest"
// 	)
//
// 	//go:embed views
// 	var views embed.FS
//
// 	func main() {
// 		t := tempest.New()
//
// 		templates, err := t.LoadFS(os.DirFS("views"))
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 			err := templates["index"].Execute(w, nil)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		})
//
// 		http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
// 			err := templates["admin/index"].Execute(w, nil)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		})
//
// 		fmt.Println("Listening on port 8080")
// 		log.Fatal(http.ListenAndServe(":8080", nil))
// 	}
//
// The above example will load all the templates in the views directory and
// subdirectories, and will group them into a map of templates to filenames.
//

package tempest

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

// Config is the configuration for the tempest instance.
// It is used to set the file extension, the directory where the includes are
// stored, and the name used for layout templates.
// Defaults are set for each of these fields if tempest is initialized
// without a config.
type Config struct {
	// The file extension of the templates.
	// Defaults to ".html".
	Ext string

	// The directory where the includes are stored.
	// Defaults to "includes".
	IncludesDir string

	// The name used for layout templates :- templates that wrap other contents.
	// Defaults to "layouts".
	Layout string
}

type Tempest struct {
	temps map[string]*template.Template
	conf  *Config
}

// New returns a new tempest instance with default configuration.
func New() *Tempest {
	return &Tempest{
		temps: make(map[string]*template.Template),
		conf: &Config{
			Ext:         ".html",
			IncludesDir: "includes",
			Layout:      "layout",
		},
	}
}

// WithConfig sets the configuration for the tempest instance.
func WithConfig(conf *Config) *Tempest {
	if conf.Ext == "" {
		conf.Ext = ".html"
	}
	if conf.IncludesDir == "" {
		conf.IncludesDir = "includes"
	}
	if conf.Layout == "" {
		conf.Layout = "layout"
	}
	return &Tempest{
		temps: make(map[string]*template.Template),
		conf:  conf,
	}
}

// LoadFS loads templates from an embedded filesystem and returns a map of
// templates to filenames.
func (tempest *Tempest) LoadFS(files fs.FS) (map[string]*template.Template, error) {

	includesDir := filepath.Clean(tempest.conf.IncludesDir)
	layoutFile := filepath.Clean(tempest.conf.Layout + tempest.conf.Ext)

	includes := make([]string, 0)
	layouts := make([]string, 0)
	rawTemps := make(map[string]string)

	templates := make(map[string]*template.Template)

	// Walk through the files and load them into the map
	fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if filepath.Base(path) == includesDir {
				includes = append(includes, path)
				return fs.SkipDir
			}
		}

		if !d.IsDir() {
			if filepath.Ext(path) == tempest.conf.Ext && filepath.Base(path) != layoutFile {
				rawTemps[path] = path
			} else {
				layouts = append(layouts, path)
			}
		}

		return nil
	})

	// sort the includes
	sort.Slice(includes, func(i, j int) bool {
		return len(includes[i]) < len(includes[j])
	})

	// sort the layouts
	sort.Slice(layouts, func(i, j int) bool {
		return len(layouts[i]) < len(layouts[j])
	})

	for _, t := range rawTemps {
		temp := template.New(layoutFile)

		// get the includes
		incls := getInclues(t, includes)
		for _, i := range incls {
			xfiles, err := fs.Glob(files, fmt.Sprintf("%s/*%s", i, tempest.conf.Ext))
			if err != nil {
				return nil, fmt.Errorf("error getting includes: %w", err)
			}
			temp, err = temp.ParseFS(files, xfiles...)
			if err != nil {
				return nil, fmt.Errorf("error parsing includes: %w", err)
			}
		}

		// get the layouts
		lyts := getLayouts(t, layouts)
		lyts = append(lyts, t)

		temp, _ = temp.ParseFS(files, lyts...)

		// remove the extension from the template name
		// views/admin/index.html -> admin/index
		key := strings.TrimPrefix(t, strings.Split(t, "/")[0]+"/")
		key = strings.TrimSuffix(key, filepath.Ext(key))
		templates[key] = temp
	}

	return templates, nil
}

func getInclues(path string, includes []string) []string {
	inc := make([]string, 0)
	for _, i := range includes {
		// fmt.Printf("i: %s\n", i)
		if strings.HasPrefix(path, filepath.Dir(i)) || filepath.Dir(i) == "." {
			inc = append(inc, i)
		}
	}
	return inc
}

func getLayouts(path string, layouts []string) []string {
	lay := make([]string, 0)
	for _, l := range layouts {
		if strings.HasPrefix(path, filepath.Dir(l)) || filepath.Dir(l) == "." {
			lay = append(lay, l)
		}
	}
	return lay
}

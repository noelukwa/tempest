## tempest 
*Made out of neccessity and frustration* 😩

### Lets you
- Use go templates in your app without repeating the parsing logic over and over.
- Use any template supported by go html/template package.
- Use `go:embed` to embed template files in your binary.
- Parse templates once.

### How
In order for tempest to parse templates, three conditions must be met.
1. Templates must be embeded 
2. The name of the template used for layouts should be `layouts.<extention>`, otherwise, it should be stated with custom config.
3. The name of the folder containing partial templates should be "inludes", otherwise, it should be stated with custom config

*📝 For requirements 2 and 3, see [examples/with-conf](https://github.com/noelukwa/tempest)*



**Requirements** 
- Fair knowledge of go [html/template](https://pkg.go.dev/html/template) package.
- Go version >= 1.16


### Example
Lets say you have a folder structure like this
```
.
├── main.go
└── templates
    ├── includes
    │   └── footer.html
    |   └── header.html
    └── admin
    │   └── dash.html
    │   └── layout.html
    ├── layout.html
    ├── index.html
    └── about.html
```

In your main.go file, you can do something like this
```go
package main

import (
    "embed"
    "log"

    "github.com/noelukwa/tempest"
)

var (
    //go:embed templates
    templates embed.FS
)

func main() {
    // Create a new tempest instance
    tempst := tempest.New()

    templs, err := tempst.ParseFS(templates)
    if err != nil {
        log.Fatal(err)
    }

    // Render a template

    mux := http.NewServeMux()

    mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
        // 🚨 Note that the template name is the file name without the extension
        // and the base folder ; in this case "templates"
        dash := templs["admin/dash"]
        dash.Execute(w, nil)
    })
}
```

### Template Directory Parsing
The template files in the `templates` directory above will be grouped as follows

```
- templates/admin/dash.html
    ├── templates/layout.html
    ├── templates/admin/layout.html 
    ├── templates/admin/dash.html 
    ├── templates/includes/footer.html
    └── templates/includes/header.html


- templates/index.html
    ├── templates/layout.html
    ├── templates/index.html 
    ├── templates/includes/footer.html
    └── templates/includes/header.html

- templates/about.html
    ├── templates/layout.html
    ├── templates/about.html 
    ├── templates/includes/footer.html
    └── templates/includes/header.html

```

### html/template basics
*When using  nested layouts, the child layout's `define` block name should correspond to the parent layout's `block` name.*

```html
<!-- templates/layout.html -->
<main>
    {{ block "content" . }}{{ end }}
</main>
```
```html
<!-- templates/admin/layout.html -->
{{ define "content" }}
<section>
    {{ block "admin-content" . }}{{ end }}
</section>
{{ end }}
```

Further Read: [Go html/template package](https://pkg.go.dev/html/template)


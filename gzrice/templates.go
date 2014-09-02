package main

import (
	"fmt"
	"os"
	"text/template"
)

// gzip support @todo Content must be transformed to byte instead of string

var tmplEmbeddedBox *template.Template

func init() {
	var err error

	// parse embedded box template
	tmplEmbeddedBox, err = template.New("embeddedBox").Parse(`package {{.Package}}

import (
	"github.com/SchumacherFM/wanderlust/gzrice/embedded"
	"time"
)

func init() {

	// define files
	{{range .Files}}{{.Identifier}} := &embedded.EmbeddedFile{
		Filename:    ` + "`" + `{{.FileName}}` + "`" + `,
		FileModTime: time.Unix({{.ModTime}}, 0),
		Content:     {{.Content | printf "%#v"}},
		IsGzip:		{{.IsGzip | printf "%#v"}}
	}
	{{end}}

	// define dirs
	{{range .Dirs}}{{.Identifier}} := &embedded.EmbeddedDir{
		Filename:    ` + "`" + `{{.FileName}}` + "`" + `,
		DirModTime: time.Unix({{.ModTime}}, 0),
		ChildFiles:  []*embedded.EmbeddedFile{
			{{range .ChildFiles}}{{.Identifier}}, // {{.FileName}}
			{{end}}
		},
	}
	{{end}}

	// link ChildDirs
	{{range .Dirs}}{{.Identifier}}.ChildDirs = []*embedded.EmbeddedDir{
		{{range .ChildDirs}}{{.Identifier}}, // {{.FileName}}
		{{end}}
	}
	{{end}}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(` + "`" + `{{.BoxName}}` + "`" + `, &embedded.EmbeddedBox{
		Name: ` + "`" + `{{.BoxName}}` + "`" + `,
		Time: time.Unix({{.UnixNow}}, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			{{range .Dirs}}"{{.FileName}}": {{.Identifier}},
			{{end}}
		},
		Files: map[string]*embedded.EmbeddedFile{
			{{range .Files}}"{{.FileName}}": {{.Identifier}},
			{{end}}
		},
	})
}`)
	if err != nil {
		fmt.Printf("error parsing embedded box template: %s\n", err)
		os.Exit(-1)
	}
}

type boxDataType struct {
	Package string
	BoxName string
	UnixNow int64
	Files   []*fileDataType
	Dirs    map[string]*dirDataType
}

type fileDataType struct {
	Identifier string
	FileName   string
	Content    []byte
	ModTime    int64
	IsGzip     bool
}

type dirDataType struct {
	Identifier string
	FileName   string
	Content    []byte
	ModTime    int64
	ChildDirs  []*dirDataType
	ChildFiles []*fileDataType
}
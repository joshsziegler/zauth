package httpserver

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gobuffalo/packr"
)

// templateHelpers allows us to use these custom functions in our templates
var templateHelpers = template.FuncMap{
	"inc":      func(i int) int { return i + 1 },
	"multiply": func(x int, y int) int { return x * y },
	"marshal": func(v interface{}) template.JS {
		a, _ := json.Marshal(v)
		return template.JS(a)
	},
	"FormatTimeAsRFC822": func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05 MST")
	},
	"HumanizeTime": humanize.Time,
	"ToLower":      strings.ToLower,
}

// MustLoadBoxedTemplates walks through a packr box, loading all templates
// ending in .html
func MustLoadBoxedTemplates(b packr.Box) *template.Template {
	t := template.New("").Funcs(templateHelpers)
	err := b.Walk(func(p string, f packr.File) error {
		if p == "" {
			return nil
		}
		var err error
		finfo, err := f.FileInfo()
		if err != nil {
			return err
		}
		// skip directory path
		if finfo.IsDir() {
			return nil
		}
		// skip all files except .html
		if !strings.HasSuffix(p, ".html") {
			return nil
		}
		// Normalize template name
		n := p
		if strings.HasPrefix(p, "\\") || strings.HasPrefix(p, "/") {
			n = n[1:] // don't want names to start with / ie. /index.html
		}
		// replace windows path seperator \ to normalized /
		n = strings.Replace(n, "\\", "/", -1)
		var h string
		if h, err = b.FindString(p); err != nil {
			return err
		}
		if _, err = t.New(n).Parse(h); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		panic("error loading template")
	}
	return t
}

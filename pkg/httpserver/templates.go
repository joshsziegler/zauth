package httpserver

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
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

// MustLoadTemplates parses all .html templates from the given filesystem,
// naming each by its base name (e.g. login.html).
func MustLoadTemplates(fsys fs.FS, pattern string) *template.Template {
	return template.Must(template.New("").Funcs(templateHelpers).ParseFS(fsys, pattern))
}

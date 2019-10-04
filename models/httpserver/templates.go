package httpserver

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/gobuffalo/packr"
)

// MustLoadBoxedTemplates walks through a packr box, loading all templates
// ending in .html
func MustLoadBoxedTemplates(b packr.Box) *template.Template {
	t := template.New("").Funcs(templateHelpers)
	err := b.Walk(func(p string, f packr.File) error {
		if p == "" {
			return nil
		}
		var err error
		var csz int64
		finfo, err := f.FileInfo()
		if err != nil {
			return err
		}

		// skip directory path
		if finfo.IsDir() {
			return nil
		}
		csz = finfo.Size()

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

		var h = make([]byte, 0, csz)

		if h, err = b.MustBytes(p); err != nil {
			return err
		}

		if _, err = t.New(n).Parse(string(h)); err != nil {
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

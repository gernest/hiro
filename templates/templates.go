package templates

import (
	"io"
	"log"
	"sort"
	"text/template"
)

// common keys
const (
	HomeBanner    = "home_banner"
	MisssingToken = "missing_token"
)

var tpl *template.Template

func init() {
	t, err := New()
	if err != nil {
		log.Fatal(err)
	}
	tpl = t
}

// New returns all bq specific templates.
func New() (*template.Template, error) {
	t := template.New("bq").Funcs(funcs())
	for _, v := range AssetNames() {
		data, err := Asset(v)
		if err != nil {
			return nil, err
		}
		_, err = t.New(v).Parse(string(data))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

// Write executes template k and writes the result into o.
func Write(o io.Writer, k string, ctx interface{}) error {
	return tpl.ExecuteTemplate(o, k, ctx)
}

func funcs() template.FuncMap {
	return template.FuncMap{
		"sortedKeys": sortKeys,
	}
}

func sortKeys(m map[string]string) []string {
	var o []string
	for k := range m {
		o = append(o, k)
	}
	sort.Strings(o)
	return o
}

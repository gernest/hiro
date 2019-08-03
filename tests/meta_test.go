package tests

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/gernest/hiro/meta"
	"github.com/gernest/hiro/templates"
)

func TestTwitter(t *testing.T) {
	tw := &meta.Twitter{
		Card:        "summary",
		Site:        "@bqservice",
		Creator:     "@gernest",
		Title:       "high performance qrcode service",
		Description: " Create, scan, manage and integrate qrcode into your business workflow",
	}
	og := &meta.OpenGraph{
		Title:       tw.Title,
		Description: tw.Description,
		Type:        "website",
	}
	var buf bytes.Buffer
	mm := &meta.Meta{OpenGraph: og, Twitter: tw}
	m, err := mm.Map()
	if err != nil {
		t.Fatal(err)
	}
	data := map[string]interface{}{
		"meta": m,
	}
	err = templates.Write(&buf, "html/meta.html", data)
	if err != nil {
		t.Fatal(err)
	}
	// ioutil.WriteFile("fixture/expected_twitter_meta.html", buf.Bytes(), 0600)
	b, err := ioutil.ReadFile("fixture/expected_twitter_meta.html")
	if err != nil {
		t.Fatal(err)
	}
	want := string(b)
	got := buf.String()
	if got != want {
		t.Errorf("expected %s got %s", want, got)
	}
}

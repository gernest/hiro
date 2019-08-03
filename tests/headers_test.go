package tests

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gernest/hiro/headers"
)

func TestHeaders(t *testing.T) {
	e := "image/png"
	g := headers.PNG()
	if g != e {
		t.Errorf("expected %s got %s", e, g)
	}

	e = "image/jpeg"
	g = headers.JPEG()
	if g != e {
		t.Errorf("expected %s got %s", e, g)
	}
}

func TestLastModified(t *testing.T) {
	sample := "Wed, 21 Oct 2015 07:28:00 GMT"
	ts, err := time.Parse(headers.LastModifiedTimeFormat, sample)
	if err != nil {
		t.Fatal(err)
	}
	h := headers.LastModifiedTime(ts)
	var buf bytes.Buffer
	h.Write(&buf)
	e := "Last-Modified: Wed, 21 Oct 2015 07:28:00 GMT"
	g := buf.String()
	g = strings.TrimSpace(g)
	if g != e {
		t.Errorf("expected %s got %s", e, g)
	}
}

func TestMerge(t *testing.T) {
	h1 := make(http.Header)
	h1.Add(headers.ContentType, headers.ApplicationJSON)
	sample := "Wed, 21 Oct 2015 07:28:00 GMT"
	ts, err := time.Parse(headers.LastModifiedTimeFormat, sample)
	if err != nil {
		t.Fatal(err)
	}
	h2 := headers.LastModifiedTime(ts)
	m := headers.Merge(h1, h2)

	buf := bytes.Buffer{}
	err = m.Write(&buf)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadFile("fixture/merge.txt")
	if err != nil {
		t.Fatal(err)
	}
	expect := string(b)
	got := buf.String()
	if got != expect {
		t.Errorf("expected %s got %s", expect, got)
	}
}

func TestIsJSONContent(t *testing.T) {
	h := make(http.Header)
	h.Add(headers.ContentType, headers.ApplicationJSON)
	if !headers.IsJSONContent(h) {
		t.Error("expected to be true")
	}
}
func TestIsForm(t *testing.T) {
	h := make(http.Header)
	h.Add(headers.ContentType, headers.ApplicationForm)
	if !headers.IsForm(h) {
		t.Error("expected to be true")
	}
}
func TestIsMultipartForm(t *testing.T) {
	h := make(http.Header)
	h.Add(headers.ContentType, headers.ApplicationMultipartForm)
	if !headers.IsMultipartForm(h) {
		t.Error("expected to be true")
	}
}

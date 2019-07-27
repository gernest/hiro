package headers

import (
	"mime"
	"net/http"
	"time"
)

const (

	// ContentType is browser content type header.
	ContentType = "Content-Type"

	// TextPlain plain text request body.
	TextPlain = "text/plain"

	// TextHTML html body.
	TextHTML = "text/html"

	// ApplicationForm is for url encoded models.
	ApplicationForm = "application/x-www-form-urlencoded"

	// ApplicationJSON is for json body.
	ApplicationJSON = "application/json"

	// ApplicationMultipartForm is for multipart models.
	ApplicationMultipartForm = "multipart/form-data"

	// LastModified header on response containing the date which the resource was
	// last modified.
	LastModified = "Last-Modified"

	lastModifiedTimeFormat = "Mon, _2 Jan 2006 15:04:05 MST"

	// Authorization header
	Authorization = "Authorization"
)

// PNG returns mime type for png.
func PNG() string {
	return mime.TypeByExtension(".png")
}

// JPEG returns mime type for jpeg.
func JPEG() string {
	return mime.TypeByExtension(".jpeg")
}

// LastModifiedTime returns http header with Last-Modified header set.
func LastModifiedTime(t time.Time) http.Header {
	h := make(http.Header)
	h.Set(LastModified, t.Format(lastModifiedTimeFormat))
	return h
}

// Merge returns a new header which contains merged properties of all h headers.
func Merge(h ...http.Header) http.Header {
	a := make(http.Header)
	for _, head := range h {
		for key := range head {
			value := head.Get(key)
			a.Add(key, value)
		}
	}
	return a
}

// IsJSONContent returns true if the Content-Type header is application/json.
func IsJSONContent(h http.Header) bool {
	return h.Get(ContentType) == ApplicationJSON
}

// IsForm returns true if the Content-Type header is
// application/x-www-form-urlencoded
func IsForm(h http.Header) bool {
	return h.Get(ContentType) == ApplicationForm
}

// IsMultipartForm returns true if the Content-Type header is
// "multipart/form-data".
func IsMultipartForm(h http.Header) bool {
	return h.Get(ContentType) == ApplicationMultipartForm
}

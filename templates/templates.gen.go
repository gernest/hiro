// Code generated by go-bindata.
// sources:
// backend/templates/data/bad_first_arg
// backend/templates/data/changelog
// backend/templates/data/download
// backend/templates/data/home_banner
// backend/templates/data/html/footer.html
// backend/templates/data/html/header.html
// backend/templates/data/html/home.html
// backend/templates/data/html/meta.html
// backend/templates/data/html/print.html
// backend/templates/data/html/privacy.html
// backend/templates/data/html/qr.html
// backend/templates/data/info
// backend/templates/data/missing_token
// DO NOT EDIT!

package templates

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _bad_first_arg = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\xce\x41\x4e\xc4\x30\x0c\x05\xd0\x75\x73\x8a\xaf\x1c\xa0\x07\x98\x2d\xe2\x06\xb3\xec\xc6\xd3\xba\x13\x23\x27\x2e\x89\x61\x90\x4a\xef\x8e\xc2\x48\x08\x98\x9d\x2d\x3f\x7f\xfd\x73\x62\xac\x52\x9b\x83\xea\xf5\x2d\x73\x71\xb8\x61\xdf\xc7\x27\xcb\x99\xca\x72\x1c\x98\xef\x13\x30\x53\x01\x8b\x27\xae\xb8\x30\x08\xef\xa4\xb2\xe0\xa5\x59\x41\xf3\x2a\xe5\x8a\x60\x15\x1b\x79\xea\x19\x74\xbf\xac\xa2\x3c\x86\x21\x9c\x93\xb4\x9f\x28\x52\xb5\x5b\xc3\x26\x5b\xff\x5a\xab\x65\x34\x5f\xa4\xc0\xcd\x3a\x7e\xfe\xa0\xbc\x29\xe3\x84\x30\x5c\x5e\xff\xd6\x89\xfb\x14\x13\xab\xda\x14\x4f\x53\xbc\x59\xd5\x65\x8a\x47\x7c\x84\xdf\x68\xec\x25\xc2\x30\x93\xff\xda\x3f\xff\xd1\xaf\x00\x00\x00\xff\xff\x27\xbc\x7e\xf4\x05\x01\x00\x00")

func bad_first_argBytes() ([]byte, error) {
	return bindataRead(
		_bad_first_arg,
		"bad_first_arg",
	)
}

func bad_first_arg() (*asset, error) {
	bytes, err := bad_first_argBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "bad_first_arg", size: 261, mode: os.FileMode(420), modTime: time.Unix(1519042478, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _changelog = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x8e\xc1\x4a\xc4\x30\x10\x86\xef\x79\x8a\x61\xb7\x07\x05\x1b\x16\x0f\x1e\x0a\xbd\xb8\xeb\xa2\x1e\x44\xa4\xe8\x41\x3c\x0c\xcd\x6c\x9b\xdd\x26\x91\x26\x94\x95\x30\xef\x2e\x6d\x22\x2c\x9e\x26\x99\xf9\xff\x8f\x6f\xbd\x86\x18\xe5\x1b\x0d\x84\x9e\x14\x33\x5c\x4d\x34\x7a\xed\xec\xbc\x7e\x4f\x4f\xe6\x6b\x21\x3e\x68\x68\x9d\x21\x08\x0e\x42\x4f\x17\xa5\x1d\x06\x92\x7b\x37\x1a\x0c\xb0\x7a\x46\x0b\xb7\x9b\xcd\xdd\x8a\x19\xc6\x74\x07\x77\x98\xd3\x2f\x68\x88\x59\x88\x18\xe5\x8e\x7c\xcb\x2c\x62\x1c\xd1\x76\x04\x72\xdb\xcf\x53\xbd\x8e\xee\x48\x6d\xf0\xcb\xa9\xf8\x4e\xbf\xaa\x96\xcc\x62\x9d\x34\x33\x23\x23\xa0\xbc\x80\x14\xa7\x9b\x62\xaa\x6a\x90\x5b\x67\x8c\x0e\xfe\xfe\xa7\xc1\x2e\xa3\x02\x76\x55\xfd\x07\x94\x7b\x6d\x55\x83\x1d\x14\x27\x66\x51\xc2\xe7\xf9\x0b\x62\xd4\x56\xd1\x19\xe6\xa4\x7c\x30\xee\xa8\x99\x21\x15\xff\xcb\x16\x13\xb3\x00\x48\xbd\xd9\xa9\xd1\x61\xa0\x25\x2e\x9f\xec\xc1\xc9\x47\xf4\xfd\x92\x27\xab\xb2\x21\x59\xb5\xcc\x12\xf2\xee\x37\x00\x00\xff\xff\x77\x0c\x2a\x8d\x75\x01\x00\x00")

func changelogBytes() ([]byte, error) {
	return bindataRead(
		_changelog,
		"changelog",
	)
}

func changelog() (*asset, error) {
	bytes, err := changelogBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "changelog", size: 373, mode: os.FileMode(420), modTime: time.Unix(1519042478, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _download = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x91\x41\x8f\xd3\x30\x10\x85\xef\xfe\x15\x23\xe5\x02\x15\x2e\x02\x09\xee\x15\x2b\x8e\x20\xc1\x6a\xcf\x75\x92\x71\x3c\xc2\xb1\xab\xb1\x43\xb7\x72\xfc\xdf\x91\xdd\xa4\xb4\x17\xb4\xa8\xda\x1c\xa2\x51\xf2\xe6\xcd\x7c\x6f\x1a\xe8\xfd\xd1\x59\xaf\xfa\x00\xda\x33\xa4\xb4\x7d\x42\x0e\xe4\x5d\xce\x42\x94\x2f\xbf\x91\x49\x9f\xc8\x0d\xd0\x19\xec\x7e\x85\x69\x0c\xef\xe0\xe4\x27\x38\x92\xb5\xe0\x10\x7b\xd0\xc4\x21\x42\xf4\x17\x2f\x88\x06\xa1\x25\xa7\xf8\x04\x5e\x17\x35\x8b\xce\x78\xea\x10\xde\x40\x40\x84\x16\xad\x3f\xbe\x85\x47\x83\xee\xb6\x69\x9d\x01\x9a\x2c\x82\x66\x3f\x82\x41\xc6\xb2\xd7\xc3\xa2\xfb\xb2\x48\x72\x16\xe4\x6a\x53\x50\x23\x42\x4f\x8c\x5d\xf4\x7c\x02\x15\xae\xe6\x6f\x85\xf0\x01\x9e\xc5\x7e\xbf\x17\xc1\xa8\x0f\xc5\x5b\x76\xc5\x6f\xf5\xf9\x4a\x16\x73\xae\x02\xe1\xa3\x41\x06\x47\xcf\x9b\xb5\xe1\xe3\xa7\xcf\xff\x6e\x69\x1a\xe8\xb1\x25\xe5\xa0\x4e\x9a\x41\x71\x67\x60\x86\xe1\x30\xc0\xbc\xb2\x09\x29\xe5\x2c\xeb\x33\x5f\xbd\xa5\x14\x29\x49\x60\xe5\x06\x84\xed\x03\xb6\x21\x67\x91\xd2\xf6\x7b\xc8\x19\xe6\x32\x71\xc7\x9d\x59\xeb\x9f\x34\xb8\x52\x5f\x65\x51\xe5\x12\xd0\x95\x4a\xec\x74\x44\xbe\xe4\x59\x4e\x56\xd3\xa1\xc1\xa9\x38\x31\xbe\xbf\x84\x7b\xa4\x68\xea\xbf\x9e\x42\x64\x6a\xa7\x48\xde\xc1\xc1\xa2\x0a\xb8\x1c\x5c\x44\xa3\x62\xbd\xf3\xe0\x63\xd5\x32\x0d\x26\x42\x88\x93\xd6\xa2\xa2\xa7\x74\xbb\x78\x53\x96\xfc\xa6\x46\x3c\x43\x3c\x55\x9f\xf3\xd2\x22\xa5\xf3\x8e\x35\xb2\x92\x99\xe2\xf1\xee\xc0\x76\x3c\xde\x17\xd8\x66\x73\x86\xdd\x6c\x6e\x81\x16\xe3\x17\x03\x35\x0d\xf0\xe1\x7e\x9e\x1f\x87\x57\xe2\x59\x8c\xff\x87\x27\x2a\x6e\x95\xb5\xe1\x5e\xa6\xc7\xea\xf3\x3a\x58\x7f\xbd\x5f\x4a\xf6\x27\x00\x00\xff\xff\x8c\x35\x30\xcc\xef\x04\x00\x00")

func downloadBytes() ([]byte, error) {
	return bindataRead(
		_download,
		"download",
	)
}

func download() (*asset, error) {
	bytes, err := downloadBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "download", size: 1263, mode: os.FileMode(420), modTime: time.Unix(1519042478, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _home_banner = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4c\x8f\xc1\x4e\x2b\x31\x0c\x45\xd7\xc9\x57\xdc\xe5\x7b\x12\x9a\x0f\x60\x07\x14\x89\x35\x54\xb0\x4e\x33\x2e\x63\xa9\x33\x6e\x1d\x0f\xa8\x44\xf9\x77\x94\x49\x40\x2c\x8f\x8e\xed\xeb\x0b\xbc\xd1\x29\xca\x4c\x30\xc1\xe1\x02\xef\x81\xfd\x44\xf8\x22\x15\x4c\xd7\x33\xdd\x60\xe2\xf7\x09\x67\xd2\xa3\xe8\x1c\x96\x48\xb8\x68\x94\x91\x90\x48\x3f\x38\xd2\xe0\x9d\xf3\xc0\xdd\x6a\x93\x28\x70\x8b\x9c\x87\x06\xa5\xe0\x5f\xce\xc3\xe3\x1c\xf8\x54\xca\x7f\x5f\xb3\x0e\x89\x8d\xda\x50\x87\x52\x3c\xb0\x93\x98\xe0\xfa\x76\x85\xbf\xf2\x45\x56\x8d\xd4\xe5\x33\x9d\x25\xb1\x89\x5e\x37\xb7\xff\x64\x33\xd2\xe6\x3a\x6c\xe2\x7e\xe5\xd3\x08\x5e\x8e\xe2\xdd\x2b\x69\x62\x59\x80\x7e\xa3\x73\x29\xde\xb5\xb1\x31\xfc\x3c\xb5\xf1\x2e\x6c\xc9\xee\x41\xe6\x99\x0d\xbf\x7b\x8d\xab\xa9\x8d\x6b\x7f\x52\x24\x0b\x6a\x34\xa2\xde\xcf\x79\x78\x92\x64\xa5\x80\x97\xad\xc9\xaa\xc1\x5a\x12\xf0\x1d\x00\x00\xff\xff\xcf\xdd\xbd\xa0\x6b\x01\x00\x00")

func home_bannerBytes() ([]byte, error) {
	return bindataRead(
		_home_banner,
		"home_banner",
	)
}

func home_banner() (*asset, error) {
	bytes, err := home_bannerBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "home_banner", size: 363, mode: os.FileMode(420), modTime: time.Unix(1519042478, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _htmlFooterHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x49\xcb\xcf\x2f\x49\x2d\xb2\xe3\x52\x50\x50\x50\xb0\x29\x4e\x2e\xca\x2c\x28\x51\x28\x2e\x4a\xb6\x55\xd2\x2f\x2e\x49\x2c\xc9\x4c\xd6\xcf\x2a\xd6\xcf\x4d\xcc\xcc\xd3\xcb\x2a\x56\xb2\xb3\xd1\x87\xa8\xb0\xe3\xb2\xd1\x87\x6a\x04\x04\x00\x00\xff\xff\x49\x1c\x72\x08\x41\x00\x00\x00")

func htmlFooterHtmlBytes() ([]byte, error) {
	return bindataRead(
		_htmlFooterHtml,
		"html/footer.html",
	)
}

func htmlFooterHtml() (*asset, error) {
	bytes, err := htmlFooterHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "html/footer.html", size: 65, mode: os.FileMode(420), modTime: time.Unix(1520029941, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _htmlHeaderHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x90\xc1\x6a\x2b\x31\x0c\x45\xf7\xf3\x15\x42\xab\xf7\xa0\xb5\x7f\x60\x32\x10\x4a\x16\x59\x15\xd2\xb4\xdb\x32\x38\x4a\x46\x89\x3d\x0e\x96\xd2\x50\x8c\xff\xbd\xb8\x29\xcd\xd0\x40\xe9\x4e\x1c\xfb\x08\xdd\xdb\x2a\xab\xa7\x2e\x67\xb3\xae\x43\x29\xad\xbd\x90\x26\x67\xde\x82\x09\xa4\x7d\x29\x90\xb3\x52\x38\xfa\x5e\x09\x70\xd0\xe0\x6d\xe5\xa6\x4e\x08\xe6\xf3\x9d\xc6\x4d\x29\x4d\xeb\x79\x3c\x40\x22\x3f\x43\xd1\x77\x4f\x32\x10\x29\xc2\x90\x68\x3b\x43\x2b\xda\x2b\x3b\xeb\x44\xec\x89\x0f\xac\x26\xf0\x68\x9c\x08\x82\xed\xfe\xae\x86\xfe\xc6\x9a\xa8\x43\x4c\xea\x4e\x0a\xec\xe2\xf8\xd3\xe6\xb0\xb3\x95\x9b\xe3\xb8\xc3\x0e\x2e\x09\xe7\x4e\xf9\x8d\x56\xf1\xa4\x54\x03\x88\x4b\x7c\xd4\xae\x01\x00\x38\xf3\xb8\x89\x67\xe3\xa3\xeb\xfd\x93\xc6\xd4\xef\xc8\x08\xe9\x52\x29\xfc\xc3\xf9\xc3\x7a\xf9\xb2\x78\x5d\x3d\x3e\xaf\x17\x78\x07\x98\xf3\x74\x15\x94\x82\xff\x9b\xd6\x7e\xad\xbb\x16\x74\x01\x20\xc9\x5d\xef\xda\x4f\xfb\xd8\x0b\x76\xdf\xde\xef\xff\xef\x6b\x18\xb9\xb5\x3e\x02\x00\x00\xff\xff\x8d\x25\xa3\xbb\xd5\x01\x00\x00")

func htmlHeaderHtmlBytes() ([]byte, error) {
	return bindataRead(
		_htmlHeaderHtml,
		"html/header.html",
	)
}

func htmlHeaderHtml() (*asset, error) {
	bytes, err := htmlHeaderHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "html/header.html", size: 469, mode: os.FileMode(420), modTime: time.Unix(1520029941, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _htmlHomeHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x51\x74\xf1\x77\x0e\x89\x0c\x70\x55\xc8\x28\xc9\xcd\xb1\xe3\xb2\x81\x50\x5c\x36\x19\xa9\x89\x29\x76\x5c\x0a\x0a\x0a\x0a\xd5\xd5\x25\xa9\xb9\x05\x39\x89\x25\xa9\x0a\x4a\x20\x59\x7d\x90\x54\x6a\x91\x1e\x88\xad\xa4\xa0\x57\x5b\xcb\x65\xa3\x0f\x51\xcd\x65\x93\x94\x9f\x52\x09\xd1\x65\x93\x92\x59\xa6\x50\x9a\xad\x5b\x5c\x90\x99\x97\x97\x5a\x64\x67\xa3\x9f\x92\x59\x66\xc7\x65\xa3\x0f\x51\x82\x61\x68\x5a\x7e\x7e\x09\x8a\xa1\x20\x53\x41\x6e\x01\x04\x00\x00\xff\xff\x8b\x59\xed\x73\xa2\x00\x00\x00")

func htmlHomeHtmlBytes() ([]byte, error) {
	return bindataRead(
		_htmlHomeHtml,
		"html/home.html",
	)
}

func htmlHomeHtml() (*asset, error) {
	bytes, err := htmlHomeHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "html/home.html", size: 162, mode: os.FileMode(420), modTime: time.Unix(1522353707, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _htmlMetaHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x24\xca\x31\x0a\x02\x31\x14\x84\xe1\xde\x53\x0c\xcb\xd6\xb1\x17\xd7\x0b\x78\x8a\x60\x06\xb1\xd8\x17\x48\x5e\xa1\x0c\xef\xee\xf2\x48\xf5\x33\x7c\x23\xed\x27\xbd\xde\x0e\x94\x6c\x84\x34\xaa\xbd\x89\xd9\x87\xb3\x3d\xf9\x9b\xd8\x97\x5c\xee\x59\x58\x3d\x79\x6c\x52\x89\xd8\xf0\xea\xe6\x34\xcf\xfd\xb1\xc6\xef\xfa\xa2\x20\xf1\xfa\x80\x44\x6b\x11\xff\x00\x00\x00\xff\xff\x81\x50\x47\x46\x66\x00\x00\x00")

func htmlMetaHtmlBytes() ([]byte, error) {
	return bindataRead(
		_htmlMetaHtml,
		"html/meta.html",
	)
}

func htmlMetaHtml() (*asset, error) {
	bytes, err := htmlMetaHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "html/meta.html", size: 102, mode: os.FileMode(420), modTime: time.Unix(1519042478, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _htmlPrintHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\x90\xcf\x4a\x03\x41\x0c\x87\xef\xf3\x14\x31\x67\x77\xab\x37\x0f\x33\x0b\xd2\xea\x41\x84\x8a\xb4\xa0\xc7\x74\x37\xec\x06\xe6\xcf\xba\x9b\xb6\xc8\x30\xef\x2e\x75\x11\xea\x29\x7c\x24\xf9\x42\x7e\xf6\x66\xb3\x5d\xef\x3e\xdf\x9e\x60\xd0\xe0\x1b\x63\x2f\x05\x3c\xc5\xde\x21\x47\x6c\x8c\xb1\x03\x53\xd7\x18\x00\x00\x1b\x58\x09\xda\x81\xa6\x99\xd5\xe1\x7e\xf7\x5c\x3d\xe0\x75\x2b\x52\x60\x87\x27\xe1\xf3\x98\x26\x45\x68\x53\x54\x8e\xea\xf0\x2c\x9d\x0e\xae\xe3\x93\xb4\x5c\xfd\xc2\x2d\x48\x14\x15\xf2\xd5\xdc\x92\x67\x77\x5f\xdf\xfd\x53\x0d\xaa\x63\xc5\x5f\x47\x39\x39\xfc\xa8\xf6\x8f\xd5\x3a\x85\x91\x54\x0e\x9e\xaf\xbc\xc2\x8e\xbb\x9e\xff\x36\x55\xd4\x73\xb3\x49\xed\x31\x70\x54\xbb\x5a\xd8\xd8\xd5\xf2\x83\xb1\x87\xd4\x7d\x2f\xb3\x39\x4f\x14\x7b\x86\xfa\x25\x1d\xea\xed\xa8\x92\xe2\x5c\xca\xa2\x91\xd0\xc3\x3c\xb5\x0e\x73\xae\x37\xa4\xb4\x7f\x7f\x2d\x05\x81\xbc\x3a\xc4\x06\x72\xe6\xd8\x95\x62\x8c\x5d\x2d\xbe\xcb\x81\x4b\x7a\x3f\x01\x00\x00\xff\xff\xc0\xac\xbc\xf2\x4d\x01\x00\x00")

func htmlPrintHtmlBytes() ([]byte, error) {
	return bindataRead(
		_htmlPrintHtml,
		"html/print.html",
	)
}

func htmlPrintHtml() (*asset, error) {
	bytes, err := htmlPrintHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "html/print.html", size: 333, mode: os.FileMode(420), modTime: time.Unix(1520078210, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _htmlPrivacyHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x59\xdd\x6e\x23\xc7\xd1\xbd\xd7\x53\xd4\x27\xe0\x83\x6d\x80\x22\x6d\xe7\xce\xa6\x05\x24\x1b\xc7\x11\x60\x38\x8b\xec\x3a\x41\x20\x08\x41\x73\xba\x86\xd3\x51\x4f\xd7\xb8\xba\x87\xdc\x89\x65\x60\x5f\x23\x80\xfd\x72\xfb\x24\x41\x55\xf7\x0c\x67\x28\xc9\x56\x74\x25\x91\xfd\x53\x3f\xa7\x4e\x9d\x6a\x6d\x9b\xd4\xfa\xeb\x8b\x8b\x6d\x83\xc6\x5e\x5f\x00\x00\xfc\xf8\x63\xc2\xb6\xf3\x26\x21\x5c\xca\xb7\x1b\xf9\x0a\x79\x2d\xbf\x5f\xc2\xfa\xa7\x9f\x2e\xb6\x9b\xbc\xfa\x62\xbb\x23\x3b\xe4\x5d\xdb\xe6\xf3\xeb\xd7\xec\x0e\xa6\x1a\xe0\x35\x79\x57\x0d\xdb\x4d\xf3\x79\xf9\xae\xbb\xfe\x06\xa9\x66\x1c\xe0\x6b\x0e\x18\x13\xec\x7a\xe7\x13\xa4\x06\x61\xf7\x43\xac\x4c\x00\xd3\x75\x60\x22\x18\xa0\x0e\x03\x44\xea\xb9\x42\x30\xc1\x42\xcd\x88\xf2\xed\x1a\xde\x36\x2e\xc2\x9b\xaf\xff\xfa\xb7\x9b\x57\x5f\x83\x8b\xd0\x31\x1d\x9c\x45\x0b\xbb\x01\xce\x8e\xbf\x35\x09\x02\xe9\xdd\xf2\x53\x51\x4c\x77\x7a\x9a\x8b\xe0\x42\xc2\x20\xdb\x6a\x62\xe8\x23\xca\xb5\x2e\xae\xb7\x9b\x6e\x32\x56\x6f\xea\xcc\x1e\x65\x7d\x1f\xd1\x42\x22\x70\xa1\x26\x6e\xe1\x88\xbb\xe8\x12\xc2\xc1\x45\x97\x88\x23\x30\xee\x0d\x5b\x17\xf6\x70\xdb\x0e\x0f\xd4\xf3\x1d\x74\xe2\xbf\xc3\x08\x47\x97\x1a\xf5\xb2\x22\xef\xb1\x4a\x8e\xc2\x4a\x0e\x5c\xa9\x31\xd6\xc5\xca\x53\xec\x19\x81\xea\xc9\xd8\xd7\xc8\x91\x82\xf1\x70\xa3\x17\x1a\xd9\x04\xae\x06\x13\x06\x0a\x08\x16\x2b\x75\x3a\x91\x1a\x3f\xdd\xf9\x06\xf9\xe0\x2a\x5c\xf8\x71\x53\xc3\x40\x3d\x54\x0d\x51\xc4\xe7\x76\xac\xc4\xc0\xa0\x0b\xcd\x5e\x82\x9d\xe8\xcc\x64\x35\x56\xb6\x52\x5d\xa2\x50\x8c\x0a\xc0\xe8\xf3\xef\xc5\x53\x17\x27\x3f\x34\x08\x83\xa4\x0d\x9f\xf6\x29\x35\x26\xc1\xed\xcd\xc3\x11\xef\xc6\xcb\xc0\x30\xe6\x80\x4b\x72\x72\x82\x25\xb2\x9a\xba\x56\xff\x0e\x7b\xb5\x6e\x74\x57\x0e\xf8\x3b\xde\xc1\xd1\x79\x3f\x5d\x1d\x28\x65\x7b\x19\x62\x23\x47\x0e\xd4\xf3\xc2\x74\x35\xb7\x44\x14\xdf\x55\xd8\x25\x81\x81\xc5\x58\xb1\xdb\xa1\x15\xd7\xc4\x19\x58\x22\xfa\x0c\x24\x08\x09\xb9\x2d\x08\x79\x7a\x07\x34\xe6\x80\x6a\x70\x34\x2d\x42\x8b\x26\xb8\xb0\x8f\x8a\xb9\x00\x62\xd5\x5b\x3d\x42\x1c\x7c\x45\xc1\x3a\xb1\x2e\xae\xe0\xd8\xb8\xaa\x11\xf4\x99\xaa\xc2\x18\xdd\xce\x23\x48\xb4\x7e\xdf\x75\x93\x97\xdf\x99\x16\xef\x56\xd0\x07\x8f\x31\x02\xa5\x06\xf9\xe8\xa2\x20\xa4\x76\xe1\x37\x7c\x18\x9d\x98\x0e\xdb\xc6\xc4\x14\xf6\xd7\xf3\x04\xbd\x5a\x22\xe0\xfb\x88\xdb\x4d\x59\x96\xf7\xcf\xa2\xf1\x27\x62\x30\xb0\xc3\x94\x90\x01\xdf\x75\xc8\x0e\x43\x85\xe2\x88\x97\x94\x4a\xde\xc4\xdd\x09\x74\x25\xf1\xad\x19\x80\xf1\x87\xde\xe5\x2c\x09\xf8\x4a\x5d\x43\x5f\xea\xa7\x42\x4e\xc6\x05\xe8\x0a\x88\xfc\x00\xce\x62\x48\xae\x76\x66\xe7\x71\x72\x61\x96\xdf\x15\xb8\x50\xf9\x5e\xb1\xb3\xeb\x93\x02\xc2\xbb\xd6\xa5\x5c\x38\xb7\xc6\x5a\x38\x36\x26\xe1\x41\xac\xf5\x31\xdf\x3d\x82\xb0\x41\xc6\x15\xe0\x7a\xbf\x96\xd4\x72\x84\x20\xb9\x7b\x00\x63\x2d\x4b\xa8\x1f\xc0\x53\xa5\xf7\x4c\x77\x3f\x40\xe7\xaa\xd4\x33\xc6\xbb\x8c\x78\xf7\x0c\xd0\xc5\x57\x61\x28\x17\xe1\x96\x51\xfc\x42\x0b\x14\x32\x42\x2d\x4a\x6c\x46\xa2\x12\xa3\x8b\x49\x99\xe4\x6e\x5b\x7c\xe8\xe3\x9d\x24\xd6\x84\x61\xba\xfb\x68\x86\xbb\x87\x5b\x29\x00\xd8\x21\x4c\x87\xee\x06\x89\x60\x29\x5d\xfb\x34\xbc\xbb\x02\x8d\xee\x69\x78\x0b\x29\x5b\x42\x45\xb8\x6c\x60\x0b\x9d\xe1\x34\x40\xcc\x59\x8c\xd9\x35\xc9\xe1\x18\xbb\xb9\xdf\x13\x73\xe6\x6c\x0d\xe2\xe4\x1a\x6e\xff\x21\xa1\x36\x01\x5a\xf9\x94\x02\x7c\x43\xb4\xf7\x38\x9d\x39\xf9\x25\x69\x50\xfb\xbd\x0b\xf7\x72\x4e\x5e\xf8\xd1\xb9\xd9\x42\x8e\x92\xbe\xa3\x09\xe9\x6e\x7d\x71\x82\xf4\x0c\xde\x73\x88\x7f\x4b\x7b\xf8\xa3\x49\x66\x09\xe6\x05\xa0\xcb\xf6\x91\x5b\x4c\x48\xb3\x0e\xa0\x30\x15\xbf\x8f\x0d\x06\x45\x90\x7c\xf2\x34\xb5\xba\x00\x95\xc9\xcc\x69\x02\x20\x33\x71\x0e\x7e\x0e\xee\x19\xfb\x59\x93\xcc\xc2\x64\x45\xc2\x2c\xa2\x1f\xa7\x86\xa9\xdf\x37\x8b\x64\x74\x4c\xb6\xaf\x52\xfc\x64\xc2\x51\xd7\x08\xaf\x55\xc6\x7b\xb4\x30\xba\x5b\x1a\xe8\xf8\xa7\x66\x2d\x57\xc9\x02\xad\x8b\xfb\x63\x5f\x35\x82\x9c\x19\x3a\xe3\x87\xf7\xff\x89\x70\x13\x12\x72\xc0\x04\xaf\x99\x12\x55\xe4\xe1\xe3\x0f\xef\x7f\xbe\x79\xfd\xe1\xfd\x2f\x9f\x8c\x75\xb2\x1a\xf1\x2c\xe5\xb3\x92\xc6\xce\x26\x49\x45\xc6\x21\x26\x6c\xe1\x80\x1c\xb5\x56\x2b\x0a\xb5\xdb\xf7\xfc\xf8\x7e\xaa\xa7\x50\x49\xb0\xa1\x4f\xce\xbb\xb8\xe8\xb5\xf3\x2e\x06\xc9\xb5\x19\x32\x56\xd4\x0b\xd5\xd9\xf0\xd2\xba\x66\x6d\x23\xf7\x5f\xa5\x4c\x88\xc9\x24\x17\x93\xab\xe2\x8c\x1a\x7f\x0d\x3f\xaf\x88\xee\x1d\xc6\xdf\x86\x4f\x59\xa8\x4d\xad\x76\x7e\xd4\x03\xb1\x35\xde\x83\x69\xa9\x0f\x49\xec\x92\xac\x67\x40\xb9\x08\x15\xb5\x2d\x05\x3f\x94\xa2\x95\x4a\xa7\x30\xb4\xd4\x47\xe8\x83\xfb\xa1\xc7\x89\xfb\x90\x95\x69\x44\xc0\xb0\x54\x4f\x48\x0b\x6b\x13\x65\xe7\x77\x4c\xc7\x88\x0c\x35\x53\xab\x21\x18\x05\x8c\xde\x28\xc8\x55\x25\xa3\x01\xd1\x83\x12\xf1\x23\x4a\xca\x49\x77\x9a\x74\xe3\xa1\xc5\x96\x78\x46\x17\xc5\xdf\xeb\xff\xbb\xba\x82\x57\x0d\x56\xf7\x52\x92\x4a\x30\x2e\x42\xe2\x1e\xb5\x9d\xeb\x79\xa6\xeb\x56\xf2\x6d\x1f\x44\xf7\xac\xe0\x5f\x7d\x94\xd6\x1b\xfb\x76\x66\x91\xcd\x4a\xa5\x1a\xe3\x17\x2c\xb4\x64\x85\x41\xf4\xd0\x80\xef\x84\xce\x03\x2e\x1c\xbe\xba\xca\xba\xed\xcd\xc8\x4d\xca\x5c\x45\x0a\x08\x53\x49\xa8\x3e\xbc\xff\xb9\x9c\xfa\xe1\xfd\x2f\xd2\xa7\x44\xa8\x25\x3f\xac\xe1\xcf\x74\x94\x5a\x5e\x4d\x80\x93\xfa\x38\xe7\xbd\x8a\xec\x48\x49\x3b\x36\xec\x66\x6c\xa5\x31\x17\x07\xfa\x47\xd7\x24\x9a\x8b\xa9\x79\x3d\xcb\x51\x42\x2c\x2a\x6d\x54\x27\x38\x9e\x88\x70\x0d\xc2\x94\x93\x80\xa0\x2e\x37\x13\x02\x74\x0a\x5c\x91\x06\xdd\x32\xe9\xc4\xc0\x58\x67\xab\xf1\x14\xc1\x0c\xf7\xfb\x40\xc7\x5c\x46\xa6\x7c\x21\xe9\xd9\xa1\x96\x24\x66\x8a\x9b\xe5\x7c\x0d\x8f\xe4\x63\x39\x5b\xd6\x8c\x27\x2f\xae\x97\xd5\x12\x36\x09\xfa\x0e\x41\xba\xf3\xa8\x3a\x23\xb5\x08\x1d\xb1\xea\x9b\x5c\x8d\xa7\x5c\xbd\xb0\xf0\xca\x6a\xe1\x1c\x11\x08\xfc\x9b\x25\x28\x78\x2c\x5a\x9e\x13\x04\x44\x0b\x11\xb1\x05\xef\xee\x11\x5c\xfa\x28\xc3\x43\x3e\x47\xbb\x52\xa1\x50\x5a\x89\x8e\x05\x61\x38\x6f\x4d\x2b\x09\xb0\x7c\x9e\x99\x63\x8e\x8c\x09\x10\x2b\xa8\x1a\x13\xaa\x5c\xf6\x2b\x3d\x4d\x2f\x16\x7f\xd7\x02\xd2\x27\x1a\x8c\xc4\x4c\x46\x2e\x1a\xf2\x99\x57\x23\xda\xda\xce\x84\xb1\x02\x5c\xb0\xee\xe0\x6c\x6f\x7c\x04\xdb\x4f\x32\xbd\x26\xef\xe9\x28\x39\x64\x34\x91\x42\xfc\x62\x19\x83\xde\x9f\x05\xd3\xbb\xeb\xb7\x04\xb5\xa9\x9c\x77\x49\x69\xf2\x24\xcb\xbe\xdc\x6e\xbc\x7b\x72\xfd\x28\xca\x66\x24\x2a\x1c\xa1\x0c\x83\x8d\xf1\xf5\xaf\x6c\x45\xd6\xbe\x59\xb6\x5d\xe9\xcc\xa0\xa9\xc8\x51\xfd\x12\x88\x9f\xdd\x6c\x62\x74\x51\x8a\x2a\xab\x1e\xe3\x87\x7f\x8b\xaf\x0d\x1d\xe7\x76\x8f\x73\xda\x7a\x79\xce\x76\x33\xf7\xfe\xd9\x86\x9e\x45\xde\x19\x26\x73\x31\xe7\x32\x3a\x25\x5a\x92\xa1\x15\x99\x75\xf9\x54\x32\x4f\x0d\x38\x4a\xcf\x0b\x9f\x72\x86\x94\x14\x4f\x61\xd1\xbe\x65\xe2\x7d\x54\x5f\xf7\x21\x6b\xa6\xd4\x60\xbb\x0c\xf0\x92\xa2\x06\x25\x6b\xda\x79\xb7\xd7\x60\x0a\x90\x13\x8d\x23\xa5\x8c\x3e\x8b\x9b\x0b\x25\x2c\xc8\xa7\x5e\x80\xb9\xeb\xb9\xa3\xf8\xf2\x52\xac\x7a\x76\x69\x78\xb1\x86\x3a\x18\xdf\x97\x39\x2c\xb1\x70\xbe\x08\xfa\x69\xbe\xeb\xe3\xf3\x61\x14\x77\x65\x14\x28\x7d\x2e\xb1\xcb\x43\xe0\xd8\x23\xda\x16\xb9\x72\x3a\x16\x64\x46\x5c\x0c\x05\xf2\x23\x53\x97\xe6\xb7\x63\x4a\xc2\xc2\x61\x0f\x2e\xad\xe1\x0f\x7d\x02\xc6\x16\xdb\x9d\x16\xb3\xbe\x19\x40\x8b\xa9\x21\xab\x68\x60\x13\x62\xeb\xa2\x88\x14\xa0\x83\xae\xc1\xd2\x04\x31\x29\x17\x9c\x16\xa3\xd0\x3b\x53\x70\xd5\x52\x40\x25\xe2\xf2\x8e\xf0\xd9\xa7\x9f\xfe\x3f\x44\x89\x5b\x6e\x21\x8c\x5e\xe7\x97\x4c\xce\xa3\x12\x34\x41\x32\xb9\xef\x0d\x9b\x90\x50\x58\x2a\x82\xd9\x45\xf2\x7d\xc2\xbc\xdb\xa5\xe1\x85\x49\xfa\xd6\x85\x7b\x85\xda\x5f\x34\xc3\x6f\x5c\x7a\x89\x6a\x99\x37\xd1\x22\xed\x83\xce\x5f\x7e\x3c\xae\xe8\x26\x39\xee\xd4\x24\xbc\xab\xee\x05\xb0\x66\x41\x61\xb2\x27\x93\xe0\x38\x9f\x58\xc7\x79\xa0\x51\x90\x9b\xa4\xe7\xac\x17\xf6\x7f\x47\x69\x51\x7e\xf8\xae\xe8\x0e\xbd\x52\x61\x20\x41\xca\x82\x72\x31\x1a\x69\xc9\x31\xd6\x24\xd4\x7b\x03\xd9\x55\x01\x86\x3d\xb8\x38\x0d\x98\x8c\x07\x87\x47\xcd\x66\x99\x8d\x17\xd7\x97\xc9\x3d\x2b\xc6\x38\x09\x26\xf1\x35\xd7\x7e\x20\x8d\x09\x93\x57\x58\xe4\x04\x16\x0d\x13\xe4\xf8\xd8\x51\x88\x6e\x27\x04\x3b\x68\x99\xe5\xf7\x94\x90\x30\xa4\xd5\x72\x7a\x79\xd4\x43\xf5\xdd\xc3\x54\x49\x25\x8c\x8e\x0d\xcb\xa6\x90\x63\x40\x33\x99\xf0\x42\xd5\xda\x38\x6f\x19\x83\x2a\xb9\xe2\xf6\xff\x06\x06\x51\x54\x1a\xf8\x71\x0e\x2e\xef\x27\x7d\xb0\xa5\x36\x04\xe9\x54\xc3\x67\xbf\x9b\x9e\x65\xca\x0e\x91\x1e\x4e\x33\x31\x8e\x3a\xe3\x30\xbf\x18\xe5\x9f\x9d\x43\x54\xbc\x56\xc5\x83\x72\xa1\xdc\x72\x93\x47\x29\x1d\xb1\x4a\x05\x09\x03\x96\x62\x35\x49\x94\x8e\xec\x9a\xb6\x40\x63\x66\x0f\x87\xe3\x40\xad\xb2\x7c\xb4\x68\x71\xf1\xe2\x4d\xa1\xdc\xe0\xda\x16\xad\x33\x09\xfd\x00\x16\x3d\x2a\x58\x5d\xcc\x46\x0a\x8b\x49\x66\x90\x4f\xb5\x21\x80\x35\xd2\x3e\x44\x64\x11\x6b\x75\x5b\x67\xb2\xfc\x9b\x16\x1c\x0d\xe3\x63\x35\xa9\xb4\x98\x7d\x58\x98\x3e\x3e\x8e\x9c\xe2\x38\xb7\xb4\xf3\x68\x62\x86\x9c\xa9\xd2\xe4\x67\xa4\xc5\x5b\xc4\x58\x92\x8f\x08\x53\x1a\x09\x41\x40\x69\x71\x86\x85\x57\x55\xb6\xbd\x18\x69\x26\xec\x51\x99\xe2\xed\xe3\xf7\xa7\x17\x77\x0c\xd5\xdf\x9d\x1d\x15\xca\xd9\xbb\x5a\x1e\x67\x64\xdc\x13\x1e\x71\x2d\x4a\xe9\xf7\x71\x75\x8a\xa7\x96\xbc\x5d\xd4\xfb\xf8\xaa\xdb\x21\x3b\xb2\xae\xd2\xa6\x51\x9f\xb5\x4a\x29\xb8\x2a\xbb\xb0\x78\x5d\x14\x20\x97\x47\x8c\xb1\x2e\xcb\x32\x21\xa0\x8e\x62\x1a\xdf\x26\x03\x1e\xcf\xcd\xa5\x70\xba\x7e\x1c\xdb\xca\xee\xe5\xe5\x8c\x80\x75\x2d\x4d\xea\x80\x73\xa4\xad\xc0\xd4\x29\x97\x59\x6e\xfd\x72\x61\x9e\xd3\x4e\x07\xbf\x74\x7e\xcd\xa8\xf8\xfe\x05\xcd\xa0\x20\x38\x4b\x9e\x30\x80\x3e\x5e\x65\x09\xcf\x10\xfb\xfd\x7e\xfc\xd3\xec\xa8\x4f\xa7\xa9\x7c\xe9\xfe\x6a\xe4\x81\x06\x63\xd6\x9c\x89\xce\xc1\xb9\xfe\x75\x0e\x3a\x8b\xa7\xa6\xf1\x68\x22\x54\x8c\xda\x03\xcc\x72\x10\xda\x1a\x68\x18\xeb\xaf\x2e\x9b\x94\xba\xf8\xc5\x66\x53\x48\x37\xbf\x18\x8d\xff\xdf\x58\x07\x4c\x97\x90\x0c\xef\x31\x7d\x75\xf9\xcf\x9d\x37\xe1\xfe\xf2\xfa\xd9\xa5\xdb\x8d\xb9\xce\x31\xde\x6e\xf2\x3f\x3e\x2e\xb6\x1b\xfd\xff\xc9\x7f\x03\x00\x00\xff\xff\x5f\x9f\xbf\x03\x46\x19\x00\x00")

func htmlPrivacyHtmlBytes() ([]byte, error) {
	return bindataRead(
		_htmlPrivacyHtml,
		"html/privacy.html",
	)
}

func htmlPrivacyHtml() (*asset, error) {
	bytes, err := htmlPrivacyHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "html/privacy.html", size: 6470, mode: os.FileMode(420), modTime: time.Unix(1521067761, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _htmlQrHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x91\xc1\x6e\xc3\x20\x0c\x86\xef\x3c\x85\x85\x72\x5c\xe1\x3e\xb9\xb9\xec\x09\xda\x37\xa0\xc5\x6d\x51\x52\x98\xa8\x1b\x6d\x42\xbc\xfb\x04\xac\x5b\xb2\x6c\xbb\xd9\xf9\x3f\xff\x31\xfe\xf1\xc2\xd7\xb1\x17\x02\x2f\x64\x6c\x2f\x00\x00\x52\x62\xba\xbe\x8e\x86\x09\x64\x51\x75\x91\x28\xaa\x52\x4b\x50\x39\x0b\xd4\x8d\x16\x78\x08\xf6\xfd\x31\xe5\x4e\xa0\x76\xfb\x9c\x6b\x8b\xd6\x4d\x70\x1c\xcd\xed\xb6\x95\xf7\x61\x73\x0c\x9e\x8d\xf3\x14\x65\xa3\xbf\x90\xfb\xb0\x39\x47\x67\xbf\xbf\x56\x85\xcd\x61\xa4\xd9\x78\xed\xe5\x12\x6a\x3f\x8d\xc6\x9f\x09\xba\xe1\xa9\x9b\x9e\xb7\x6a\xb7\x57\x2f\xc1\x33\xbd\xf1\xe7\x1a\x4b\xd7\xb8\x76\x68\x82\xed\x53\xea\x86\x9c\x51\xb3\xfd\x97\x99\xfe\x62\x50\xff\xe6\x9e\x12\x79\xfb\x63\x15\xd4\xf5\x35\xb3\x3b\x68\xeb\xa6\xd6\xce\xca\xc7\x28\xea\x76\xe3\x55\x2a\xa7\x10\x78\x91\x4a\x89\xa5\x84\xf9\x11\x00\x00\xff\xff\xd1\xcd\x10\x9b\xd3\x01\x00\x00")

func htmlQrHtmlBytes() ([]byte, error) {
	return bindataRead(
		_htmlQrHtml,
		"html/qr.html",
	)
}

func htmlQrHtml() (*asset, error) {
	bytes, err := htmlQrHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "html/qr.html", size: 467, mode: os.FileMode(420), modTime: time.Unix(1519042478, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _info = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x52\x50\x08\x4b\x2d\x2a\xce\xcc\xcf\x53\x50\x50\x50\xb0\x52\xa8\xae\xd6\x83\xf2\x6b\x6b\xb9\x14\x14\x9c\x4a\x33\x73\x52\x14\x5c\x12\x4b\x52\x21\x72\x60\x3e\x88\x0b\x96\x75\xce\xcf\xcd\xcd\x2c\x51\x80\xeb\x84\xf0\xc1\x52\xee\xf9\x0a\x49\x60\xbd\x99\x79\x69\xf9\x56\x5c\x20\x25\xfe\xc1\x48\x66\x78\xe6\xa5\xe5\xeb\xb9\xfb\xfb\x07\x83\x55\x2b\x28\x38\x06\x39\x7b\x60\x4a\x83\x44\xa1\x0a\xc2\x5c\x83\x82\x3d\xfd\xfd\x30\xd5\x40\x25\x6a\x6b\x01\x01\x00\x00\xff\xff\xfa\x09\x79\x55\xca\x00\x00\x00")

func infoBytes() ([]byte, error) {
	return bindataRead(
		_info,
		"info",
	)
}

func info() (*asset, error) {
	bytes, err := infoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "info", size: 202, mode: os.FileMode(420), modTime: time.Unix(1519042478, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _missing_token = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\xcd\xb1\xad\xc3\x30\x0c\x00\xd1\xfa\x6b\x8a\x9b\xc0\x4b\xfc\x15\xb2\x80\x22\x33\x16\x01\x99\x8c\x25\x0a\x29\x0c\xef\x9e\x22\x4d\xd2\xdf\xe1\xc1\x79\x2e\xff\xbe\xef\xd9\xd6\xeb\xa2\xcb\x31\xb5\xcb\x20\xcf\xa8\x62\xa1\x25\x87\xba\x2d\x29\xc1\xad\x0a\xa5\xa9\x58\xa0\x03\xf3\xf8\x8e\x64\x65\x38\xcf\x26\x79\x08\x7d\x1a\x51\x85\x87\xb7\xe6\x2f\xb5\x8d\xf2\x11\x12\x84\xff\x6c\x0b\xe9\x2f\x01\xdc\x0f\x9a\x6f\x6a\xef\x00\x00\x00\xff\xff\x98\x50\x02\xae\x91\x00\x00\x00")

func missing_tokenBytes() ([]byte, error) {
	return bindataRead(
		_missing_token,
		"missing_token",
	)
}

func missing_token() (*asset, error) {
	bytes, err := missing_tokenBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "missing_token", size: 145, mode: os.FileMode(420), modTime: time.Unix(1519042478, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"bad_first_arg": bad_first_arg,
	"changelog": changelog,
	"download": download,
	"home_banner": home_banner,
	"html/footer.html": htmlFooterHtml,
	"html/header.html": htmlHeaderHtml,
	"html/home.html": htmlHomeHtml,
	"html/meta.html": htmlMetaHtml,
	"html/print.html": htmlPrintHtml,
	"html/privacy.html": htmlPrivacyHtml,
	"html/qr.html": htmlQrHtml,
	"info": info,
	"missing_token": missing_token,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"bad_first_arg": &bintree{bad_first_arg, map[string]*bintree{}},
	"changelog": &bintree{changelog, map[string]*bintree{}},
	"download": &bintree{download, map[string]*bintree{}},
	"home_banner": &bintree{home_banner, map[string]*bintree{}},
	"html": &bintree{nil, map[string]*bintree{
		"footer.html": &bintree{htmlFooterHtml, map[string]*bintree{}},
		"header.html": &bintree{htmlHeaderHtml, map[string]*bintree{}},
		"home.html": &bintree{htmlHomeHtml, map[string]*bintree{}},
		"meta.html": &bintree{htmlMetaHtml, map[string]*bintree{}},
		"print.html": &bintree{htmlPrintHtml, map[string]*bintree{}},
		"privacy.html": &bintree{htmlPrivacyHtml, map[string]*bintree{}},
		"qr.html": &bintree{htmlQrHtml, map[string]*bintree{}},
	}},
	"info": &bintree{info, map[string]*bintree{}},
	"missing_token": &bintree{missing_token, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}


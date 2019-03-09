// Code generated by go-bindata. (@generated) DO NOT EDIT.
// sources:
// static-source/scripts.js
// static-source/standard.css
package main

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

var _scriptsJs = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x54\xdf\x6f\x22\x37\x10\x7e\xe7\xaf\x98\xf3\x0b\xf6\x65\xb3\x81\xa8\xed\x43\xd0\xde\xe9\xb8\xa4\xba\x4a\x49\x4e\x4a\xa8\xaa\x0a\xf1\x60\xf0\xb0\xb8\x35\x36\xb2\xbd\xc0\x92\xe3\x7f\xaf\x6c\x2f\x3f\x36\x6d\x22\x35\x0f\x59\xfc\xf9\xf3\x37\xf3\x79\x3c\xb3\xe6\x16\xbc\x5c\xa2\xa9\xbc\x83\x02\xc6\x93\x4e\x27\x40\xce\x54\x76\x86\x50\x80\xc6\x0d\xdc\xad\x51\xfb\xe7\x88\x50\x72\x85\x61\xe5\xae\x08\x1b\x74\x12\x2b\x37\x7a\x89\xce\xf1\x32\xf0\xe7\x95\x9e\x79\x69\x34\x8d\x34\x06\x2f\x1d\x00\x80\x20\x19\x81\x5b\xf4\x5c\xaa\x10\x29\x2e\x73\xc1\x3d\xcf\xdd\x4a\x49\x4f\xbb\x59\x97\x1d\xc9\x9e\xdb\x12\xfd\xd0\x6c\xa1\x00\x61\x66\xd5\x32\x90\x4b\xf4\x77\x0a\xc3\xcf\x61\xfd\x9b\xa0\xe7\x82\xe3\xde\x84\x0d\xe2\xe9\xd9\x82\xeb\x12\xbf\x28\xb4\xfe\x1e\xd7\xa8\xe8\x51\x2a\x6b\xa5\x30\xee\x4f\x5e\x01\xd7\x41\x62\x3f\xe8\x74\x3a\x07\x13\x30\x35\xdb\xaf\x4a\xce\xfe\xa6\x52\x9c\x5b\xf1\x50\xc0\x2d\xf7\x98\x6b\xb3\xa1\xa7\xa4\x57\x01\x5f\xd6\x23\xb9\x44\xea\x5f\x7b\x79\xc7\x88\x14\x89\xcc\x43\xd2\x8f\xe6\xf7\x95\xe0\x1e\xa9\x14\x19\xfc\xd2\x63\xef\x99\xca\x80\x94\x16\x51\x93\x0c\xc8\x08\x9d\x97\xba\x04\xa3\x55\x4d\x58\x67\x7f\x6e\xe3\x98\x53\xe3\x41\xce\x81\x82\x87\x0f\x05\xe8\x4a\x29\x38\xc0\xe1\xcf\x36\x35\x0f\xf6\x0e\x26\xf6\x80\xca\xe1\x5b\xa4\x86\x13\xff\x5b\xf4\x95\xd5\x60\x83\xc3\x5f\x2b\xa5\xfe\x44\x6e\x29\x83\x0b\x58\x71\x41\x23\xfa\x60\xb4\x5f\x44\xa8\x9f\x5d\xb7\x76\x6e\x79\x4d\x59\xc2\xc8\x88\x9c\xef\x7c\x33\x95\x75\xcd\xde\x29\x87\xa8\x26\x75\xe5\xd1\xb5\x43\x3c\xe3\xcc\x68\x91\x0e\xb4\xef\x21\x50\x74\x06\x1b\x29\xfc\x22\x83\x5d\xf2\xbd\x83\x02\x76\xf0\xe3\x07\x74\x7b\xdd\xf0\x84\x74\x30\x07\x17\xd0\x8d\xab\xc6\x91\xce\x15\xea\xd2\x2f\xe0\x53\x91\x8e\xc3\x67\xd0\x70\x13\x6f\xe1\x8b\xb5\xbc\xa6\x09\xbd\x3c\x31\x2f\xa0\xcf\xf2\xbf\x8c\xd4\x74\x17\xd2\xd3\x83\x76\x2e\x6f\x96\x54\x85\x55\x06\x4d\x4b\xb5\x4a\x36\x26\x7c\x39\x45\x4b\xb2\x43\xd9\xc3\xb7\x26\x19\xb1\x28\xc8\x24\x97\x5a\xe0\xf6\xfb\x9c\x46\x05\x06\x45\x01\x97\x7d\x06\x2f\x49\x11\x0a\x48\xec\xa6\x54\x29\x5c\x3e\x53\xdc\xb9\x7b\xe9\x7c\x6e\x71\x69\xd6\x48\x0f\x21\xce\x9e\x56\x0a\x02\x31\x0a\xfb\xef\xc3\x5c\x88\x26\xec\xf9\xfe\xe9\x9d\xbb\x61\xfd\x35\x90\x1f\xf9\x12\x29\x69\xbc\x11\x36\xee\x85\xac\x35\xda\x6f\xa3\x87\xfb\xd0\x3c\x69\xa3\x7d\x53\x56\x96\x0b\xff\x2c\x77\x38\x94\xe5\xd0\x6c\xe9\x79\x27\xf2\x35\x97\x8a\x4f\x15\xfe\x11\xaf\xbf\x80\x07\xee\x17\xf9\x5c\x19\x63\xe9\x46\x6a\x61\x36\x49\x3f\x6d\x5f\xc1\xcf\xfd\x6b\x06\x1f\xc3\x27\x4a\xc4\xa2\xa5\x29\x43\x5f\x49\x7d\x2a\xa0\xdf\xbb\xfe\x89\xc1\xe7\x57\x41\x6e\x02\x1c\x4f\xbf\xd5\xcf\xdd\xa9\x2c\x2f\xa7\x66\xdb\x65\xb9\xf3\xb5\xc2\x7c\xd3\x24\x77\x0c\x77\x01\x64\xb5\x25\x6d\x9b\xff\xee\xfe\x30\x97\x4f\xf5\xa7\xbe\x5e\xa1\x99\x1f\xa7\xf5\x58\x8a\x09\x7c\x28\x0a\x20\x95\x16\x38\x97\x3a\x54\x07\x5e\x60\xa6\x90\xdb\x51\x22\xd1\x73\x32\xdb\x0f\xfe\xff\x4c\x6a\x45\x2b\xc0\xa1\x3f\x48\x1f\x27\x3d\x7b\x79\x7b\x3c\x85\x27\x93\x01\xb9\x7b\x7a\xfa\xfe\x74\x03\x8f\x06\xaa\x68\xcf\xc1\xdc\x58\x08\x4d\x1e\xf4\xc3\x75\xb8\x9c\xb0\x7d\xb2\x0c\x1f\xa1\xdf\xeb\xf5\x42\xef\xfe\x13\x00\x00\xff\xff\x6f\xed\x47\xe3\x9e\x06\x00\x00")

func scriptsJsBytes() ([]byte, error) {
	return bindataRead(
		_scriptsJs,
		"scripts.js",
	)
}

func scriptsJs() (*asset, error) {
	bytes, err := scriptsJsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "scripts.js", size: 1694, mode: os.FileMode(420), modTime: time.Unix(1552072328, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _standardCss = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\xd2\xdb\x6e\xa3\x30\x10\x06\xe0\x7b\x9e\xc2\xd2\x6a\x2f\x41\x06\x96\x6c\xe3\x3c\x8d\xb1\x07\x63\xc5\x87\xc8\x4c\x1a\xda\xaa\xef\x5e\x35\x50\x70\x02\x39\xf4\x32\x93\xff\xff\x0c\x8c\xa5\x7e\xcd\x54\x00\x70\xe4\x23\x21\x84\x90\x9a\x8b\xbd\x0a\xfe\xe8\x64\x2a\xbc\xf1\x81\xfd\x29\x4b\x21\xca\x72\x97\x7c\x26\xc9\x77\x38\x80\xbc\x19\x15\x82\x52\x4a\xa7\x28\xb7\x35\x84\x9b\xe1\xa6\x79\x89\xc3\x2a\xc0\xdb\xcd\x2c\xcf\xb7\xb0\x15\x53\xb6\x37\x3c\x28\x18\xd3\x27\x2d\xb1\x65\x55\x4e\x0f\xfd\xee\x3c\x68\x41\xab\x16\xe3\x49\xe3\x1d\xa6\x9d\x7e\x07\x56\x80\x9d\x14\xf9\x8c\x52\x54\xff\x96\x4a\x9e\x55\x91\xb3\x64\xa2\xd2\xf3\x8c\xb4\x20\xf5\xd1\x3e\x80\xf2\x62\xb3\x06\xe5\x45\x4c\xad\x48\x51\xef\x37\x92\xec\x2c\x37\xe6\x81\xb4\x29\x96\x10\xcd\xfe\xc7\xce\x92\x99\x4b\xcf\x2b\xd2\x6a\x11\xfc\x7d\xa6\x5c\xd9\x3a\xbd\xf8\xd0\x4b\xa4\xbc\xde\xfa\x63\xa4\xf6\xfd\x48\x58\x1e\x94\x76\x69\x18\xf6\xfc\xd3\x1b\xa7\xb5\x47\xf4\x76\x1e\x37\xc6\x73\x64\x06\x1a\x1c\x7e\x23\xf4\x98\x72\xa3\x95\x63\x44\x80\x43\x08\xc3\x7c\xb8\xf7\xe4\xd4\x6a\x84\xf9\x4c\xad\xd2\xeb\x73\x19\xe1\x47\xf4\xe7\xc8\x21\x43\x8d\x06\x2e\x1f\xab\x36\x5e\xec\xd3\x0e\x79\x40\x46\xaa\xbf\xbb\xe5\x7f\xe0\xe4\x05\x62\xa1\xeb\xf8\x74\xa1\xe7\xf7\xdf\xd0\xb5\xf6\x28\x0f\xfd\xfb\xf6\x57\x00\x00\x00\xff\xff\x38\x14\xf0\xe4\x6b\x04\x00\x00")

func standardCssBytes() ([]byte, error) {
	return bindataRead(
		_standardCss,
		"standard.css",
	)
}

func standardCss() (*asset, error) {
	bytes, err := standardCssBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "standard.css", size: 1131, mode: os.FileMode(420), modTime: time.Unix(1551633529, 0)}
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
	"scripts.js":   scriptsJs,
	"standard.css": standardCss,
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
	"scripts.js":   &bintree{scriptsJs, map[string]*bintree{}},
	"standard.css": &bintree{standardCss, map[string]*bintree{}},
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

// Code generated by go-bindata. DO NOT EDIT.
// sources:
// static-source/scripts.js (2.074kB)
// static-source/standard.css (1.962kB)

package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
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
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
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

var _scriptsJs = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x55\x4b\x6f\x1b\x37\x10\xbe\xef\xaf\x98\xf0\x62\x32\x5e\x6d\x64\x17\xbe\x58\xd8\x04\xb1\xeb\x22\x2d\xfc\x28\x6c\x17\x45\x11\xf4\x40\x2d\x47\x2b\xd6\x14\x29\x90\x5c\xbd\x1c\xfd\xf7\x82\xcb\x7d\x48\x82\x65\xe4\x24\x71\xe6\x9b\x6f\x86\xdf\xcc\x70\x17\xdc\x82\x33\x95\x2d\x10\x72\xd0\xb8\x84\x9b\x05\x6a\xff\x54\x5b\x28\xf9\x84\xe1\xe4\x3e\x11\x36\x4a\x22\x2a\x33\x7a\x86\xce\xf1\x32\xe0\x27\x95\x2e\xbc\x34\x9a\xd6\x30\x06\xaf\x09\x00\x40\xa0\xac\x0d\x90\xc3\x1f\x4f\x0f\xf7\xd9\x9c\x5b\x87\x11\x93\x09\xee\x39\xab\x61\x6e\x29\x7d\x31\x6d\xcc\x7e\x3d\xc7\x36\x1e\xa0\xe0\x0e\x81\xbc\x20\xce\xb9\x92\x0b\x24\x97\x8d\x1d\xa0\xb3\x51\x96\x74\xc6\xb1\x45\xfe\x32\x4a\xf6\x82\xab\xb9\xe0\x1e\xaf\xcc\x6a\x27\x38\x14\xe6\xb9\x2d\xd1\x5f\x99\x15\xe4\x20\x4c\x51\xcd\x42\xf6\x12\xfd\x8d\xc2\xf0\xf7\x6a\xfd\xbb\x68\x4a\x92\x82\x8d\xba\xd0\x62\xca\x75\x89\x5f\x15\x5a\x7f\x8b\x0b\x54\xb4\xe3\x49\xe3\x5d\x33\xe7\xb9\xaf\x5c\x7b\x52\xdc\xf9\xbb\xa8\x13\x1b\xbd\x5f\xa9\x40\x85\x87\x95\x76\xb6\xbe\x96\xf7\x49\x2c\x2a\xc3\xc5\x9f\xbc\xdc\x15\x4b\x99\x82\x87\xf6\x64\xd1\x4b\xa3\xee\xdb\x64\x3b\x4a\x92\xa4\xed\x1d\x8c\xcd\xea\x5a\xc9\xe2\x85\x4a\x11\x3a\xb0\xdd\x71\xf5\x65\x44\xdf\xae\x84\xef\xe8\x17\xaa\x85\x06\x16\x9a\x8f\xda\xdf\x1b\x81\x99\xc5\x99\x59\xe0\xf5\x54\x2a\xd1\xe8\xc7\xf6\xd2\xed\x34\xf7\xe7\x93\x91\x28\xfc\x60\xcc\x2d\x09\x69\x8f\x74\x2a\x05\x52\x5a\x44\x4d\x52\x20\x35\x4e\x4e\x68\x18\x3a\x33\x81\x17\x0e\x1f\xf2\x1c\x48\xa5\x05\x4e\xa4\x46\x41\x18\xbc\x42\xa1\x90\xdb\x67\x39\x43\x53\x79\xfa\xc2\x19\x6c\x13\x08\xd0\x1c\x1c\xfa\xd6\xde\x2d\x00\x7b\x3d\x9e\xd7\xa2\x08\x59\x6f\x1e\x1f\x1f\x1e\x2f\xe1\xde\xf4\xf7\x74\x30\x31\x16\x2e\x5c\x46\xd8\x36\x85\x0b\xf8\x08\x67\xc3\xe1\xb0\x16\xa5\x57\x65\xb6\x0e\xd9\x68\xb7\x5e\x72\x02\x14\x3c\x7c\xc8\x41\x57\x4a\x41\xbf\x35\x00\xb6\xd9\xe1\x5f\xb9\x0f\x01\xb1\xdd\x80\xca\xe1\x31\x50\x3b\x12\xb5\x03\x7d\x65\x35\xd8\xa0\xf0\x6f\x95\x52\xff\x20\xb7\x94\xc1\x29\x90\x01\x81\x53\x98\x73\x41\x6b\xdf\x9d\xd1\x7e\x5a\x3b\xce\xd2\xf3\xce\xdf\xf1\x77\xb8\x98\xa0\x81\x3c\xef\x51\x7c\x33\x95\x75\x9d\xef\xf2\xcd\xf0\x3b\xa9\x2b\x8f\xfb\xa8\xde\xfb\x84\x85\xd1\x22\x7a\xf7\xf5\x0a\x10\x9d\xc2\x52\x0a\x3f\x4d\x61\x13\xf5\xd9\x40\x0e\x1b\xf8\xf1\x03\x4e\x86\x27\x61\xab\x75\x10\x01\x4e\xe1\xa4\x3e\x35\x37\xd7\x99\x42\x5d\xfa\x29\x7c\xce\x63\x38\x7c\x01\x0d\x97\xb5\x5a\x5f\xad\xe5\x6b\x1a\xad\x83\x1e\x79\x0a\x67\x2c\xfb\xcf\x48\x4d\x37\xa1\x48\x3d\xda\xaf\xe5\xe8\x50\xb4\xaf\x45\xf3\x96\xee\xf5\xf6\x3b\xe1\xb3\x31\x5a\x92\xb6\x13\x1b\x7e\xd7\x24\x25\xda\xfc\x55\xbf\x6a\x24\xad\x67\xea\xdf\x4c\x6a\x81\xab\x87\x09\x8d\x6c\x0c\xf2\x1c\x06\x67\x61\x78\xa3\x01\x72\x88\xa1\x4d\x83\x9b\x8d\x2c\x14\x77\xee\x56\x3a\xdf\x2c\x24\x6d\xf3\xed\xac\x48\xcc\x08\x3b\x29\xe3\x1c\xb3\xb7\x89\xb8\x10\x6d\x11\xbb\x80\x7e\x55\xdd\xd5\xfa\x3a\xa0\xef\xf9\x0c\x29\x69\x6e\x4d\xd8\xf7\x61\xb8\x84\x46\xfb\xed\xf9\xee\x16\xf2\x56\x8e\x9f\xe1\x08\x2f\x6c\xac\x4d\xbc\xc1\x13\x77\xe6\x60\x32\xac\x2c\xa7\xfe\x49\x6e\xf0\x4a\x96\xe1\x51\xdb\xfd\x60\xf1\x05\x97\x8a\x8f\x15\xfe\x5d\xb7\x38\x87\x3b\xee\xa7\xd9\x44\x19\x63\x29\x5d\x4a\x2d\xcc\x32\x66\x88\xfe\xc1\x2f\x43\x06\x9f\xe0\xe2\xec\x9c\xc1\xc7\xf0\x53\x13\xd5\xe3\x11\xbf\x2d\xf4\x80\xf0\x73\x1e\xc1\x5f\x0e\x32\x5d\xb6\xb1\xc7\xde\xb7\x93\xb1\x2c\x07\x63\xb3\x3a\x61\x99\xf3\x6b\x85\xd9\xb2\x29\xb0\x4b\x76\x0a\x64\xbe\x22\xef\x93\xf4\x8f\xe4\x21\x0f\xed\x89\x06\xe7\xc0\x5a\xb6\x6d\xf2\x7f\x00\x00\x00\xff\xff\xfb\x59\x5f\x4c\x1a\x08\x00\x00")

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

	info := bindataFileInfo{name: "scripts.js", size: 2074, mode: os.FileMode(0644), modTime: time.Unix(1568801885, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xd8, 0x76, 0xd, 0x4, 0xc1, 0xa9, 0x6b, 0x1b, 0x9e, 0xcd, 0x6a, 0xcc, 0x48, 0x5f, 0xa9, 0x4d, 0x1d, 0x34, 0x68, 0x91, 0xae, 0x9e, 0xdf, 0x91, 0x81, 0x82, 0x76, 0xb, 0x37, 0xf5, 0x60, 0xf6}}
	return a, nil
}

var _standardCss = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x94\xc1\x72\x9b\x30\x10\x86\xef\x3c\x85\x26\x99\xde\x82\x07\x4c\x9c\x34\xca\xad\xd3\x37\x68\xd3\xfb\x82\x16\xac\x46\x48\x8c\xb4\x38\x76\x33\x79\xf7\x8e\x2d\xd9\xe0\x00\x36\x9d\xe9\xcd\x2c\xff\xff\xfd\xcb\x6a\xe5\x86\xbd\x97\x46\x53\x5c\x42\x2d\xd5\x8e\xb3\x9b\xef\xf8\x1b\x36\x2d\xfb\x01\xda\xdd\xdc\xb1\x5f\x68\x05\x68\xb8\x63\x0e\xb4\x8b\x1d\x5a\x59\x3e\x7f\x44\x91\x90\x9b\x45\x65\x11\x35\x7b\x8f\x18\x63\x2c\x87\xe2\xb5\xb2\xa6\xd5\x22\x2e\x8c\x32\x96\xdf\x66\x59\x51\x64\xd9\x73\x14\xc4\x16\xc5\x1d\xdb\xff\xd0\xe6\xa5\x11\x40\x38\x69\x2c\x8a\x24\x49\x92\x93\x11\xea\x1c\xed\xa4\xb8\x2c\xbf\xf6\xc5\x95\xc5\xdd\xa4\x16\xd2\x27\x7c\x2a\x4e\xda\xad\x02\x5b\x1d\xdb\x78\x93\x82\xd6\x7c\x95\x26\xcd\xf6\xf9\x50\x58\xa3\xac\xd6\xd4\xaf\x1c\xa6\xe4\xe4\x1f\xe4\x4b\xac\x4f\x14\x31\x87\xb2\x5c\xdd\x0f\x29\xe9\x62\xd5\xe3\x0c\x31\x3d\xd3\x7c\x8c\xa8\x51\xc8\xb6\xbe\x02\x4a\x97\x0f\x63\xa0\x74\xd9\x47\x8d\x90\x7a\xbe\x7f\x21\x09\x57\x83\x52\x57\x48\x0f\xcb\x21\x28\x59\x3c\xf6\x39\x43\x4c\x67\x9a\x4f\x11\xb5\x2c\xac\xb9\x8c\xc9\x46\x4e\x3d\x39\x1b\xf4\x10\x92\x7d\x3e\xf5\xeb\x10\x47\x40\xad\x0b\x94\xe0\x62\x63\xb6\xb4\x33\xdd\x7a\x53\x9c\x83\x65\xcd\x82\x24\xa9\xe3\xda\x74\x72\x96\x26\xc9\x17\xcf\x68\x8c\x93\x24\x8d\xe6\x0c\x72\x67\x54\x4b\xe8\xeb\x0a\xcb\x7d\x96\x7f\xc8\x0d\x91\xa9\x39\x4b\xc7\x43\x14\x38\xf2\x97\x56\x0c\xa3\x1e\x57\x57\x92\x6c\xf8\xac\x89\xa8\x45\x6e\xb6\x81\xea\xef\x29\x7b\x5b\xcb\xa3\xb7\x54\x06\x88\xef\x7b\xf5\xcf\x35\xd8\x4a\xea\x38\x40\x4e\xa7\x16\xca\x3e\xe9\x54\xed\x1a\xb2\xa8\x80\xe4\x26\x40\x09\xb7\x14\x83\x92\x95\xe6\xac\x40\x4d\x68\xbb\x5e\x64\x15\x77\xfd\x78\x2c\x67\xd0\x92\x39\x48\xce\xe7\x7d\x6c\x46\x99\xe2\x35\x76\x04\x96\x38\x3b\x0e\xe3\xec\x1d\x6a\xe1\x21\x9e\x51\xa3\x73\x50\x0d\x4f\xed\x21\x19\x33\x07\x70\x58\x9c\x69\x76\x68\xb0\x86\xed\xcf\x6f\x2f\x03\xf6\xfd\x25\xb6\xb7\x5f\x42\xef\xdf\x6d\xa4\x93\xb9\x54\x92\x76\x9c\xad\xa5\x10\xa8\x43\xe4\xa5\xf5\xf8\x0f\xb9\x53\x7b\xd5\xdb\xa4\xf3\x35\xfb\x88\xa2\x70\x3d\x3f\x8d\x5a\x48\xd7\x28\xd8\x71\xa6\x8d\x46\xaf\x13\x73\x85\x41\x37\xfc\xd6\x0b\xd0\x59\x62\xff\xa7\x36\x47\xfb\x37\x00\x00\xff\xff\xe6\xf2\x69\x67\xaa\x07\x00\x00")

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

	info := bindataFileInfo{name: "standard.css", size: 1962, mode: os.FileMode(0644), modTime: time.Unix(1578476198, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x97, 0xb4, 0xa5, 0xe9, 0x5a, 0x49, 0x2b, 0x4c, 0x32, 0x10, 0xa1, 0xcc, 0xac, 0x69, 0x55, 0xb4, 0x2c, 0x87, 0x86, 0xff, 0xf9, 0x4a, 0x80, 0x2e, 0xb6, 0x4b, 0xe2, 0x90, 0x10, 0x63, 0xdd, 0x48}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
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

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
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
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
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

// RestoreAsset restores an asset under the given directory.
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
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}

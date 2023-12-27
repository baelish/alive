// Code generated by go-bindata. DO NOT EDIT.
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
	info  fileInfoEx
}

type fileInfoEx interface {
	os.FileInfo
	MD5Checksum() string
}

type bindataFileInfo struct {
	name        string
	size        int64
	mode        os.FileMode
	modTime     time.Time
	md5checksum string
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
func (fi bindataFileInfo) MD5Checksum() string {
	return fi.md5checksum
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _bindataScriptsJs = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x57\x5d\x6f\xdb\xb8\x12\x7d\xf7\xaf\x98\xf2\xa1\x96\x6a\x5b\x76\x7a\xd1\xfb\x10\x47\x2d\x92\xdc\x02\xed\xdd\x26\x29\x92\x14\x7d\x28\x8a\x05\x2d\x8d\x64\xd6\x34\x29\x90\xf4\x57\x5a\xff\xf7\x05\x3f\x64\xc9\x89\xd3\x4d\xb1\xdd\x3c\xc4\x36\x39\x3c\x73\x38\x9a\x39\x33\x1a\x0e\xe1\x1a\xb9\xa4\x39\xb0\x02\x14\x0e\x96\x4c\x33\xc3\x44\x09\x0b\x6d\xff\x4f\x68\x36\x1b\x16\x52\xad\xa8\xca\x61\xb2\x30\x46\x0a\x9d\x74\x58\x01\xd1\x8d\x51\x4c\x94\xd1\x8a\x89\x5c\xae\x92\x0a\x55\x21\xd5\x9c\x8a\x0c\x93\x12\xcd\x5b\x61\x14\x43\x7d\xb6\xb9\xdd\x54\x18\x11\x41\x97\xac\xa4\x86\x49\x41\xe2\x2f\xa3\xaf\x89\xd9\x54\x18\x43\x9a\xa6\x40\xac\x83\x3f\x83\x03\x12\x7f\xef\x00\x70\x99\x39\xd3\x44\x39\x5e\x51\x3c\xee\x6c\x3b\x9d\x8e\x23\x5a\x32\x6d\x50\xc1\x8a\x99\x29\x4c\xe4\x1a\x70\x89\xc2\x80\x96\x0b\x95\x61\x87\x63\xfd\x15\x52\x10\xb8\x82\xb7\x76\xf7\xc6\xad\x44\x64\xe8\x6c\xf5\x90\xc4\xe3\x8e\xb7\x4a\xa4\x98\xa3\xd6\xb4\xb4\xf6\xc5\x42\x64\xd6\x6b\xe4\xcc\x62\xb0\x4c\x20\xe0\xa7\xf0\xff\x9b\xab\xcb\xa4\xa2\x4a\xa3\xdf\x4f\x72\x6a\x68\x3c\x76\x36\x7a\xc5\x4c\x36\x0d\xeb\xfe\x66\xfe\x30\x40\x46\x35\x02\x99\x21\x56\x94\xb3\x25\x92\xe3\xb0\x0e\xb0\x5b\xb3\xd7\xdb\xad\x4e\x14\xd2\xd9\xee\xb7\x3f\xbd\xa8\x72\x6a\xf0\x4c\xae\x5b\xa7\xed\x4d\x0d\x55\x25\x9a\x33\xb9\x86\x14\x72\x99\x2d\xe6\xd6\xbd\x0d\x3d\x47\xfb\xf5\x6c\xf3\x3e\x0f\x9c\x58\xde\xf6\x61\x9f\x5d\x73\xf6\x59\x9a\x82\x58\x70\xde\x70\x76\x9e\xa7\x54\x94\x78\xca\x51\x99\x0f\xb8\x44\xde\x1c\xe8\xfb\x90\x24\xda\x50\xb3\xd0\xf5\x2f\x4e\xb5\xb9\xf0\xa1\x0c\x41\xb1\x7f\xdb\x7d\xa7\xde\x74\x4e\xd7\xb7\x67\x9f\x9c\x5f\x42\xf6\xbd\xda\x5b\x29\xb9\x82\xb4\xb9\x5b\xeb\x42\xfa\x6c\x73\xce\xa9\xd6\x97\x74\x8e\x11\xf1\x30\x2e\x9b\x5a\x08\x4a\xae\xf6\x4f\xdc\xd2\xd2\xd9\x77\x4d\xde\x75\x99\xc7\x84\x40\xf5\xee\xf6\xe2\x03\xa4\xd0\x26\x34\x6e\xa1\x3c\x20\xeb\x12\x75\x64\xd9\x5a\x07\xda\x6c\x38\x26\x39\xd3\x15\xa7\x1b\x48\x81\x08\x29\x90\x6c\x01\xb9\xc6\x47\x2c\x0c\x9d\x70\x1c\x28\xb9\x22\xdb\x9f\x47\x07\xd7\x15\x53\x78\x5a\xd8\x24\xff\xa7\x21\x6a\x61\xfd\x8e\x38\xb5\xe0\x0e\x07\xab\xcd\xfd\xdf\x8e\xd8\xa1\x42\xc9\x91\xe3\xfd\x42\xd9\xad\x1d\x2c\x85\x43\x28\x5e\x74\x3e\xd2\xb2\x5d\xad\x0f\x24\xa9\xe3\xf9\x6c\xc7\x5e\x98\x6c\x25\x19\x29\xb9\x61\x55\xa7\x16\x12\xab\x4f\xef\xe4\x12\x55\x64\x58\x55\x3f\xc5\xa6\x70\x7f\x52\xb5\x24\x40\x11\xef\xc6\xdb\xef\x3d\x11\xeb\xa7\xb5\xd5\x8a\xdb\x84\xcb\x6c\x46\x9c\x60\xb6\x89\x5c\x2d\x4c\xf4\x9b\x39\x10\xf2\x08\x85\x29\xcb\x73\x14\xa4\x16\x6d\x1b\x9b\x8c\xb3\x6c\xb6\x47\xe8\xdc\xae\x44\x2c\xaf\x49\x85\x36\xb2\x0b\xf4\x54\x61\x61\xc1\x86\x13\xb9\x1e\x12\xe8\x01\xcb\x5b\x6d\x60\x2e\x97\x68\x51\x1a\xc8\xe6\x41\xd7\x98\x4f\xba\xa6\xcb\x87\xdd\x25\x2a\xaa\x50\x98\x4b\x99\x63\xa2\x9c\x8f\xf3\x29\xe3\x79\x10\xbf\xa6\x0d\xed\xb4\xdb\x75\x1c\xab\x7d\x7f\xd0\x71\x43\xa5\x25\xed\x3b\x22\x99\x09\x2d\xe9\x7f\xd4\x60\x14\x5b\x1a\xb7\x6c\x8e\x2e\x91\x6c\x05\x05\x14\x78\xfe\x1c\xa2\xf0\xb5\x07\xff\x1d\x8d\x46\xa3\x18\x4e\xec\xe9\x18\xbe\x3f\xcc\x42\xb0\xb5\x11\xcc\x53\xc8\xcc\x53\x6f\x4d\xbc\x7a\x0f\x26\x54\x91\xf6\xfd\x33\x2b\x1e\x1f\x98\x36\xe1\xfa\x11\xa1\xf3\x09\x2a\xd2\x27\xa5\x42\x14\xfe\x73\x43\xfa\x44\xc8\x4f\xae\x2b\x91\x3e\x51\x98\x3f\x82\x41\xf3\x3c\x0a\x07\xdb\x06\x8f\x2a\xba\xef\x1f\xe4\x81\xfc\x10\x32\x76\x51\x8a\x6c\x6f\x95\x05\xcc\xa8\x97\xc6\x85\xc8\xb1\x60\xc2\xfa\x87\xef\x90\x71\xa4\xca\xc6\x54\x2e\x4c\x34\xa3\xf1\x18\xb6\xf6\xd8\xcc\x46\x46\xfb\x68\xdb\x1d\x97\x6c\xbb\x5e\x1f\x37\xda\xfa\x3b\x42\xf0\x08\x92\x0b\xc4\xce\xfe\xa1\xe5\xaf\x47\xe4\xed\xf5\xf5\xd5\xf5\x31\x5c\xca\x26\xd9\x34\x68\x26\x32\x04\x5b\x2b\xf3\x8d\x4b\x2e\x9f\x19\x31\xf4\x80\x24\xa4\xf6\xba\xed\xc3\x2b\x78\x01\x47\x36\xb5\xea\x84\xfe\xa8\x98\x30\x60\xd8\x1c\x81\x09\x98\x6f\xa0\x52\x58\xa0\x52\x98\x83\x1b\xea\x4c\x93\xdc\x01\xda\xb4\xc5\x24\xf4\x03\x37\x58\xc0\xb3\x87\xf3\x84\x6a\xe7\xbe\x09\xda\x19\xc4\xff\xb0\x51\xad\xaf\x6e\x03\xcd\x42\x09\x50\x89\x91\xef\x6f\xae\xc2\xd4\xd9\x50\xa7\x39\x50\xd0\x6e\xb5\x61\x59\xd1\x3c\x12\x7d\x58\xb1\xdc\x4c\xfb\x70\xe7\xc9\xdc\x41\x0a\x77\xf0\xe3\x07\x74\x47\x5d\xcb\x58\x58\x8f\xd0\x83\xae\xfb\x15\xdc\x88\x84\xa3\x28\xcd\x14\x5e\xa7\xfe\x38\xbc\x01\x01\xc7\x8e\xda\xa9\x52\x74\x13\xf9\xd5\x41\x63\xd9\x83\xa3\x38\xf9\x26\x99\x88\xee\x6c\xac\xc5\x4e\x27\xce\xdd\x14\x05\xd4\x8e\x51\xc0\xed\x1c\x05\xb2\x00\xba\x2f\x5c\x8f\x8c\x5a\x7d\xa8\x27\xac\x90\x05\x75\x40\x9d\x5a\x7c\x79\x62\x56\xda\xb4\xc9\x71\x7d\x55\x44\x1e\xcd\xcf\xdc\x83\x23\x5b\x30\x7e\xc5\x66\x93\x3b\x1b\xa2\xfd\x77\x65\x00\xb5\xc7\x70\xac\x0f\x2d\x9f\xd0\x2e\x85\x83\x65\x10\x68\xec\x59\xfc\x72\xfa\x87\x8d\x27\x81\xd8\x1a\xf0\xf4\xf2\x03\x40\x3e\x99\x9f\x46\xa7\x52\xb8\x64\x72\xa1\xc3\x9c\xab\x6b\x38\x8d\xca\x9c\xe6\xdf\x68\x86\xc2\x58\xdc\xa8\x4b\xed\x08\x34\xc1\x92\x89\x6e\x1f\xc8\x09\x67\xaf\x5b\x45\xe9\xca\xf1\xd8\x95\xa9\x0f\x46\x62\xe4\xa7\xaa\x42\x75\x4e\x75\xd8\x85\xc8\xd9\x87\x57\x93\x1e\x90\xf8\x64\x68\x41\x9a\x0e\x74\x41\x67\x08\x13\x56\xba\x77\x20\x23\xa1\x60\x06\xa8\x86\x39\x15\x1b\xbb\x5c\xa2\x36\x76\x0b\xb5\x5d\x5d\x31\xce\x9d\x85\x99\x22\x64\x0b\x65\x7b\x5c\xdd\x6f\x9b\x3c\x54\xac\x9c\x9a\x1b\x76\x87\x67\xac\xb4\x6d\xb4\x5d\xe0\x74\x49\x19\xb7\x03\xd9\x67\x97\xfb\x29\x5c\x50\x33\x4d\x0a\x2e\xa5\x8a\xea\x17\x40\x17\x57\xbf\x3f\xf8\xcf\x28\x86\x21\xbc\x3a\x7a\x19\xc3\x0b\xfb\x31\xde\x21\xb9\xe2\xf1\xef\x2c\xd1\x3d\xd4\xd7\xa9\x3f\xf1\xe6\x9e\xbb\xe3\x1d\xc0\x63\xed\xac\x3b\x61\xe5\x60\x22\xd7\xdd\x38\xcc\x91\xab\x40\x73\xe7\xad\x07\xa4\x5a\x93\x86\x46\xb1\xe0\xfc\x73\xd8\x44\x7d\xb8\x53\xee\x3d\xfd\xae\x3d\xe1\xe0\xba\x21\x5b\x0a\xa9\x20\xb2\x58\x0c\x52\x18\x8d\x81\xc1\xc9\x3d\xd8\xa0\x0f\x63\x60\xbd\x5e\x5b\x0e\xf7\xad\xbe\xb0\xaf\xf7\x58\x47\x3b\xda\x03\x78\x19\xef\x71\xdf\x76\xb6\x9d\xbf\x02\x00\x00\xff\xff\x66\x08\xf7\x66\xb1\x0f\x00\x00")

func bindataScriptsJsBytes() ([]byte, error) {
	return bindataRead(
		_bindataScriptsJs,
		"/scripts.js",
	)
}



func bindataScriptsJs() (*asset, error) {
	bytes, err := bindataScriptsJsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "/scripts.js",
		size: 4017,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1703686311, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}

var _bindataStandardCss = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x56\x51\x53\xc3\x36\x0c\x7e\xef\xaf\xd0\xc1\xf1\xd2\x6b\xda\xb4\xa5\x0c\xc2\xd3\x76\xfb\x07\x1b\x7b\x57\x62\x25\xf5\x70\xec\x9c\xad\x94\x76\x1c\xff\x7d\x97\x38\x29\x6e\x49\x60\x94\x3d\xfa\x8b\xf4\x59\xfa\x24\x4b\x59\x4c\x41\xa1\xe3\xa7\x4a\x20\x13\x64\x0a\x9d\x03\x98\x2e\x26\xd5\xfc\x1d\x16\xf0\x3a\x01\x00\xc8\x8d\xe6\xc8\xc9\x7f\x28\xb9\x8d\x6f\x1e\x5b\xa8\x44\x5b\x48\x1d\xa5\xca\x64\xcf\x11\x69\x91\x00\xd6\x6c\x06\xbe\x39\x46\xcb\xe1\xd7\xca\x38\xc9\xd2\xe8\x04\x30\x75\x46\xd5\x4c\x1e\xb7\xb2\xd8\x72\x02\xeb\x6a\xef\xcf\xa9\x61\x36\x65\x02\xcb\x06\x78\x9b\x4c\x84\xdc\x5d\x3b\x46\xae\x5d\x94\xa2\x85\x4f\xc3\x84\x5f\x36\x37\xad\xd3\xbc\x94\x99\x35\x43\xc6\x42\xba\x4a\xe1\x21\x01\x6d\x34\x79\x5b\xf1\x1d\x63\x57\xa2\x52\xff\xcd\x76\xb2\x98\x42\x49\xce\x61\xd1\x0b\xdd\xea\xdc\x43\xe7\xc1\xdf\x5d\xaa\x71\x3c\xdf\x50\x79\x92\xf6\xe9\x15\x9f\xa4\xfc\x95\x61\x93\x02\x4b\x56\x41\x02\x73\x7f\x7e\xfd\x41\x3b\x0c\x54\x35\x24\x3d\xe9\x80\xb3\x02\x2f\xe3\x5e\x25\x45\x79\xd8\x36\x43\xed\xd5\x25\xe0\x88\x59\xea\x02\x04\xe5\x58\x2b\x6e\xe9\x20\xc7\x52\xaa\x43\x5b\x91\x19\xa4\x46\x1c\x66\xc0\x98\x36\x31\xb4\xb7\xf9\xcf\x09\x5c\xfd\x4e\x7f\xe3\xae\x86\x3f\x50\xbb\xab\x19\xfc\x45\x56\xa0\xc6\x19\x38\xd4\x2e\x72\x64\x65\xfe\xd8\xdd\x92\x9a\x3d\x64\x46\x19\xeb\xa5\x22\x2f\x56\x61\x89\x74\x9f\x17\x66\xcf\x85\x35\xb5\x16\x51\x6b\x98\x5c\xaf\xd7\x59\xb6\x5e\xfb\x92\x58\x12\x33\x98\x6b\xd3\xbd\xcc\x31\x97\x2c\x8b\xe3\x38\xf6\x2e\x58\xa6\x64\x47\x2d\xf3\xfc\xfe\x68\x59\x58\x3a\x8c\x1a\xe2\xf2\x81\x1e\xb2\xa3\x5c\x4d\x22\x8d\xd8\x27\x79\xec\x15\xda\x81\x9e\x5d\x35\x8d\xd7\x40\x5b\x6a\x9f\xf1\x66\x19\xf7\x05\x79\x91\x82\xb7\x3d\xd0\xf6\xdc\x30\xc5\xd2\x77\x6f\x40\xb2\xda\xdc\x8e\x92\x5c\xcc\xd1\x01\xbe\xf9\x49\xc8\xba\x1c\x60\x59\xae\xce\x79\x96\xab\xbb\x51\x9e\x9f\xd0\x74\x40\x1b\x8e\x1f\x29\xe7\x34\xf1\xfc\xfe\x8c\xe4\x6e\x35\xca\x71\x29\x85\x3f\x07\x13\xe1\x6b\x8a\x75\x3c\x46\x71\x29\xc3\xfa\x58\x5d\x3f\x11\x3e\x2a\x7a\x46\x00\x47\x8f\xe3\xc3\x0b\x5a\xb5\x39\x7b\x06\xdf\xdf\xf0\xb2\x95\xfd\xb2\xc9\x95\x41\x4e\x9a\xd1\x71\x3a\xa3\xfc\xc4\x39\x8a\xd3\xc1\x7e\x35\xad\x3e\x8e\x18\x4b\x0a\x59\xee\x3a\x52\xa6\x3d\x47\xa8\x64\xa1\x13\xc8\x48\x33\x59\x9f\x4d\x2a\x8b\xe8\x3d\x18\xcf\x19\x8c\xc0\x76\x3d\xe0\xfe\xcf\xdf\x9e\x4e\xb6\x83\x47\xc6\x06\xf2\x62\x0a\xb4\xaf\xa4\xa5\x5f\x73\x26\x1b\x3a\x86\xf0\x27\xde\x52\xe7\xa6\x9b\x74\x8d\x58\xc1\x71\x4c\xb3\x0a\x85\x90\xba\x88\xfc\xc0\xdd\xf4\xb5\xf2\x9e\xfd\xf6\x0b\x35\xf0\xf2\xbe\xdb\x6c\x07\x6c\x5a\x69\x3d\xff\x8e\x2c\xcb\x0c\x55\xff\x89\x4d\x15\x78\xd7\x7d\x5b\x2b\xe9\x38\x72\x7c\x50\x14\xf1\xa1\xa2\x3e\xab\x30\x42\xa9\x95\xd4\x74\x5c\x8b\x67\xd5\x64\x53\x7d\x04\xfb\x5d\xd3\x37\x54\xb3\xef\x8c\x51\x2c\xab\x70\xe3\x75\xc8\x77\xd6\xd3\xff\xff\xa3\xe4\x1f\x4b\x77\xc5\xdb\xe4\xdf\x00\x00\x00\xff\xff\x2b\xec\xfb\x9d\xc9\x09\x00\x00")

func bindataStandardCssBytes() ([]byte, error) {
	return bindataRead(
		_bindataStandardCss,
		"/standard.css",
	)
}



func bindataStandardCss() (*asset, error) {
	bytes, err := bindataStandardCssBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "/standard.css",
		size: 2505,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1703686736, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}


//
// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
//
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

//
// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
// nolint: deadcode
//
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

//
// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or could not be loaded.
//
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

//
// AssetNames returns the names of the assets.
// nolint: deadcode
//
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

//
// _bindata is a table, holding each asset generator, mapped to its name.
//
var _bindata = map[string]func() (*asset, error){
	"/scripts.js":   bindataScriptsJs,
	"/standard.css": bindataStandardCss,
}

//
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
//
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, &os.PathError{
					Op: "open",
					Path: name,
					Err: os.ErrNotExist,
				}
			}
		}
	}
	if node.Func != nil {
		return nil, &os.PathError{
			Op: "open",
			Path: name,
			Err: os.ErrNotExist,
		}
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

var _bintree = &bintree{Func: nil, Children: map[string]*bintree{
	"": {Func: nil, Children: map[string]*bintree{
		"scripts.js": {Func: bindataScriptsJs, Children: map[string]*bintree{}},
		"standard.css": {Func: bindataStandardCss, Children: map[string]*bintree{}},
	}},
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
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
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

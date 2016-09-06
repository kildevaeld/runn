// Code generated by go-bindata.
// sources:
// docker.js
// DO NOT EDIT!

package docker

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

var _dockerJs = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xac\x91\xc1\x4e\xc3\x30\x0c\x86\xcf\xc9\x53\xf8\x96\x56\xaa\xfa\x02\x53\x6f\xf0\x1c\x55\x59\x3d\x35\x22\xc4\xc1\x71\x36\x10\xda\xbb\x93\x26\xa3\x0c\xa4\x72\xe2\x66\x5b\xbf\xff\xef\x4f\xac\xcf\x13\x43\x5c\x60\x00\xc6\xd7\x64\x19\x1b\x13\x17\x74\xce\xb4\x07\xad\x4f\xc9\x1f\xc5\x92\x87\x07\x3a\x3e\x23\x37\x14\xd6\x2e\xb6\xf0\xa1\x95\x2c\x36\xf6\xb7\x41\xde\xbe\x55\x07\x7d\xd5\xc5\x32\x30\x09\xe5\x79\xdd\xec\x4b\x2b\xef\x01\xb3\x6b\xa9\xfb\xa7\x64\xdd\x9c\x05\x1b\xe3\x87\xbb\x3d\x6d\x3d\x0c\x03\xf8\xe4\x5c\x0b\xb2\x30\x5d\xc0\xe3\x05\x1e\x99\x89\x1b\xe3\xe9\x8b\xbb\xc6\x55\xe3\x18\xd8\x9e\x27\xc1\x71\xae\x79\xef\x33\x76\x60\x0a\xd2\x74\xb0\x4d\xbe\xd9\xc8\xdc\xe5\x1f\x88\xc9\xc9\xca\xd7\xea\xda\x96\xa7\xd4\xac\x9c\xfc\x6e\xd2\xff\x83\x2a\xf5\x0b\x1b\x85\xc2\x1e\xf7\x5e\x35\xb1\xfc\x25\x7b\xa1\x39\x39\xec\xf1\x2d\x50\x11\xd6\x93\xe4\x43\x7c\x06\x00\x00\xff\xff\x22\x65\x78\xa3\xfc\x01\x00\x00")

func dockerJsBytes() ([]byte, error) {
	return bindataRead(
		_dockerJs,
		"docker.js",
	)
}

func dockerJs() (*asset, error) {
	bytes, err := dockerJsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "docker.js", size: 508, mode: os.FileMode(420), modTime: time.Unix(1472772261, 0)}
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
	"docker.js": dockerJs,
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
	"docker.js": &bintree{dockerJs, map[string]*bintree{}},
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


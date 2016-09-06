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

var _dockerJs = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xbc\x55\x4d\x8f\xdb\x36\x10\x3d\x8b\xbf\x62\xba\x40\x41\xaa\xb6\xf8\x07\x04\x1f\x8a\x7e\x00\xe9\xa1\x05\xd2\xe3\x66\x0f\x8c\x4c\xdb\x4c\x24\x4a\x25\x29\x77\x17\xb6\xfe\x7b\x67\xa8\x0f\x53\xf6\x16\x08\x50\xa0\x41\x20\x59\x9c\xf7\xde\xbc\x21\x67\xb8\x67\xe5\xc0\x9f\x60\x07\x4e\xff\xd5\x1b\xa7\x05\xf7\x27\x9e\x6f\x59\x76\x68\x5d\xa3\x42\x1a\xe8\x83\xa9\x79\x2e\xc7\x00\x63\x91\x59\x9b\x4a\x23\xe6\x47\xe7\xd4\x9b\xec\x5c\x1b\xda\xf0\xd6\x69\x19\xd7\x4b\xc6\xd8\xa1\xb7\x55\x30\xad\x85\x9f\xdb\xea\xab\x76\x20\x72\xb8\x30\x36\x30\xa6\x5f\xbb\xd6\x05\x2f\xa7\xf5\xdd\x04\x28\x13\x4a\xa5\xea\x5a\x54\xcd\x7e\x0b\xca\x1d\x3d\x11\xb3\x98\x13\xc1\xa3\x07\xc1\xf7\x23\xfb\x7b\x8f\xff\xf9\x16\x16\xb0\xfc\xd2\x1a\x2b\x38\xf0\x3c\x07\x96\x39\x1d\x7a\x67\xb1\x4c\xa9\x5f\x75\x25\x7c\x2e\x7d\xd8\xb7\x7d\x20\x1f\xa4\x18\x7d\x2f\x16\x6e\x65\x30\x16\x7f\x22\x5a\x39\xda\x89\xc5\x99\xb0\xaa\xd1\x5b\x30\x8d\x3a\xe2\xab\xed\x68\x31\x1a\x04\xfc\x37\x7d\x22\x7e\xfa\x75\xbd\x5e\x86\x32\x86\xe2\x83\x32\x1e\x6a\x75\x24\xc4\xf3\x4b\x5c\xc2\x72\x40\xd0\xfa\x57\xfd\x06\xc6\xde\x2b\xce\xac\xb3\xaa\x6f\xaa\xcf\x88\x7d\x29\x97\xb8\x39\x80\x20\xf6\x6e\x07\x5c\xdb\xf3\xaf\xa6\xd6\xfc\x7a\x85\x64\x89\x43\xf2\xdd\xf5\x9f\x6b\x43\x27\x9d\xa4\x98\x65\x16\x3e\x21\x91\x90\x8f\xac\x28\x52\x1c\x28\xf0\x40\xf9\x6e\x6c\x00\xe3\xe3\x1b\x4b\xa9\xf3\x7b\xe5\xb1\x08\x2a\xe0\x19\x5f\x2f\xab\xd8\xc0\x1e\x51\xf8\x94\x8d\xea\xc4\x6d\xd3\xcf\x28\x09\xd3\x59\xf2\xa2\xe0\xb0\x89\xc6\x36\x80\xe7\x8c\xcf\x73\x09\x43\x5e\xae\x94\xe2\x36\xcb\xae\xf7\x27\xb2\x94\x34\xc5\x0d\x36\x80\xae\xbd\xbe\xf3\x9a\xf0\x28\xd1\x66\x95\x06\x6b\x4b\xe8\x8c\x8d\xaf\xd8\x9a\xd8\x7f\x25\xcb\xe2\x86\x84\x93\xf1\xf2\xa4\xfc\x4f\xad\x0d\xca\x58\xed\x62\xcf\xc4\x5d\xc9\x32\xc4\x3d\x36\xb1\xeb\x2d\x14\x7b\xea\xe5\xa2\x20\xec\xd2\xd5\xa3\x9b\xc5\xfc\x16\x92\xee\x43\x2b\xd9\x54\x02\xa5\x8d\x59\x8d\xff\xd8\x5b\x6b\xec\x31\x4d\x79\x37\x03\x73\x52\x63\x7d\xa7\xab\x00\xc5\x01\x2e\x17\xf9\x61\x3f\x0c\xb1\xc8\x48\x9c\xa6\x44\x06\x67\x1a\x91\x2f\x79\x6e\x05\x3c\x4d\x22\xe3\x7c\x3c\x4d\x3c\x74\x94\xcd\x69\x10\x47\x0e\x27\xfc\x37\x24\xc5\x2c\x2c\xab\xb0\xbd\xdb\x5a\xcb\xba\x3d\x0a\xee\xc6\x5a\x60\xe2\x8e\x53\x9e\x3f\x0c\x35\x2d\xae\x0d\x97\x34\xde\xf3\x00\xb7\xdd\x6a\x7e\xe3\x9e\x4c\x0a\xf1\x9a\xe1\x04\x41\xed\x78\x71\xc9\xb8\x84\xf7\x48\xdf\x68\x1b\x7c\x9e\x4a\xb9\xe6\xe1\x22\x78\x14\x43\x50\x71\x40\xb5\x67\x0a\xbf\xac\xe9\xe6\x9b\xf8\xe6\x5d\xf6\xe7\xde\xd4\xfb\x15\xbf\x53\xe1\xb4\x85\xa0\x8e\x8f\x22\x11\x4c\x32\xd8\xc2\x08\xe0\x11\xb6\x05\x62\xac\x44\x97\x76\x79\xdf\x18\xf5\xb5\xd3\x14\xb4\xfa\x6f\xf8\xa8\x8f\xbf\xbc\x76\x82\x7f\xfa\xe4\x7f\x10\x7c\x43\xa8\x0d\xcf\xf1\x0b\x9f\x4b\x7e\x84\xcb\xa0\x7d\x10\xf3\xe1\xcc\x8d\xd2\xf9\xa7\xf9\x90\x72\x36\x64\xb3\x83\x74\x4c\xfe\x0f\x13\x50\xa8\xd4\x47\x62\xe3\x03\xcd\xd4\x7f\xb7\x40\xd8\x36\xfe\xe9\xfc\x17\x0f\x71\x76\xd3\xcd\x98\x5d\xb7\xe9\xc1\x4c\x63\xf2\xae\x9d\xe0\xde\xe2\x24\x52\xaa\x2f\x1e\x43\xbb\xe9\xd8\x27\x56\xda\x3f\xcb\xf4\xff\xf6\xe7\x1f\xbf\xcb\x4e\x39\xaf\x05\x71\xc6\xcb\xa3\x52\xa1\x3a\x81\xd0\xe9\x85\x9d\x4e\xa0\xce\x6f\x02\xb6\xaf\x6b\x22\x91\xc9\x7f\x02\x00\x00\xff\xff\xf9\xe2\xb8\x72\x3a\x08\x00\x00")

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

	info := bindataFileInfo{name: "docker.js", size: 2106, mode: os.FileMode(420), modTime: time.Unix(1473191200, 0)}
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


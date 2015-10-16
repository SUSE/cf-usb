// Code generated by go-bindata.
// sources:
// schemas/config.json
// schemas/dials.json
// DO NOT EDIT!

package driverdata

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

var _schemasConfigJson = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x9c\x90\xc1\xea\xc2\x30\x0c\xc6\xcf\xed\x53\xfc\xc9\xf9\xff\x04\x5e\x1d\x78\xf7\x2a\x32\xdc\x1a\x46\x14\x93\x9a\x56\x41\x64\xef\x6e\x90\xb2\x0a\x6e\x1e\x84\x36\x87\x2f\xf9\xe5\x6b\xbf\x87\x77\x90\xef\x11\x61\xf5\x07\xd2\x1d\xb1\xcf\x60\x4a\x54\x89\xa8\x99\x30\x99\x6e\x23\x0e\x3a\x95\x13\xea\x46\xb6\x18\x28\x35\x4a\x37\xd4\xd2\xaa\x7c\xca\x4a\x3c\x18\xef\xc6\xff\xca\x34\xd2\x5b\x5d\x0b\xb3\x2d\x27\xe1\x0f\x6c\xb2\x75\x33\xc6\xa6\x85\xd7\x82\x16\x39\x44\x21\xce\x53\x63\xce\xb8\x58\x57\x8a\xce\x87\x01\x7f\x40\x5a\xfb\x60\x7a\x7b\xed\x12\xea\x4b\xb1\x6b\x07\x14\x2f\x57\x52\x0c\x36\xb7\x5b\x8a\xed\x5b\x36\xde\xed\xfd\xf8\x0c\x00\x00\xff\xff\xbc\xda\xb8\xab\x94\x01\x00\x00")

func schemasConfigJsonBytes() ([]byte, error) {
	return bindataRead(
		_schemasConfigJson,
		"schemas/config.json",
	)
}

func schemasConfigJson() (*asset, error) {
	bytes, err := schemasConfigJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schemas/config.json", size: 404, mode: os.FileMode(436), modTime: time.Unix(1444987012, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemasDialsJson = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xaa\xae\x05\x04\x00\x00\xff\xff\x43\xbf\xa6\xa3\x02\x00\x00\x00")

func schemasDialsJsonBytes() ([]byte, error) {
	return bindataRead(
		_schemasDialsJson,
		"schemas/dials.json",
	)
}

func schemasDialsJson() (*asset, error) {
	bytes, err := schemasDialsJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schemas/dials.json", size: 2, mode: os.FileMode(436), modTime: time.Unix(1444908984, 0)}
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
	"schemas/config.json": schemasConfigJson,
	"schemas/dials.json":  schemasDialsJson,
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
	"schemas": &bintree{nil, map[string]*bintree{
		"config.json": &bintree{schemasConfigJson, map[string]*bintree{}},
		"dials.json":  &bintree{schemasDialsJson, map[string]*bintree{}},
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

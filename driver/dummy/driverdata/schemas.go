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
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path/filepath"
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
	name string
	size int64
	mode os.FileMode
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

var _schemasConfigJson = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xaa\xe6\xe2\x54\x2a\xa9\x2c\x48\x55\x52\xb0\x52\x50\xca\x4f\xca\x4a\x4d\x2e\x51\xd2\x01\x8a\x15\x14\xe5\x17\xa4\x16\x95\x64\xa6\x16\x83\x64\x80\xaa\xe0\x42\x95\xf1\xf9\x79\xa9\x4a\x10\x31\x24\xbd\xc5\x25\x45\x99\x79\xe9\x4a\x40\xc1\x5a\x1d\x14\xd5\x25\xe5\xf9\x78\x55\x73\x81\x35\x28\x15\xa5\x16\x96\x66\x16\xa5\xa6\x80\xa4\xa3\x51\xed\xd2\x51\x40\x35\x2d\x96\xab\x96\x0b\x10\x00\x00\xff\xff\xb9\x92\x44\x6e\xb8\x00\x00\x00")

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

	info := bindataFileInfo{name: "schemas/config.json", size: 184, mode: os.FileMode(420), modTime: time.Unix(1444902332, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _schemasDialsJson = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xaa\x56\x50\xe0\x52\x50\x50\x50\x2a\xa9\x2c\x48\x55\xb2\x52\xca\x4f\xca\x4a\x4d\x2e\x51\xd2\x01\x8b\x15\x14\xe5\x17\xa4\x16\x95\x64\xa6\x16\x2b\x59\x41\xd5\x81\x84\x73\x13\x2b\xe2\x53\x92\x8a\x33\xab\x52\xe3\x73\x93\x90\x65\x90\xcc\xc9\xcc\x2b\x49\x4d\x4f\x2d\x82\x18\x04\xd3\x97\x99\x97\x99\x5b\x9a\xab\x64\x65\x00\x15\xac\x05\xd1\xb5\x10\xbb\x8a\x52\x0b\x4b\x33\x8b\x52\x53\x94\xac\xa2\x71\xd9\x04\x12\x8d\xe5\xaa\x05\x04\x00\x00\xff\xff\xca\xaa\x24\x1c\xb2\x00\x00\x00")

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

	info := bindataFileInfo{name: "schemas/dials.json", size: 178, mode: os.FileMode(420), modTime: time.Unix(1444902332, 0)}
	a := &asset{bytes: bytes, info:  info}
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
	if (err != nil) {
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
	"schemas/dials.json": schemasDialsJson,
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
	Func func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"schemas": &bintree{nil, map[string]*bintree{
		"config.json": &bintree{schemasConfigJson, map[string]*bintree{
		}},
		"dials.json": &bintree{schemasDialsJson, map[string]*bintree{
		}},
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


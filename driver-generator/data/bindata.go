// Code generated by go-bindata.
// sources:
// templates/Makefile.template
// templates/config.go.template
// templates/driver.go.template
// templates/main.go.template
// DO NOT EDIT!

package data

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

var _templatesMakefileTemplate = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x74\x90\x5f\x4b\xf3\x30\x18\xc5\xaf\x9b\x4f\x71\xde\x52\xb6\xf6\x22\x2b\xef\xed\xa4\xe2\x5f\x14\x95\x55\xd4\xcb\x82\xa4\xcd\x63\x1b\x6c\x92\xb2\x46\x99\x8c\x7d\x77\xd3\xb5\x63\x28\x78\x15\x72\xce\x93\x93\xf3\x7b\x98\x68\xdb\x25\x6a\x32\xa5\x32\x52\x38\x81\xf2\x43\xb5\x92\xb1\xfd\xb1\x64\xc1\x19\x55\x8d\x45\x18\xc5\xf9\xfd\xeb\x65\xfe\x90\x3f\x25\x59\x76\x8a\x8b\xc1\x55\xa6\x8e\xe2\x55\x3e\xc9\x21\x0b\x68\xd3\xd9\xb5\xc3\x4d\xfe\x78\xfe\x72\x9b\x45\x71\xdf\x50\xdb\xa2\xb6\x92\x3a\x74\xc2\x35\xc9\xf2\xa0\xed\x43\xa3\x68\x9c\x4c\x30\x9b\x15\x2c\xa8\xed\xf8\x37\xb8\xc5\x76\xbb\xb8\x5a\xab\x4f\x5a\xaf\x84\xa6\xdd\x0e\x95\x96\xa9\xdc\x0b\xe9\x2f\x2b\xd5\x42\x99\x45\x6d\x19\x3b\x32\xfc\x59\xfb\x5a\x97\x24\x87\xde\xb8\x7b\xce\x57\xe8\xab\x86\xb4\xe8\xa1\x8c\xb3\x18\xe3\xfb\x9f\x44\x6f\x3e\x70\x72\x52\x70\x2d\x36\x1e\xc5\x35\xf8\x0f\xee\xbe\x3a\x82\x44\x11\xe3\x1f\xb8\xf1\x55\xa6\x31\x14\x09\x38\x6d\xa8\x82\x47\x2a\x45\xdf\x80\x57\x08\xe3\x4a\x62\xbe\xdd\xcd\x3d\xa9\xdf\x07\x3f\xec\x9a\x77\xef\x75\x16\x8e\x0f\x07\x21\x1c\xd8\x8f\xd7\x74\x2a\xe8\xf1\x0e\x5d\x53\x24\x21\x8a\x93\x80\xb1\xe0\x3b\x00\x00\xff\xff\xac\x5a\xeb\x77\xba\x01\x00\x00")

func templatesMakefileTemplateBytes() ([]byte, error) {
	return bindataRead(
		_templatesMakefileTemplate,
		"templates/Makefile.template",
	)
}

func templatesMakefileTemplate() (*asset, error) {
	bytes, err := templatesMakefileTemplateBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/Makefile.template", size: 442, mode: os.FileMode(420), modTime: time.Unix(1456843431, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesConfigGoTemplate = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x7c\xcc\xb1\x0a\x02\x31\x0c\x80\xe1\xf9\xfa\x14\xe1\x1e\xc0\x07\x70\xd5\xc1\x49\x5c\xdc\xaf\xd4\x58\xa3\xd8\x94\xa4\xa7\x48\xe9\xbb\x1b\x2a\x0e\x3a\x5c\xa6\x90\x7c\xfc\xd9\x87\x9b\x8f\x08\x81\xd3\x99\xa2\x73\xe5\x95\x11\x6a\x5d\x6d\x85\x1e\x28\x7b\x7f\xc7\xd6\x3e\xfb\xa6\x0b\xd0\x22\x73\x28\x50\xdd\x70\x54\x14\xf8\x1d\x7b\x52\x8a\x30\x5d\x95\xd3\x7a\x9c\x0d\xd0\x69\x9c\xdc\x70\xf0\xaa\x8b\x34\x1b\x78\xb2\x74\xbc\x63\x2d\x8b\xf8\x62\xa0\x57\x59\xfe\x21\xa5\x7e\xf9\x56\x0d\x18\x6c\xee\x1d\x00\x00\xff\xff\xad\xb3\x78\x78\xe6\x00\x00\x00")

func templatesConfigGoTemplateBytes() ([]byte, error) {
	return bindataRead(
		_templatesConfigGoTemplate,
		"templates/config.go.template",
	)
}

func templatesConfigGoTemplate() (*asset, error) {
	bytes, err := templatesConfigGoTemplateBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/config.go.template", size: 230, mode: os.FileMode(420), modTime: time.Unix(1456504424, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesDriverGoTemplate = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xcc\x95\xcf\x6e\xa3\x3a\x14\xc6\xd7\xe1\x29\x7c\x59\x41\x14\x60\x5f\xa9\x8b\x7b\xcb\xd5\x28\xd2\x4c\x55\xb5\x4f\xe0\x9a\x13\xea\x29\xb1\x19\xdb\xa4\x33\x8a\x78\xf7\xf1\xbf\x10\xda\x10\x42\x92\x8e\x34\xab\x52\xfb\xf3\xe7\xef\xfc\x7c\xec\xd4\x98\xbc\xe2\x12\xd0\x76\x9b\xe6\x82\x6e\x40\xdc\xe3\x35\xb4\x6d\x10\xd0\x75\xcd\x85\x42\x51\x30\x0b\x81\x11\x5e\x50\x56\x66\xdf\x25\x67\x61\xa0\x47\x4a\xaa\x5e\x9a\xe7\x94\xf0\x75\xf6\x52\x93\x8a\x37\x45\x46\x56\x49\x23\x9f\xb3\xc2\x9a\x84\x5a\xa3\x1d\xff\xc3\x12\x96\xd6\xa7\x6d\xfd\x4c\xf6\x61\x23\x3f\x5c\x60\x85\xc3\x09\xc6\x99\x54\x58\x35\xf2\x83\xb4\xa6\x1b\xae\x70\x95\x94\xbc\xc2\x3a\x66\xa5\x0b\xd2\x11\xe2\x20\x50\xbf\xea\x83\xd2\xdc\x37\x92\x4a\x34\x44\xa1\x6d\x30\xab\x78\xa9\xf5\xc8\xae\x4a\xbf\xda\x7f\x02\x0d\x60\xd5\x30\x82\xee\xe1\x6d\x70\x79\x34\xb0\x28\x46\x2e\xa3\x97\x1b\x6b\x01\xaa\x11\x6c\x38\xc1\xd6\x59\xdc\x20\xf7\x37\x7d\x02\x29\x29\x67\x51\xf8\x41\x9d\x78\xa4\x71\xdb\xa5\x8a\x8a\x61\xcb\x18\x3d\xe8\x63\x8a\x04\xfc\x68\x40\x2a\x34\x37\xe7\x95\x3e\xe2\xb7\x6f\xda\x5a\x07\x5d\x20\x01\xb2\xe6\x4c\x02\x9a\x3f\x73\x5e\xc5\x08\x84\xe0\x36\x68\x91\xfa\x14\x4b\xb6\xe2\x51\x58\x6b\x9b\xc4\xdb\x84\x0b\x5f\x65\xae\xcf\x68\x1b\xee\x46\x6f\x0c\x41\xb3\xdb\xdc\x8f\xc4\xad\x06\x3e\xab\x31\xa3\x24\x0a\xef\xb9\x42\xba\x83\x2a\x58\x03\x53\x50\x84\xf1\x2c\x98\xcd\xbb\xdd\x6f\x91\x86\x0f\x41\xc7\x87\xd1\x6a\x42\x6d\x5f\x40\xe5\x98\x56\xf2\x89\xbc\xc0\x1a\x77\x55\xba\x18\xfd\xda\xdc\xc8\xf1\xea\x4a\x50\x49\x61\x9c\x12\x69\xad\x4e\x56\xea\xbf\x6c\x81\xc5\x3e\xc2\xc2\xec\x80\x6e\x6e\xd1\xbe\x87\xd3\x7f\xa5\x04\x15\x85\xce\x58\x66\x05\xc5\x95\x4c\xed\xbd\x89\x83\x19\x5d\xd9\x15\xff\xdc\x9a\x8a\x4d\xb0\x1d\x00\x3d\x1a\xcc\x34\x81\x77\x8c\x3c\xdf\xde\x86\xf1\x05\xc8\xee\x38\x5b\xd1\xf2\x53\x98\x11\x6b\x75\x09\x34\xd2\x0b\x31\x81\x9a\x93\x5f\x87\xad\xbf\xe5\xd9\xdc\x1e\x04\xdf\x50\x73\x1b\x97\x4c\x3f\x38\x8c\x40\x47\xce\x5f\xf1\x03\xc1\xa3\x9b\xef\x23\xf5\xd2\x9d\x62\xe4\xb6\xed\xcc\x12\xea\xb5\x47\xe0\x76\xd3\xb4\xd8\x03\xee\x36\x58\xe6\x0b\x14\xba\xb2\x0f\x6f\x67\xea\xfa\x20\xd6\x12\xdb\x93\x03\x8a\xdc\x8c\x8f\xdf\x62\x83\xd1\x95\x97\x3e\xd9\x87\xd8\xf2\x36\x1f\xe9\x9d\x00\xac\x35\x17\x74\xe8\x31\xc6\xbd\xa9\x6b\xe8\x9a\xce\xfd\xf3\x5c\x5b\xd7\x63\x47\xe0\xe4\x1c\xa4\x86\xf9\xff\x4f\x2a\xd5\x18\xdd\xb3\xe1\x31\x10\x9a\xba\x66\x5f\x68\x0f\x73\x7e\x87\x10\x0f\x24\x03\x30\xa9\x0e\x20\x56\x98\xc0\xb6\x1d\x03\xe9\xac\x12\xb2\xf7\xba\x02\x68\xcf\xe4\x9d\xaa\x97\x74\x3a\xf9\xa3\x48\x2f\x78\x30\xc7\x60\xaa\x71\x8e\x5e\xd7\x13\x1d\xc7\xd9\x07\x00\xa6\x31\xfe\x1a\x98\xd3\xda\xf8\xf3\xa0\x3f\xc2\x86\xbf\x8e\x35\xf1\x81\xe0\x4a\xf4\xc2\xfa\x4d\xe8\xe3\xc9\x5c\x4f\x9d\xd1\x09\xae\x50\x81\x7b\x3b\xc7\x1e\xde\x33\x88\xe6\x50\x9f\xfa\x05\x1b\x90\x5c\xf3\xca\x16\x7b\xbb\x73\xda\x78\x1a\x17\x34\xfd\xc9\xfc\x1d\x00\x00\xff\xff\x97\x3b\xda\xd5\xd0\x0c\x00\x00")

func templatesDriverGoTemplateBytes() ([]byte, error) {
	return bindataRead(
		_templatesDriverGoTemplate,
		"templates/driver.go.template",
	)
}

func templatesDriverGoTemplate() (*asset, error) {
	bytes, err := templatesDriverGoTemplateBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/driver.go.template", size: 3280, mode: os.FileMode(420), modTime: time.Unix(1456504424, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesMainGoTemplate = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x64\x51\xcd\x6e\xe2\x30\x10\x3e\xc7\x4f\xe1\xf5\xc9\x91\x20\xb9\xef\x8a\x0b\x4b\x5b\x55\xaa\x50\x55\x54\xf5\xec\x86\xc1\x4c\x49\x6c\x6b\x62\xd2\x03\xca\xbb\x77\xec\x00\x87\xf4\x64\xeb\x9b\xef\x6f\xec\x60\x9a\x93\xb1\x20\x3b\x83\x4e\x08\xec\x82\xa7\x28\xb5\x28\x94\x83\x58\x53\x68\xea\xaf\xde\x3b\x3e\x15\x43\xbe\x57\x82\x8f\xcb\xa5\x5a\x9b\x1e\x9e\x33\x77\x1c\xeb\x3d\xe1\x00\x54\x33\xbc\xc9\xb7\xad\xe9\x60\x1c\x93\xc0\x62\x3c\x9e\x3f\xab\xc6\x77\xb5\x33\x11\x0e\xe8\x9a\x63\x1d\x10\x66\xb3\x80\x83\x8f\xa6\x5d\x5a\xdf\x1a\x67\xeb\x96\xfb\x90\x12\xa5\x10\x87\xb3\x6b\x72\x33\x5d\xca\x8b\x28\x06\x43\xb2\xf5\x96\xa7\x72\x25\x33\xab\xda\xc2\xf7\x4b\x46\xb4\x9a\xe5\x2f\xa7\x5a\x8a\x6d\x8a\x49\x54\xbd\x81\xc5\x3e\x02\xed\xd0\x9d\xf4\x5d\xff\x41\x78\xc3\x7c\x5f\xed\xe2\x1e\x88\x16\x57\xfb\xcd\xc3\xfa\xfd\xa9\x4c\x16\x41\xfe\x5d\x49\xae\x9e\x14\xaf\xe4\x07\x64\x9a\x4e\x83\x29\x26\x4d\x67\x05\x12\x73\x06\x4d\x77\x3d\xd5\x29\x45\x81\x07\xc9\x61\xd9\xf9\x5e\x2e\x31\x7f\x2d\xa3\x16\x72\xca\x29\xff\x65\xc5\x9f\x95\x74\xd8\xa6\x37\xb9\xed\xf6\x68\xf8\x05\xb5\xa2\xab\xcb\x32\xb4\x67\x8b\x8e\x75\x4c\xe7\xa4\x31\xad\x50\xed\x80\x06\xf8\xef\xf7\xd0\xe8\xeb\xb7\xa6\x92\x19\xa5\x0c\xf3\x42\xa3\xf8\x09\x00\x00\xff\xff\xea\x88\xda\x1f\x13\x02\x00\x00")

func templatesMainGoTemplateBytes() ([]byte, error) {
	return bindataRead(
		_templatesMainGoTemplate,
		"templates/main.go.template",
	)
}

func templatesMainGoTemplate() (*asset, error) {
	bytes, err := templatesMainGoTemplateBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/main.go.template", size: 531, mode: os.FileMode(420), modTime: time.Unix(1456504424, 0)}
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
	"templates/Makefile.template": templatesMakefileTemplate,
	"templates/config.go.template": templatesConfigGoTemplate,
	"templates/driver.go.template": templatesDriverGoTemplate,
	"templates/main.go.template": templatesMainGoTemplate,
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
	"templates": &bintree{nil, map[string]*bintree{
		"Makefile.template": &bintree{templatesMakefileTemplate, map[string]*bintree{}},
		"config.go.template": &bintree{templatesConfigGoTemplate, map[string]*bintree{}},
		"driver.go.template": &bintree{templatesDriverGoTemplate, map[string]*bintree{}},
		"main.go.template": &bintree{templatesMainGoTemplate, map[string]*bintree{}},
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

// Code generated by go-bindata.
// sources:
// swagger-spec/api.json
// DO NOT EDIT!

package data

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

var _swaggerSpecApiJson = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x9d\x5f\x93\x9b\x38\x12\xc0\xdf\xf3\x29\x28\xdf\x3d\x4d\x5d\xcc\xec\xd6\x3d\xed\xd3\xcd\x66\xb6\xae\x5c\x77\xb5\x9b\xaa\xa9\x3c\xdd\xa5\x66\x64\x23\xdb\xda\xc3\xc0\x4a\x22\x93\xdc\x54\xbe\xfb\x4a\x02\x6c\xfe\x48\x20\x40\x38\x30\xe9\xbc\x64\x8c\x5b\x8d\xd4\x74\xff\x68\xa9\x65\x78\x79\xe3\x89\x7f\x2b\xf6\x8c\x0e\x07\x4c\x57\x3f\x79\xab\x1f\xd7\xb7\xab\xbf\x65\x47\x49\xb4\x8f\xc5\xa1\x4c\x46\x1d\x09\x30\xdb\x51\x92\x70\x12\x47\x52\xf6\x43\x44\x3e\x61\xca\x50\xe8\x31\x4c\x3f\x91\x1d\xf6\xb6\x34\xfe\x1f\xa6\xde\x09\x45\xe8\x80\x4f\x38\xe2\xde\xdd\xfb\x4d\xae\x4f\x69\x90\xf2\x79\xeb\xdb\xf5\xed\xfa\x87\xf2\x77\x9c\xf0\x10\x2b\xbd\x0f\x3f\x97\x55\xa0\x84\x54\xc4\x30\x3d\xb1\xdf\xf6\x0f\xd9\x29\xa5\xfc\x91\xf3\x84\xfd\xe4\xfb\x07\xc2\x8f\xe9\x76\xbd\x8b\x4f\xfe\x31\xd9\x85\x71\x1a\xf8\xbb\xfd\xdb\x94\x6d\xcb\xcd\x77\x71\xc4\xd1\x8e\x57\x06\xa6\xbe\xc0\x27\x44\x42\xa5\x4e\x34\x0a\xf0\xa7\x7f\x1c\x13\x2c\x75\xad\xce\x62\x5f\x4b\x6a\x42\x71\xee\x88\xe1\xa6\x9a\x08\x9d\x54\xa7\xee\x12\xb4\x3b\x62\xef\x62\xcf\xb3\x44\x4a\xc3\xa2\xd7\xa2\xd3\xcf\xcf\xcf\x6b\xa4\x64\xd7\x31\x3d\xf8\xb9\x5e\xe6\xff\x7b\xf3\xee\x97\x5f\x1f\x7e\x79\x2b\x14\xac\x8f\xfc\x14\x96\xba\xf1\xa6\xd4\x19\x39\x1e\x96\x9e\x30\x13\x2a\xff\x73\xe9\x1e\x4a\x12\xa1\x09\xc9\x2b\xe5\xff\xce\x84\xc1\xd5\x57\x1f\xf3\x36\x09\x8d\x83\x74\xd7\xb3\x0d\xe2\x47\x56\x75\x07\x3f\xbb\xdc\x4d\x1b\x1c\x70\xd3\xbe\x3a\x07\xfa\x27\xe6\xcc\xe3\xc7\xb3\xdf\x88\x2b\xed\x89\xf1\x44\x78\x27\x25\x3c\xe9\x80\xff\x8d\x6a\xe6\x53\x7a\xe2\x04\x53\xd5\xd1\x4d\x20\xf5\x88\xf3\xfd\xac\x34\x6c\xa4\xcb\x6a\xe4\x13\x44\xc5\x65\x11\x9e\x53\x1d\x72\xf9\x5f\xb3\xb7\xe7\xd6\xc5\x25\x45\x29\x3f\xc6\x94\xfc\x5f\x9d\x59\x73\x1a\xd3\x30\xef\xca\xed\x3c\x2e\x3a\xda\xda\x9a\xa8\x46\x47\x8c\x02\x61\xdb\x16\x39\x8a\xff\x48\x09\xc5\xd2\x02\x9c\xa6\xb8\x45\x92\x7f\x49\xd4\x00\x18\xa7\x24\x3a\xac\xb4\x82\x5f\x1b\x47\x3f\x6a\x0c\x49\x31\x4b\x62\xe9\x9f\xda\xeb\xab\x44\x7e\xbc\xbd\x35\x7e\xa9\x04\x6a\xd6\x79\x48\x77\xc2\x15\xd9\x3e\x0d\xbd\x42\x7b\xdb\xa8\x99\x88\x94\x13\x6a\x3d\x83\x92\xfb\x2b\xc5\x7b\xa9\xfe\x2f\x7e\x80\xf7\x24\x22\xf2\x74\xac\xf0\x58\x63\xd3\xa6\x15\xd4\x51\x7d\x7f\xc4\x48\xf6\x28\x0d\xf5\xbe\x6e\x1a\xee\x87\x08\x7f\x4e\x84\x83\xe3\xc0\xc3\x94\xc6\xad\x57\x78\xdc\x58\x33\xf5\x7d\x87\xda\x38\x5a\x3d\x52\x33\xc5\x2a\x89\x99\x5d\xac\x3f\x38\x8a\xf5\x34\x09\x10\xc7\x10\xee\xf3\x09\xf7\x1f\x20\xdc\xcd\xc3\x5d\x7a\xb8\x6b\x93\x9f\xdc\xae\x7e\x16\x8b\x8f\x22\x67\x40\x61\x7c\x68\xe6\x01\xd6\x70\xf8\xa0\x14\x55\xf8\x90\x2b\xed\xc1\x84\x77\x79\x37\x00\x08\x75\xc9\xc5\xdd\xff\x21\x2c\xdb\x8f\x18\xc2\xb2\x31\x63\x53\x87\x7b\x65\xe3\x52\x05\x3d\x65\x9e\x8a\xb6\x71\xca\x55\x4c\x8a\x19\xd9\xda\x36\x13\x87\x9b\xb2\x56\x72\x71\x31\x38\x59\x44\x28\x2f\x85\x5b\xb2\xe3\xd8\x0f\xa8\x5a\x8c\x99\x24\xfc\xd1\x27\x44\x42\xb4\x0d\xb1\xf7\x94\x9f\xe6\xc9\x96\x06\xf7\x79\xb7\x00\x08\x75\xc9\xef\x11\x08\xc5\x98\x11\xa5\xe8\x4b\x8b\x42\x25\x4c\x38\x3e\x99\x87\x54\x11\xd5\x87\x5f\xe6\xab\xe6\xf8\x93\xff\xf4\x31\x68\xfe\x06\x40\x74\xf9\x34\x74\x29\xe0\x1d\xc5\x22\x5d\xf7\x90\x17\xe1\x67\x2f\xbb\x48\x56\x38\xd9\xa9\x76\x19\x51\x00\x28\x0d\x49\x2b\xa0\xe8\xdb\x5b\x98\x25\x30\x99\xbd\x36\xa2\x6d\x1c\xb4\x45\x76\xdd\x6a\xd9\xd5\x14\xe6\xf2\xb6\xd8\xcb\x2e\x70\xe0\xc6\x1c\xe3\x22\xa6\x8b\x1d\xb6\x21\x73\xa5\x55\x95\xdc\x8c\x85\x01\x01\x5d\xed\x47\xda\x73\x28\xff\x25\xfb\xe3\x91\x04\x5f\x47\xe5\x53\x9e\x31\x68\x8c\x99\x12\x70\xad\x21\x79\x0d\xae\x89\x4b\xdd\x3d\x28\x59\x02\xeb\x8f\xb6\xcd\xfd\x22\xf3\xc5\x41\xfc\xc9\x16\xe2\x5a\x4d\xf9\x4d\xb0\x0c\xe4\x3b\x7f\xaa\x27\x6d\x69\x9f\x15\x5a\x7b\xa0\x65\x9e\x00\x4c\x33\x48\x02\xd3\x94\xe4\xf2\x53\x56\x0b\xe2\x41\xca\x9a\x09\xb4\xdf\x32\x00\xdc\xed\x47\xea\xe0\x0e\x70\x28\xf0\x39\x09\xbb\x33\xd5\xc0\x6e\x83\x24\xb0\x5b\x49\x5e\x39\x1f\xfd\xfb\x10\xb8\x64\x9e\x0c\x70\x71\x37\x1f\xf6\xb7\x84\x6b\x8a\x0c\xf6\x69\x64\x18\xa3\x20\x47\x91\xa7\x74\xd9\xe4\x92\xb2\x11\xf0\xc8\x20\x09\x3c\x52\x92\xd3\x9a\x81\x1d\x51\xb7\x01\x64\xed\xec\x1e\xf1\x56\xc9\x9a\x11\xf6\x24\xc4\x5e\x87\xf2\xb9\xd8\x40\xf6\x75\x12\x23\xe4\x9e\x80\x3f\xe3\x5d\xca\x65\x9d\xd1\xad\x35\x54\xbf\x47\xdc\x9f\xb4\x5b\xac\x2b\x12\x27\x71\xbf\x20\x02\x43\xdc\x97\xa3\x7f\x1b\xc8\xe1\x43\x5e\xdd\x6c\xba\xec\x5b\x9f\xa1\xf7\x03\x56\x84\xe5\xce\xd7\x3d\x39\x78\xb9\xc6\x1e\xeb\xc3\x0f\xc6\x16\x70\x17\x84\xbb\xe0\xcc\x77\x15\xfc\xf6\xaf\xa9\x57\x86\x1f\x73\x25\x7d\x21\x01\x3c\x6c\x7e\xd2\xf1\xf0\x91\x44\x8c\xa3\x68\xa7\xf1\x8a\x7e\x1b\x8d\x50\x18\x56\x59\x98\x66\xb8\x63\xfd\x36\x16\x6d\xce\xdd\x01\x22\xd6\x25\x67\x40\xc4\x21\xac\xcb\x0c\xf1\x47\x8a\x69\xeb\xa2\x6c\xc9\x0e\x7b\x14\xb2\xef\x9c\x89\xdf\x62\x63\xd5\x99\x05\x8e\x77\x58\x75\x01\x6a\xf4\xb6\xa3\x9c\x3b\xe7\xee\xf7\xdb\x79\xb4\x69\x69\x06\xc4\xb9\x02\x71\xcc\xf6\xaf\x0d\xad\xab\xae\x73\xdd\xc2\x8c\x45\xb8\xcc\xac\x42\x53\x78\xba\xcd\xbe\xac\x6f\x6b\x21\x48\xdf\x9a\x9f\xda\xd3\xb7\xcb\xc4\x36\x3f\x32\x7a\xb3\x13\x73\x99\xcf\x01\x5c\x1b\x92\x57\x85\xeb\x94\x33\xdd\x82\x2a\x30\xe5\x05\x6a\xf6\x1f\xf4\x1c\x77\x45\x35\x52\xca\x1e\xbf\x5a\x05\xee\x01\xf7\xca\x8d\x67\x66\x98\xec\x66\xee\x7c\x0f\xd5\x5d\x50\x14\x81\x2d\xce\x00\x79\x7a\x4d\xa0\xf9\x63\xad\xec\xb7\x5a\xf0\x0b\xea\x96\xa3\xce\xb6\x53\xdd\x2b\xd1\x41\xeb\x08\xe5\x9d\x55\x80\x7c\xa3\x24\x20\xdf\x20\x09\x7b\xae\x5e\x11\x81\xc6\xcc\xd4\xfd\x44\x5e\xff\xc1\xd3\xf5\xf7\xa2\x75\xf6\xb4\x95\x1e\x3f\xbf\x94\xa7\x04\x74\x01\xba\xca\x8d\x67\x8b\xae\xde\xb3\x74\xc0\x55\xfb\x11\x03\xae\xf2\x27\x7f\x8e\x2d\x07\x4b\x14\xe1\xcf\x84\x71\xe1\x1b\xde\x53\xae\xd4\xfa\x21\x13\x0f\x45\x27\x80\x47\x75\xc9\xf9\xf1\xa8\xb3\xb2\x5b\xeb\xb2\x9b\x1a\xb0\x9e\x72\xa4\x44\xb9\x79\x11\x49\x37\xa7\x5b\x5c\x9d\x38\x8f\x63\x78\x02\xc7\xac\x9f\xc0\x71\xc6\xad\x77\xe9\xb0\x75\x4d\xbc\x78\x0c\x33\x90\xb7\x2e\x39\xf1\x7e\x74\xa3\xdd\x6b\x43\x1a\xb0\x30\x57\x3c\xcd\xbb\xd5\x1d\xce\x0a\xae\xb5\x3a\xd7\x09\x93\x05\xac\xca\xcd\xce\x38\xc0\xd2\xe6\x27\x6d\x7e\xeb\xbf\xe4\x7f\x39\xa8\x93\xcb\x64\xf7\x02\xdd\x67\xc2\x8f\xea\x10\x09\xbc\x9b\x9b\xfc\xf0\xe6\xfe\xe6\xa6\x57\xee\x0b\x00\x6e\x48\x5e\x05\xc0\x53\x4c\xc1\x37\xf7\x5e\xbc\x57\x2e\x61\x41\xf9\x57\x32\x03\x07\x5a\x2e\x8b\x96\xa3\x0a\xe6\x2e\x11\x98\x55\xd1\x81\x82\x26\x49\xa0\x60\x55\x72\x99\x69\x79\xbe\xd1\xc4\xa9\x25\x20\x1d\x07\xca\xb7\x1f\x19\x5d\x24\x67\x1e\xea\xb9\x9a\x9b\x9d\x02\x70\x6e\x92\x04\x9c\x57\x25\x67\x5d\x11\x07\xe6\xd8\x1c\xb5\x9a\x87\x27\x21\x8a\xc6\x16\x99\x9e\x94\x92\xa7\x5e\xf3\xea\xf7\xea\xbc\xc0\xa1\xba\xe4\xc4\x75\x25\x82\xc2\x05\x15\x93\x44\x6f\xa1\x82\xa4\x13\x1e\x5d\x41\x92\x11\x0b\xe5\xa3\xd9\x96\x8f\xd4\xe5\xe9\x59\x2a\x7a\x6f\x68\x03\x40\x9d\x10\xa8\x86\x0b\x55\x1b\xcf\x80\x5a\x51\xe6\x03\xde\x3e\xa6\xe2\xcf\xae\xb5\x9b\xb3\xa2\x6b\x4d\x52\xdb\xf1\x01\x33\xd4\xe5\x23\xac\x25\x5b\xf4\x5f\xe4\x7f\xae\xea\x35\x52\x57\x69\xa5\xf2\xe6\x46\x1e\x90\x7e\xae\x7c\x5f\x1d\xba\x38\x7f\xef\x14\x13\x80\xd8\x90\x9c\x1e\x88\xd3\x4e\x73\xbb\x90\x3b\xdb\x39\xee\x2c\xd2\x4c\x87\x3c\x07\x5a\x9e\x3f\x39\xa8\xda\xa8\xdb\x7d\xb5\x60\xd3\x20\x61\xb1\x81\x64\x6c\x39\x07\xc8\xa8\x95\x04\x32\x16\x92\x0b\x4c\x9a\xf3\x4a\x8e\x3b\x1b\x40\x86\xac\xe9\x21\x30\xff\xfc\x69\x64\x0d\xc7\x3e\xfd\xed\x49\xfa\x4a\xa5\x07\x48\xaf\x95\x04\xd2\x17\x92\x50\xe7\x79\x6d\x5c\xd2\xcf\xdc\xe5\xca\xff\xe8\x3a\x8f\x52\x62\x59\xe7\xb9\x0b\xc3\x7b\x75\x4e\xa0\x4f\x5d\x12\x7e\x3b\x94\x89\xc2\x6f\x87\x94\xdc\x95\x9f\x31\x29\x82\x12\x2a\x3f\xb3\xad\xfc\xc8\xcb\x23\x53\x3f\xeb\xea\x8f\x84\x2c\x30\xb6\x21\xb9\xc0\x19\xec\xaf\xf2\xad\xcd\xfa\xab\x79\x6e\x02\xd3\x57\x0f\xd2\xc4\xda\xa0\x9c\xa7\x89\xfe\x4b\xbe\x4f\xc4\x49\x81\x47\xea\xaa\xcc\x70\x73\xe5\xb6\x3f\xc4\x01\xc0\x69\x25\x67\xb2\x51\x68\xf0\x14\xd6\x1d\xe9\xa0\x8c\x63\x97\xdb\x41\xf6\xd6\x37\x7b\x1b\x50\xc6\x51\xf9\x5b\xb5\x8c\xd3\x87\x77\xf9\xb3\x2b\x01\x79\x3a\x49\x40\x5e\x21\x39\xbd\x25\x26\xaa\xcf\x04\xb3\xc9\x70\x07\xb1\x12\x32\x5c\x43\xd3\x05\xd0\x7c\x60\x81\x66\x54\xfa\x9a\x3f\x97\x12\x70\xae\x93\x04\x9c\x17\x92\x50\x84\x79\x6d\xec\xb1\x9a\x5d\xbb\x79\xf1\xa1\x4c\x38\x7b\xbd\xef\x50\x34\x80\xb7\x1d\x1a\x24\x81\x49\x85\x24\xcc\xaa\xcb\x99\x22\xbc\xff\xb0\x73\xa8\xf6\x48\x54\x7f\xe5\xa6\x59\x95\x4e\x5c\xe9\xef\x6a\x4b\x45\xf8\xd3\x26\x1c\x0b\xc7\x8c\xb7\xbf\x0b\x6b\xd4\xac\x50\x76\xf0\x26\xbd\x56\x21\x61\x5c\x8b\x14\x59\x4c\x09\x70\xc4\x55\xc5\xba\xf2\xed\xc7\xfa\xa2\x00\x95\x48\xe5\xc4\xe0\xce\xc5\x29\x8c\xae\xde\x15\x57\x1a\x97\xa9\x74\xce\xa8\xb8\x75\xe4\x67\xa9\x94\x61\xaa\x70\xd4\xe2\x3d\x09\x62\xec\x39\xa6\x86\x47\x4d\x6b\x42\xbc\x32\x30\xed\x75\xb9\xe8\x6e\x37\x5f\xb3\x9f\xd6\xe5\xd3\x36\x50\xc9\x7f\x86\x58\xcc\x7a\x55\x8c\xd8\xdd\xd9\x1c\xe6\x0d\x24\xda\xc7\x6e\xe3\x40\x3e\x09\x5b\x62\x65\x94\xa7\x17\x4a\x86\xbb\xba\xc5\xd8\x8b\x1f\xb5\x3b\x1d\xbe\xd5\x16\x89\xec\xae\x3d\xca\x42\xc4\xec\x50\x83\x38\xa0\xe9\xb7\x53\xfd\x5b\x12\x05\xea\x35\xf1\x9d\x5a\xb7\x71\x1c\x62\x5d\x89\x50\xa7\xb6\x35\x8e\x87\x19\xa2\x72\x6b\x74\xa8\x98\xa3\x43\x0b\x62\x6d\x76\x6a\x58\xec\xd0\x18\x98\x5b\xe9\xfa\x2b\x52\x75\xa4\xde\x8a\xdf\xd9\xe7\x3c\x4e\x06\x05\xa1\xaa\x06\xbb\x8d\x40\x63\x26\xec\x20\xec\xdc\xfb\x9b\xf3\x40\xce\x87\xbf\x8c\xa0\xd8\x53\x3c\x0a\x09\x16\x1e\xa6\xe6\x2d\x53\x33\xfe\xbb\x62\x79\xe5\x6d\x91\x23\xf0\xa0\x5b\x23\xd0\xe1\xc0\xbe\xb3\x36\xee\x50\x7b\x63\xd4\x24\x9e\xb1\xb0\x9b\xbe\x63\xb5\xee\x21\x39\xa1\xc7\xe9\x77\x2f\x37\x54\xce\xe7\xbe\x6c\x4a\x5b\xad\x4f\x65\x1f\x25\x6e\x83\x63\x7e\xfe\x3f\x41\xf6\x98\x45\x54\xde\x74\x3a\xa6\xcf\xd7\x63\x2d\xef\xc8\x8f\xa6\x05\x6a\x53\x0c\x6b\x18\x3e\x46\x45\xb6\x10\xd5\xd3\xbf\xbb\xfc\xf3\x84\x19\x63\xe8\xe0\x1c\x7e\x81\x85\x46\x12\x71\x7c\xc0\x9a\xb5\xb5\x8e\xb5\xb2\x37\x5f\xff\x0c\x00\x00\xff\xff\x79\xf5\x57\xa1\x52\xba\x00\x00")

func swaggerSpecApiJsonBytes() ([]byte, error) {
	return bindataRead(
		_swaggerSpecApiJson,
		"swagger-spec/api.json",
	)
}

func swaggerSpecApiJson() (*asset, error) {
	bytes, err := swaggerSpecApiJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "swagger-spec/api.json", size: 47698, mode: os.FileMode(436), modTime: time.Unix(1445596741, 0)}
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
	"swagger-spec/api.json": swaggerSpecApiJson,
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
	"swagger-spec": &bintree{nil, map[string]*bintree{
		"api.json": &bintree{swaggerSpecApiJson, map[string]*bintree{
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

package commands

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/hpcloud/cf-usb/driver-generator/data"
)

func GenerateCommand(c *cli.Context) {
	validateArgsCount(c, 1)

	di := driverInfo{}

	di.DriverName = c.Args().First()

	di.GeneratedPath = c.String("path")

	baseImp, err := baseImport(di.DriverName)
	if err != nil {

	}

	di.BaseImport = baseImp

	if _, err := os.Stat(di.GeneratedPath); err == nil {
		log.Println(di.GeneratedPath, "exists. Please chose an empty folder")
		return
	}

	fmt.Println("Generateing:", di.DriverName, "to", di.GeneratedPath)

	//	//Create folder structure
	//	os.MkdirAll(path.Join(di.GeneratedPath, "cmd", "driver", di.DriverName), 0777)
	//	os.MkdirAll(path.Join(di.GeneratedPath, "driver", di.DriverName, "driverdata"), 0777)
	//	os.MkdirAll(path.Join(di.GeneratedPath, "driver", di.DriverName, "schemas"), 0777)
	err = writeTemplate("config", filepath.Join("driver", di.DriverName, "config"), di)
	if err != nil {
		log.Fatalln("cannot generate config.go:", err)
	}

	err = writeTemplate("main", filepath.Join("cmd", "driver", di.DriverName), di)
	if err != nil {
		log.Fatalln("cannot generate main.go:", err)
	}

	err = writeTemplate("driver", filepath.Join("driver", di.DriverName), di)
	if err != nil {
		log.Fatalln("cannot generate driver.go:", err)
	}

	os.MkdirAll(path.Join(di.GeneratedPath, "driver", di.DriverName, "schemas"), 0777)
}

func baseImport(tgt string) (string, error) {
	p, err := filepath.Abs(tgt)
	if err != nil {
		log.Fatalln(err)
	}

	var pth string
	for _, gp := range filepath.SplitList(os.Getenv("GOPATH")) {
		pp := filepath.Join(gp, "src")
		if strings.HasPrefix(p, pp) {
			pth, err = filepath.Rel(pp, p)
			if err != nil {
				return "", err
			}
			break
		}
	}

	if pth == "" {
		err := errors.New("target must reside inside a location in the $GOPATH/src")
		return "", err
	}
	return pth, nil
}

func writeTemplate(templateName, relativePath string, di driverInfo) error {
	asset, err := data.Asset(fmt.Sprintf("templates/%s.template", templateName))
	if err != nil {
		return err
	}

	templ, err := template.New(templateName).Parse(string(asset))
	if err != nil {
		return err
	}

	parentFolder := filepath.Join(di.GeneratedPath, relativePath)

	os.MkdirAll(parentFolder, 0777)

	file, err := os.Create(filepath.Join(parentFolder, fmt.Sprintf("%s.go", templateName)))
	if err != nil {
		return err
	}

	defer file.Close()
	err = templ.Execute(io.Writer(file), di)

	return err
}

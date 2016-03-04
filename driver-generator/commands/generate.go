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

	di.GeneratedPath = strings.TrimSpace(c.String("path"))

	if di.GeneratedPath == "" {
		log.Fatalln("Incorrect Usage, Requires --path argument")
	}

	baseImp, err := baseImport(di.GeneratedPath)
	if err != nil {

	}

	di.BaseImport = baseImp

	if _, err := os.Stat(di.GeneratedPath); err == nil {
		log.Println(di.GeneratedPath, "exists. Please chose an empty folder")
		return
	}

	log.Println("Generateing:", di.DriverName, "to", di.GeneratedPath)

	err = writeTemplate("config.go", filepath.Join("driver", di.DriverName, "config"), di)
	if err != nil {
		log.Fatalln("cannot generate config.go:", err)
	}

	err = writeTemplate("main.go", filepath.Join("cmd", "driver", di.DriverName), di)
	if err != nil {
		log.Fatalln("cannot generate main.go:", err)
	}

	err = writeTemplate("driver.go", filepath.Join("driver", di.DriverName), di)
	if err != nil {
		log.Fatalln("cannot generate driver.go:", err)
	}

	err = writeTemplate("Makefile", "", di)
	if err != nil {
		log.Fatalln("cannot generate makefile:", err)
	}

	schemasPath := path.Join(di.GeneratedPath, "driver", di.DriverName, "schemas")
	os.MkdirAll(schemasPath, 0777)

	configFile, err := os.Create(path.Join(schemasPath, "config.json"))
	if err != nil {
		log.Fatalln("cannot write config json file", err)
	}
	defer configFile.Close()

	configFile.WriteString("{}")

	dialsFile, err := os.Create(path.Join(schemasPath, "dials.json"))
	defer dialsFile.Close()

	dialsFile.WriteString("{}")

	log.Println("Done!")
	log.Println("You can build the driver by running make in the generated directory.")

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

	file, err := os.Create(filepath.Join(parentFolder, templateName))
	if err != nil {
		return err
	}

	defer file.Close()
	err = templ.Execute(io.Writer(file), di)

	return err
}

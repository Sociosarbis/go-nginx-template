package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	template "github.com/sociosarbis/go/template/src"
)

var (
	templatePath string
	dest         string
	upstream     string
)

func initFlags() {
	flag.StringVar(&templatePath, "template-path", "", "source path of template")
	flag.StringVar(&dest, "dest", "", "output path of generated file.default is the template path")
	flag.StringVar(&upstream, "upstream", "", "host of upstream api")
	flag.Parse()
}

func main() {
	initFlags()
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("can not get working directory")
	}
	if templatePath == "" {
		log.Fatalf("template-path is not defined")
	}

	if !filepath.IsAbs(templatePath) {
		templatePath = filepath.Join(pwd, templatePath)
	}

	if dest == "" {
		dest = templatePath
	}

	if !filepath.IsAbs(dest) {
		dest = filepath.Join(pwd, dest)
	}

	if upstream == "" {
		log.Fatalf("upstream is not defined")
	}

	if success := template.GenerateFile(&template.Config{
		TemplatePath: templatePath,
		Dest:         dest,
	}, &template.Context{
		Upstream: upstream,
	}); !success {
		fmt.Printf("file generation failed")
	}
}

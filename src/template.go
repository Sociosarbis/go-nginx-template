package template

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type Context struct {
	Upstream string
}

type Config struct {
	TemplatePath string
	Dest         string
}

func GenerateFile(conf *Config, data *Context) bool {
	tpl, err := template.New(filepath.Base(conf.TemplatePath)).ParseFiles(conf.TemplatePath)
	if err != nil {
		log.Fatalf("Unable to parse template: %s", err)
	}
	buf := new(bytes.Buffer)
	err = tpl.ExecuteTemplate(buf, filepath.Base(conf.TemplatePath), data)
	if err != nil {
		log.Fatalf("Template error: %s\n", err)
	}
	contents := buf.Bytes()
	destDir := filepath.Dir(conf.Dest)
	if err = os.MkdirAll(destDir, os.ModeDir); err != nil {
		log.Fatalf("can not create output directory")
	}
	dest, err := ioutil.TempFile(destDir, "docker-gen")
	defer func() {
		dest.Close()
		os.Remove(dest.Name())
	}()
	if err != nil {
		log.Panicf("Unable to create temp file: %s\n", err)
	}
	if n, err := dest.Write(contents); n != len(contents) || err != nil {
		log.Panicf("Failed to write to temp file: wrote %d, exp %d, err=%v", n, len(contents), err)
	}
	dest.Close()
	if fi, err := os.Stat(conf.Dest); err == nil || os.IsNotExist(err) {
		if err != nil && os.IsNotExist(err) {
			emptyFile, err := os.Create(conf.Dest)
			if err != nil {
				log.Fatalf("Unable to create empty destination file: %s\n", err)
			} else {
				emptyFile.Close()
				fi, _ = os.Stat(conf.Dest)
			}
		}
		if err := Chmod(dest, fi); err != nil {
			log.Panicf("Unable to chmod temp file: %s\n", err)
		}
		if err := Chown(dest, fi); err != nil {
			log.Panicf("Unable to chown temp file: %s\n", err)
		}
	}
	err = os.Rename(dest.Name(), conf.Dest)
	if err != nil {
		log.Panicf("Unable to create dest file %s: %s\n", conf.Dest, err)
	}
	log.Printf("Generated '%s'", conf.Dest)
	return true
}

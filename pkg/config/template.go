package config

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var (
	templates              = make(map[string]*template.Template)
	globalStorage *Storage = nil
)

func LoadTemplates(path string) error {
	files, err := filepath.Glob(filepath.Join(path, "*.conf.tmpl"))
	if err != nil {
		return err
	}

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		fullName := filepath.Base(file)
		tmpl, err := template.New(fullName).Parse(string(content))
		if err != nil {
			return err
		}
		trimmedName := strings.TrimSuffix(fullName, ".conf.tmpl")
		templates[trimmedName] = tmpl
	}

	// Update metrics
	activeTemplates.Set(float64(len(templates)))

	return nil
}

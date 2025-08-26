package internal

import (
	"fmt"
	"log/slog"
	"os"
	"text/template"
)

type TemplateInfo struct {
	TmplString string
	FilePath   string
}

type Generator struct {
	name            string
	flags           ParsedFlags
	templateInfoMap map[string]TemplateInfo
	dirs            []string
}

func NewGenerator(name string, flags ParsedFlags, templateInfoMap map[string]TemplateInfo, dirs ...string) Generator {
	return Generator{
		name:            name,
		flags:           flags,
		templateInfoMap: templateInfoMap,
		dirs:            dirs,
	}
}

func (g Generator) Generate() error {
	if g.flags.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	slog.Info("Generating " + g.name + "...")

	pathToTemplate := map[string]*template.Template{}

	for templateName, tmplInfo := range g.templateInfoMap {

		tmpl, err := template.New(templateName).Parse(tmplInfo.TmplString)
		if err != nil {
			return fmt.Errorf("error parsing %s template: %w", templateName, err)
		}
		pathToTemplate[tmplInfo.FilePath] = tmpl
		slog.Debug("Template " + templateName + " ok.")
	}

	// create directory structure
	err := createDirsAll(g.dirs...)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	for path, tmpl := range pathToTemplate {
		err = createFile(tmpl, path, g.flags)
		if err != nil {
			return fmt.Errorf("error creating %s file: %w", path, err)
		}
	}
	return nil
}

func createFile(tmpl *template.Template, filePath string, data ParsedFlags) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}
	slog.Debug("File created.", "file", filePath)
	return nil
}

func createDirsAll(dirs ...string) error {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating %s directory: %w", dir, err)
		}
		slog.Debug("Directory created.", "dir", dir)
	}
	return nil
}

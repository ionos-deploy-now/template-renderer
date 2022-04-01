package cmd

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/bmatcuk/doublestar/v4"
	"io/fs"
	"os"
	"strings"
	"syscall"
	"text/template"
)

type ConfigurationTemplate struct {
	Path     string
	Filename string
	Owner    int
	Group    int
	Mode     fs.FileMode
	Template *template.Template
}

var templateFunctions = sprig.TxtFuncMap()

func init() {
	delete(templateFunctions, "env")
	delete(templateFunctions, "expandenv")
}

func LoadTemplateFiles(templateDir string, templateExtension string) ([]ConfigurationTemplate, error) {
	var templates []ConfigurationTemplate
	files, err := doublestar.Glob(os.DirFS(templateDir), "**/*"+templateExtension)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		subPaths := strings.Split(file, "/")
		templateName := subPaths[len(subPaths)-1]
		filename := strings.TrimSuffix(templateName, templateExtension)
		if filename == "" {
			break
		}

		var fileInfo syscall.Stat_t
		if err := syscall.Stat(joinPath(templateDir, file), &fileInfo); err != nil {
			return nil, err
		}
		templates = append(templates, ConfigurationTemplate{
			Path:     joinPath(subPaths[:len(subPaths)-1]...),
			Filename: filename,
			Owner:    int(fileInfo.Uid),
			Group:    int(fileInfo.Gid),
			Mode:     fs.FileMode(fileInfo.Mode),
			Template: template.Must(template.New(templateName).
				Funcs(templateFunctions).
				Option("missingkey=error").
				ParseFiles(joinPath(templateDir, file))),
		})
	}
	return templates, nil
}

func (t ConfigurationTemplate) Render(data Data, outputDir string, copyPermissions bool) error {
	var buffer bytes.Buffer
	if err := t.Template.Execute(&buffer, data); err != nil {
		return err
	}
	if err := os.MkdirAll(joinPath(outputDir, t.Path), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(joinPath(outputDir, t.Path, t.Filename))
	if err != nil {
		return err
	}
	if _, err = file.Write(buffer.Bytes()); err != nil {
		return err
	}

	if copyPermissions {
		if err := file.Chown(t.Owner, t.Group); err != nil {
			return err
		}
		if err := file.Chmod(t.Mode); err != nil {
			return err
		}
	}
	return nil
}

func joinPath(s ...string) string {
	var values []string
	for _, value := range s {
		if value != "" {
			values = append(values, value)
		}
	}
	return strings.Join(values, "/")
}

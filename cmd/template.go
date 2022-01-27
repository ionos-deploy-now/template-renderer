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

func LoadTemplateFiles(templateDir string, templateExtension string) []ConfigurationTemplate {
	var templates []ConfigurationTemplate
	files, err := doublestar.Glob(os.DirFS(templateDir), "**/*"+templateExtension)
	handleError(err)
	for _, file := range files {
		subPaths := strings.Split(file, "/")
		templateName := subPaths[len(subPaths)-1]
		var fileInfo syscall.Stat_t
		handleError(syscall.Stat(joinPath(templateDir, file), &fileInfo))
		templates = append(templates, ConfigurationTemplate{
			Path:     joinPath(subPaths[:len(subPaths)-1]...),
			Filename: strings.TrimSuffix(templateName, templateExtension),
			Owner:    int(fileInfo.Uid),
			Group:    int(fileInfo.Gid),
			Mode:     fs.FileMode(fileInfo.Mode),
			Template: template.Must(template.New(templateName).
				Funcs(templateFunctions).
				Option("missingkey=error").
				ParseFiles(joinPath(templateDir, file))),
		})
	}
	return templates
}

func (t ConfigurationTemplate) Render(data Data, outputDir string, copyPermissions bool) {
	var buffer bytes.Buffer
	handleError(t.Template.Execute(&buffer, data))

	handleError(os.MkdirAll(joinPath(outputDir, t.Path), os.ModePerm))
	file, err := os.Create(joinPath(outputDir, t.Path, t.Filename))
	handleError(err)
	_, err = file.Write(buffer.Bytes())
	handleError(err)

	if copyPermissions {
		handleError(file.Chown(t.Owner, t.Group))
		handleError(file.Chmod(t.Mode))
	}
}

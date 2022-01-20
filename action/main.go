package action

import (
	"bufio"
	"github.com/Masterminds/sprig"
	"github.com/bmatcuk/doublestar/v4"
	"io/fs"
	"log"
	"os"
	"strings"
	"syscall"
	"text/template"
)

var (
	templateDir       string
	templateExtension string
	outputDir         string
	envVarPrefix      string
	copyPermissions   bool
)

var templateFunctions = sprig.TxtFuncMap()

func init() {
	delete(templateFunctions, "env")
	delete(templateFunctions, "expandenv")
}

type ConfigurationTemplate struct {
	Path     string
	Filename string
	Owner    int
	Group    int
	Mode     fs.FileMode
	Template *template.Template
}

func loadTemplateFiles() []ConfigurationTemplate {
	var templates []ConfigurationTemplate
	files, err := doublestar.Glob(os.DirFS(templateDir), "**/*.template")
	handleError(err)
	for _, file := range files {
		subPaths := strings.Split(file, "/")
		filename := subPaths[len(subPaths)-1]
		var fileInfo syscall.Stat_t
		handleError(syscall.Stat(joinPath(templateDir, file), &fileInfo))
		templates = append(templates, ConfigurationTemplate{
			Path:     joinPath(subPaths[:len(subPaths)-1]...),
			Filename: filename,
			Owner:    int(fileInfo.Uid),
			Group:    int(fileInfo.Gid),
			Mode:     fs.FileMode(fileInfo.Mode),
			Template: template.Must(template.New(filename).
				Funcs(templateFunctions).
				Option("missingkey=error").
				ParseFiles(joinPath(templateDir, file))),
		})
	}
	return templates
}

func (t ConfigurationTemplate) Fill(data map[string]interface{}) {
	handleError(os.MkdirAll(joinPath(outputDir, t.Path), os.ModePerm))
	file, err := os.Create(strings.TrimSuffix(joinPath(outputDir, t.Path, t.Filename), ".template"))
	handleError(err)
	writer := bufio.NewWriter(file)
	handleError(t.Template.Execute(writer, data))
	handleError(writer.Flush())
	if copyPermissions {
		handleError(file.Chown(t.Owner, t.Group))
		handleError(file.Chmod(t.Mode))
	}
}

func getDataFromEnvironment() map[string]interface{} {
	data := make(map[string]string)
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, envVarPrefix) {
			name := strings.TrimPrefix(strings.Split(env, "=")[0], envVarPrefix)
			value := strings.Split(env, "=")[1]
			data[name] = value
		}
	}
	return map[string]interface{}{"data": data}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func joinPath(s ...string) string {
	return strings.Join(s, "/")
}

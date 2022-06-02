package cmd

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"io/fs"
	"os"
	"strings"
	"syscall"
	"template-renderer/cmd/types"
	"text/template"
)

type ConfigurationTemplate struct {
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

func ReadTemplateFile(templateFilePath types.Path, templateFileName string, templateExtension string) (*ConfigurationTemplate, error) {
	fullPath := templateFilePath.Append(templateFileName).String()
	var fileInfo syscall.Stat_t
	if err := syscall.Stat(fullPath, &fileInfo); err != nil {
		return nil, err
	}
	return &ConfigurationTemplate{
		Filename: strings.TrimSuffix(templateFileName, templateExtension),
		Owner:    int(fileInfo.Uid),
		Group:    int(fileInfo.Gid),
		Mode:     fs.FileMode(fileInfo.Mode),
		Template: template.Must(template.New(templateFileName).
			Funcs(templateFunctions).
			Option("missingkey=error").
			ParseFiles(fullPath)),
	}, nil
}

func (t ConfigurationTemplate) Render(data *Data, outputDir types.Path, copyPermissions bool, runtimeVariableFiles *[]string) error {
	var buffer bytes.Buffer
	if err := t.Template.Execute(&buffer, data.Values); err != nil {
		return err
	}
	if data.RuntimeValuesUsed {
		*runtimeVariableFiles = append(*runtimeVariableFiles, outputDir.Append(t.Filename).String())
		data.ResetUsedRuntimeValues()
	}

	file, err := os.Create(outputDir.Append(t.Filename).String())
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

package cmd

import (
	"io/fs"
	"os"
	"strings"
	"syscall"
	"template-renderer/cmd/types"
)

type Directory struct {
	Name           string
	Owner          int
	Group          int
	Mode           fs.FileMode
	SubDirectories []*Directory
	Templates      []*ConfigurationTemplate
}

func ReadDirectory(baseDir types.Path, name types.Path, templateExtension string) (*Directory, error) {
	currentDir := baseDir.Append(name)
	entries, err := os.ReadDir(currentDir.String())
	if err != nil {
		return nil, err
	}
	var subDirectories []*Directory
	var templates []*ConfigurationTemplate
	for _, entry := range entries {
		if entry.IsDir() {
			subDirectory, err := ReadDirectory(currentDir, types.Path(entry.Name()), templateExtension)
			if err != nil {
				return nil, err
			}
			subDirectories = append(subDirectories, subDirectory)
		} else if strings.HasSuffix(entry.Name(), templateExtension) && entry.Name() != templateExtension {
			template, err := ReadTemplateFile(currentDir, entry.Name(), templateExtension)
			if err != nil {
				return nil, err
			}
			templates = append(templates, template)
		}
	}
	var fileInfo syscall.Stat_t
	if err := syscall.Stat(currentDir.String(), &fileInfo); err != nil {
		return nil, err
	}
	return &Directory{
		Name:           name.String(),
		Owner:          int(fileInfo.Uid),
		Group:          int(fileInfo.Gid),
		Mode:           fs.FileMode(fileInfo.Mode),
		SubDirectories: subDirectories,
		Templates:      templates,
	}, nil
}

func (d Directory) Render(data *Data, outputDir types.Path, copyPermissions bool, runtimeVariableFiles *[]string) error {
	currentDir := outputDir.Append(d.Name)
	_, err := os.Stat(currentDir.String())
	if os.IsNotExist(err) {
		if copyPermissions {
			if err := os.MkdirAll(currentDir.String(), d.Mode); err != nil {
				return err
			}
			if err := os.Chown(currentDir.String(), d.Owner, d.Group); err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(currentDir.String(), os.ModePerm); err != nil {
				return err
			}
		}
	}
	for _, directory := range d.SubDirectories {
		if err := directory.Render(data, currentDir, copyPermissions, runtimeVariableFiles); err != nil {
			return err
		}
	}
	for _, template := range d.Templates {
		if err := template.Render(data, currentDir, copyPermissions, runtimeVariableFiles); err != nil {
			return err
		}
	}
	return nil
}

package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "templater",
		Short:        "",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			templateDir, templateExtension, inputData, outputDir, copyPermissions, err := readFlags(cmd)

			templates, err := LoadTemplateFiles(templateDir, templateExtension)
			if err != nil {
				return err
			}

			data, err := ParseInputData(inputData)
			if err != nil {
				return err
			}

			for _, template := range templates {
				err = template.Render(data, outputDir, copyPermissions)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	rootCmd.Flags().StringP("template-dir", "i", "./", "Set the input directory.")
	rootCmd.Flags().StringP("template-extension", "t", ".template", "Set a file extension to detect templates.")
	rootCmd.Flags().StringP("output-dir", "o", "./", "Set the output directory.")
	rootCmd.Flags().StringArrayP("data", "d", []string{}, "Data to use for rendering templates as yaml or json objects. Multiple objects will be merged before rendering.")
	rootCmd.Flags().Bool("copy-permissions", false, "Copy the user, group and mode of the template.")
	return rootCmd
}
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func readFlags(cmd *cobra.Command) (templateDir string, templateExtension string, inputData []string, outputDir string, copyPermissions bool, err error) {
	if templateDir, err = cmd.Flags().GetString("template-dir"); err != nil {
		return
	}
	if templateExtension, err = cmd.Flags().GetString("template-extension"); err != nil {
		return
	}
	if inputData, err = cmd.Flags().GetStringArray("data"); err != nil {
		return
	}
	if outputDir, err = cmd.Flags().GetString("output-dir"); err != nil {
		return
	}
	if copyPermissions, err = cmd.Flags().GetBool("copy-permissions"); err != nil {
		return
	}
	return
}

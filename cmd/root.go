package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

type rootConfig struct {
	templateDir       string
	templateExtension string
	inputData         []string
	runtimeData       []string
	outputDir         string
	copyPermissions   bool
	githubAction      bool
}

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "templater",
		Short:        "",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := readFlags(cmd)
			if err != nil {
				return err
			}

			templates, err := LoadTemplateFiles(config.templateDir, config.templateExtension)
			if err != nil {
				return err
			}

			var usedValues []string
			data, err := ParseInputData(config.inputData, config.runtimeData, &usedValues)
			if err != nil {
				return err
			}

			for _, template := range templates {
				err = template.Render(data, config.outputDir, config.copyPermissions)
				if err != nil {
					return err
				}
			}

			var message string
			if config.githubAction {
				message = fmt.Sprintf("::set-output name=used_runtime_values::%s", strings.Join(usedValues, ","))

			} else if len(usedValues) > 0 {
				message = "Runtime values used while rendering templates:\n"
				message += strings.Join(usedValues, "\n")
			}
			if _, err = cmd.OutOrStdout().Write([]byte(message)); err != nil {
				return err
			}

			return nil
		},
	}
	rootCmd.Flags().StringP("template-dir", "i", "./", "Set the input directory.")
	rootCmd.Flags().StringP("template-extension", "t", ".template", "Set a file extension to detect templates.")
	rootCmd.Flags().StringP("output-dir", "o", "./", "Set the output directory.")
	rootCmd.Flags().StringArrayP("data", "d", []string{}, "Data to use for rendering templates as yaml or json objects. Multiple objects will be merged before rendering.")
	rootCmd.Flags().StringArrayP("runtime-data", "r", []string{}, "Same as --data but used values will be printed out after rendering to use them as placeholders for another template engine.")
	rootCmd.Flags().Bool("copy-permissions", false, "Copy the user, group and mode of the template.")
	rootCmd.Flags().Bool("github-action", false, "Use github action output format.")
	return rootCmd
}
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func readFlags(cmd *cobra.Command) (config rootConfig, err error) {
	if config.templateDir, err = cmd.Flags().GetString("template-dir"); err != nil {
		return
	}
	if config.templateExtension, err = cmd.Flags().GetString("template-extension"); err != nil {
		return
	}
	if config.inputData, err = cmd.Flags().GetStringArray("data"); err != nil {
		return
	}
	if config.runtimeData, err = cmd.Flags().GetStringArray("runtime-data"); err != nil {
		return
	}
	if config.outputDir, err = cmd.Flags().GetString("output-dir"); err != nil {
		return
	}
	if config.copyPermissions, err = cmd.Flags().GetBool("copy-permissions"); err != nil {
		return
	}
	if config.githubAction, err = cmd.Flags().GetBool("github-action"); err != nil {
		return
	}
	return
}

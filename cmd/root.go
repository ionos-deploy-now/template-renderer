package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"strings"
	"template-renderer/cmd/types"
)

type rootConfig struct {
	templateDir                   string
	templateExtension             string
	secrets                       string
	runtimeData                   string
	outputDir                     string
	copyPermissions               bool
	outputRuntimePlaceholderFiles bool
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

			directory, err := ReadDirectory(types.Path(config.templateDir), "", config.templateExtension)
			if err != nil {
				return err
			}

			data, err := ParseInputData(config.secrets, config.runtimeData, args)
			if err != nil {
				return err
			}

			var runtimeVariableFiles []string
			if err := directory.Render(data, types.Path(config.outputDir), config.copyPermissions, &runtimeVariableFiles); err != nil {
				return err
			}

			if config.outputRuntimePlaceholderFiles {
				filePath, present := os.LookupEnv("GITHUB_OUTPUT")
				if present {
					f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
					if err != nil {
						return err
					}

					defer f.Close()

					if _, err = f.WriteString("runtime-placeholder-files=" + strings.Join(runtimeVariableFiles, ",") + "\n"); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
	rootCmd.Flags().StringP("template-dir", "i", "./", "Set the input directory.")
	rootCmd.Flags().StringP("template-extension", "t", ".template", "Set a file extension to detect templates.")
	rootCmd.Flags().StringP("output-dir", "o", "./", "Set the output directory.")
	rootCmd.Flags().StringP("secrets", "s", "", "Data to use for rendering templates. Will be prefixed with \"secrets\".")
	rootCmd.Flags().StringP("runtime", "r", "", "Data to use for rendering templates. Will be prefixed with \"runtime\".")
	rootCmd.Flags().Bool("output-runtime-placeholder-files", false, "Print file that contain runtime placeholder as github-action output to replace them with another template engine later.")
	rootCmd.Flags().Bool("copy-permissions", false, "Copy the user, group and mode of the template.")
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
	if config.secrets, err = cmd.Flags().GetString("secrets"); err != nil {
		return
	}
	if config.runtimeData, err = cmd.Flags().GetString("runtime"); err != nil {
		return
	}
	if config.outputDir, err = cmd.Flags().GetString("output-dir"); err != nil {
		return
	}
	if config.copyPermissions, err = cmd.Flags().GetBool("copy-permissions"); err != nil {
		return
	}
	if config.outputRuntimePlaceholderFiles, err = cmd.Flags().GetBool("output-runtime-placeholder-files"); err != nil {
		return
	}
	return
}

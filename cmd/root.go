package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "template-action",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		templateDir := getStringFlag(cmd, "template-dir")
		templateExtension := getStringFlag(cmd, "template-extension")
		inputData := getStringsFlag(cmd, "data")
		outputDir := getStringFlag(cmd, "output-dir")
		copyPermissions := getBoolFlag(cmd, "copy-permissions")

		templates := LoadTemplateFiles(templateDir, templateExtension)
		data := ParseInputData(inputData)
		for _, template := range templates {
			template.Render(data, outputDir, copyPermissions)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("template-dir", "i", "./", "Set the input directory.")
	rootCmd.Flags().StringP("template-extension", "t", ".template", "Set a file extension to detect templates.")
	rootCmd.Flags().StringP("output-dir", "o", "./", "Set the output directory.")
	rootCmd.Flags().StringArrayP("data", "d", []string{}, "Data to use for rendering templates as yaml or json objects. Multiple objects will be merged before rendering.")
	rootCmd.Flags().Bool("copy-permissions", false, "Copy the user, group and mode of the template.")
}

func getStringFlag(cmd *cobra.Command, name string) string {
	value, err := cmd.Flags().GetString(name)
	handleError(err)
	return value
}

func getStringsFlag(cmd *cobra.Command, name string) []string {
	value, err := cmd.Flags().GetStringArray(name)
	handleError(err)
	return value
}

func getBoolFlag(cmd *cobra.Command, name string) bool {
	value, err := cmd.Flags().GetBool(name)
	handleError(err)
	return value
}

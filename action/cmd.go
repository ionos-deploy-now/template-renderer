package action

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
		templates := loadTemplateFiles()
		data := getDataFromEnvironment()
		for _, t := range templates {
			t.Fill(data)
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
	rootCmd.PersistentFlags().StringVarP(&templateDir, "input", "i", "./", "Set the input directory.")
	rootCmd.PersistentFlags().StringVarP(&templateExtension, "template-extension", "t", ".template", "Set a file extension to detect templates.")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "./", "Set the output directory.")
	rootCmd.PersistentFlags().StringVarP(&envVarPrefix, "env-var-prefix", "e", "data.", "Specify a prefix to select environment variables as input values.")
}

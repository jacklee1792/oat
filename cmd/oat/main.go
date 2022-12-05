package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "oat",
	Short: "oat - tools to manipulate OpenAPI specifications",
}

var (
	inputPath  string
	outputPath string
)

func init() {
	rootCmd.AddCommand(filterOpsCmd)
	rootCmd.AddCommand(cleanSchemasCmd)
	rootCmd.AddCommand(validateCmd)
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}

package main

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "oat",
	Short: "oat - tools to manipulate OpenAPI specifications",
}

var (
	inputPath  string
	outputPath string
)

func init() {
	rootCmd.PersistentFlags().IntVarP(
		&openapi3.CircularReferenceCounter,
		"circular-reference-counter",
		"c",
		openapi3.CircularReferenceCounter, // use default as set in openapi3
		"Sets max depth for circular references (openapi3.CircularReferenceCounter)",
	)

	rootCmd.AddCommand(filterOpsCmd)
	rootCmd.AddCommand(cleanSchemasCmd)
	rootCmd.AddCommand(validateCmd)
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jacklee1792/oat"
	"github.com/spf13/cobra"
)

var cleanSchemasCmd = &cobra.Command{
	Use:   "clean-schemas [exceptSchemaName ...]",
	Short: "removed unused schemas in components: schemas (with exceptions)",
	PreRun: func(cmd *cobra.Command, args []string) {
		if outputPath == "" {
			outputPath = makeOutputPath(inputPath, "cleaned")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		spec, err := openapi3.NewLoader().LoadFromFile(inputPath)
		cobra.CheckErr(err)
		sc := oat.SchemaCleaner{}
		n := sc.Clean(spec, args)
		data, err := json.MarshalIndent(spec, "", "  ")
		cobra.CheckErr(err)
		err = os.WriteFile(outputPath, data, 0666)
		cobra.CheckErr(err)
		fmt.Printf("Done! Removed %d schemas\n", n)
	},
}

func init() {
	cleanSchemasCmd.Flags().StringVarP(&inputPath, "input", "i", "", "Input path of OpenAPI specification")
	cobra.CheckErr(cleanSchemasCmd.MarkFlagRequired("input"))
	cleanSchemasCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output path of cleaned OpenAPI specification")
}

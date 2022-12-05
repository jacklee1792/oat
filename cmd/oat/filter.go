package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jacklee1792/oat"
	"github.com/spf13/cobra"
)

var filterOpsCmd = &cobra.Command{
	Use:   "filter-ops [operationId ...]",
	Short: "select operations to keep by operation ID",
	PreRun: func(cmd *cobra.Command, args []string) {
		if outputPath == "" {
			outputPath = makeOutputPath(inputPath, "filtered")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		spec, err := openapi3.NewLoader().LoadFromFile(inputPath)
		cobra.CheckErr(err)
		n := oat.FilterOperationsById(spec, args)
		data, err := json.MarshalIndent(spec, "", "\t")
		cobra.CheckErr(err)
		err = os.WriteFile(outputPath, data, 0666)
		cobra.CheckErr(err)
		fmt.Printf("Done! Removed %d operations\n", n)
	},
}

func init() {
	filterOpsCmd.Flags().StringVarP(&inputPath, "input", "i", "", "Input path of OpenAPI specification")
	cobra.CheckErr(filterOpsCmd.MarkFlagRequired("input"))
	filterOpsCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output path of filtered OpenAPI specification")
}

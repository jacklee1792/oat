package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "check validity of the provided OpenAPI specification",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := openapi3.NewLoader().LoadFromFile(inputPath)
		cobra.CheckErr(err)
		fmt.Println("Spec is valid")
	},
}

func init() {
	validateCmd.Flags().StringVarP(&inputPath, "input", "i", "", "Input path of OpenAPI specification")
	cobra.CheckErr(validateCmd.MarkFlagRequired("input"))
}

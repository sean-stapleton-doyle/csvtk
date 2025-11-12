package cmd

import (
	"github.com/spf13/cobra"
)

func getDelimiter(cmd *cobra.Command) rune {
	delimiterStr, _ := cmd.Flags().GetString("delimiter")

	if delimiterStr == "" {
		return ','
	}

	if delimiterStr == "\\t" || delimiterStr == "\t" {
		return '\t'
	}

	if len(delimiterStr) > 0 {
		return rune(delimiterStr[0])
	}

	return ','
}

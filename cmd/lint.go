package cmd

import (
	"fmt"
	"os"

	"sean-stapleton-doyle/csvtk/pkg/csvlint"

	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint [file]",
	Short: "Validate a CSV file against RFC 4180",
	Long: `Validate a CSV file according to RFC 4180 standards.
Reports any parsing errors or inconsistent field counts.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		delimiter := getDelimiter(cmd)
		lazyQuotes, _ := cmd.Flags().GetBool("lazy-quotes")

		errors, hasFatalError, err := csvlint.ValidateFile(filename, delimiter, lazyQuotes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating file: %v\n", err)
			os.Exit(1)
		}

		if hasFatalError {
			fmt.Println("❌ CSV validation failed with fatal errors:")
			for _, e := range errors {
				fmt.Printf("  %s\n", e.Error())
			}
			os.Exit(1)
		}

		if len(errors) > 0 {
			fmt.Printf("⚠️  CSV validation completed with %d warning(s):\n", len(errors))
			for _, e := range errors {
				fmt.Printf("  %s\n", e.Error())
			}
			os.Exit(1)
		}

		fmt.Println("✓ CSV file is valid")
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)
	lintCmd.Flags().StringP("delimiter", "d", ",", "Field delimiter")
	lintCmd.Flags().Bool("lazy-quotes", false, "Allow lazy quotes (less strict parsing)")
}

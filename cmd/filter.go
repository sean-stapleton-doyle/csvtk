package cmd

import (
	"fmt"
	"os"

	"sean-stapleton-doyle/csvtk/pkg/csveditor"
	"sean-stapleton-doyle/csvtk/pkg/csvparser"

	"github.com/spf13/cobra"
)

var filterCmd = &cobra.Command{
	Use:   "filter [column] [value] [file]",
	Short: "Filter rows based on column values",
	Long: `Filter rows in a CSV file based on column values.
By default, filters rows where the column value equals the given value.
Use flags to specify different filter operations including regex and numeric comparisons.

Examples:
  csvtk filter City "New York" data.csv
  csvtk filter Age 30 data.csv --operator ">"
  csvtk filter Email "@gmail.com" data.csv --regex
  cat data.csv | csvtk filter Name "John" -  # from stdin`,
	Args: cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		columnName := args[0]
		value := args[1]
		filename := ""
		if len(args) == 3 {
			filename = args[2]
		} else {
			filename = "-"
		}

		config := csvparser.DefaultConfig()
		config.Delimiter = getDelimiter(cmd)

		csv, err := csvparser.ParseFromFileOrStdin(filename, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing CSV: %v\n", err)
			os.Exit(1)
		}

		operator, _ := cmd.Flags().GetString("operator")
		regex, _ := cmd.Flags().GetBool("regex")

		if regex {
			operator = "regex"
		}

		if operator == "" {
			if contains, _ := cmd.Flags().GetBool("contains"); contains {
				operator = "contains"
			} else if startsWith, _ := cmd.Flags().GetBool("starts-with"); startsWith {
				operator = "starts-with"
			} else if endsWith, _ := cmd.Flags().GetBool("ends-with"); endsWith {
				operator = "ends-with"
			} else if notEquals, _ := cmd.Flags().GetBool("not-equals"); notEquals {
				operator = "not-equals"
			} else {
				operator = "equals"
			}
		}

		strategy := csveditor.NewFilterStrategy(operator)

		filtered, err := csveditor.FilterWithStrategy(csv, columnName, value, strategy)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error filtering CSV: %v\n", err)
			os.Exit(1)
		}

		output, _ := cmd.Flags().GetString("output")
		if output == "" || output == "-" {

			err = filtered.Write(os.Stdout, config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing CSV: %v\n", err)
				os.Exit(1)
			}
		} else {
			err = filtered.WriteToFile(output, config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing CSV: %v\n", err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "Filtered %d rows to %s\n", len(filtered.Records), output)
		}
	},
}

func init() {
	rootCmd.AddCommand(filterCmd)
	filterCmd.Flags().StringP("delimiter", "d", ",", "Field delimiter")
	filterCmd.Flags().StringP("output", "o", "", "Output file (defaults to stdout)")
	filterCmd.Flags().StringP("operator", "p", "", "Filter operator: equals, contains, starts-with, ends-with, not-equals, regex, >, <, >=, <=, ==, !=")
	filterCmd.Flags().Bool("regex", false, "Use regex matching")

	filterCmd.Flags().Bool("contains", false, "Filter rows where column contains the value")
	filterCmd.Flags().Bool("starts-with", false, "Filter rows where column starts with the value")
	filterCmd.Flags().Bool("ends-with", false, "Filter rows where column ends with the value")
	filterCmd.Flags().Bool("not-equals", false, "Filter rows where column does not equal the value")
}

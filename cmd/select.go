package cmd

import (
	"fmt"
	"os"
	"strings"

	"sean-stapleton-doyle/csvtk/pkg/csveditor"
	"sean-stapleton-doyle/csvtk/pkg/csvparser"

	"github.com/spf13/cobra"
)

var selectCmd = &cobra.Command{
	Use:   "select [columns] [file]",
	Short: "Select specific columns from a CSV file",
	Long: `Select and output only specific columns from a CSV file.
Columns should be specified as a comma-separated list of column names.

Example: csvtk select "Name,Email,Age" data.csv
         cat data.csv | csvtk select "Name,Age" -`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		columnsStr := args[0]
		filename := "-"
		if len(args) == 2 {
			filename = args[1]
		}

		columnNames := strings.Split(columnsStr, ",")
		for i := range columnNames {
			columnNames[i] = strings.TrimSpace(columnNames[i])
		}

		config := csvparser.DefaultConfig()
		config.Delimiter = getDelimiter(cmd)

		csv, err := csvparser.ParseFromFileOrStdin(filename, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing CSV: %v\n", err)
			os.Exit(1)
		}

		selected, err := csveditor.SelectColumns(csv, columnNames)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error selecting columns: %v\n", err)
			os.Exit(1)
		}

		output, _ := cmd.Flags().GetString("output")
		if output == "" {

			err = selected.Write(os.Stdout, config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing CSV: %v\n", err)
				os.Exit(1)
			}
		} else {
			err = selected.WriteToFile(output, config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing CSV: %v\n", err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "Selected %d columns to %s\n", len(columnNames), output)
		}
	},
}

func init() {
	rootCmd.AddCommand(selectCmd)
	selectCmd.Flags().StringP("delimiter", "d", ",", "Field delimiter")
	selectCmd.Flags().StringP("output", "o", "", "Output file (defaults to stdout)")
}

package cmd

import (
	"fmt"
	"os"

	"sean-stapleton-doyle/csvtk/pkg/csveditor"
	"sean-stapleton-doyle/csvtk/pkg/csvparser"

	"github.com/spf13/cobra"
)

var sortCmd = &cobra.Command{
	Use:   "sort [column] [file]",
	Short: "Sort a CSV file by a column",
	Long: `Sort a CSV file by the values in a specified column.

Examples:
  csvtk sort Age data.csv -o sorted.csv
  cat data.csv | csvtk sort Name - > sorted.csv`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		columnName := args[0]
		filename := "-"
		if len(args) == 2 {
			filename = args[1]
		}

		config := csvparser.DefaultConfig()
		config.Delimiter = getDelimiter(cmd)

		csv, err := csvparser.ParseFromFileOrStdin(filename, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing CSV: %v\n", err)
			os.Exit(1)
		}

		descending, _ := cmd.Flags().GetBool("descending")

		sortConfig := csveditor.SortConfig{
			ColumnName: columnName,
			Descending: descending,
		}

		err = csveditor.Sort(csv, sortConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sorting CSV: %v\n", err)
			os.Exit(1)
		}

		output, _ := cmd.Flags().GetString("output")
		direction := "ascending"
		if descending {
			direction = "descending"
		}

		if output == "" || output == "-" {

			err = csv.Write(os.Stdout, config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing CSV: %v\n", err)
				os.Exit(1)
			}
		} else {
			err = csv.WriteToFile(output, config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing CSV: %v\n", err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "Sorted by column '%s' (%s) in %s\n", columnName, direction, output)
		}
	},
}

func init() {
	rootCmd.AddCommand(sortCmd)
	sortCmd.Flags().StringP("delimiter", "d", ",", "Field delimiter")
	sortCmd.Flags().StringP("output", "o", "", "Output file (defaults to input file)")
	sortCmd.Flags().BoolP("descending", "r", false, "Sort in descending order")
}

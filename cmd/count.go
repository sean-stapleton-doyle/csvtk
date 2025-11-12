package cmd

import (
	"fmt"
	"os"

	"sean-stapleton-doyle/csvtk/pkg/csvparser"

	"github.com/spf13/cobra"
)

var countCmd = &cobra.Command{
	Use:   "count",
	Short: "Count rows or columns in a CSV file",
	Long:  `Count the number of rows or columns in a CSV file.`,
}

var countRowsCmd = &cobra.Command{
	Use:   "rows [file]",
	Short: "Count the number of rows in a CSV file",
	Long:  `Count the number of data rows in a CSV file (excluding the header). Use "-" or omit file to read from stdin.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := "-"
		if len(args) > 0 {
			filename = args[0]
		}

		config := csvparser.DefaultConfig()
		config.Delimiter = getDelimiter(cmd)

		csv, err := csvparser.ParseFromFileOrStdin(filename, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing CSV: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(csv.CountRows())
	},
}

var countColumnsCmd = &cobra.Command{
	Use:   "columns [file]",
	Short: "Count the number of columns in a CSV file",
	Long:  `Count the number of columns in a CSV file. Use "-" or omit file to read from stdin.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := "-"
		if len(args) > 0 {
			filename = args[0]
		}

		config := csvparser.DefaultConfig()
		config.Delimiter = getDelimiter(cmd)

		csv, err := csvparser.ParseFromFileOrStdin(filename, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing CSV: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(csv.CountColumns())
	},
}

func init() {
	rootCmd.AddCommand(countCmd)
	countCmd.AddCommand(countRowsCmd)
	countCmd.AddCommand(countColumnsCmd)

	countRowsCmd.Flags().StringP("delimiter", "d", ",", "Field delimiter")
	countColumnsCmd.Flags().StringP("delimiter", "d", ",", "Field delimiter")
}

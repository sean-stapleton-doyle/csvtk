package cmd

import (
	"fmt"
	"os"
	"strconv"

	"sean-stapleton-doyle/csvtk/pkg/csveditor"
	"sean-stapleton-doyle/csvtk/pkg/csvparser"

	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move rows or columns in a CSV file",
	Long:  `Move rows or columns to a different position in a CSV file.`,
}

var moveColumnCmd = &cobra.Command{
	Use:   "column [column-name] [target-index] [file]",
	Short: "Move a column to a different position",
	Long: `Move a column to a different position in the CSV file.
The target index is 0-based.`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		columnName := args[0]
		targetIndexStr := args[1]
		filename := args[2]

		targetIndex, err := strconv.Atoi(targetIndexStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid target index: %v\n", err)
			os.Exit(1)
		}

		config := csvparser.DefaultConfig()
		config.Delimiter = getDelimiter(cmd)

		csv, err := csvparser.ParseFile(filename, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing CSV: %v\n", err)
			os.Exit(1)
		}

		err = csveditor.MoveColumn(csv, columnName, targetIndex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error moving column: %v\n", err)
			os.Exit(1)
		}

		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = filename
		}

		err = csv.WriteToFile(output, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing CSV: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Moved column '%s' to index %d in %s\n", columnName, targetIndex, output)
	},
}

var moveRowCmd = &cobra.Command{
	Use:   "row [old-index] [new-index] [file]",
	Short: "Move a row to a different position",
	Long: `Move a row to a different position in the CSV file.
Both indices are 0-based and refer to data rows (excluding the header).`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		oldIndexStr := args[0]
		newIndexStr := args[1]
		filename := args[2]

		oldIndex, err := strconv.Atoi(oldIndexStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid old index: %v\n", err)
			os.Exit(1)
		}

		newIndex, err := strconv.Atoi(newIndexStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid new index: %v\n", err)
			os.Exit(1)
		}

		config := csvparser.DefaultConfig()
		config.Delimiter = getDelimiter(cmd)

		csv, err := csvparser.ParseFile(filename, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing CSV: %v\n", err)
			os.Exit(1)
		}

		err = csveditor.MoveRow(csv, oldIndex, newIndex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error moving row: %v\n", err)
			os.Exit(1)
		}

		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = filename
		}

		err = csv.WriteToFile(output, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing CSV: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Moved row %d to index %d in %s\n", oldIndex, newIndex, output)
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)
	moveCmd.AddCommand(moveColumnCmd)
	moveCmd.AddCommand(moveRowCmd)

	moveColumnCmd.Flags().StringP("delimiter", "d", ",", "Field delimiter")
	moveColumnCmd.Flags().StringP("output", "o", "", "Output file (defaults to input file)")

	moveRowCmd.Flags().StringP("delimiter", "d", ",", "Field delimiter")
	moveRowCmd.Flags().StringP("output", "o", "", "Output file (defaults to input file)")
}

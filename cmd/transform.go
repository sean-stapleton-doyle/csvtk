package cmd

import (
	"fmt"
	"os"

	"sean-stapleton-doyle/csvtk/pkg/csveditor"
	"sean-stapleton-doyle/csvtk/pkg/csvparser"

	"github.com/spf13/cobra"
)

var transformCmd = &cobra.Command{
	Use:   "transform",
	Short: "Transform CSV data",
	Long:  `Transform CSV data by applying various operations to columns or all cells.`,
}

var transformLowerCmd = &cobra.Command{
	Use:   "lower [column] [file]",
	Short: "Convert text to lowercase",
	Long: `Convert text in a column to lowercase. If no column is specified, transforms all cells.

Examples:
  csvtk transform lower Name data.csv
  csvtk transform lower --all data.csv
  cat data.csv | csvtk transform lower Name -`,
	Args: cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		runTransform(cmd, args, csveditor.ToLower, "lowercase")
	},
}

var transformUpperCmd = &cobra.Command{
	Use:   "upper [column] [file]",
	Short: "Convert text to uppercase",
	Long: `Convert text in a column to uppercase. If no column is specified, transforms all cells.

Examples:
  csvtk transform upper Name data.csv
  csvtk transform upper --all data.csv
  cat data.csv | csvtk transform upper City -`,
	Args: cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		runTransform(cmd, args, csveditor.ToUpper, "uppercase")
	},
}

var transformReplaceCmd = &cobra.Command{
	Use:   "replace [old] [new] [column] [file]",
	Short: "Replace text in a column",
	Long: `Replace occurrences of text in a column.

Examples:
  csvtk transform replace "old" "new" Name data.csv
  csvtk transform replace "@example.com" "@newdomain.com" Email data.csv --all
  cat data.csv | csvtk transform replace "," ";" City -`,
	Args: cobra.RangeArgs(2, 4),
	Run: func(cmd *cobra.Command, args []string) {
		old := args[0]
		new := args[1]

		var columnName string
		var filename string

		if len(args) == 4 {
			columnName = args[2]
			filename = args[3]
		} else if len(args) == 3 {
			columnName = args[2]
			filename = "-"
		} else {

			allCols, _ := cmd.Flags().GetBool("all")
			if !allCols {
				fmt.Fprintf(os.Stderr, "Error: column name required or use --all flag\n")
				os.Exit(1)
			}
			filename = "-"
		}

		config := csvparser.DefaultConfig()
		config.Delimiter = getDelimiter(cmd)

		csv, err := csvparser.ParseFromFileOrStdin(filename, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing CSV: %v\n", err)
			os.Exit(1)
		}

		allCols, _ := cmd.Flags().GetBool("all")
		transformFunc := csveditor.ReplaceAll(old, new)

		if allCols || columnName == "" {
			err = csveditor.TransformAll(csv, transformFunc)
		} else {
			err = csveditor.TransformColumn(csv, columnName, transformFunc)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error transforming CSV: %v\n", err)
			os.Exit(1)
		}

		output, _ := cmd.Flags().GetString("output")
		writeOutput(csv, output, config, "Replaced text in")
	},
}

var transformTrimCmd = &cobra.Command{
	Use:   "trim [column] [file]",
	Short: "Trim whitespace from cells",
	Long: `Trim leading and trailing whitespace from cells in a column or all cells.

Examples:
  csvtk transform trim Name data.csv
  csvtk transform trim --all data.csv`,
	Args: cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		runTransform(cmd, args, csveditor.TrimSpace, "trimmed")
	},
}

func runTransform(cmd *cobra.Command, args []string, transform csveditor.TransformFunc, operation string) {
	var columnName string
	var filename string

	allCols, _ := cmd.Flags().GetBool("all")

	if len(args) == 2 {
		columnName = args[0]
		filename = args[1]
	} else if len(args) == 1 {
		if allCols {
			filename = args[0]
		} else {
			columnName = args[0]
			filename = "-"
		}
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

	if allCols || columnName == "" {
		err = csveditor.TransformAll(csv, transform)
	} else {
		err = csveditor.TransformColumn(csv, columnName, transform)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error transforming CSV: %v\n", err)
		os.Exit(1)
	}

	output, _ := cmd.Flags().GetString("output")
	writeOutput(csv, output, config, fmt.Sprintf("Transformed to %s", operation))
}

func writeOutput(csv *csvparser.CSV, output string, config *csvparser.Config, successMsg string) {
	if output == "" || output == "-" {

		err := csv.Write(os.Stdout, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing CSV: %v\n", err)
			os.Exit(1)
		}
	} else {
		err := csv.WriteToFile(output, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing CSV: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "%s %s\n", successMsg, output)
	}
}

func init() {
	rootCmd.AddCommand(transformCmd)
	transformCmd.AddCommand(transformLowerCmd)
	transformCmd.AddCommand(transformUpperCmd)
	transformCmd.AddCommand(transformReplaceCmd)
	transformCmd.AddCommand(transformTrimCmd)

	for _, cmd := range []*cobra.Command{transformLowerCmd, transformUpperCmd, transformReplaceCmd, transformTrimCmd} {
		cmd.Flags().StringP("delimiter", "d", ",", "Field delimiter")
		cmd.Flags().StringP("output", "o", "", "Output file (defaults to stdout)")
		cmd.Flags().Bool("all", false, "Apply transformation to all columns")
	}
}

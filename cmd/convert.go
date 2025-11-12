package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sean-stapleton-doyle/csvtk/pkg/csvparser"

	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert [file]",
	Short: "Convert between CSV and TSV formats",
	Long: `Convert a CSV file to TSV (tab-delimited) or vice versa.
Use --to-tsv or --to-csv flags to specify the conversion direction.
If neither flag is specified, the tool will infer based on the input file extension.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		toTSV, _ := cmd.Flags().GetBool("to-tsv")
		toCSV, _ := cmd.Flags().GetBool("to-csv")

		if !toTSV && !toCSV {
			ext := strings.ToLower(filepath.Ext(filename))
			if ext == ".tsv" || ext == ".tab" {
				toCSV = true
			} else {
				toTSV = true
			}
		}

		var inputDelimiter, outputDelimiter rune
		var outputExt string

		if toTSV {
			inputDelimiter = ','
			outputDelimiter = '\t'
			outputExt = ".tsv"
		} else {
			inputDelimiter = '\t'
			outputDelimiter = ','
			outputExt = ".csv"
		}

		if cmd.Flags().Changed("delimiter") {
			inputDelimiter = getDelimiter(cmd)
		}

		config := csvparser.DefaultConfig()
		config.Delimiter = inputDelimiter

		csv, err := csvparser.ParseFile(filename, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
			os.Exit(1)
		}

		output, _ := cmd.Flags().GetString("output")
		if output == "" {

			ext := filepath.Ext(filename)
			output = strings.TrimSuffix(filename, ext) + outputExt
		}

		outputConfig := csvparser.DefaultConfig()
		outputConfig.Delimiter = outputDelimiter

		err = csv.WriteToFile(output, outputConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			os.Exit(1)
		}

		format := "TSV"
		if toCSV {
			format = "CSV"
		}
		fmt.Printf("Converted to %s: %s\n", format, output)
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringP("delimiter", "d", "", "Input delimiter (auto-detected if not specified)")
	convertCmd.Flags().StringP("output", "o", "", "Output file (auto-generated if not specified)")
	convertCmd.Flags().Bool("to-tsv", false, "Convert to TSV format")
	convertCmd.Flags().Bool("to-csv", false, "Convert to CSV format")
}

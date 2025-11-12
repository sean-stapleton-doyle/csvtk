package cmd

import (
	"fmt"
	"os"

	"sean-stapleton-doyle/csvtk/pkg/csveditor"
	"sean-stapleton-doyle/csvtk/pkg/csvparser"

	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename [old-name] [new-name] [file]",
	Short: "Rename a column header",
	Long: `Rename a column header in a CSV file.

Examples:
  csvtk rename "Old Name" "New Name" data.csv
  csvtk rename Email EmailAddress data.csv -o updated.csv
  cat data.csv | csvtk rename City Location -`,
	Args: cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		oldName := args[0]
		newName := args[1]
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

		err = csveditor.RenameHeader(csv, oldName, newName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error renaming header: %v\n", err)
			os.Exit(1)
		}

		output, _ := cmd.Flags().GetString("output")
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
			fmt.Fprintf(os.Stderr, "Renamed header '%s' to '%s' in %s\n", oldName, newName, output)
		}
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
	renameCmd.Flags().StringP("delimiter", "d", ",", "Field delimiter")
	renameCmd.Flags().StringP("output", "o", "", "Output file (defaults to stdout)")
}

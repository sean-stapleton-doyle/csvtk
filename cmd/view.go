package cmd

import (
	"fmt"
	"os"

	"sean-stapleton-doyle/csvtk/pkg/csvparser"
	"sean-stapleton-doyle/csvtk/pkg/csvviewer"

	"github.com/spf13/cobra"
)

var viewCommand = &cobra.Command{
	Use:   "view [file]",
	Short: "View a CSV file in an interactive terminal viewer",
	Long: `Open a CSV file in an interactive terminal viewer with keyboard navigation.

Keyboard shortcuts:
  ↑/k: Move up one row
  ↓/j: Move down one row
  PgUp: Move up one page
  PgDn: Move down one page
  g/home: Go to first row
  G/end: Go to last row
  q: Quit viewer`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		config := csvparser.DefaultConfig()
		config.Delimiter = getDelimiter(cmd)

		csv, err := csvparser.ParseFile(filename, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing CSV: %v\n", err)
			os.Exit(1)
		}

		err = csvviewer.Run(csv, filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running viewer: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(viewCommand)
	viewCommand.Flags().StringP("delimiter", "d", ",", "Field delimiter")
}

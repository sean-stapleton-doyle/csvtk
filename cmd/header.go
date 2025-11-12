package cmd

import (
	"fmt"
	"os"
	"strings"

	"sean-stapleton-doyle/csvtk/pkg/csvparser"

	"github.com/spf13/cobra"
)

var headerCmd = &cobra.Command{
	Use:   "header [file]",
	Short: "Display the header row of a CSV file",
	Long:  `Display the header row of a CSV file. Use "-" or omit file to read from stdin.`,
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

		numbered, _ := cmd.Flags().GetBool("numbered")

		if numbered {
			for i, col := range csv.Header {
				fmt.Printf("%d: %s\n", i, col)
			}
		} else {
			fmt.Println(strings.Join(csv.Header, ", "))
		}
	},
}

func init() {
	rootCmd.AddCommand(headerCmd)
	headerCmd.Flags().StringP("delimiter", "d", ",", "Field delimiter")
	headerCmd.Flags().BoolP("numbered", "n", false, "Show column numbers")
}

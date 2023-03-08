package cmd

import (
	"encoding/csv"
	"errors"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the latest value set for a key",
	Long: "Get the latest value set for a key",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires exactly one argument")
		}

		key := args[0]
		return validateKey(key)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		filepath := rootCmd.PersistentFlags().Lookup("database").Value.String()
		db, err := os.Open(filepath)

		cobra.CheckErr(err)
		defer db.Close()

		lookupKey := args[0]

		csvReader := csv.NewReader(db)
		keyFound := false
		var answer string

		for {
			record, err := csvReader.Read()

			if err == io.EOF {
				break
			} else {
				cobra.CheckErr(err)
			}

			key := record[0]

			if key == lookupKey {
				keyFound = true
				answer = record[1]
			}
		}

		if keyFound {
			cmd.Printf("Answer: %s", answer)
			return nil
		}	else {
			return errors.New("the key is not present in the database")
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

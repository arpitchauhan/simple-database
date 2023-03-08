package cmd

import (
	"encoding/csv"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Add a key-value pair to the database",
	Long:  "Add a key-value pair to the database",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("requires exactly 2 args")
		}

		key := args[0]

		return validateKey(key)
	},
	Run: func(cmd *cobra.Command, args []string) {
		filepath := rootCmd.PersistentFlags().Lookup("database").Value.String()
		db, err := os.OpenFile(
			filepath,
			os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0o644,
		)
		cobra.CheckErr(err)
		defer db.Close()

		csvWriter := csv.NewWriter(db)
		err = csvWriter.Write(args)

		cobra.CheckErr(err)
		csvWriter.Flush()

		err = csvWriter.Error()
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}

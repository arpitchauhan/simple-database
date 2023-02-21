/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/package cmd

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the value of a variable",
	Long:  "Set the value of a variable",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("requires exactly 2 args")
		}

		key := args[0]

		if len(strings.TrimSpace(key)) == 0 {
			return errors.New("key cannot be empty")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db, err := os.OpenFile(
			rootCmd.PersistentFlags().Lookup("database").Value.String(),
			os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0o644,
		)

		check(err)

		csvWriter := csv.NewWriter(db)
		err = csvWriter.Write(args)

		check(err)
		csvWriter.Flush()

		err = csvWriter.Error()
		check(err)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	rootCmd.PersistentFlags().String("database", "database.csv", "database contents file")
}

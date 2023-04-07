package cmd

import (
	"github.com/arpitchauhan/simple-database/client"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var setValueForKey = client.SetValueForKey

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Add a key-value pair to the database",
	Long:  "Add a key-value pair to the database",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		err := setValueForKey(key, value)

		if err != nil {
			status, _ := status.FromError(err)
			if status.Code() == codes.Unavailable {
				cmd.Printf("Error: the server is not running")
				return
			}

			cobra.CheckErr(err)
		}

		cmd.Printf("Successful!")
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}

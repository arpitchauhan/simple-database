package cmd

import (
	"github.com/arpitchauhan/simple-database/client"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var getValueForKey = client.GetValueForKey

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the latest value set for a key",
	Long: "Get the latest value set for a key",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lookupKey := args[0]
		answer, err := getValueForKey(lookupKey)

		if err != nil {
			status, _ := status.FromError(err)

			if status.Code() == codes.Unavailable {
				cmd.Printf("Error: the server is not running")
				return
			} else if status.Code() == codes.NotFound {
				cmd.Printf("Error: the key was not found")
				return
			}

			cobra.CheckErr(err)
		}

		cmd.Printf("Answer: %s", answer)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

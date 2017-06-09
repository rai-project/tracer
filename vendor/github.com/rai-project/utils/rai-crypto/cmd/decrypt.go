package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rai-project/utils"
	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:          "decrypt",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.Errorf("expected one argument, got %v", len(args))
		}
		val, err := utils.DecryptStringBase64(appsecret, args[0])
		if err != nil {
			return err
		}
		fmt.Print(val)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(decryptCmd)
}

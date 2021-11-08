package cmd

import (
	"github.com/isayme/go-shadowsocks/cmd/local"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(localCmd)
}

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "shadowsocks local client",
	Run: func(cmd *cobra.Command, args []string) {
		local.Run()
	},
}

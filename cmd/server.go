package cmd

import (
	"github.com/isayme/go-shadowsocks/cmd/server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(severCmd)
}

var severCmd = &cobra.Command{
	Use:   "server",
	Short: "shadowsocks server",
	Run: func(cmd *cobra.Command, args []string) {
		server.Run()
	},
}

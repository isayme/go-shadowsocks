package cmd

import (
	"fmt"
	"os"

	"github.com/isayme/go-shadowsocks/util"
	"github.com/spf13/cobra"
)

var versionFlag bool

func init() {
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "wstunnel version")
}

var rootCmd = &cobra.Command{
	Use: "shadowsocks",
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			util.PrintVersion()
			os.Exit(0)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

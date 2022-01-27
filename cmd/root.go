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
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(localCmd)
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

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "shadowsocks server",
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "shadowsocks local",
	Run: func(cmd *cobra.Command, args []string) {
		runLocal()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

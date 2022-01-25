/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"

	"github.com/Joe-Degs/kake/cmd/kaked"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "start the kake server daemon (kaked)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		daemon := kaked.DefaultDaemon()
		// TODO(Joe-Degs):
		// do the config loading thing here then pass to daemon
		if err := daemon.Init(nil); err != nil {
			log.Printf("error starting kaked daemon: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kakkle",
	Short: "kakkle tool for managing kake services",
	Long: `Kakkle provides a interface for starting kake server daemon(kaded)
and serves as a command line interface for dealing communicating with the
daemon`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kake.yaml)")
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

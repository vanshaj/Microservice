/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vanshaj/Microservice/cliApps/cobra/scan"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan the hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}
		hl := &scan.HostsList{}
		fr, err := os.Open(f)
		if err != nil {
			return err
		}
		if err = hl.Load(fr); err != nil {
			return err
		}
		ports, err := cmd.Flags().GetIntSlice("ports")
		if err != nil {
			return err
		}
		scanAction(os.Stdout, hl, ports)
		return nil
	},
}

func scanAction(out io.Writer, hl *scan.HostsList, ports []int) {
	results := scan.Run(hl, ports)
	var message strings.Builder
	for _, result := range results {
		for _, portState := range result.PortStates {
			fmt.Fprintf(&message, "Host: %s has port %d in %s stage\n", result.Host, portState.Port, portState.Open.String())
		}
	}
	fmt.Fprint(out, message.String())
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	scanCmd.PersistentFlags().IntSliceP("ports", "p", []int{80}, "list of ports to scan")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

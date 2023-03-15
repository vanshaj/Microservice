/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/vanshaj/Microservice/cliApps/cobra/scan"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:          "add",
	Short:        "add hosts",
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		hosts_file, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}
		return addAction(os.Stdout, hosts_file, args)
	},
}

func addAction(w io.Writer, hosts_file string, args []string) error {
	f, err := os.OpenFile(hosts_file, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	h := &scan.HostsList{}
	if err := h.Load(f); err != nil {
		return err
	}
	err = os.Truncate(hosts_file, 0)
	if err != nil {
		return err
	}
	for _, eachHost := range args {
		if err := h.Add(eachHost); err != nil {
			return err
		}
		fmt.Fprintln(w, "added host: ", eachHost)
	}
	return h.Save(f)
}

func init() {
	hostsCmd.AddCommand(addCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

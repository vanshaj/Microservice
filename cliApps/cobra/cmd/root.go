/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "Pscan",
	Short:   "Pscan",
	Long:    `pScan allows you to add, list and delete hosts from the scan list`,
	Version: "0.1",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		// if config file is passed via the flag --config set the file
		viper.SetConfigFile(cfgFile)
	} else {
		// get the home directory
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// this is the path to look for the config file in
		viper.AddConfigPath(home)
		// this is the name of the config with extension i.e if there is no --config file it will use $HOME/.pscan.yaml as the configfile
		viper.SetConfigName(".pscan")
	}
	// Will check for an environment variable every time viper.Get request is made
	viper.AutomaticEnv()
	// ReadInConfig will discover and load the configuration file from disk and key/value stores, searching in one of the defined paths.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file: ", viper.ConfigFileUsed())
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pscan.yaml)")
	rootCmd.PersistentFlags().StringP("hosts-file", "f", "pScan.hosts", "pScan hosts file")
	replacer := strings.NewReplacer("-", "_")
	// will replace all the env variables with - in it to _
	viper.SetEnvKeyReplacer(replacer)
	//SetEnvPrefix defines a prefix that ENVIRONMENT variables will use. E.g. if your prefix is "spf", the env registry will look for env variables that start with "SPF_".
	viper.SetEnvPrefix("PSCAN")
	// bind --hosts-file flag to an environment variable i.e PSCAN_HOSTS_FILE
	viper.BindPFlag("hosts-file", rootCmd.PersistentFlags().Lookup("hosts-file"))
	versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`
	rootCmd.SetVersionTemplate(versionTemplate)
}

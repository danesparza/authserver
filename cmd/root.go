package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/logutils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile               string
	ProblemWithConfigFile bool
	loglevel              string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "authserver",
	Short: "A service to manage, issue and validate OAuth2 grants",
	Long: `A service to manage, issue and validate OAuth2 grants.
	
The service has a management UI, an API layer, and OAuth2 helper methods`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/authserver.yml)")

	rootCmd.PersistentFlags().StringVarP(&loglevel, "loglevel", "l", "WARN", "Log level: DEBUG/INFO/WARN/ERROR")

	//	Bind config flags for optional config file override:
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.SetConfigName("authserver") // name of config file (without extension)
		viper.AddConfigPath(home)         // adding home directory as first search path
		viper.AddConfigPath(".")          // also look in the working directory
	}

	viper.AutomaticEnv() // read in environment variables that match

	//	Set our defaults
	viper.SetDefault("loglevel", "INFO")
	viper.SetDefault("apiservice.port", "3000")
	viper.SetDefault("uiservice.port", "3001")
	viper.SetDefault("apiservice.allowed-origins", "*")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		ProblemWithConfigFile = true
	}

	//	Set the log level from config (if we have it)
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(viper.GetString("loglevel")),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)

	//	If we have a config file, report it:
	if viper.ConfigFileUsed() != "" {
		log.Println("[DEBUG] Using config file:", viper.ConfigFileUsed())
	}
}

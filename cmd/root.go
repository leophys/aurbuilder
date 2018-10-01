package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var logLevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "AUR-builder",
	Short: "An helper for building AUR packages",
	Long: `This is a cli helper designed to work inside a docker
    container. The full source code can be found at

    https://github.com/leophys/AUR-builder`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		pacmanList := exec.Command("pacman", "-Q")
		var out bytes.Buffer
		pacmanList.Stdout = &out
		err := pacmanList.Run()
		if err != nil {
			log.WithFields(log.Fields{
				"subcommand": nil,
				"error":      err,
			}).Panic("Execution terminated abruptly")
			// panic(err)
		} else {
			log.Println("%s\n", out.String())
		}
	},
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.AUR-builder.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log", "panic", "Set the log level. Default: panic.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
			log.WithFields(log.Fields{
				"subcommand": nil,
				"err":        err,
			}).Fatal("Home dir not found")
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".AUR-builder" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".AUR-builder")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.WithFields(log.Fields{
			"subcommand":  nil,
			"config file": viper.ConfigFileUsed(),
		}).Info("Using config file")
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	switch logLevel {
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	}
}

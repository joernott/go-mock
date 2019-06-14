package cmd

import (
	"fmt"
	"os"

	"bufio"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/joernott/mock/rules"
	"github.com/joernott/mock/server"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ConfigFile string
var Port int
var Rules string
var LogLevel int
var LogFile string

var rootCmd = &cobra.Command{
	Use:   "mock",
	Short: "Mock is a simple generic REST API mock",
	Long:  `A configurable REST API mock written in go`,
	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		if LogFile == "" {
			log.SetOutput(os.Stdout)
		} else {
			f, err := os.Create(LogFile)
			if err != nil {
				fmt.Println("Could not create logfile '" + LogFile + "'")
				os.Exit(10)
			}
			w := bufio.NewWriter(f)
			log.SetOutput(w)
		}
		switch LogLevel {
		case 0:
			log.SetLevel(log.PanicLevel)
		case 1:
			log.SetLevel(log.FatalLevel)
		case 2:
			log.SetLevel(log.ErrorLevel)
		case 3:
			log.SetLevel(log.WarnLevel)
		case 4:
			log.SetLevel(log.InfoLevel)
		case 5:
			log.SetLevel(log.DebugLevel)
		default:
			log.SetLevel(log.DebugLevel)
		}
		spew.Dump(LogLevel)
		log.WithFields(log.Fields{
			"LogFile":  LogFile,
			"LogLevel": LogLevel,
		}).Debug("Logging configured")

		if ConfigFile != "" {
			log.Debug("Read config from " + ConfigFile)
			viper.SetConfigFile(ConfigFile)
		} else {
			log.Debug("Read config from home directory")
			home, err := homedir.Dir()
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			viper.AddConfigPath(home)
			ex, err := os.Executable()
			if err != nil {
				log.Error(err)
				panic(err)
			}
			pwd := filepath.Dir(ex)
			viper.AddConfigPath(pwd)
			viper.SetConfigName("mock")
		}

		if err := viper.ReadInConfig(); err != nil {
			log.Error("Can't read config" + err.Error())
		}
		log.Debug("PersistentPreRun finished")
	},
	Run: func(cmd *cobra.Command, args []string) {
		Rules, err := rules.LoadRules(viper.GetString("rules"))
		if err != nil {
			log.Error(err)
			os.Exit(2)
		}
		server.Router(viper.GetInt("port"), Rules)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "c", "./mock.json", "config file (default is ./mock.json)")
	rootCmd.PersistentFlags().IntVarP(&Port, "port", "P", 8000, "Network port")
	rootCmd.PersistentFlags().StringVarP(&Rules, "rules", "r", "./rules.json", "rule file")
	rootCmd.PersistentFlags().IntVarP(&LogLevel, "loglevel", "l", 5, "log level")
	rootCmd.PersistentFlags().StringVarP(&LogFile, "logfile", "L", "", "logfile, when empty, log to stdout")

	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("rules", rootCmd.PersistentFlags().Lookup("rules"))
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))

	viper.SetDefault("port", 8000)
	viper.SetDefault("rules", "./rules.json")
	viper.SetDefault("loglevel", 5)
	viper.SetDefault("logfile", "")
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var Settings *Config
var HomePath string
var ConfigFilePath string
var ConfigFileName string = ".dockenv"
var ConfigFileType string = "yaml"

type Service struct {
	Name                 string `mapstructure:"name"`
	DockerComposeFile    string `mapstructure:"file"`
	DockerComposeService string `mapstructure:"service"`
	Port                 int    `mapstructure:"port"`
}

type Config struct {
	Domain   string    `mapstructure:"domain"`
	Services []Service `mapstructure:"services"`
}

func initConfig() {
	home, homeError := os.UserHomeDir()

	if homeError != nil {
		fmt.Fprintln(os.Stderr, "Undefined user home path")
		return
	}

	HomePath = home
	ConfigFilePath = HomePath + "/" + ConfigFileName + "." + ConfigFileType

	createConfigFile()
	Settings = loadSettings()
	createInitialSettings()
}

func loadSettings() (config *Config) {
	viper.AddConfigPath(HomePath)
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Print("Error reading env file", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Print(err)
	}

	return
}

func createConfigFile() {
	_, configFileError := os.Stat(ConfigFilePath)

	if os.IsNotExist(configFileError) {
		configFile, configFileError := os.Create(ConfigFilePath)

		if configFileError != nil {
			fmt.Fprintln(os.Stderr, "Unable to create config file")
		}

		defer configFile.Close()
	}
}

func createInitialSettings() {
	if Settings.Domain == "" {
		Settings.Domain = "local"

		viper.Set("domain", Settings.Domain)
		viper.WriteConfig()
	}
}

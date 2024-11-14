package config

import (
	"fmt"

	"github.com/go-importexportExcelCRUD/models"
	"os"

	"github.com/spf13/viper"
)

var Config *models.Configurations

func GetConfig() *models.Configurations {
	if Config != nil {
		return Config
	}

	viper.SetConfigName("config") // name of config file (without extension)

	if os.Getenv("GO_ENV") == "local" {
		viper.AddConfigPath("./config") // path to look for the config file in
	} else {
		viper.AddConfigPath("./config") // path to look for the config file from machine
	}

	viper.SetConfigType("json")

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	var configurations models.Configurations
	err := viper.Unmarshal(&configurations)
	fmt.Println(configurations)
	if err != nil {
		panic(fmt.Errorf("Fatal error while decoding config file: %s \n", err))
	}

	validate := validator.New()
	if err := validate.Struct(&configurations); err != nil {
		panic(fmt.Errorf("Fatal error while validating config file: %s \n", err))
	}

	Config = &configurations
	return Config
}

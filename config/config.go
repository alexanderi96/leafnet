package config

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/go-yaml/yaml"
)

var (
	conf_file_path string = ".config/leafnet/config.yaml"

	Config map[string]string
)

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	homeDir := usr.HomeDir

	conf_file_path = fmt.Sprintf("%s/%s", homeDir, conf_file_path)
	if !checkPathExists(conf_file_path) {
		log.Fatal("Config file not found: ", conf_file_path)
	}

	//read leafnet config file from ~/.config/leafnet/config.yaml and put the variables in the config map
	Config, err = readConfig(conf_file_path)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Unable to read the config file: %s", err))
	}
	log.Println("Config file loaded")
}

func checkPathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// the readConfig function, it reads the yaml file and puts the variables in the config map
func readConfig(path string) (map[string]string, error) {

	config := make(map[string]string)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// read the yaml file
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

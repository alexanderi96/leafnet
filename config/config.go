package config

import (
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

var (
	Config map[string]string
)

// the readConfig function, it reads the yaml file and puts the variables in the config map
func ReadConfig(path string) error {
	log.Println("Reading config file: ", path)
	Config = make(map[string]string)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// read the yaml file
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		return err
	}

	return nil
}

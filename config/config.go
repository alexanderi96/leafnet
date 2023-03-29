package config

import (
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

var (
	Config map[string]string
)

// the readConfig function, it reads the toml file and puts the variables in the config map
func ReadConfig(path string) error {
	log.Println("Reading config file: ", path)
	Config = make(map[string]string)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// read the toml file
	decoder := toml.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		return err
	}

	return nil
}

package config

import (
	"log"
	"os"
	"time"
)

type Updater func(configPath string) error

var readConfig = func(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	return nil
}

func InfiniteUpdater(configPath string) error {
	for {
		err := readConfig(configPath)
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
}

func UpdateConfig(configPath string, updater Updater) error {
	return updater(configPath)
}

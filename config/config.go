package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/AmogusAzul/AutoShutdown/utils"
	"github.com/fsnotify/fsnotify"
)

const (
	ConfigFile = "config.json"
)

type Config struct {
	Path string

	Target       utils.Hour
	AlertAdvance time.Duration
	Postpone     time.Duration

	KillChan   chan bool
	UpdateChan chan bool
}

type StringConfigs struct {
	ShutdownTime     string `json:"ShutdownTime"`
	AlertAdvantage   string `json:"AlertAdvantage"`
	PostponeDuration string `json:"PostponeDuration"`
}

func CreateConfigReader() *Config {
	path, pathErr := filepath.Abs(ConfigFile)
	if pathErr != nil {
		log.Fatalf("Couldn't find %s (which is obligatory for working)", path)
		os.Create(ConfigFile)
		path, pathErr = filepath.Abs(ConfigFile)
	}

	config := &Config{
		Path:       path,
		UpdateChan: make(chan bool),
		KillChan:   make(chan bool),
	}

	config.updateFromJson()

	return config

}

func (config *Config) Watch() {

	go func() {

		var killed bool

		// creates a new file watcher
		watcher, err := fsnotify.NewWatcher()

		if err != nil {
			fmt.Println("ERROR", err)
		}

		if err := watcher.Add(config.Path); err != nil {
			fmt.Println("ERROR", err)
		}

		for {

			if killed {
				break
			}

			select {

			case <-watcher.Events:
				time.Sleep(50 * time.Millisecond)
				config.updateFromJson()

				go func() { config.UpdateChan <- true }()

			case watcherErr := <-watcher.Errors:
				log.Fatalf("Error at file watcher: %s", watcherErr)

			case killed = <-config.KillChan:
			}

		}

	}()

}

func (config *Config) updateFromJson() {

	// Read File as byte slice
	data, errRead := os.ReadFile(config.Path)
	if errRead != nil {
		log.Panicf("Couldn't read file %s\n", config.Path)
		return
	}

	// Parsing to custom struct
	fields := StringConfigs{}

	errParse := json.Unmarshal(data, &fields)
	if errParse != nil {

		log.Panicf("Couldn't get json data. Aborting config change.\n")
		return
	}

	target, errTarget := utils.ParseHour(fields.ShutdownTime)
	advantage, errAdvantage := time.ParseDuration(fields.AlertAdvantage)
	postponeT, errPostpone := time.ParseDuration(fields.PostponeDuration)

	if errTarget != nil {
		log.Panicf("Couldn't get errTarget json data field. Aborting config change.%v\n", target)
		return
	}
	if errAdvantage != nil {
		log.Panicf("Couldn't get errAdvantage json data field. Aborting config change.\n")
		return
	}
	if errPostpone != nil {
		log.Panicf("Couldn't get errAdvantage json data field. Aborting config change.\n")
		return
	}

	config.Target = target
	config.AlertAdvance = advantage
	config.Postpone = postponeT

	config.UpdateChan <- true

}

package main

import (
	"fmt"

	"github.com/AmogusAzul/AutoShutdown/alarmer"
	"github.com/AmogusAzul/AutoShutdown/config"
)

func main() {

	config := config.CreateConfigReader()
	config.Watch()

	alarmer := alarmer.GetAlarmer(func() { fmt.Println("TRIGGERED") }, func() { fmt.Printf("WILL TRIGGER IN %d seconds\n", config.AlertAdvance/1000000000) }, config)

	alarmer.Activate()

	fmt.Scanln()

}

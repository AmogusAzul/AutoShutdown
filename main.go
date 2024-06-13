package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/AmogusAzul/AutoShutdown/alarmer"
	"github.com/AmogusAzul/AutoShutdown/config"
)

func main() {

	config := config.CreateConfigReader()
	config.Watch()

	targetFunc := func() {
		if runtime.GOOS != "windows" {
			fmt.Println("Not In Windows!")
			return
		}

		if err := exec.Command("cmd", "/C", "shutdown", "/p").Run(); err != nil {
			fmt.Println("Failed to initiate shutdown:", err)
		}
	}
	alertFunc := func() {
		fmt.Printf("WILL TRIGGER IN %d seconds\n",
			config.AlertAdvance/time.Second)
	}

	alarmer := alarmer.GetAlarmer(
		&targetFunc,
		&alertFunc,
		config)

	alarmer.Activate()

	fmt.Println("Waiting")

	fmt.Scanln()

}

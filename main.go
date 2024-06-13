package main

import (
	"fmt"
	"os/exec"
	"runtime"

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
		fmt.Printf("WILL TRIGGER IN %v\n",
			config.AlertAdvance)
	}

	alarmer := alarmer.GetAlarmer(
		&targetFunc,
		&alertFunc,
		config)

	alarmer.Activate()

	fmt.Println("Waiting")

	response := ""
	killed := false

	for !killed {
		fmt.Scanln(&response)

		fmt.Println("response:", response)

		switch response {

		case "k":
			killed = true
		case "p":
			alarmer.Postpone()
		}
	}

}

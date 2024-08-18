package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

const (
	markName  = "Golang-cli-project"
	markValue = "1"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <hh:mm> <Text Message>\n", os.Args[0])
		os.Exit(1)
	}

	now := time.Now()
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	t, err := w.Parse(os.Args[1], now)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		os.Exit(2)
	}

	if t == nil {
		fmt.Println("Unable to parse time")
		os.Exit(2)
	}

	if now.After(t.Time) {
		fmt.Println("Please set a future time")
		os.Exit(3)
	}

	if os.Getenv(markName) == markValue {
		// This is the reminder execution block
		err = beeep.Alert("Reminder", strings.Join(os.Args[2:], " "), "assets/information.png")
		if err != nil {
			fmt.Println("Error displaying reminder:", err)
			os.Exit(4)
		}
		os.Exit(0)
	} else {
		// Calculate the time difference and schedule the reminder
		diff := t.Time.Sub(now)
		fmt.Println("Reminder will be displayed after", diff.Round(time.Second))
		time.Sleep(diff)

		// Re-run the program with the environment variable set
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", markName, markValue))
		if err := cmd.Start(); err != nil {
			fmt.Println("Error restarting process:", err)
			os.Exit(5)
		}
		os.Exit(0)
	}
}

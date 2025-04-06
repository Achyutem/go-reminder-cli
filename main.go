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
		fmt.Printf("Usage: %s <time expression> <Text Message>\nExample: %s \"in 10 minutes\" Take a break\n", os.Args[0], os.Args[0])
		os.Exit(1)
	}

	now := time.Now()
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	t, err := w.Parse(os.Args[1], now)
	if err != nil || t == nil {
		fmt.Println("Could not parse time. Try something like 'in 10 minutes' or '10:30am'")
		os.Exit(2)
	}

	if now.After(t.Time) {
		fmt.Println("Time must be in the future.")
		os.Exit(3)
	}

	message := strings.Join(os.Args[2:], " ")

	if os.Getenv(markName) == markValue {
		err = beeep.Notify("Reminder", message, "")
		if err != nil {
			fmt.Println("Reminder:", message)
		}
		os.Exit(0)
	} else {
		diff := t.Time.Sub(now)
		fmt.Println("Reminder scheduled in", diff.Round(time.Second))

		time.Sleep(diff)

		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", markName, markValue))
		if err := cmd.Start(); err != nil {
			fmt.Println("Error launching reminder:", err)
			os.Exit(5)
		}
		os.Exit(0)
	}
}

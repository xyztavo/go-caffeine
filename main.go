package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
)

var (
	duration = flag.Duration("t", 0, "Duration to prevent sleep (e.g. 1h, 30m, 5h). 0 means indefinitely")
	version  = "1.0.2"
)

func main() {
	flag.Parse()

	// Print version if requested
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("go-caffeine version %s ✨\n", version)
		return
	}

	fmt.Printf("Starting go-caffeine ☕ (Press Ctrl+C to exit)\n")
	if *duration > 0 {
		fmt.Printf("System will stay awake for %v ⏰\n", *duration)
	} else {
		fmt.Println("System will stay awake indefinitely 🔋")
	}

	// Create and configure spinner
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = " Keeping system awake..."
	s.Start()

	// Set up channel to handle interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create timer channel if duration is specified
	var timer *time.Timer
	var startTime time.Time
	if *duration > 0 {
		timer = time.NewTimer(*duration)
		startTime = time.Now()
	}

	// Start the keep-awake loop based on OS
	keepAwake := getKeepAwakeFunc()
	ticker := time.NewTicker(time.Second * 30)
	updateTicker := time.NewTicker(time.Second)
	defer ticker.Stop()
	defer updateTicker.Stop()
	defer s.Stop()

	for {
		if *duration > 0 {
			select {
			case <-sigChan:
				fmt.Println("\n👋 Exiting go-caffeine...")
				return
			case <-ticker.C:
				keepAwake()
			case <-updateTicker.C:
				remaining := *duration - time.Since(startTime)
				if remaining > 0 {
					s.Suffix = fmt.Sprintf(" Keeping system awake... %v remaining", remaining.Round(time.Second))
				}
			case <-timer.C:
				fmt.Println("\n⌛ Duration expired, exiting go-caffeine...")
				return
			}
		} else {
			select {
			case <-sigChan:
				fmt.Println("\n👋 Exiting go-caffeine...")
				return
			case <-ticker.C:
				keepAwake()
			}
		}
	}
}

func getKeepAwakeFunc() func() {
	switch runtime.GOOS {
	case "darwin":
		return func() {
			exec.Command("caffeinate", "-i", "-t", "60").Start()
		}
	case "windows":
		return func() {
			exec.Command("powershell", "-Command", "Add-Type -TypeDefinition '@using System; using System.Runtime.InteropServices; public class Sleep { [DllImport(\"kernel32.dll\")] public static extern uint SetThreadExecutionState(uint esFlags); }'; [Sleep]::SetThreadExecutionState([uint32]\"0x80000003\")").Start()
		}
	case "linux":
		return func() {
			exec.Command("xdg-screensaver", "reset").Start()
		}
	default:
		return func() {
			fmt.Printf("Unsupported operating system: %s\n", runtime.GOOS)
			os.Exit(1)
		}
	}
}

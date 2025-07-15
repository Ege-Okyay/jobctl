package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/Ege-Okyay/jobctl/internal/cli"
)

// clearScreen clears the terminal screen.
func clearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Print("\033[H\033[2J")
	}
}

// LaunchInteractiveShell starts an interactive shell for managing jobs.
func LaunchInteractiveShell() {
	// Gracefully handle Ctrl+C to exit the shell.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("\nReceived interrupt - shutting down.")
		os.Exit(0)
	}()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Interactive mode - type 'help' for commands, 'exit' to quit.")

	for {
		fmt.Print("jobctl> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nGoodybe!")
				return
			}

			fmt.Fprintln(os.Stderr, "read error:", err)
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch line {
		case "exit", "quit":
			fmt.Println("Goodbye!")
			return
		case "clear", "cls":
			clearScreen()
			continue
		}

		// To execute a command, we replace os.Args with the new command and call the CLI setup again.
		parts := strings.Fields(line)
		os.Args = append([]string{"jobctl"}, parts...)

		cli.Setup()
	}
}

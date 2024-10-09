package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

var variables = make(map[string]string)

func main() {
	exitCMDMaps := map[string]struct{}{
		"exit":   {},
		"quit":   {},
		"exit()": {},
		"quit()": {},
		".quit":  {},
		".exit":  {},
	}
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Println(blue("IShell 1.0.0 (main, May 29 2024, 09:33:50) [xcli 1.0.0]"))
	fmt.Println("Type 'copyright', 'credits' or 'license' for more information")
	fmt.Println("MyInteractiveShell 1.0.0 -- An enhanced Interactive Shell. Type '?' for help.")
	fmt.Println()

	// Initialize readline
	rl, err := readline.New("")
	if err != nil {
		fmt.Println("Error initializing readline:", err)
		os.Exit(1)
	}
	defer rl.Close()

	index := 1
	for {
		tipmsg := fmt.Sprintf("In [%d]: ", index)
		rl.SetPrompt(green(tipmsg))
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				fmt.Println("Bye!")
				break
			}
			fmt.Println("Readline error:", err)
			break
		}

		input := strings.TrimSpace(line)
		if _, ok := exitCMDMaps[input]; ok {
			fmt.Println("Bye!")
			break
		}

		handleCommand(input, yellow, red, cyan)
		fmt.Println()
		index++
	}
}

func executeCommand(cmd string) (string, error) {
	// Execute the command using sh -c
	command := exec.Command("sh", "-c", cmd)
	output, err := command.CombinedOutput() // Get output
	return string(output), err
}

func handleCommand(command string, yellow, red, cyan func(a ...interface{}) string) {
	// Check for assignment syntax
	if strings.Contains(command, "=") {
		parts := strings.SplitN(command, "=", 2)
		varName := strings.TrimSpace(parts[0])
		cmd := strings.TrimSpace(parts[1])

		if output, err := executeCommand(cmd); err == nil {
			variables[varName] = output
			// Do not print the output here
		} else {
			fmt.Println(red("Error executing command:", err))
		}
	} else if strings.HasPrefix(command, "print(") && strings.HasSuffix(command, ")") {
		varName := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(command, "print("), ")"))
		if value, exists := variables[varName]; exists {
			fmt.Println(value)
		} else {
			fmt.Println(red("Variable not found:", varName))
		}
	} else if value, exists := variables[command]; exists {
		// Print the value if the input is a variable name
		fmt.Println(value)
	} else {
		switch command {
		case "hello":
			fmt.Println(yellow("Hello world!"))
		case "help":
			fmt.Println(cyan("CMD: hello, help, exit, print(<variable>)"))
		default:
			fmt.Printf(red("Unknown: %s\n"), command)
		}
	}
}

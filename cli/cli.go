package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CLI struct {
	prompt   string
	commands map[string]func([]string)
}

func NewCLI(prompt string) *CLI {
	return &CLI{
		prompt:   prompt,
		commands: make(map[string]func([]string)),
	}
}

func (c *CLI) AddCommand(name string, handler func([]string)) {
	c.commands[name] = handler
}

func (c *CLI) Start() {
	fmt.Println("Type 'help' for a list of available commands.")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(c.prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		inputParts := strings.Fields(input)

		if len(inputParts) == 0 {
			continue
		}

		command := strings.ToLower(inputParts[0])
		args := inputParts[1:]

		if command == "exit" {
			if c.confirmExit() {
				fmt.Println("Terminating the node...")
				os.Exit(0)
				return
			}
		} else {
			handler, exists := c.commands[command]
			if !exists {
				fmt.Println("Unknown command. Type 'help' for available commands.")
				continue
			}
			handler(args)
		}
	}
}

func (c *CLI) confirmExit() bool {
	fmt.Print("Are you sure you want to exit (y/n)? ")
	input := strings.ToLower(c.readInput())
	return input == "y"
}

func (c *CLI) readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

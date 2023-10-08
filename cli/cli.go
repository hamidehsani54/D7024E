package cli

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// CLI is a struct representing a command-line interface.
type CLI struct {
    prompt   string                        // The command prompt displayed to the user.
    commands map[string]func([]string)    // A map of command names to their handler functions.
}

// NewCLI is a constructor function for creating a new CLI instance.
func NewCLI(prompt string) *CLI {
    return &CLI{
        prompt:   prompt,
        commands: make(map[string]func([]string)),
    }
}

// AddCommand associates a command name with its handler function.
func (c *CLI) AddCommand(name string, handler func([]string)) {
    c.commands[name] = handler
}

// Start initiates the CLI and enters a loop to process user input.
func (c *CLI) Start() {
    fmt.Println("Type 'help' for a list of available commands.")
    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Print(c.prompt)                 // Display the command prompt to the user.
        input, _ := reader.ReadString('\n')  // Read user input.
        input = strings.TrimSpace(input)     // Remove leading/trailing whitespace.
        inputParts := strings.Fields(input)  // Split input into words.

        if len(inputParts) == 0 {
            continue
        }

        command := strings.ToLower(inputParts[0])  // Extract the first word as the command.
        args := inputParts[1:]                      // Extract the remaining words as arguments.

        if command == "exit" {
            if c.confirmExit() {  // Check if the user wants to exit and confirm.
                fmt.Println("Terminating the node...")
                os.Exit(0)
                return
            }
        } else {
            handler, exists := c.commands[command]  // Look up the handler for the command.
            if !exists {
                fmt.Println("Unknown command. Type 'help' for available commands.")
                continue
            }
            handler(args)  // Execute the handler function with arguments.
        }
    }
}

// confirmExit prompts the user to confirm if they want to exit the CLI.
func (c *CLI) confirmExit() bool {
    fmt.Print("Are you sure you want to exit (y/n)? ")
    input := strings.ToLower(c.readInput())  // Read and convert user input to lowercase.
    return input == "y"  // Return true if the user confirmed with 'y'.
}

// readInput reads and returns user input from the standard input.
func (c *CLI) readInput() string {
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')
    return strings.TrimSpace(input)  // Remove leading/trailing whitespace and return the input.
}
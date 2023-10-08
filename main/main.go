package main

import (
	"D7024E/logic"
	"D7024E/cli"
	"fmt"
	"net"
	"os"
	"time"
	"strings"
)

func main() {
	// Get the role from environment variables
	role := os.Getenv("master")

	// Retrieve the  IP
	ip, err := getContainerIP()
	if err != nil {
		panic(err)
	}
	fmt.Println("ip:", ip)

	// Generate the KademliaID based on role
	var nodeID *logic.KademliaID
	if role == "true" {
		fmt.Println("I am master")
		nodeID = logic.NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51")
	} else {
		fmt.Println("I am NOT master")
		nodeID = logic.NewRandomKademliaID()
	}
	fmt.Println("Node ID:", nodeID)
	port := 4000

	// Initialize the network with the node's ID and IP
	netInstance := logic.InitNetwork(nodeID, ip)

	// Start listening for incoming messages
	go netInstance.Listen(ip, port)

	if role != "true" {
		contact := logic.NewContact(logic.NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51"), "masterNode")
		netInstance.SendPingMessage(&contact)

		netInstance.Kademlia.JoinNetwork()
		time.Sleep(60 * time.Second)
		netInstance.Kademlia.Store([]byte("test"))
		time.Sleep(20 * time.Second)
		fmt.Println("QQQCXXXX", netInstance.Kademlia.PrintData())
	}

	exitCh := make(chan struct{})
	go startCLI(exitCh, netInstance)

	for {
		time.Sleep(time.Hour)
	}
}

// Function to get the container IP address
func getContainerIP() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	addr, err := net.LookupIP(hostname)
	if err != nil {
		return "", err
	}

	if len(addr) == 0 {
		return "", err
	}

	return addr[0].String(), nil
}

func startCLI(exitCh chan struct{}, network *logic.Network) {
    myCLI := cli.NewCLI("> ")

	myCLI.AddCommand("ping", func(args []string) {
        // Handle 'put' command
        fmt.Println("Executing 'ping' command with args:", args)
		contact := logic.NewContact(logic.NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51"), "masterNode")
		network.SendPingMessage(&contact)
    })

    myCLI.AddCommand("put", func(args []string) {
		// Handle 'put' command
		fmt.Println("Executing 'put' command with args:", args)
	
		if len(args) < 1 {
			fmt.Println("No arguments provided for 'put' command!")
			return
		}
	
		// Convert the first argument to a byte slice and pass to Store
		dataToStore := strings.Join(args, " ")  // Join with spaces
		network.Kademlia.Store([]byte(dataToStore))
	})
	
    myCLI.AddCommand("get", func(args []string) {
        // Handle 'get' command
        fmt.Println("Executing 'get' command with args:", args)
    })

	myCLI.AddCommand("print", func(args []string) {
        // Handle 'get' command
        fmt.Println("Executing 'print' command with args:", args)
		fmt.Println(network.Kademlia.PrintData())
    })

    myCLI.AddCommand("help", func(args []string) {
        // Handle 'help' command
        fmt.Println("Available commands: put, get, exit")
    })

    // Add an exit command to gracefully exit the CLI and terminate the node
    myCLI.AddCommand("exit", func(args []string) {
        fmt.Println("Terminating the node...")
        // Send a signal to the network Goroutine to exit
        close(exitCh)
    })

    myCLI.Start()
}
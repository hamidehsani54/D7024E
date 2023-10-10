package main

import (
	"D7024E/cli"
	"D7024E/logic"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
	//"crypto/sha1"
	//"encoding/hex"
)

func main() {
	role := os.Getenv("master")

	ip, err := getContainerIP()
	if err != nil {
		panic(err)
	}
	fmt.Println("ip:", ip)

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

	netInstance := logic.InitNetwork(nodeID, ip)
	go netInstance.Listen(ip, port)

	if role != "true" {
		contact := logic.NewContact(logic.NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51"), "masterNode")
		netInstance.SendPingMessage(&contact)

		netInstance.Kademlia.JoinNetwork()
		time.Sleep(60 * time.Second)
	}

	exitCh := make(chan struct{})
	go startCLI(exitCh, netInstance)

	for {
		time.Sleep(time.Hour)
	}
}

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

		fmt.Println("Executing 'ping' command with args:", args)
		contact := logic.NewContact(logic.NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51"), "masterNode")
		network.SendPingMessage(&contact)
	})

	myCLI.AddCommand("put", func(args []string) {
		fmt.Println("Executing 'put' command with args:", args)
		if len(args) < 1 {
			fmt.Println("No arguments provided for 'put' command!")
			return
		}

		dataToStore := strings.Join(args, " ")
		fmt.Println("The hash: ", network.Kademlia.Store([]byte(dataToStore)))
	})

	myCLI.AddCommand("get", func(args []string) {
		fmt.Println("Executing 'get' command with args:", args)
		data := strings.Join(args, " ")
		var t1, t2 string
		var contacts []logic.Contact
		t1, t2, contacts = network.Kademlia.LookupData(data)
		if contacts != nil {
			fmt.Println("Contacts:", contacts)
		}
		fmt.Println("The data: ", t1, " From Node: ", t2)

	})

	myCLI.AddCommand("print", func(args []string) {
		fmt.Println("Executing 'print' command with args:", args)
		fmt.Println(network.Kademlia.PrintData())
	})

	myCLI.AddCommand("help", func(args []string) {
		fmt.Println("Available commands: put, get, exit")
	})

	myCLI.Start()
}

package main

import (
    "D7024E/logic"
    "fmt"
    "net"
    "os"
    "time"
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

    //kademliaInstance := logic.InitKademlia(netInstance)
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
    // If the node is not a master, it should join the network
    

    if role != "true" {
        // Assuming a function to initialize Kademlia
        //time.Sleep(60 * time.Second)
        contact := logic.NewContact(logic.NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51"), "masterNode")
        netInstance.SendPingMessage(&contact)
        //time.Sleep(60 * time.Second)
		kademliaInstance := logic.InitKademlia(netInstance)
        kademliaInstance.JoinNetwork()
    }

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
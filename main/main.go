package main

import (
    "D7024E/logic"
    "os"
    "time"
    //"fmt"
    "net"
)

func main() {
    // Get the role from environment variables
    role := os.Getenv("master")

    // Retrieve the container's IP
    ip, err := getContainerIP()
    if err != nil {
        panic(err)  // Handle this error gracefully in a production setting
    }

    // Generate the KademliaID based on role (master or regular node)
    var nodeID *logic.KademliaID
    if role == "true" {
        nodeID = logic.NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51")
    } else {
        nodeID = logic.NewRandomKademliaID()
    }

    port := 4000

    // Initialize the network with the node's ID and IP
    netInstance := logic.InitNetwork(nodeID, ip)

    // Start listening for incoming messages
    go netInstance.Listen(ip, port)

    // If the node is not a master, it should join the network
    if role != "true" {
        kademliaInstance := logic.InitKademlia(netInstance) // Assuming a function to initialize Kademlia
        kademliaInstance.JoinNetwork()
    }

    // For demonstration, continuously send pings to the master node
    /*targetIP := "10.10.0.10"  // Assuming this is the master node or bootstrap node
    message := logic.PingMessage{
        Type:       "ping",
        KademliaID: nodeID.String(),
        IP:         ip,
    }*/

    for {
        //logic.SendDial(targetIP, &message)
        time.Sleep(1 * time.Second)
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
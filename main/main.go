package main

import (
    "D7024E/logic"
    "os"
    "time"
    "fmt"
    "net"
)

func main() {
    // Get the role from environment variables
    role := os.Getenv("master")

    // Retrieve the  IP
    ip, err := getContainerIP()
    if err != nil {
        panic(err) 
    }

    // Generate the KademliaID based on role
    var nodeID *logic.KademliaID
    if role == "true" {
        fmt.Println("I am master")
        nodeID = logic.NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51")
    } else {
        fmt.Println("I am NOT master")
        nodeID = logic.NewRandomKademliaID()
    }

    port := 4000

    // Initialize the network with the node's ID and IP
    netInstance := logic.InitNetwork(nodeID, ip)

    // Start listening for incoming messages
    go netInstance.Listen(ip, port)
    
    // If the node is not a master, it should join the network
    /*
    if role != "true" {
         // Assuming a function to initialize Kademlia
        
        kademliaInstance.JoinNetwork()
    }
    */

    // For demonstration, continuously send pings to the master node
    /*targetIP := "10.10.0.10"  // Assuming this is the master node or bootstrap node
    message := logic.PingMessage{
        Type:       "ping",
        KademliaID: nodeID.String(),
        IP:         ip,
    }*/
    
    contact := logic.NewContact(logic.NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51"), "10.10.0.10")
    for {
        //fmt.Println(contact.String())
        netInstance.SendPingMessage(&contact)
        time.Sleep(1 * time.Second)
        //fmt.Println("Yoo my man", netInstance.SendFindContactMessage(&contact, logic.NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51")))
        
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
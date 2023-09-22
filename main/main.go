package main

import (
    "D7024E/logic"
    "os"
    "time"
    "net"
)

func main() {
    // Generate a random KademliaID for the node
    nodeID := logic.NewRandomKademliaID()

    // Get the container IP address
    ip, err := getContainerIP()
    if err != nil {
        panic(err)  // Handle this error gracefully in a production setting
    }

    port := 4000
    targetIP := "10.10.0.10"  // Assuming this is the master node or bootstrap node

    message := logic.PingMessage{
        Type:       "ping",
        KademliaID: nodeID.String(),
        IP:         ip,
    }

    netInstance := logic.InitNetwork(nodeID, ip)

    // Now use the instance to call Listen
    go netInstance.Listen(ip, port)

    for {
        logic.SendDial(targetIP, &message)
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
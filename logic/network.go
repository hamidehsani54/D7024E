package logic

import (
    "net"
    "fmt"
    "log"
    //"strings"
    "encoding/json"
    "time"
)

type Network struct {
    Node *Contact
	rt *RoutingTable
}

func InitNetwork(id *KademliaID, address string) *Network {
    node := &Contact{
        ID:      id,
        Address: address,
    }
    
    return &Network{
        Node: node,
        rt: NewRoutingTable(*node),
    }
}


type Message struct {
    Type       string `json:"type"`
    KademliaID string `json:"kademliaID"`
    IP         string `json:"ip"`
    target     *KademliaID
    Contacts    []Contact
}

func (network *Network) Listen(ip string, port int) {
    // TODO
    addr := fmt.Sprintf("%s:%d", ip, port)
    udpAddr, err := net.ResolveUDPAddr("udp", addr)
    listener, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()
    fmt.Printf("Listening on %s...\n", udpAddr.String())

    buffer := make([]byte, 4096)
    for{
        n, addr, err := listener.ReadFromUDP(buffer)
        if err != nil{
            log.Print(err)
            continue
        }
        receivedData := string(buffer[:n])
        handleInput(buffer[:n],addr, network)
        log.Printf("Received data from %s: %s\n", addr, receivedData)
    }
}

func SendDial(targetIp string, message *Message) Message {
    err_message := Message{
        Type:       "error",
    }
    addr := fmt.Sprintf("%s:%d", targetIp, 4000)
    conn, err := net.Dial("udp", addr)
    if err != nil{
        log.Fatal(err)
    }
    defer conn.Close()

    msgBytes, err := json.Marshal(message)
    if err != nil {
        log.Println("Error marshalling message:", err)
        return err_message
    }
    _, err = conn.Write(msgBytes)
    if err != nil {
        fmt.Println("Writing error: ", err, " in sendMsg")
        return err_message
    }

    buffer := make([]byte, 4096)
	setReadDeadlineError := conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	if setReadDeadlineError != nil {
		return err_message
	}
	end, readError := conn.Read(buffer)
	if readError != nil {
		return err_message
	}

    var pingMsg Message
    err2 := json.Unmarshal(buffer[:end], &pingMsg)
    if err2 != nil {
        log.Println("Error unmarshalling received message:", err)
        return err_message
    }
	return pingMsg
}

func (network *Network) SendPingMessage(contact *Contact) {
    message := Message{
        Type:       "ping",
        KademliaID: network.Node.ID.String(),
        IP:         network.Node.Address,
    }
    
    SendDial(contact.Address ,&message)
}



func (network *Network) SendFindContactMessage(contact *Contact, t *KademliaID) Message {
    fmt.Println("printing contacts", contact.Address)
    message := Message{
        Type:       "FindContact",
        KademliaID: network.Node.ID.String(),
        IP:         network.Node.Address,
        target:     t,
    }
    
    return SendDial(contact.Address ,&message)
}

func (network *Network) SendFindDataMessage(hash string) {
    // TODO
}

func (network *Network) SendStoreMessage(data []byte) {
    // TODO
}

func handleInput(message []byte, addr *net.UDPAddr, network *Network){
    var pingMsg Message
    err := json.Unmarshal(message, &pingMsg)
    if err != nil {
        log.Println("Error unmarshalling received message:", err)
        return
    }
    // Now you can use pingMsg.Type, pingMsg.KademliaID, etc.
  
    switch{
        case pingMsg.Type == "ping":
            handlePing(pingMsg, addr, network)
        case pingMsg.Type == "pong":
            handlePong(pingMsg)
        case pingMsg.Type == "FindContact":
            fmt.Println("Tst: ")
            handleFindContact(pingMsg, network)
        case pingMsg.Type == "err":
            fmt.Println("Error in error")
    }
}

func handlePing(message Message, addr *net.UDPAddr, network *Network) Message{
    //Respnse to sender
    fmt.Println("Ping from: ", message)

	senderContact := NewContact(NewKademliaID(message.KademliaID), message.IP)
    network.rt.AddContact(senderContact) 

	sendMsg := &Message{
        Type: "pong",
        KademliaID: network.Node.ID.String(),
        IP: network.Node.Address,//
    }
    return SendDial(message.IP, sendMsg)
}


func handlePong(message Message){
    fmt.Println("Pong from: ", message)
}


func handleFindContact(message Message, network *Network) Message{
    var K = 10
    contacts := network.rt.FindClosestContacts(message.target, K)

    sendContacts:= &Message{
        Type:       "ReturnContact",
        KademliaID: network.Node.ID.String(),
        IP:         network.Node.Address,
        Contacts: contacts,
    }

    return SendDial(message.IP, sendContacts)
}
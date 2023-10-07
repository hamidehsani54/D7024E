package logic

import (
	"fmt"
	"log"
	"net"

	//"strings"
	"encoding/json"
	"time"
)

type Network struct {
	Node *Contact
	rt   *RoutingTable
	Kademlia *Kademlia
}

func InitNetwork(id *KademliaID, address string) *Network {
	node := &Contact{
		ID:      id,
		Address: address,
	}

	net :=  &Network{
		Node: node,
		rt:   NewRoutingTable(*node),
	}

	net.Kademlia = &Kademlia{
        Network: net,
        // Initialize other necessary fields for Kademlia, e.g., DataList
        DataList: make([]DataStore, 0),
    }
	
	return net
}

type Message struct {
	Type       string `json:"type"`
	KademliaID string `json:"kademliaID"`
	IP         string `json:"ip"`
	Target     string `json:"target"`
	Contacts   []Contact
	Data 	   []byte
	Key		   string `json:"Key"`
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
	for {
		n, addr, err := listener.ReadFromUDP(buffer)
		if err != nil {
			log.Print(err)
			continue
		}
		receivedData := string(buffer[:n])
		network.handleInput(buffer[:n], addr)
		log.Printf("Received data from %s: %s\n", addr, receivedData)
	}
}

func SendDial(targetIp string, message *Message) Message {
	err_message := Message{
		Type: "error",
	}
	addr := fmt.Sprintf("%s:%d", targetIp, 4000)
	conn, err := net.Dial("udp", addr)
	if err != nil {
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

func (network *Network) SendPingMessage(contact *Contact) Message{
	fmt.Println("SendPingMessage")
	message := Message{
		Type:       "ping",
		KademliaID: network.Node.ID.String(),
		IP:         network.Node.Address,
	}

	return SendDial(contact.Address, &message)
}

func (network *Network) SendFindContactMessage(contact *Contact, t string) Message {
	fmt.Println("SendFindContactMessage")
	message := Message{
		Type:       "FindContact",
		KademliaID: network.Node.ID.String(),
		IP:         network.Node.Address,
		Target:     t,
	}

	fmt.Println("this is the address: ", contact.Address)
	
	return SendDial(contact.Address, &message)
}

func (network *Network) SendFindDataMessage(contact *Contact, hash string) Message{
	fmt.Println("SendFindDataMessage")
	
	sendMsg := &Message{
		Type:       "SendFindMessage",
		KademliaID: network.Node.ID.String(),
		IP:         network.Node.Address, 
		Key: 		hash,
	}

	return SendDial(contact.Address, sendMsg)
}

func (network *Network) handleFindDataMessage(message Message) Message {
	fmt.Println("handleStoreMessage")
	hash, data := network.Kademlia.FindLocalData(message.Key)
	// TODO
	//Send data to target
	senderContact := NewContact(NewKademliaID(message.KademliaID), message.IP)
	network.rt.AddContact(senderContact)
	var sendData []byte
	if message.Key == hash {
		sendData = data
	}else {
		sendData = nil
	}
	sendMsg := &Message{
		Type:       "SendDataMessage",
		KademliaID: network.Node.ID.String(),
		IP:         network.Node.Address, 
		Data: 		sendData,
	}

	return SendDial(message.IP, sendMsg)
}

func (network *Network) SendStoreMessage(contact *Contact, data []byte) Message{
	fmt.Println("SendStoreMessage: ", string(data))
	
	sendMsg := &Message{
		Type:       "StoreMessage",
		KademliaID: network.Node.ID.String(),
		IP:         network.Node.Address, 
		Data: 		data,
	}

	return SendDial(contact.Address, sendMsg)
}

func (network *Network) handleStoreMessage(message Message) {
	fmt.Println("handleStoreMessage: ", message.Data, " ", string(message.Data))
	// TODO
	//Send data to target
	network.Kademlia.addData(message.Data)
	fmt.Println("handleStore2")
	senderContact := NewContact(NewKademliaID(message.KademliaID), message.IP)
	network.rt.AddContact(senderContact)
}

func (network *Network) handleInput(message []byte, addr *net.UDPAddr) {
	var pingMsg Message
	err := json.Unmarshal(message, &pingMsg)
	if err != nil {
		log.Println("Error unmarshalling received message:", err)
		return
	}
	// Now you can use pingMsg.Type, pingMsg.KademliaID, etc.

	switch {
	case pingMsg.Type == "ping":
		network.handlePing(pingMsg, addr)
	case pingMsg.Type == "pong":
		handlePong(pingMsg)
	case pingMsg.Type == "FindContact":
		fmt.Println("Tst: ")
		fmt.Println(pingMsg.Target)
		fmt.Println("pingMsg: ", pingMsg)
		network.handleFindContact(pingMsg)
	case pingMsg.Type == "StoreMessage":
		fmt.Println("Tst222: ")
		network.handleStoreMessage(pingMsg)
	case pingMsg.Type == "SendDataMessage":
		network.handleFindDataMessage(pingMsg)
	case pingMsg.Type == "err":
		fmt.Println("Error in error")
	}
}

func (network *Network) handlePing(message Message, addr *net.UDPAddr) Message {
	fmt.Println("handlePing")
	//Respnse to sender
	fmt.Println("Ping from: ", message)

	senderContact := NewContact(NewKademliaID(message.KademliaID), message.IP)
	network.rt.AddContact(senderContact)

	sendMsg := &Message{
		Type:       "pong",
		KademliaID: network.Node.ID.String(),
		IP:         network.Node.Address, //
	}
	return SendDial(message.IP, sendMsg)
}

func handlePong(message Message) {
	fmt.Println("Pong from: ", message)
}

func (network *Network) handleFindContact(message Message) Message {
	fmt.Println("handleFindContact")
	var K = 10
	fmt.Println("Target ", message.Target)
	fmt.Println("Xasd: ", message)
	contacts := network.rt.FindClosestContacts(NewKademliaID(message.Target), K)
	//self := network.rt.FindClosestContacts(network.Node.ID, K)
	
	senderContact := NewContact(NewKademliaID(message.KademliaID), message.IP)
	network.rt.AddContact(senderContact)

	//fmt.Println("This is the self printout: \n", self)
	sendContacts := &Message{
		Type:       "ReturnContact",
		KademliaID: network.Node.ID.String(),
		IP:         network.Node.Address,
		Contacts:   contacts,
	}

	return SendDial(message.IP, sendContacts)
}

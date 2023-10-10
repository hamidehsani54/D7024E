package logic

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"net"

	//"strings"
	"encoding/json"
	"time"
)

type Network struct {
	Node                 *Contact
	rt                   *RoutingTable
	Kademlia             *Kademlia
	pendingRequests      map[string]chan []Contact
	pendingDataRequests  map[string]chan Message
	pendingStoreRequests map[string]chan []byte
}

func InitNetwork(id *KademliaID, address string) *Network {
	node := &Contact{
		ID:      id,
		Address: address,
	}

	net := &Network{
		Node:                 node,
		rt:                   NewRoutingTable(*node),
		pendingRequests:      make(map[string]chan []Contact), // Initialize the map here
		pendingDataRequests:  make(map[string]chan Message),
		pendingStoreRequests: make(map[string]chan []byte),
	}

	net.Kademlia = &Kademlia{
		Network:  net,
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
	Data       []byte
	Key        string `json:"Key"`
	RequestID  string `json:"requestID"`
}

const (
	BufferSize  = 4096
	DialTimeout = (5 * time.Second)
	DefaultPort = 4000
)

func (network *Network) Listen(ip string, port int) error {
	addr := fmt.Sprintf("%s:%d", ip, port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	listener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Printf("Listening on %s...\n", udpAddr.String())

	buffer := make([]byte, BufferSize)
	for {
		n, _, err := listener.ReadFromUDP(buffer)
		if err != nil {
			log.Print(err)
			continue
		}
		//receivedData := string(buffer[:n])
		network.handleInput(buffer[:n])
		//log.Printf("Received data from %s: %s\n", addr, receivedData)
	}
}

func SendDial(targetIp string, message *Message) error {
	addr := fmt.Sprintf("%s:%d", targetIp, DefaultPort)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return fmt.Errorf("Failed to dial to %s: %v", addr, err)
	}
	defer conn.Close()

	msgBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("Error marshalling message: %v", err)
	}

	_, err = conn.Write(msgBytes)
	if err != nil {
		return fmt.Errorf("Writing error to %s: %v in sendMsg", addr, err)
	}

	return nil
}

func (network *Network) handleInput(message []byte) {
	var receivedMessage Message
	err := json.Unmarshal(message, &receivedMessage)
	if err != nil {
		log.Println("Error unmarshalling received message:", err)
		return
	}
	// Now you can use pingMsg.Type, pingMsg.KademliaID, etc.

	switch receivedMessage.Type {
	case "ping":
		network.handlePing(receivedMessage)
	case "pong":
		network.handlePong(receivedMessage)
	case "FindContact":
		network.handleFindContact(receivedMessage)
	case "FindContactResponse":
		network.handleFindContactResponse(receivedMessage)
	case "StoreMessage":
		network.handleStoreMessage(receivedMessage)
	case "SendDataMessage":
		network.handleFindDataMessage(receivedMessage)
	case "SendDataResponse":
		network.handleDataResponse(receivedMessage)
	case "SendStoreResponse":
		network.handleStoreResponse(receivedMessage)
	case "err":
		log.Println("Received error message")
	default:
		log.Printf("Unknown message type: %s", receivedMessage.Type)
	}
}

func generateUniqueRequestID() string {
	//time stamp uniqueid
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func (network *Network) SendPingMessage(contact *Contact) error {
	//fmt.Println("SendPingMessage")
	message := Message{
		Type:       "ping",
		KademliaID: network.Node.ID.String(),
		IP:         network.Node.Address,
		RequestID:  generateUniqueRequestID(),
	}

	return SendDial(contact.Address, &message)
}

func (network *Network) handlePing(message Message) error {
	//fmt.Println("handlePing")

	senderContact := NewContact(NewKademliaID(message.KademliaID), message.IP)
	network.rt.AddContact(senderContact)

	sendMsg := &Message{
		Type:       "pong",
		KademliaID: network.Node.ID.String(),
		IP:         network.Node.Address, //
	}
	return SendDial(message.IP, sendMsg)
}

func (network *Network) SendFindContactMessage(contact *Contact, target string) (chan []Contact, error) {
	//fmt.Println("SendFindContactMessage")

	requestID := generateUniqueRequestID()

	ch := make(chan []Contact)
	network.pendingRequests[requestID] = ch

	message := Message{
		Type:       "FindContact",
		KademliaID: network.Node.ID.String(),
		IP:         network.Node.Address,
		Target:     target,
		RequestID:  requestID,
	}

	err := SendDial(contact.Address, &message)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (network *Network) handleFindContact(message Message) {
	//fmt.Println("Handling Find Contact Request...")
	targetID := NewKademliaID(message.Target)

	contacts := network.rt.FindClosestContacts(targetID, 10)
	sendMsg := &Message{
		Type:       "FindContactResponse",
		KademliaID: network.Node.ID.String(),
		IP:         network.Node.Address,
		Contacts:   contacts,
		RequestID:  message.RequestID,
	}

	err := SendDial(message.IP, sendMsg)
	if err != nil {
		log.Printf("Couldnt send FindContact to %s", message.IP)
	}
}

func (network *Network) handleFindContactResponse(message Message) {
	ch, ok := network.pendingRequests[message.RequestID]
	if !ok {
		log.Printf("Unknown request ID: %s", message.RequestID)
		return
	}
	ch <- message.Contacts

	delete(network.pendingRequests, message.RequestID)
	close(ch)
}

func (network *Network) SendFindDataMessage(contact *Contact, hash string) (chan Message, error) {
	//fmt.Println("Requesting Data...")

	requestID := generateUniqueRequestID()
	ch := make(chan Message)
	network.pendingDataRequests[requestID] = ch

	sendMsg := &Message{
		Type:       "SendDataMessage",
		KademliaID: network.Node.ID.String(),
		IP:         network.Node.Address,
		Key:        hash,
		RequestID:  requestID,
	}

	err := SendDial(contact.Address, sendMsg)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (network *Network) handleFindDataMessage(message Message) {
	//fmt.Println("Handling Data Request...")

	// Look for the data based on the hash
	hash, data := network.Kademlia.FindLocalData(message.Key)

	sendData := []byte{}
	if message.Key == hash {
		sendData = data
	}

	sendMsg := &Message{
		Type:       "SendDataResponse",
		KademliaID: network.Node.ID.String(),
		IP:         network.Node.Address,
		Data:       sendData,
		RequestID:  message.RequestID,
	}

	SendDial(message.IP, sendMsg)
}

func (network *Network) handleDataResponse(message Message) {
	//fmt.Println("Handling Data Response...")

	ch, exists := network.pendingDataRequests[message.RequestID]
	if !exists {
		log.Printf("Unknown request ID: %s", message.RequestID)
		return
	}

	ch <- message
	close(ch)
	delete(network.pendingDataRequests, message.RequestID)
}

func (network *Network) SendStoreMessage(contact *Contact, data []byte) (chan []byte, error) {
	requestID := generateUniqueRequestID()

	ch := make(chan []byte)
	network.pendingStoreRequests[requestID] = ch

	sendMsg := &Message{
		Type:       "StoreMessage",
		KademliaID: network.Node.ID.String(),
		IP:         network.Node.Address,
		Data:       data,
		RequestID:  requestID,
	}

	err := SendDial(contact.Address, sendMsg)

	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (network *Network) handleStoreMessage(message Message) {
	//fmt.Println("handleStoreMessage: ", message.Data, " ", string(message.Data))

	network.Kademlia.addData(message.Data)
	senderContact := NewContact(NewKademliaID(message.KademliaID), message.IP)
	network.rt.AddContact(senderContact)

	hash := sha1.Sum(message.Data)
	hashString := hex.EncodeToString(hash[:])

	sendMsg := &Message{
		Type:       "SendStoreResponse",
		KademliaID: network.Node.ID.String(),
		IP:         network.Node.Address,
		Data:       []byte(hashString),
		RequestID:  message.RequestID,
	}

	SendDial(message.IP, sendMsg)
}

func (network *Network) handleStoreResponse(message Message) {
	//fmt.Println("Handling Store Response...")

	ch, exists := network.pendingStoreRequests[message.RequestID]
	if !exists {
		log.Printf("Unknown request ID: %s", message.RequestID)
		return
	}

	ch <- message.Data
	close(ch)
	delete(network.pendingDataRequests, message.RequestID)
}

func (network *Network) handlePong(message Message) {
	//fmt.Println("Pong from: ", message)
	senderContact := NewContact(NewKademliaID(message.KademliaID), message.IP)
	network.rt.AddContact(senderContact)
}

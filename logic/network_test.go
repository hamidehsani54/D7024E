package logic

import (
	"encoding/json"
	"testing"
)

func TestGenerateUniqueRequestID(t *testing.T) {
	id1 := generateUniqueRequestID()
	id2 := generateUniqueRequestID()

	if id1 == id2 {
		t.Fatalf("Generated IDs are not unique: %s, %s", id1, id2)
	}
}

func TestHandleFindContactResponse(t *testing.T) {
	network := InitNetwork(NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51"), "localhost:8000")
	requestID := "testRequestID1"
	ch := make(chan []Contact, 1)
	network.pendingRequests[requestID] = ch
	message := Message{
		RequestID: requestID,
	}
	network.handleFindContactResponse(message)

	if _, exists := network.pendingRequests[requestID]; exists {
		t.Fatalf("Expected request ID %s to be removed from pending requests", requestID)
	}
}

func TestHandleDataResponse(t *testing.T) {
	network := InitNetwork(NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51"), "localhost:8000")
	requestID := "testRequestID2"
	ch := make(chan Message, 1)
	network.pendingDataRequests[requestID] = ch

	message := Message{
		RequestID: requestID,
	}
	network.handleDataResponse(message)

	if _, exists := network.pendingDataRequests[requestID]; exists {
		t.Fatalf("Expected request ID %s to be removed from pending data requests", requestID)
	}
}

func TestHandleInput(t *testing.T) {
	net := InitNetwork(NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51"), "localhost:8000")
	input := []byte("ping")
	net.handleInput(input)

	msg := Message{
		Type:       "ping",
		KademliaID: "27f2d5effb3dcfe4d7bdd17e64a3101226648a51",
		IP:         "192.168.1.1",
		Target:     "targetNodeID",
		Contacts:   []Contact{{}},
		Data:       []byte("Some data"),
		Key:        "exampleKey",
		RequestID:  "requestID1234",
	}
	data, _ := json.Marshal(msg)
	net.handleInput([]byte(data))

	msg2 := Message{
		Type:       "pong",
		KademliaID: "27f2d5effb3dcfe4d7bdd17e64a3101226648a51",
		IP:         "192.168.1.1",
		Target:     "targetNodeID",
		Contacts:   []Contact{{}},
		Data:       []byte("Some data"),
		Key:        "exampleKey",
		RequestID:  "requestID1234",
	}

	data2, _ := json.Marshal(msg2)
	net.handleInput([]byte(data2))

}

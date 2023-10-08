package logic

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"testing"
)

func TestKademliaDataFunctions(t *testing.T) {
	// Initialize a mock Kademlia instance
	kademlia := &Kademlia{
		Network:  nil, // As it's not used in the functions we're testing
		DataList: []DataStore{},
	}

	// Test addData
	data := []byte("Test Data")
	kademlia.addData(data)

	if len(kademlia.DataList) != 1 {
		t.Fatalf("expected 1 item in DataList, got %d", len(kademlia.DataList))
	}

	hash := sha1.Sum(data)
	hashString := hex.EncodeToString(hash[:])

	if kademlia.DataList[0].Hash != hashString || string(kademlia.DataList[0].Data) != string(data) {
		t.Fatalf("stored data or hash doesn't match. Expected hash %s with data %s, got hash %s with data %s",
			hashString, string(data), kademlia.DataList[0].Hash, string(kademlia.DataList[0].Data))
	}

	// Test PrintData
	expectedOutput := `Stored Data:
Item 1:
  Hash: ` + hashString + `
  Data: Test Data
`
	printedData := kademlia.PrintData()

	if !strings.EqualFold(printedData, expectedOutput) {
		t.Fatalf("Printed data doesn't match expected output. Expected:\n%s\nGot:\n%s", expectedOutput, printedData)
	}
}

func TestKademliaFindLocalData(t *testing.T) {
	// Initialize a mock Kademlia instance with some data
	data := []byte("Test Data")
	hash := sha1.Sum(data)
	hashString := hex.EncodeToString(hash[:])

	kademlia := &Kademlia{
		Network: nil,
		DataList: []DataStore{
			{Data: data, Hash: hashString},
		},
	}

	// Test FindLocalData
	foundHash, foundData := kademlia.FindLocalData(hashString)
	if foundHash != hashString || string(foundData) != string(data) {
		t.Fatalf("expected hash %s with data %s, got hash %s with data %s", hashString, string(data), foundHash, string(foundData))
	}

	// Test for non-existent hash
	_, notFoundData := kademlia.FindLocalData("nonexistenthash")
	if notFoundData != nil {
		t.Fatal("expected to not find data for nonexistent hash")
	}
}

func TestAllContactsAdd(t *testing.T) {
	// Mock data
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	contact2 := NewContact(NewKademliaID("EEEEEEEE00000000000000000000000000000000"), "localhost:8002")
	target := NewKademliaID("DDDDDDDD00000000000000000000000000000000")

	// Initialize allContacts instance
	allContactsInstance := &allContacts{
		Contacts: []Contact{},
		Seen:     make(map[string]bool),
	}

	// Add contacts
	allContactsInstance.Add(contact1, target)
	allContactsInstance.Add(contact2, target)

	// Test order based on distance
	if !allContactsInstance.Contacts[0].ID.Equals(contact2.ID) || !allContactsInstance.Contacts[1].ID.Equals(contact1.ID) {
		t.Fatal("Contacts not in expected order based on distance")
	}
}

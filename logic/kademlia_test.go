package logic

import (
	"crypto/sha1"
	"encoding/hex"
	"testing"
)

func TestKademliaFindLocalData(t *testing.T) {
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
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	contact2 := NewContact(NewKademliaID("EEEEEEEE00000000000000000000000000000000"), "localhost:8002")
	target := NewKademliaID("DDDDDDDD00000000000000000000000000000000")

	// Initialize allContacts instance
	allContactsInstance := &allContacts{
		Contacts: []Contact{},
		Seen:     make(map[string]bool),
	}

	allContactsInstance.Add(contact1, target)
	allContactsInstance.Add(contact2, target)

	// Test distance
	if !allContactsInstance.Contacts[0].ID.Equals(contact2.ID) || !allContactsInstance.Contacts[1].ID.Equals(contact1.ID) {
		t.Fatal("Contacts not in expected order based on distance")
	}
}

func TestJoinNetwork(t *testing.T) {
	net := InitNetwork(NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51"), "localhost:8000")
	k := net.Kademlia
	k.JoinNetwork()
	contacts := k.LookupContact(net.Node.ID)
	if len(contacts) == 0 {
		t.Fatal("Failed to join network")
	}
}

/*func TestLookUpData(t *testing.T) {
	node1 := InitNetwork(NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51"), "localhost:8000")

	data := []byte("Test")

	hashString := node1.Kademlia.Store(data)

	if hashString == "Failed to store data" || hashString == "Failed to return data" {
		t.Fatalf("Failed to store data: %s", hashString)
	}
}
*/

package logic

import (
	"testing"
)

func TestNewContact(t *testing.T) {
	id := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	address := "localhost:8001"
	contact := NewContact(id, address)

	if contact.ID != id {
		t.Fatalf("expected id to be %s, got %s", id, contact.ID)
	}

	if contact.Address != address {
		t.Fatalf("expected address to be %s, got %s", address, contact.Address)
	}
}

func TestCalcDistance(t *testing.T) {
	contact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	target := NewKademliaID("0000000000000000000000000000000000000000")

	contact.CalcDistance(target)

	if !contact.distance.Equals(NewKademliaID("FFFFFFFF00000000000000000000000000000000")) {
		t.Fatalf("expected distance to be FFFFFFFF00000000000000000000000000000000, got %s", contact.distance)
	}
}

func TestContactLess(t *testing.T) {
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	contact2 := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8002")

	contact1.CalcDistance(NewKademliaID("1000000000000000000000000000000000000000"))
	contact2.CalcDistance(NewKademliaID("1000000000000000000000000000000000000000"))

	if !contact2.Less(&contact1) {
		t.Fatal("expected contact2 to be less than contact1")
	}
}

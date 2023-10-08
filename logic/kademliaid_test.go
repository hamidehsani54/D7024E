package logic

import (
	"testing"
)

func TestKademliaID(t *testing.T) {
	// Test the creation of a new KademliaID
	id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	if id1.String() != "ffffffff00000000000000000000000000000000" {
		t.Fatalf("expected id1 to be FFFFFFFF00000000000000000000000000000000, got %s", id1.String())
	}

	// Test random generation of a KademliaID
	id2 := NewRandomKademliaID()
	if len(id2.String()) != 40 { // 20 bytes = 40 characters in hex
		t.Fatalf("expected length of id2 to be 40, got %d", len(id2.String()))
	}

	// Validate the Less function
	id3 := NewKademliaID("EEEEEEEE00000000000000000000000000000000")
	if !id3.Less(id1) {
		t.Fatal("expected id3 to be less than id1")
	}

	// Validate the Equals function
	id4 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	if !id4.Equals(id1) {
		t.Fatal("expected id4 to be equal to id1")
	}

	// Test CalcDistance function
	distance := id1.CalcDistance(id3)
	expectedDistance := NewKademliaID("1111111100000000000000000000000000000000")
	if !distance.Equals(expectedDistance) {
		t.Fatalf("expected distance to be 1111111100000000000000000000000000000000, got %s", distance.String())
	}

	// Test string representation
	if id1.String() != "ffffffff00000000000000000000000000000000" {
		t.Fatalf("expected id1 string representation to be FFFFFFFF00000000000000000000000000000000, got %s", id1.String())
	}

}

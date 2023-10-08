package logic

import (
	"testing"
)

func TestLenBucket(t *testing.T) {
	bucket := newBucket()

	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000001"), "localhost:8001")
	contact2 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000002"), "localhost:8002")
	contact3 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000003"), "localhost:8003")
	contact4 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000004"), "localhost:8004")
	contact5 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000005"), "localhost:8005")
	contact6 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000006"), "localhost:8006")

	// Add contacts to bucket
	bucket.AddContact(contact1)
	bucket.AddContact(contact2)
	bucket.AddContact(contact3)
	bucket.AddContact(contact4)
	bucket.AddContact(contact5)
	bucket.AddContact(contact6)

	// Ensure contacts are added
	if bucket.Len() != 6 {
		t.Fatalf("expected bucket length to be 6, got %d", bucket.Len())
	}
}

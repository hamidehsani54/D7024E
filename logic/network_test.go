package logic

import (
	"testing"
)

func TestGenerateUniqueRequestID(t *testing.T) {
	id1 := generateUniqueRequestID()
	id2 := generateUniqueRequestID()

	if id1 == id2 {
		t.Fatalf("Generated IDs are not unique: %s, %s", id1, id2)
	}
}

package logic

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

type Kademlia struct {
	Network  *Network
	DataList []DataStore
}

type DataStore struct {
	Data []byte
	Hash string
}

type allContacts struct {
	Contacts []Contact
	Seen     map[string]bool
}

const alpha = 3

func (kademlia *Kademlia) JoinNetwork() {
	kademlia.Network.rt.AddContact(NewContact(NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51"), "masterNode"))

	contacts := kademlia.LookupContact(kademlia.Network.Node.ID)
	if len(contacts) == 0 {
		fmt.Printf("Failed to join network")
	}
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) []Contact {
	contacts := kademlia.iterativeFindNode(target)
	return contacts
}

func (kademlia *Kademlia) iterativeFindNode(nodeID *KademliaID) []Contact {
	kNearest := kademlia.Network.rt.FindClosestContacts(nodeID, alpha)
	nodeList := &allContacts{
		Contacts: []Contact{},
		Seen:     make(map[string]bool),
	}

	for _, contact := range kNearest {
		nodeList.Add(contact, nodeID)
	}
	return kademlia.lookupContactHelp(nodeID, kNearest, nodeList)
}

func (kademlia *Kademlia) lookupContactHelp(nodeID *KademliaID, earlierContacts []Contact, nodeList *allContacts) []Contact { // Added error return type
	resultsChan := make(chan []Contact)
	var wg sync.WaitGroup
	for _, contact := range earlierContacts {
		if _, found := nodeList.Seen[contact.ID.String()]; !found {
			wg.Add(1)
			nodeList.Seen[contact.ID.String()] = true
			go kademlia.queryNodeForClosestContacts(contact, nodeID.String(), resultsChan, &wg)
		}
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for contacts := range resultsChan {
		for _, contact := range contacts {
			nodeList.Add(contact, nodeID)
		}
	}

	contactClosest := kademlia.Network.rt.FindClosestContacts(nodeID, alpha)
	foundContacts := 0
	for _, contact := range contactClosest {
		for _, previousContact := range earlierContacts {
			if contact.ID.Equals(previousContact.ID) {
				foundContacts++
				break
			}
		}
	}

	if foundContacts == len(contactClosest) {
		if len(nodeList.Contacts) < alpha {
			return nodeList.Contacts
		}
		return nodeList.Contacts[:alpha]
	} else {
		return kademlia.lookupContactHelp(nodeID, contactClosest, nodeList)
	}
}

func (kademlia *Kademlia) queryNodeForClosestContacts(contact Contact, target string, ch chan []Contact, wg *sync.WaitGroup) {
	defer wg.Done()

	responseCh, err := kademlia.Network.SendFindContactMessage(&contact, target)
	if err != nil {
		log.Printf("Failed to send find contact message to %s: %v", contact.Address, err)
		return
	}

	select {
	case contacts := <-responseCh:
		ch <- contacts
	case <-time.After(20 * time.Second):
		log.Printf("Timed out waiting for response from %s", contact.Address)
	}

}

func (s *allContacts) Add(contact Contact, target *KademliaID) {

	if _, found := s.Seen[contact.ID.String()]; found {
		return
	}
	//s.Seen[contact.ID.String()] = true

	contact.CalcDistance(target)

	position := sort.Search(len(s.Contacts), func(i int) bool {
		return s.Contacts[i].distance.Less(contact.distance) || s.Contacts[i].distance.Equals(contact.distance)
	})

	s.Contacts = append(s.Contacts, Contact{})
	copy(s.Contacts[position+1:], s.Contacts[position:])
	s.Contacts[position] = contact
}

func (kademlia *Kademlia) LookupData(hash string) (string, string, []Contact) {
	closestContacts := kademlia.LookupContact(NewKademliaID(hash))
	responseCh := make(chan Message, len(closestContacts))
	errorMsg := &Message{
		Type: "Error",
		Data: []byte("Error in request"),
	}

	var wg sync.WaitGroup

	for _, contact := range closestContacts {
		wg.Add(1)
		go func(contact Contact) {
			defer wg.Done()

			dataCh, err := kademlia.Network.SendFindDataMessage(&contact, hash)
			if err != nil {
				log.Printf("Failed to send find data message to %s: %v", contact.Address, err)
				responseCh <- *errorMsg
				return
			}

			select {
			case data := <-dataCh:
				responseCh <- data
			case <-time.After(10 * time.Second):
				log.Printf("Timed out waiting for data response from %s", contact.Address)
				responseCh <- *errorMsg
			}
		}(contact)
	}

	wg.Wait()
	close(responseCh)

	for data := range responseCh {
		if data.Type != "Error" {
			return string(data.Data), string(data.KademliaID), nil
		}
	}
	return "Could not find data", " ", closestContacts
}

func (kademlia *Kademlia) FindLocalData(hash string) (string, []byte) {
	for _, datastore := range kademlia.DataList {
		if datastore.Hash == hash {
			return datastore.Hash, datastore.Data
		}
	}
	return "Error", nil
}

func (kademlia *Kademlia) Store(data []byte) string {
	hash := sha1.Sum(data)
	hashString := hex.EncodeToString(hash[:])
	dataTarget := kademlia.LookupContact(NewKademliaID(hashString))

	responseCh := make(chan []byte, len(dataTarget))

	//dataTarget = dataTarget[:10]
	for _, contact := range dataTarget {
		if kademlia.Network.Node.ID == contact.ID {
			kademlia.addData(data)
			continue
		}
		dataCh, err := kademlia.Network.SendStoreMessage(&contact, data)
		if err != nil {
			log.Printf("Failed to send find data message to %s: %v", contact.Address, err)
			responseCh <- nil
			return "Failed to store data"
		}
		select {
		case data := <-dataCh:
			responseCh <- data
		case <-time.After(10 * time.Second):
			log.Printf("Timed out waiting for data response from %s", contact.Address)
			responseCh <- nil
		}
	}

	for data := range responseCh {
		if data != nil {
			return hashString
		}
	}
	return "Failed to return data"
}

func (kademlia *Kademlia) addData(data []byte) {
	hash := sha1.Sum(data)

	println("Data: ", string(data), " Hash: ", string(hash[:]))
	hashString := hex.EncodeToString(hash[:])

	dataStore := DataStore{
		Data: data,
		Hash: hashString,
	}

	kademlia.DataList = append(kademlia.DataList, dataStore)
}

func (kademlia *Kademlia) PrintData() string {
	var result strings.Builder
	result.WriteString("Stored Data:\n")

	for i, dataStore := range kademlia.DataList {
		result.WriteString(fmt.Sprintf("Item %d:\n", i+1))
		result.WriteString(fmt.Sprintf("  Hash: %s\n", dataStore.Hash))
		result.WriteString(fmt.Sprintf("  Data: %s\n", string(dataStore.Data)))
	}

	return result.String()
}

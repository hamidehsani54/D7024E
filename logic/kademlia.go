package logic

import (
  "sort"
  "fmt"
)

type Kademlia struct {
	Network *Network
}

type allContacts struct {
    Contacts []Contact
    Seen     map[string]bool
}


const alpha = 3
func (kademlia *Kademlia) JoinNetwork() {
	fmt.Println("Joining started...")
	// add bootstrap node to routing table
	kademlia.Network.rt.AddContact(NewContact(NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51"), "10.10.0.10"))
    
	// lookup on itself
	res := kademlia.LookupContact(kademlia.Network.Node.ID)

	fmt.Println("Join lookup result:")
	for _, c := range res {
		fmt.Println("  ", c.String())
	}
}

func InitKademlia(network *Network) *Kademlia {
	return &Kademlia{
		Network: network,
	}
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) []Contact {
    fmt.Println("LookupContact started...")
	contacts := kademlia.iterativeFindNode(target)
	return contacts
}


func (kademlia *Kademlia) iterativeFindNode(nodeID *KademliaID) []Contact {
    kNearest := kademlia.Network.rt.FindClosestContacts(nodeID, alpha)
    //queriedNodes := make(map[string]bool)
    resultsChan := make(chan []Contact)
    
    nodeList := &allContacts{
        Contacts: []Contact{},
        Seen:     make(map[string]bool),
    }

    for _, contact := range kNearest {
        fmt.Println("iterativeFindNode ---------------------...")
        nodeList.Add(contact, nodeID)
        go kademlia.queryNodeForClosestContacts(contact, nodeID, resultsChan)
    }

    for contacts := range resultsChan {
        for _, contact := range contacts {
            if _, found := nodeList.Seen[contact.ID.String()]; !found{
                nodeList.Add(contact, nodeID)
                go kademlia.queryNodeForClosestContacts(contact, nodeID, resultsChan)
            }  
        }
        /*
        if allContacts.HasConverged() { // Some condition to check if you're no longer finding closer nodes
            close(resultsChan)
        }*/
    }

    return nodeList.Contacts[:alpha] // Returns top alpha closest contacts
}

func (kademlia *Kademlia) queryNodeForClosestContacts(contact Contact, target *KademliaID, ch chan []Contact) {
    contacts := kademlia.Network.SendFindContactMessage(&contact).Contact
    //Ta bara K närmaste
    ch <- contacts
}

func (s *allContacts) Add(contact Contact, target *KademliaID) {
    // Check if contact already exists
    if _, found := s.Seen[contact.ID.String()]; found {
        return
    }

    // Mark as seen
    s.Seen[contact.ID.String()] = true

    contact.CalcDistance(target)

    // Find the position to insert using binary search
    position := sort.Search(len(s.Contacts), func(i int) bool {
        return s.Contacts[i].distance.Less(contact.distance) || s.Contacts[i].distance.Equals(contact.distance)
    })

    // Insert at the found position
    s.Contacts = append(s.Contacts, Contact{})
    copy(s.Contacts[position+1:], s.Contacts[position:])
    s.Contacts[position] = contact
}


func (kademlia *Kademlia) sortTakeK(contacts []Contact, alpha int) []Contact {
    // Compute distances to the target (each contact is its own target)
    for i := range contacts {
        contacts[i].CalcDistance(contacts[i].ID)
    }

    // Sort based on computed distance
    sort.Slice(contacts, func(i, j int) bool {
        return contacts[i].Less(&contacts[j])
    })

    // Take the first alpha contacts
    if len(contacts) > alpha {
        return contacts[:alpha]
    }

    return contacts
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
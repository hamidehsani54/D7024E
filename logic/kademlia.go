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
	
	// Add master node to contacts, with the known values
	kademlia.Network.rt.AddContact(NewContact(NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51"), "10.10.0.10"))

	kademlia.LookupContact(kademlia.Network.Node.ID)
}

func InitKademlia(network *Network) *Kademlia {
	return &Kademlia{
		Network: network,
	}
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) []Contact {
    fmt.Println("LookupContact-------")
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
        Add condition so that we dont continue for too long.
        Take the K first from our sorted list and query the nodes that has not yet been visited/queried
        insert the new nodes in the list and keep it sorted, repeat until the top K is visited.
        }*/
    }

    return nodeList.Contacts[:alpha] // Returns top alpha closest contacts
}

func (kademlia *Kademlia) queryNodeForClosestContacts(contact Contact, target *KademliaID, ch chan []Contact) {
    fmt.Println(&contact)
    c := kademlia.Network.SendFindContactMessage(&contact, target)
    fmt.Println(c.Contacts)
    //Ta bara K närmaste
    ch <- c.Contacts
}

func (s *allContacts) Add(contact Contact, target *KademliaID) {
    // Check if contact already exists
    if _, found := s.Seen[contact.ID.String()]; found {
        return
    }

    // Mark as seen
    s.Seen[contact.ID.String()] = true

    contact.CalcDistance(target)

    // Find the position to insert element
    position := sort.Search(len(s.Contacts), func(i int) bool {
        return s.Contacts[i].distance.Less(contact.distance) || s.Contacts[i].distance.Equals(contact.distance)
    })

    // Insert at the found position
    s.Contacts = append(s.Contacts, Contact{})
    copy(s.Contacts[position+1:], s.Contacts[position:])
    s.Contacts[position] = contact
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
package logic

import (
  "sort"
  "fmt"
  "sync"
  "crypto/sha1"
  "encoding/hex"
)

type Kademlia struct {
	Network *Network
    DataList []DataStore
}

type DataStore struct {
    data string
    hash string
}

type allContacts struct {
    Contacts []Contact
    Seen     map[string]bool
}


const alpha = 3
func (kademlia *Kademlia) JoinNetwork() {
	
	// Add master node to contacts, with the known values
	kademlia.Network.rt.AddContact(NewContact(NewKademliaID("27f2d5effb3dcfe4d7bdd17e64a3101226648a51"), "masterNode"))

	kademlia.LookupContact(kademlia.Network.Node.ID)
}

func InitKademlia(network *Network) *Kademlia {
	return &Kademlia{
		Network: network,
	}
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) []Contact {
    fmt.Println("LookupContact-------", target)
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

    fmt.Println("return kadmelia")
    return kademlia.lookupContactHelp(nodeID, kNearest, nodeList)
}

func (kademlia *Kademlia) lookupContactHelp(nodeID *KademliaID, earlierContacts []Contact, nodeList *allContacts) []Contact {
    resultsChan := make(chan []Contact)
    var wg sync.WaitGroup
    fmt.Println("1")
    for _, contact := range earlierContacts {
        if _, found := nodeList.Seen[contact.ID.String()]; !found {
            fmt.Println("2")
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
    println("target: in KADEMLIA ", target)
    defer wg.Done()
    c := kademlia.Network.SendFindContactMessage(&contact, target)

    //self := kademlia.Network.rt.FindClosestContacts(target, 10)
    //print("Self print: ", self)
    print("\n Others print: ", c.Contacts)
    ch <- c.Contacts
}



func (s *allContacts) Add(contact Contact, target *KademliaID) {
    // Check if contact already exists
    if _, found := s.Seen[contact.ID.String()]; found {
        return
    }

    // Mark as seen
    //s.Seen[contact.ID.String()] = true

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
    // Compute the hash of the data
    hash := sha1.Sum(data)
    hashString := hex.EncodeToString(hash[:])
    dataTarget := kademlia.LookupContact(NewKademliaID(hashString))

    //dataTarget = dataTarget[:10]

    for _, contact := range dataTarget{
        if kademlia.Network.Node.ID == contact.ID{
            kademlia.addData(data)
        }
        kademlia.Network.SendStoreMessage(&contact, data)
    }
}

func (kademlia *Kademlia) addData(data []byte){
    hash := sha1.Sum(data)
    hashString := hex.EncodeToString(hash[:])

    dataStore := DataStore{
        data: string(data),
        hash: hashString,
    }

    kademlia.DataList = append(kademlia.DataList, dataStore)
}
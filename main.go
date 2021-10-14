// BOLLON hugo / RODRIGUEZ Samuel
package main

import (
	"fmt"
	"sync"
	"time"
)

const NbProcess = 4

var (
	ProcessPool []*Process
)

func main() {
	// Instanciation du bus de communication
	bus := CreateBus()

	// Initialisation de la pool de process
	ProcessPool = make([]*Process, NbProcess)

	// Instanciation du pool de process
	for i := 0; i < NbProcess; i++ {
		process := Process{
			id:       i,
			syncChan: make(chan *sync.WaitGroup, 1),
			com: Com{
				bus:       &bus,
				clock:     0,
				mutex:     sync.Mutex{},
				letterBox: make(LetterBox, 10),
			},
		}
		// Lancement d'une goroutine exécutant la focntion reader pour chaque process
		go process.reader()
		ProcessPool[i] = &process
	}

	broadcastMessage := BroadcastMessage{
		Message{
			Sender:    3,
			Msg:       "Bonjour",
			Timestamp: 0,
			Synchrone: true,
		},
	}

	dedicatedMessage := DedicatedMessage{
		Message: Message{
			Sender:    3,
			Msg:       "Bonjour",
			Timestamp: 0,
			Synchrone: true,
		},
		Receiver: 1,
	}

	fmt.Println("On fait un broadcast depuis le processus 3 :")
	ProcessPool[3].com.broadcast(broadcastMessage)
	ProcessPool[3].com.broadcastSync(broadcastMessage)
	fmt.Println("On envoie un message à 1 depuis le processus 3 :")
	ProcessPool[3].com.sendTo(dedicatedMessage)
	ProcessPool[3].com.sendToSync(dedicatedMessage)
	ProcessPool[3].SyncAll()
	time.Sleep(5 * time.Second)
}

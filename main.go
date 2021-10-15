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
			id:        i,
			syncChan:  make(chan *sync.WaitGroup, 1),
			tokenChan: make(chan Token, 1),
			com: Com{
				bus:       &bus,
				clock:     0,
				mutex:     sync.Mutex{},
				letterBox: make(LetterBox, 10),
			},
		}
		// Lancement d'une goroutine exécutant la fonction reader pour chaque process
		go process.reader()

		// Lancement d'une goroutine exécutant la fonction onToken pour chaque process
		go process.OnToken()

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

	// On mets le token sur le circuit
	ProcessPool[0].SendToken()

	fmt.Println("On fait un broadcast depuis le processus 3 :")
	ProcessPool[3].com.broadcast(broadcastMessage)
	fmt.Println("On fait un broadcast synchrone depuis le processus 3 :")
	ProcessPool[3].com.broadcastSync(broadcastMessage)
	fmt.Println("On envoie un message à 1 depuis le processus 3 :")
	ProcessPool[3].com.sendTo(dedicatedMessage)
	fmt.Println("On envoie un message synchrone à 1 depuis le processus 3 :")
	ProcessPool[3].com.sendToSync(dedicatedMessage)
	ProcessPool[3].SyncAll()
	time.Sleep(5 * time.Second)

	// Test de la section critique
	ProcessPool[2].RequestCriticalSection(func(i ...interface{}) {
		fmt.Println(i)
	}, "Execution du callback par le processus 2 car il a la section critique")
}

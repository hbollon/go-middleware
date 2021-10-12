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

	// Instanciation du pool de process
	for i := 0; i < NbProcess; i++ {
		process := Process{
			id: i,
			com: Com{
				bus:   &bus,
				clock: 0,
				mutex: sync.Mutex{},
			},
		}
		// Lancement d'une goroutine exécutant la focntion reader pour chaque process
		go process.reader()
		ProcessPool = append(ProcessPool, &process)
	}

	broadcastMessage := BroadcastMessage{
		Message{
			Sender:    3,
			Msg:       "Bonjour",
			Timestamp: 0,
			Synchrone: true,
		},
	}

	broadcasMessagetAsync := BroadcastMessage{
		Message{
			Sender:    3,
			Msg:       "Bonjour",
			Timestamp: 0,
			Synchrone: false,
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

	dedicatedMessageAsync := DedicatedMessage{
		Message: Message{
			Sender:    3,
			Msg:       "Bonjour",
			Timestamp: 0,
			Synchrone: false,
		},
		Receiver: 1,
	}

	fmt.Println("On fait un broadcast depuis le processus 3 :")
	ProcessPool[3].com.broadcast(broadcastMessage)
	ProcessPool[3].com.broadcast(broadcasMessagetAsync)
	fmt.Println("On envoie un message à 1 depuis le processus 3 :")
	ProcessPool[3].com.sendTo(dedicatedMessage)
	ProcessPool[3].com.sendTo(dedicatedMessageAsync)
	time.Sleep(3 * time.Second)
	fmt.Println("On va lire dans la boite aux lettre de 2 :")
	ProcessPool[2].readLetterBox()
	time.Sleep(5 * time.Second)
}

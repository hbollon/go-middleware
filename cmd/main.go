package main

// BOLLON hugo / RODRIGUEZ Samuel

import (
	"fmt"
	"time"

	"github.com/hbollon/go-middleware"
)

func main() {
	middleware.InitProcessPool(4)

	broadcastMessage := middleware.BroadcastMessage{
		middleware.Message{
			Sender:    3,
			Msg:       "Bonjour",
			Timestamp: 0,
			Synchrone: true,
		},
	}

	dedicatedMessage := middleware.DedicatedMessage{
		Message: middleware.Message{
			Sender:    3,
			Msg:       "Bonjour",
			Timestamp: 0,
			Synchrone: true,
		},
		Receiver: 1,
	}

	// On mets le token sur le circuit
	middleware.ProcessPool[0].SendToken()

	fmt.Println("On fait un broadcast depuis le processus 3 :")
	middleware.ProcessPool[3].Com.Broadcast(broadcastMessage)
	fmt.Println("On fait un broadcast synchrone depuis le processus 3 :")
	middleware.ProcessPool[3].Com.BroadcastSync(broadcastMessage)
	fmt.Println("On envoie un message à 1 depuis le processus 3 :")
	middleware.ProcessPool[3].Com.SendTo(dedicatedMessage)
	fmt.Println("On envoie un message synchrone à 1 depuis le processus 3 :")
	middleware.ProcessPool[3].Com.SendToSync(dedicatedMessage)
	middleware.ProcessPool[3].SyncAll()
	time.Sleep(5 * time.Second)

	// Test de la section critique
	middleware.ProcessPool[2].RequestCriticalSection(func(i ...interface{}) {
		fmt.Println(i)
	}, "Execution du callback par le processus 2 car il a la section critique")
}

package main

import (
	"sync"
	"time"
)

var (
	ProcessPool []*Process
	NbProcess   int
)

func main() {
	bus := CreateBus()
	NbProcess = 4
	for i := 0; i < NbProcess; i++ {
		process := Process{
			id: i,
			com: Com{
				bus:   &bus,
				clock: 0,
				mutex: sync.Mutex{},
			},
		}
		go process.reader()
		ProcessPool = append(ProcessPool, &process)
	}

	message := BroadcastMessage{
		Message{
			Sender:    3,
			Msg:       "Bonjour",
			Timestamp: 0,
		},
	}

	dedicatedMessage := DedicatedMessage{
		Message: Message{
			Sender:    3,
			Msg:       "Bonjour",
			Timestamp: 0,
		},
		Receiver: 1,
	}

	ProcessPool[3].com.broadcast(message)
	ProcessPool[3].com.sendTo(dedicatedMessage)
	time.Sleep(10 * time.Second)
}

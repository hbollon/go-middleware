package main

// BOLLON hugo / RODRIGUEZ Samuel

import (
	"fmt"
	"sync"
	"time"
)

// LetterBox est une liste de messages à partir d'un tableau d'interface
// (pseudo généricité)
type LetterBox chan interface{}

// Process est un processus comportant une instance de Com et un id
type Process struct {
	com       Com
	id        int
	syncChan  chan *sync.WaitGroup
	tokenChan chan Token
	state     StateSectionCritique
}

// Com est notre interface de communication
// Il contient un pointeur vers le bus de communication, un mutex,
// un tableau de messages (boite aux lettres) et une clock
type Com struct {
	clock     int
	mutex     sync.Mutex
	bus       *Bus
	letterBox LetterBox
}

// OnBroadcast est la fonction appelée lorsqu'un message de type BroadcastMessage
// est reçu
func (c *Com) OnBroadcast(receiver int, msg BroadcastMessage) {
	c.mutex.Lock()
	c.clock = Max(c.clock, msg.Timestamp) + 1
	c.mutex.Unlock()
	fmt.Printf("Process n°%d received msg: %v, clock: %d\n", receiver, msg, c.clock)
}

// Broadcast envoi un message asynchrone de type BroadcastMessage directement dans la letterbox
func (c *Com) Broadcast(msg BroadcastMessage) {
	c.mutex.Lock()
	c.clock++
	c.mutex.Unlock()

	msg.Timestamp = c.clock
	for _, p := range ProcessPool {
		p.SendLetterBox(msg)
	}
}

// SendTo envoi un message asynchrone de type DedicatedMessage dans la letterbox du receiver
func (c *Com) SendTo(msg DedicatedMessage) {
	c.mutex.Lock()
	c.clock++
	c.mutex.Unlock()

	msg.Timestamp = c.clock
	ProcessPool[msg.Receiver].SendLetterBox(msg)
}

// BroadcastSync envoi un message synchrone de type BroadcastMessage sur le bus
func (c *Com) BroadcastSync(msg BroadcastMessage) {
	c.mutex.Lock()
	c.clock++
	c.mutex.Unlock()

	msg.Timestamp = c.clock
	for i := 0; i < NbProcess; i++ {
		c.bus.Send(msg)
	}
}

// SendToSync envoi un message synchrone de type DedicatedMessage sur le bus
func (c *Com) SendToSync(msg DedicatedMessage) {
	c.mutex.Lock()
	c.clock++
	c.mutex.Unlock()

	msg.Timestamp = c.clock
	for i := 0; i < NbProcess; i++ {
		c.bus.Send(msg)
	}
}

// OnDedicatedMessage est la fonction appelée lorsqu'un message de type DedicatedMessage
// est reçu
func (c *Com) OnDedicatedMessage(receiver int, msg DedicatedMessage) {
	c.mutex.Lock()
	c.clock = Max(c.clock, msg.Timestamp) + 1
	c.mutex.Unlock()

	if msg.Receiver == receiver {
		fmt.Printf("Process n°%d received msg: %v, clock: %d\n", receiver, msg, c.clock)
	} else {
		fmt.Printf("Process n°%d: rejected message\n", receiver)
	}
}

// SendLetterBox envoi un message dans la boite aux lettres
func (p *Process) SendLetterBox(payload interface{}) {
	fmt.Println("Process n°", p.id, "sending msg in letterbox:", payload)
	p.com.mutex.Lock()
	p.com.letterBox <- payload // On depose le message
	p.com.mutex.Unlock()
}

// ReadLetterBox lit le contenu de la boite aux lettres
func (p *Process) ReadLetterBox() {
	p.com.mutex.Lock()
	select {
	case msg := <-p.com.letterBox:
		fmt.Printf("Process n°%d received msg in letterbox: %v\n", p.id, msg)
		break
	default:
		break
	}
	p.com.mutex.Unlock()
}

// IncClock incrémente la clock
func (p *Process) IncClock() {
	p.com.mutex.Lock()
	p.com.clock++
	p.com.mutex.Unlock()
}

// Reader est une fonctione exécutée sur une goroutine (thread) qui vérifie continuellement le bus,
// la letterbox et les demandes de synchronisation
func (p *Process) Reader() {
	for {
		select {
		case msg := <-*(p.com.bus):
			if m, ok := msg.(BroadcastMessage); ok {
				if !m.Synchrone {
					p.SendLetterBox(m)
				} else {
					p.com.OnBroadcast(p.id, m)
				}
			} else if m, ok := msg.(DedicatedMessage); ok {
				if !m.Synchrone {
					if m.Receiver == p.id {
						p.SendLetterBox(m)
					}
				} else {
					p.com.OnDedicatedMessage(p.id, m)
				}
			}
			break
		case syncRequest := <-p.syncChan:
			if syncRequest != nil {
				p.sync(syncRequest)
			}
			break
		default:
			p.ReadLetterBox()
			break
		}

		time.Sleep(time.Millisecond * 100)
	}
}

// SyncAll crée un "waitgroup" avec un delta egal au nombre de processus
// et appelle ensuite la fonction sync pour chaque processus
func (p *Process) SyncAll() {
	fmt.Println("Synching all processes...")
	var wg sync.WaitGroup
	wg.Add(NbProcess)

	for _, p := range ProcessPool {
		p.syncChan <- &wg
	}
	wg.Wait() // on attend que tous les processus soient sync (waitgroup == 0)
	fmt.Println("All processes are synchronized")
}

func (p *Process) sync(wg *sync.WaitGroup) {
	wg.Done() // on décrémente le waitgroup
}

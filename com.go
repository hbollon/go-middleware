// BOLLON hugo / RODRIGUEZ Samuel
package main

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
	com      Com
	id       int
	syncChan chan *sync.WaitGroup
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

// onBroadcast est la fonction appelée lorsqu'un message de type BroadcastMessage
// est reçu
func (c *Com) onBroadcast(receiver int, msg BroadcastMessage) {
	c.mutex.Lock()
	c.clock = max(c.clock, msg.Timestamp) + 1
	c.mutex.Unlock()
	fmt.Printf("Process n°%d received msg: %v, clock: %d\n", receiver, msg, c.clock)
}

// broadcast envoi un message de type BroadcastMessage sur le bus
func (c *Com) broadcast(msg BroadcastMessage) {
	c.mutex.Lock()
	c.clock++
	c.mutex.Unlock()

	msg.Timestamp = c.clock
	for _, p := range ProcessPool {
		p.sendLetterBox(msg)
	}
}

// sendTo envoi un message de type DedicatedMessage sur le bus
func (c *Com) sendTo(msg DedicatedMessage) {
	c.mutex.Lock()
	c.clock++
	c.mutex.Unlock()

	msg.Timestamp = c.clock
	ProcessPool[msg.Receiver].sendLetterBox(msg)
}

// broadcast envoi un message de type BroadcastMessage sur le bus
func (c *Com) broadcastSync(msg BroadcastMessage) {
	c.mutex.Lock()
	c.clock++
	c.mutex.Unlock()

	msg.Timestamp = c.clock
	for i := 0; i < NbProcess; i++ {
		c.bus.Send(msg)
	}
}

// sendTo envoi un message de type DedicatedMessage sur le bus
func (c *Com) sendToSync(msg DedicatedMessage) {
	c.mutex.Lock()
	c.clock++
	c.mutex.Unlock()

	msg.Timestamp = c.clock
	for i := 0; i < NbProcess; i++ {
		c.bus.Send(msg)
	}
}

// onDedicatedMessage est la fonction appelée lorsqu'un message de type DedicatedMessage
// est reçu
func (c *Com) onDedicatedMessage(receiver int, msg DedicatedMessage) {
	c.mutex.Lock()
	c.clock = max(c.clock, msg.Timestamp) + 1
	c.mutex.Unlock()

	if msg.Receiver == receiver {
		fmt.Printf("Process n°%d received msg: %v, clock: %d\n", receiver, msg, c.clock)
	} else {
		fmt.Printf("Process n°%d: rejected message\n", receiver)
	}
}

func (p *Process) sendLetterBox(payload interface{}) {
	fmt.Println("Process n°", p.id, "sending msg in letterbox:", payload)
	p.com.mutex.Lock()
	p.com.letterBox <- payload // On depose le message
	p.com.mutex.Unlock()
}

// readLetterBox lit le contenu de la boite aux lettres
func (p *Process) readLetterBox() {
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

// reader est une fonctione exécutée sur une goroutine (thread) qui lit constamment le bus
func (p *Process) reader() {
	for {
		select {
		case msg := <-*(p.com.bus):
			if m, ok := msg.(BroadcastMessage); ok {
				if !m.Synchrone {
					p.sendLetterBox(m)
				} else {
					p.com.onBroadcast(p.id, m)
				}
			} else if m, ok := msg.(DedicatedMessage); ok {
				if !m.Synchrone {
					if m.Receiver == p.id {
						p.sendLetterBox(m)
					}
				} else {
					p.com.onDedicatedMessage(p.id, m)
				}
			}
			break
		case syncRequest := <-p.syncChan:
			if syncRequest != nil {
				p.sync(syncRequest)
			}
			break
		default:
			p.readLetterBox()
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
	wg.Wait() // on attend que tous les processus soient sync
	fmt.Println("All processes are synchronized")
}

func (p *Process) sync(wg *sync.WaitGroup) {
	wg.Done() // on décrémente le waitgroup
}

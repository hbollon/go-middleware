// BOLLON hugo / RODRIGUEZ Samuel
package main

import (
	"fmt"
	"sync"
	"time"
)

// LetterBox est une liste de messages à partir d'un tableau d'interface
// (pseudo généricité)
type LetterBox []interface{}

// Process est un processus comportant une instance de Com et un id
type Process struct {
	com Com
	id  int
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
	for i := 0; i < NbProcess; i++ {
		c.bus.Send(msg)
	}
}

// sendTo envoi un message de type DedicatedMessage sur le bus
func (c *Com) sendTo(msg DedicatedMessage) {
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

// readLetterBox lit le contenu de la boite aux lettres
func (p *Process) readLetterBox() {
	for _, msg := range p.com.letterBox {
		fmt.Printf("Process n°%d received msg: %v\n", p.id, msg)
	}
	p.com.letterBox = nil
}

// incClock incrémente la clock
func (p *Process) incClock() {
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
					p.com.letterBox = append(p.com.letterBox, m)
				} else {
					p.com.onBroadcast(p.id, m)
				}
			} else if m, ok := msg.(DedicatedMessage); ok {
				if !m.Synchrone {
					if m.Receiver == p.id {
						p.com.letterBox = append(p.com.letterBox, m)
					}
				} else {
					p.com.onDedicatedMessage(p.id, m)
				}
			}
		default:
			break
		}

		time.Sleep(time.Millisecond * 100)
	}
}

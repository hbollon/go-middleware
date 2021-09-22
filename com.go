package main

import (
	"fmt"
	"sync"
	"time"
)

type Process struct {
	com Com
	id  int
}

type Com struct {
	clock int
	mutex sync.Mutex
	bus   *Bus
}

func (c *Com) onBroadcast(receiver int, msg BroadcastMessage) {
	c.mutex.Lock()
	c.clock = max(c.clock, msg.Timestamp) + 1
	c.mutex.Unlock()
	fmt.Printf("Process n°%d received msg: %v, clock: %d\n", receiver, msg, c.clock)
}

func (c *Com) broadcast(msg BroadcastMessage) {
	c.mutex.Lock()
	c.clock++
	c.mutex.Unlock()

	msg.Timestamp = c.clock
	for i := 0; i < NbProcess; i++ {
		c.bus.Send(msg)
	}
}

func (c *Com) sendTo(msg DedicatedMessage, dest int) {
	c.mutex.Lock()
	c.clock++
	c.mutex.Unlock()

	msg.Sender = dest
	msg.Timestamp = c.clock
	for i := 0; i < NbProcess; i++ {
		c.bus.Send(msg)
	}
}

func (c *Com) onDedicatedMessage(msg DedicatedMessage) {
	c.mutex.Lock()
	c.clock = max(c.clock, msg.Timestamp) + 1
	c.mutex.Unlock()
}

func (p *Process) reader() {
	for {
		select {
		case msg := <-*(p.com.bus):
			if m, ok := msg.(BroadcastMessage); ok {
				p.com.onBroadcast(p.id, m)
			} else if m, ok := msg.(DedicatedMessage); ok {
				p.com.onDedicatedMessage(m)
			}
		default:
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

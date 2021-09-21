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

func (c *Com) onBroadcast(receiver int, msg Message) {
	c.mutex.Lock()
	c.clock = max(c.clock, msg.Timestamp) + 1
	c.mutex.Unlock()
	fmt.Printf("Process n°%d received msg: %v\n", receiver, msg)
}

func (c *Com) broadcast(msg Message) {
	c.mutex.Lock()
	c.clock++
	c.mutex.Unlock()

	for i := 0; i < NbProcess; i++ {
		c.bus.Send(msg)
	}
}

func (p *Process) reader() {
	for {
		select {
		case msg := <-*(p.com.bus):
			p.com.onBroadcast(p.id, msg)
		default:
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

}

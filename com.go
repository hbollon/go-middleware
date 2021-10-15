package middleware

// BOLLON hugo / RODRIGUEZ Samuel

import (
	"fmt"
	"sync"
	"time"
)

var (
	ProcessPool []*Process
	NbProcess   int
)

// LetterBox est une liste de messages à partir d'un tableau d'interface
// (pseudo généricité)
type LetterBox chan interface{}

// Process est un processus comportant une instance de Com et un id
type Process struct {
	Com       Com
	Id        int
	SyncChan  chan *sync.WaitGroup
	TokenChan chan Token
	State     StateSectionCritique
}

// Com est notre interface de communication
// Il contient un pointeur vers le bus de communication, un mutex,
// un tableau de messages (boite aux lettres) et une clock
type Com struct {
	Clock     int
	Mutex     sync.Mutex
	Bus       *Bus
	LetterBox LetterBox
}

// InitProcessPool initialise le tableau de processus, le bus et lance les différents threads
func InitProcessPool(nbProcess int) {
	NbProcess = nbProcess

	// Instanciation du bus de communication
	bus := CreateBus()

	// Initialisation de la pool de process
	ProcessPool = make([]*Process, NbProcess)

	// Instanciation du pool de process
	for i := 0; i < NbProcess; i++ {
		process := Process{
			Id:        i,
			SyncChan:  make(chan *sync.WaitGroup, 1),
			TokenChan: make(chan Token, 1),
			Com: Com{
				Bus:       &bus,
				Clock:     0,
				Mutex:     sync.Mutex{},
				LetterBox: make(LetterBox, 10),
			},
		}
		// Lancement d'une goroutine exécutant la fonction reader pour chaque process
		go process.Reader()

		// Lancement d'une goroutine exécutant la fonction onToken pour chaque process
		go process.OnToken()

		ProcessPool[i] = &process
	}
}

// OnBroadcast est la fonction appelée lorsqu'un message de type BroadcastMessage
// est reçu
func (c *Com) OnBroadcast(receiver int, msg BroadcastMessage) {
	c.Mutex.Lock()
	c.Clock = Max(c.Clock, msg.Timestamp) + 1
	c.Mutex.Unlock()
	fmt.Printf("Process n°%d received msg: %v, clock: %d\n", receiver, msg, c.Clock)
}

// Broadcast envoi un message asynchrone de type BroadcastMessage directement dans la letterbox
func (c *Com) Broadcast(msg BroadcastMessage) {
	c.Mutex.Lock()
	c.Clock++
	c.Mutex.Unlock()

	msg.Timestamp = c.Clock
	for _, p := range ProcessPool {
		p.SendLetterBox(msg)
	}
}

// SendTo envoi un message asynchrone de type DedicatedMessage dans la letterbox du receiver
func (c *Com) SendTo(msg DedicatedMessage) {
	c.Mutex.Lock()
	c.Clock++
	c.Mutex.Unlock()

	msg.Timestamp = c.Clock
	ProcessPool[msg.Receiver].SendLetterBox(msg)
}

// BroadcastSync envoi un message synchrone de type BroadcastMessage sur le bus
func (c *Com) BroadcastSync(msg BroadcastMessage) {
	c.Mutex.Lock()
	c.Clock++
	c.Mutex.Unlock()

	msg.Timestamp = c.Clock
	for i := 0; i < NbProcess; i++ {
		c.Bus.Send(msg)
	}
}

// SendToSync envoi un message synchrone de type DedicatedMessage sur le bus
func (c *Com) SendToSync(msg DedicatedMessage) {
	c.Mutex.Lock()
	c.Clock++
	c.Mutex.Unlock()

	msg.Timestamp = c.Clock
	for i := 0; i < NbProcess; i++ {
		c.Bus.Send(msg)
	}
}

// OnDedicatedMessage est la fonction appelée lorsqu'un message de type DedicatedMessage
// est reçu
func (c *Com) OnDedicatedMessage(receiver int, msg DedicatedMessage) {
	c.Mutex.Lock()
	c.Clock = Max(c.Clock, msg.Timestamp) + 1
	c.Mutex.Unlock()

	if msg.Receiver == receiver {
		fmt.Printf("Process n°%d received msg: %v, clock: %d\n", receiver, msg, c.Clock)
	} else {
		fmt.Printf("Process n°%d: rejected message\n", receiver)
	}
}

// SendLetterBox envoi un message dans la boite aux lettres
func (p *Process) SendLetterBox(payload interface{}) {
	fmt.Println("Process n°", p.Id, "sending msg in letterbox:", payload)
	p.Com.Mutex.Lock()
	p.Com.LetterBox <- payload // On depose le message
	p.Com.Mutex.Unlock()
}

// ReadLetterBox lit le contenu de la boite aux lettres
func (p *Process) ReadLetterBox() {
	p.Com.Mutex.Lock()
	select {
	case msg := <-p.Com.LetterBox:
		fmt.Printf("Process n°%d received msg in letterbox: %v\n", p.Id, msg)
		break
	default:
		break
	}
	p.Com.Mutex.Unlock()
}

// IncClock incrémente la clock
func (p *Process) IncClock() {
	p.Com.Mutex.Lock()
	p.Com.Clock++
	p.Com.Mutex.Unlock()
}

// Reader est une fonctione exécutée sur une goroutine (thread) qui vérifie continuellement le bus,
// la letterbox et les demandes de synchronisation
func (p *Process) Reader() {
	for {
		select {
		case msg := <-*(p.Com.Bus):
			if m, ok := msg.(BroadcastMessage); ok {
				if !m.Synchrone {
					p.SendLetterBox(m)
				} else {
					p.Com.OnBroadcast(p.Id, m)
				}
			} else if m, ok := msg.(DedicatedMessage); ok {
				if !m.Synchrone {
					if m.Receiver == p.Id {
						p.SendLetterBox(m)
					}
				} else {
					p.Com.OnDedicatedMessage(p.Id, m)
				}
			}
			break
		case syncRequest := <-p.SyncChan:
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
		p.SyncChan <- &wg
	}
	wg.Wait() // on attend que tous les processus soient sync (waitgroup == 0)
	fmt.Println("All processes are synchronized")
}

func (p *Process) sync(wg *sync.WaitGroup) {
	wg.Done() // on décrémente le waitgroup
}

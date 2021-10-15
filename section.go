// BOLLON hugo / RODRIGUEZ Samuel
package main

import "time"

type Token struct {
	processID int
	clock     int
}

type StateSectionCritique int
type SCCallback func(...interface{})

const (
	StateSectionCritique_Released StateSectionCritique = iota
	StateSectionCritique_Requested
	StateSectionCritique_Critical_Section
)

func (p *Process) SendToken() {
	p.IncClock()
	token := Token{
		clock:     p.com.clock,
		processID: (p.id + 1) % NbProcess,
	}

	ProcessPool[token.processID].tokenChan <- token
}

func (p *Process) OnToken() {
	for {
		select {
		case token := <-p.tokenChan:
			if token.processID == p.id {
				p.com.clock = max(p.com.clock, token.clock) + 1
				if p.state == StateSectionCritique_Requested {
					p.state = StateSectionCritique_Critical_Section
					for p.state != StateSectionCritique_Released {
						time.Sleep(time.Millisecond * 10)
					}
				}
				p.SendToken()
				p.state = StateSectionCritique_Released
			}
			break
		default:
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (p *Process) RequestCriticalSection(callback SCCallback, args ...interface{}) {
	p.state = StateSectionCritique_Requested
	defer p.ReleaseCriticalSection() // on diffère la libération de la CS (éxécuté automatiquement en fin de programme ou de sortie de fonction)

	for p.state != StateSectionCritique_Critical_Section {
		time.Sleep(time.Millisecond * 10)
	}

	if p.state == StateSectionCritique_Critical_Section {
		callback(args)
	}

}

func (p *Process) ReleaseCriticalSection() {
	p.state = StateSectionCritique_Released
}

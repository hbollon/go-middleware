package middleware

// BOLLON hugo / RODRIGUEZ Samuel

import "time"

// Token struct
type Token struct {
	processID int
	clock     int
}

// Type énuméré pour definir les différents states
type StateSectionCritique int

// SCCallback est une fonction de callback pour la section critique acceptant
// un nombre d'arguments illimité et de n'importe quel type
type SCCallback func(...interface{})

const (
	StateSectionCritique_Released StateSectionCritique = iota
	StateSectionCritique_Requested
	StateSectionCritique_Critical_Section
)

// SendToken incrémente l'horloge et envoie un token au process suivant
func (p *Process) SendToken() {
	p.IncClock()
	token := Token{
		clock:     p.Com.Clock,
		processID: (p.Id + 1) % NbProcess,
	}

	ProcessPool[token.processID].TokenChan <- token
}

// OnToken est une routine exécuté sur une goroutine (thread) pour chaque process
// Elle vérifie la présence d'un token dans le channel et si oui, elle passe le status du process à Critical_Section
// (si il l'a demandé), attends que le token soit libéré et, enfin, envoie le token au process suivant
func (p *Process) OnToken() {
	for {
		select {
		case token := <-p.TokenChan:
			if token.processID == p.Id {
				p.Com.Clock = Max(p.Com.Clock, token.clock) + 1
				if p.State == StateSectionCritique_Requested {
					p.State = StateSectionCritique_Critical_Section
					for p.State != StateSectionCritique_Released {
						time.Sleep(time.Millisecond * 10)
					}
				}
				p.SendToken()
				p.State = StateSectionCritique_Released
			}
			break
		default:
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

// RequestCriticalSection demande la section critique et execute le callback une fois obtenue
func (p *Process) RequestCriticalSection(callback SCCallback, args ...interface{}) {
	p.State = StateSectionCritique_Requested
	defer p.ReleaseCriticalSection() // on diffère la libération de la CS (éxécuté automatiquement en fin de programme ou de sortie de fonction)

	for p.State != StateSectionCritique_Critical_Section {
		time.Sleep(time.Millisecond * 10)
	}

	if p.State == StateSectionCritique_Critical_Section {
		callback(args)
	}

}

// ReleaseCriticalSection libère la section critique
func (p *Process) ReleaseCriticalSection() {
	p.State = StateSectionCritique_Released
}

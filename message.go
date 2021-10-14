// BOLLON hugo / RODRIGUEZ Samuel
package main

type Message struct {
	Sender    int
	Msg       interface{}
	Timestamp int

	// Ce paramètre a été introduit pour permettre un envoi asynchrone
	// (donc qui arrive dans les bal des autres processus) tout en passant par le système de bus général
	// et transféré par le reader de chaque processus.
	Synchrone bool
}

type BroadcastMessage struct {
	Message
}

type DedicatedMessage struct {
	Message
	Receiver int
}

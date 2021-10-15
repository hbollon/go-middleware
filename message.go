package main

// BOLLON hugo / RODRIGUEZ Samuel

// Message est la struct principale utilisé pour créer les autres types de messages
// et permettre une pseudo-généricité
type Message struct {
	Sender    int
	Msg       interface{}
	Timestamp int

	// Ce paramètre a été introduit pour permettre un envoi asynchrone
	// (donc qui arrive dans les bal des autres processus) tout en passant par le système de bus général
	// et transféré par le reader de chaque processus.
	Synchrone bool
}

// BroadcastMessage est un message de type broadcast (simple wrapper de Message)
type BroadcastMessage struct {
	Message
}

// DeliverMessage wrap Message et ajoute un receiver
type DedicatedMessage struct {
	Message
	Receiver int
}

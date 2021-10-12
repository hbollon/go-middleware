// BOLLON hugo / RODRIGUEZ Samuel
package main

type Message struct {
	Sender    int
	Msg       interface{}
	Timestamp int
	Synchrone bool
}

type BroadcastMessage struct {
	Message
}

type DedicatedMessage struct {
	Message
	Receiver int
}

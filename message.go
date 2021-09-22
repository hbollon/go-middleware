package main

type Message struct {
	Sender    int
	Msg       interface{}
	Timestamp int
}

type BroadcastMessage struct {
	Message
}

type DedicatedMessage struct {
	Message
	Receiver int
}

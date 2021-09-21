package main

type Bus chan Message

func CreateBus() Bus {
	return make(Bus, 10)
}

func (b *Bus) Send(msg Message) {
	*b <- msg
}

func (b *Bus) Recv() Message {
	return <-*b
}

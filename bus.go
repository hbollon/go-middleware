package main

type Bus chan interface{}

func CreateBus() Bus {
	return make(Bus, 10)
}

func (b *Bus) Send(msg interface{}) {
	*b <- msg
}

func (b *Bus) Recv() interface{} {
	return <-*b
}

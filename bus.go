package middleware

// BOLLON hugo / RODRIGUEZ Samuel

// Bus générique via un chanel acceptant n'importe quel type
type Bus chan interface{}

// Initialise un bus avec un buffer de 10
func CreateBus() Bus {
	return make(Bus, 10)
}

// Envoie un message sur le bus
func (b *Bus) Send(msg interface{}) {
	*b <- msg
}

// Reçoit un message sur le bus
func (b *Bus) Recv() interface{} {
	return <-*b
}

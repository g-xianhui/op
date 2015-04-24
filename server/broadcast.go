package main

const (
	_ = iota
	BC_SUBSCRIPTE
	BC_UNSUBSCRIPTE
	BC_PUSH
)

type BCCmd struct {
	cmd   int
	ud    interface{}
	reply chan interface{}
}

type Broadcast struct {
	c         chan *BCCmd
	nextid    uint32
	listeners map[uint32]chan Msg
}

func (bc *Broadcast) subscripte(listener chan Msg) uint32 {
	id := bc.nextid
	bc.listeners[id] = listener
	bc.nextid++
	return id
}

func (bc *Broadcast) unsubscripte(id uint32) {
	delete(bc.listeners, id)
}

func (bc *Broadcast) push(m Msg) {
	for _, v := range bc.listeners {
		v <- m
	}
}

func createBroadcast() *Broadcast {
	bc := &Broadcast{}
	bc.c = make(chan *BCCmd)
	bc.listeners = make(map[uint32]chan Msg)
	go func() {
		for m := range bc.c {
			switch m.cmd {
			case BC_SUBSCRIPTE:
				listener := m.ud.(chan Msg)
				id := bc.subscripte(listener)
				m.reply <- id
			case BC_UNSUBSCRIPTE:
				id := m.ud.(uint32)
				bc.unsubscripte(id)
			case BC_PUSH:
				msg := m.ud.(Msg)
				bc.push(msg)
			}
		}
	}()
	return bc
}

func subscripte(bc *Broadcast, listener chan Msg) uint32 {
	m := &BCCmd{cmd: BC_SUBSCRIPTE, ud: listener}
	m.reply = make(chan interface{})
	bc.c <- m
	id := (<-m.reply).(uint32)
	return id
}

func unsubscripte(bc *Broadcast, id uint32) {
	m := &BCCmd{cmd: BC_UNSUBSCRIPTE, ud: id}
	bc.c <- m
}

func broadcast(bc *Broadcast, msg Msg) {
	m := &BCCmd{cmd: BC_PUSH, ud: msg}
	bc.c <- m
}

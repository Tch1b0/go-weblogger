package chans

import "sync"

type SuperChannel[T any] struct {
	mutex sync.Mutex

	Sender    chan T
	Receivers []chan T
}

func (sc *SuperChannel[T]) Start() {
	for {
		value := <-sc.Sender

		sc.mutex.Lock()
		removeQueue := []chan T{}

		for _, recv := range sc.Receivers {
			select {
			case recv <- value:
				// value accepted
			default:
				removeQueue = append(removeQueue, recv)
			}
		}

		for _, recv := range removeQueue {
			sc.unsafeRemoveReceiver(recv)
		}

		sc.mutex.Unlock()
	}
}

func (sc *SuperChannel[T]) AddReceiver(recv chan T) {
	sc.mutex.Lock()
	sc.Receivers = append(sc.Receivers, recv)
	sc.mutex.Unlock()
}

func (sc *SuperChannel[T]) unsafeRemoveReceiver(recv chan T) {
	newReceivers := make([]chan T, len(sc.Receivers)-1)

	i := 0
	for _, el := range sc.Receivers {
		if el != recv {
			newReceivers[i] = el
			i++
		}
	}

	sc.Receivers = newReceivers
}

func (sc *SuperChannel[T]) RemoveReceiver(recv chan T) {
	sc.mutex.Lock()
	sc.unsafeRemoveReceiver(recv)
	sc.mutex.Unlock()
}

func NewSuperChannel[T any]() *SuperChannel[T] {
	ch := new(SuperChannel[T])
	ch.Sender = make(chan T)

	return ch
}

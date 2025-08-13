package chans

import (
	"fmt"
	"sync"
	"time"
)

type DebouncedChan[T any] struct {
	data  []T
	mutex sync.Mutex

	Sender   chan T
	Receiver chan []T
}

func (db *DebouncedChan[T]) Flush() {
	db.mutex.Lock()
	db.Receiver <- db.data

	fmt.Println("This message should ONLY be displayed in the local terminal")

	db.data = []T{}
	db.mutex.Unlock()
}

func (db *DebouncedChan[T]) Debounce() {
	t := time.NewTicker(100 * time.Millisecond)

	for {
		select {
		case value := <-db.Sender:
			db.mutex.Lock()
			db.data = append(db.data, value)
			db.mutex.Unlock()
		case <-t.C:
			if len(db.data) > 0 {
				db.Flush()
			}
		}
	}
}

func NewDebouncedChan[T any]() *DebouncedChan[T] {
	ch := new(DebouncedChan[T])

	ch.Sender = make(chan T, 100)
	ch.Receiver = make(chan []T)

	return ch
}

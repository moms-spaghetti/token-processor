package main

import (
	"errors"
	"log"
	"time"
)

type tokenStore struct {
	ch chan struct{}
}

func newTokenStore(lim int) *tokenStore {
	ch := make(chan struct{}, lim)
	for i := 0; i < lim; i++ {
		ch <- struct{}{}
	}

	return &tokenStore{
		ch: ch,
	}
}

func (ts *tokenStore) jobProcessor(f func()) error {
	select {
	case <-ts.ch:
		log.Print("token available: new worker")
		f()
		ts.ch <- struct{}{}
		return nil
	default:
		return errors.New("no tokens remaining")
	}
}

func doWork() {
	time.Sleep(10 * time.Second)
	log.Print("worker: work done")
}

func main() {
	ts := newTokenStore(3)

	for range time.Tick(1 * time.Second) {
		go func() {
			if err := ts.jobProcessor(doWork); err != nil {
				log.Print("jobProcessor error: ", err)
			}
		}()
	}
}

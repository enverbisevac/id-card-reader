package main

import (
	"fmt"
	"github.com/ebfe/scard"
	"sync"
	"time"
)

type Reader interface {
	Release() error
	OnCardInserted(device string) error
	OnCardRemoved()
}

type readerDevice struct {
	sync.RWMutex
	*scard.Context
	*scard.Card
	readers []string
}

func NewReaderDevice() *readerDevice {
	// Establish a context
	ctx, err := scard.EstablishContext()
	if err != nil {
		panic(err)
	}

	newObject := &readerDevice{
		Context: ctx,
		readers: make([]string, 0),
	}

	newObject.readers, err = newObject.Context.ListReaders()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d readers:\n", len(newObject.readers))
	for i, reader := range newObject.readers {
		fmt.Printf("[%d] %s\n", i, reader)
	}

	go newObject.listReaders()
	if len(newObject.readers) > 0 {
		go newObject.checkCardStatus()
	}

	return newObject
}

func (r *readerDevice) listReaders() {
	var err error
	for {
		r.RLock()
		r.readers, err = r.Context.ListReaders()
		if err != nil {
			panic(err)
		}
		r.RUnlock()
		fmt.Println("Checking for new reader device")
		time.Sleep(5 * time.Second)
	}
}

func (r *readerDevice) OnCardInserted(device string) error {
	var err error
	r.Card, err = r.Connect(device, scard.ShareExclusive, scard.ProtocolAny)
	fmt.Println("Kartica ubacena")
	return err
}

func (r *readerDevice) OnCardRemoved() {
	fmt.Println("Ubaci karticu")
	r.Card = nil
}

func (r *readerDevice) Release() error {
	if r.Card != nil {
		err := r.Card.Disconnect(scard.ResetCard)
		if err != nil {
			return err
		}
	}
	return r.Context.Release()
}

func (r *readerDevice) checkCardStatus() {
	rs := make([]scard.ReaderState, len(r.readers))
	// initial
	for i := range rs {
		r.RLock()
		rs[i].Reader = r.readers[i]
		r.RUnlock()
		rs[i].CurrentState = scard.StateUnaware
	}

	initial := true
	for {
		index := -1
		for i := range rs {
			if rs[i].EventState&scard.StatePresent != 0 {
				index = i
				r.Lock()
				r.OnCardInserted(r.readers[i])
				r.Unlock()
			}
			rs[i].CurrentState = rs[i].EventState
		}

		if !initial && index == -1 {
			r.OnCardRemoved()
		}

		err := r.Context.GetStatusChange(rs, -1)
		if err != nil {
			index = -1
			break
		}
		initial = false
	}
}

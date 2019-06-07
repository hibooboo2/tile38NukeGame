package ui

import (
	"log"
	"time"
)

var (
	keyboardEventChans   []chan KeyboardEvent
	newKeyboardEventChan chan chan KeyboardEvent
	mainKeyBoardEvents   chan KeyboardEvent
)

func init() {
	mainKeyBoardEvents = make(chan KeyboardEvent, 2)
	newKeyboardEventChan = make(chan chan KeyboardEvent, 2)
	go startEventHandler()
}

func startEventHandler() {
	t := time.NewTicker(time.Millisecond * 1000)
	for {
		select {
		case evtChan := <-newKeyboardEventChan:
			keyboardEventChans = append(keyboardEventChans, evtChan)
		case evt := <-mainKeyBoardEvents:
			// log.Println("Got event!")
			for _, evtChan := range keyboardEventChans {
				select {
				case <-t.C:
					log.Println("Skipped key event:", evt.Key)
				case evtChan <- evt:
				}
			}
		}
	}
}

func GetKeyboardEvents() <-chan KeyboardEvent {
	evtChan := make(chan KeyboardEvent, 2)
	newKeyboardEventChan <- evtChan
	return evtChan
}

type KeyboardEvent struct {
	Key       KeyCode
	Modifiers uint8
	When      time.Time
	KeyType   KeyType
}

type KeyType uint8

const (
	None KeyType = iota
	KeyPress
	KeyUp
	KeyDown
)

package main

import (
	"github.com/stianeikeland/go-rpio"
	"time"
	//"fmt"
	//"strconv"
)

const (
	DATA_PIN  int = 27 //gpio 27, pin 13
	LATCH_PIN int = 22 //gpio 22, pin 15
	CLOCK_PIN int = 4  //gpio 4, pin 7
	OE_PIN    int = 17 //gpio 17, pin 11
)

func zoneController() {
	defer cleanup()

	for {
		select {
		case address := <-startChan:
			startZone(address)
			status.Watering = true
			status.Address = address
			broadcastStatus()
		case <-stopChan:
			allOff()
			status.Watering = false
			status.RunningProgram = false
			broadcastStatus()
		}
	}
}

func broadcastStatus() {
	for _, client := range clients {
		client.sendStatus()
	}
}

func runZone(timer *time.Timer) {
	//fmt.Println("runZone")
	select {
	case <-timer.C:
		stopChan <- true
	case <-quitChan:
		stopChan <- true
		timer.Stop()
	}
}

func startZone(address int) {
	//fmt.Println("startZone", address)
	for i := 0; i < 8; i++ {
		if i == address {
			rpio.Pin(DATA_PIN).High()
		} else {
			rpio.Pin(DATA_PIN).Low()
		}
		tick()
	}
	latch()
}

func initializeController() {
	err := rpio.Open()
	if err != nil {
		showError(err)
	}
	rpio.Pin(DATA_PIN).Output()
	rpio.Pin(DATA_PIN).Low()
	rpio.Pin(LATCH_PIN).Output()
	rpio.Pin(LATCH_PIN).Low()
	rpio.Pin(CLOCK_PIN).Output()
	rpio.Pin(CLOCK_PIN).Low()
	rpio.Pin(OE_PIN).Output()
	rpio.Pin(OE_PIN).Low()
}

func tick() {
	rpio.Pin(CLOCK_PIN).High()
	rpio.Pin(CLOCK_PIN).Low()
}

func latch() {
	rpio.Pin(LATCH_PIN).High()
	rpio.Pin(LATCH_PIN).Low()
}

func cleanup() {
	allOff()
	rpio.Close()
}

func allOff() {
	rpio.Pin(DATA_PIN).Low()
	for i := 0; i < 8; i++ {
		tick()
	}
	latch()
}

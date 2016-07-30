package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	//"github.com/robfig/cron" //https://godoc.org/github.com/robfig/cron
	//"bufio"
	//"github.com/gorilla/websocket"
	//"github.com/stianeikeland/go-rpio"
	//"os"
	//"strconv"
	//"strings"
	//"time"
)

var clients = make(map[string]Client)
var startChan = make(chan int)
var stopChan = make(chan bool)
var quitChan = make(chan bool)

func showError(err error) {
	fmt.Println(err.Error())
}

func main() {
	initZones()
	initializeController()

	go zoneController()
	stopChan <- true //turn everything off in cas something is stuck on

	initProgram()
	program.Enable()

	r := mux.NewRouter()

	r.HandleFunc("/socket", socketHandler)

	http.Handle("/", r)
	fmt.Println("listening")
	log.Fatal(http.ListenAndServe(":4000", r))

	// if len(os.Args) > 1 {

	// 	//if supplied with arg "cron" set up the cron job
	// 	if os.Args[1] == "cron" {
	// 		c := cron.New()
	// 		c.AddFunc("0 3 * * *", func() { initAndRunProgram() })
	// 		c.Start()
	// 	} else {
	// 		fmt.Println("arg should be 'cron' to initiate cron job")
	// 	}
	// } else {
	// 	//if no args, run the user interface
	// 	defer cleanup()
	// 	initialize()
	// 	initProgram()
	// 	askZone()
	// }
}

// func askZone() {
// 	fmt.Println("We've got these zones:")
// 	for i, _ := range zones {
// 		fmt.Println(zones[i].Num, zones[i].Name)
// 	}
// 	fmt.Println("Where would you like to water? ('bye' to quit, 'program' to run program)")
// 	z, cmd := getInput()
// 	if z < 0 {
// 		if cmd == "bye" {
// 			goodbye()
// 		} else if cmd == "program" {
// 			runProgram()
// 		} else {
// 			fmt.Println("I don't understand")
// 			askZone()
// 		}
// 	} else {
// 		fmt.Println("For how many seconds?")
// 		t, cmd := getInput()
// 		if t < 0 {
// 			if cmd == "bye" {
// 				goodbye()
// 			} else {
// 				fmt.Println("I don't understand")
// 				askZone()
// 			}
// 		} else {
// 			zoneOn(getZone(z), t)
// 			askZone()
// 		}
// 	}
// }

// func getInput() (int, string) {
// 	reader := bufio.NewReader(os.Stdin)
// 	input, _ := reader.ReadString('\n')
// 	input = strings.TrimSuffix(input, "\n")
// 	num, err := strconv.Atoi(input)
// 	if err == nil {
// 		return num, ""
// 	} else {
// 		return -1, input
// 	}
// }

// func goodbye() {
// 	fmt.Println("Goodbye")
// 	allOff()
// }

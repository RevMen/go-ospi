package main

import (
	"fmt"
	"github.com/robfig/cron" //https://godoc.org/github.com/robfig/cron
	"time"
)

var program = Program{}

func runCommand(index int) {
	fmt.Println("index:", index)
	if index == len(program.Commands) {
		stopChan <- true
	} else if program.Commands[index].Time < 1 {
		runCommand(index + 1)
	} else {
		timer := time.NewTimer(time.Second * time.Duration(program.Commands[index].Time))
		startChan <- program.Commands[index].Address
		select {
		case <-timer.C:
			runCommand(index + 1)
		case <-quitChan:
			stopChan <- true
			timer.Stop()
		}
	}
}

func startProgram() {
	fmt.Println("starting program")
	runCommand(0)
}

type Command struct {
	Address int
	Time    int
}

func initProgram() {
	program.Name = "Nightly"
	program.Cron = cron.New()

	program.Commands = make([]Command, 8)
	program.Commands[0] = Command{Address: 0, Time: 0}    //unused
	program.Commands[1] = Command{Address: 1, Time: 1500} //north side
	program.Commands[2] = Command{Address: 2, Time: 0}    // blocked
	program.Commands[3] = Command{Address: 3, Time: 1200} // back center
	program.Commands[4] = Command{Address: 4, Time: 1200} //back south
	program.Commands[5] = Command{Address: 5, Time: 1500} // south side
	program.Commands[6] = Command{Address: 6, Time: 1500} //front center
	program.Commands[7] = Command{Address: 7, Time: 0}    //unused
	program.Cron.AddFunc("0 0 6 * * *", func() { startProgram() })

	// program.Commands = make([]Command, 8)
	// program.Commands[0] = Command{Address: 0, Time: 0}  //unused
	// program.Commands[1] = Command{Address: 1, Time: 30} //north side
	// program.Commands[2] = Command{Address: 2, Time: 0}  // blocked
	// program.Commands[3] = Command{Address: 3, Time: 30} // back center
	// program.Commands[4] = Command{Address: 4, Time: 30} //back south
	// program.Commands[5] = Command{Address: 5, Time: 30} // south side
	// program.Commands[6] = Command{Address: 6, Time: 30} //front center
	// program.Commands[7] = Command{Address: 7, Time: 0}  //unused
	// program.Cron.AddFunc("0 39 * * * *", func() { startProgram() })

}

type Program struct {
	Name     string     `json:"name"`
	Commands []Command  `json:"commands"`
	Cron     *cron.Cron `json:"-"`
}

func (p *Program) Enable() {
	p.Cron.Start()
}

func (p *Program) Disable() {
	p.Cron.Stop()
}

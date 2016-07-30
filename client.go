package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var status = Status{Watering: false, RunningProgram: false}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		showError(err)
	}

	id := time.Now().Format(time.RFC3339)

	client := Client{ID: id, conn: conn, Send: make(chan Message)}

	clients[id] = client

	go client.writer()
	client.reader()
}

type Client struct {
	ID   string
	conn *websocket.Conn
	Send chan Message
}

func (c *Client) reader() {
	defer func() {
		c.conn.Close()
		close(c.Send)
		c.remove()
	}()

	for {
		msg := Message{}
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("read error on", c.ID)
			break
		}

		fmt.Println(msg)

		switch msg.Action {
		case "stop":
			c.stopNow()
		case "status":
			c.sendStatus()
		case "zones":
			c.sendZones()
		case "start":
			c.startNow(msg.Address)
		case "run":
			startChan <- msg.Address
			timer := time.NewTimer(time.Second * time.Duration(msg.Time))
			go runZone(timer)
		case "start_program":
			go startProgram()
		}
	}
}

func (c *Client) writer() {
	for {
		msg := <-c.Send
		err := c.conn.WriteJSON(msg)
		if err != nil {
			//showErr(err)
			break
		}
	}
}

func (c *Client) sendZones() {
	msg := Message{
		Action: "zones",
		Zones:  zones,
	}
	c.Send <- msg
}

func (c *Client) sendStatus() {
	msg := Message{
		Action: "status",
		Status: status,
	}
	c.Send <- msg
}

func (c *Client) startNow(address int) {
	startChan <- address
}

func (c *Client) runFor(address int, t int) {
	startChan <- address
	timer := time.NewTimer(time.Second * time.Duration(t))
	go runZone(timer)
}

func (c *Client) stopNow() {
	if status.Watering {
		quitChan <- true
	}
	stopChan <- true
}

func (c *Client) remove() {
	delete(clients, c.ID)
}

type Message struct {
	Action  string `json:"action"`
	Address int    `json:"address,omitempty"`
	Time    int    `json:"time,omitempty"`
	Status  Status `json:"status,omitempty"`
	Zones   []Zone `json:"zones,omitempty"`
}

type Status struct {
	Watering       bool `json:"watering"`
	RunningProgram bool `json:"running_program"`
	ProgramID      int  `json:"program_id,omitempty"`
	Address        int  `json:"address,omitempty"`
}

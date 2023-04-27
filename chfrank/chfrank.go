package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8889", "http service address")
var done = make(chan int)
var interaction_reader_co int = 0
var interaction_phase int = 0

func InteractionReader(c *websocket.Conn) {
	defer close(done)
	for interaction_reader_co == 0 {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		str_message := string(message)

		log.Println("recv: ", str_message)

	}

}

func main() {

	flag.Parse()
	log.SetFlags(0)

	var init_counter int = 0

	done_init := make(chan int, 1)

	var init_phase int = 0

	var init_phase_co int = 0

	var id string

	var pw string

	done_cred := make(chan int, 1)

	var cred_phase int = 0

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/enter"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	fmt.Println("ID: ")
	fmt.Scanln(&id)
	fmt.Println("PW: ")
	fmt.Scanln(&pw)

	go func() {
		defer close(done)
		for init_phase_co == 0 {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)

			str_message := string(message)

			if str_message != "ALIVE" && init_counter < 3 {

				log.Printf("recv: %s", "Terminate: Server Dead")
				return

			} else {
				init_counter += 1

			}

			if init_counter == 3 {

				close(done_init)

			}

		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for init_phase == 0 {
		select {
		case <-done:
			return
		case t := <-ticker.C:

			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}

		case <-done_init:

			init_phase = 1

		case <-interrupt:
			log.Println("interrupt")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
	select {
	case <-done:
	case <-time.After(time.Second):
	}

	init_phase_co = 1

	log.Println("Dial Successful")

	u = url.URL{Scheme: "ws", Host: *addr, Path: "/access"}
	log.Printf("connecting to %s", u.String())

	c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done = make(chan int)

	go func() {

		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		str_message := string(message)

		if str_message != "ALIVE" {

			log.Printf("recv: %s", "Terminate")
			return

		} else if str_message == "ALIVE" {
			log.Println("Auth Successful")
			close(done_cred)
			return
		} else {

			log.Println("Unkown Error")
			return

		}

	}()

	ticker = time.NewTicker(time.Second)
	defer ticker.Stop()

	credentials := id + ":" + pw

	var write_message string = ""

	for cred_phase == 0 {

		write_message = credentials

		select {
		case <-done:
			return
		case <-ticker.C:

			log.Println("Waiting Auth...")
			err := c.WriteMessage(websocket.TextMessage, []byte(write_message))
			if err != nil {
				log.Println("write:", err)
				return
			}

		case <-done_cred:

			cred_phase = 1

		}

	}

	go InteractionReader(c)

	for interaction_phase == 0 {

		write_message = "Tell me if I'm in"

		select {
		case <-done:
			return
		case <-ticker.C:

			err := c.WriteMessage(websocket.TextMessage, []byte(write_message))
			if err != nil {
				log.Println("write:", err)
				return
			}

		case <-interrupt:
			log.Println("interrupt")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

}

package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8889", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var testid string = "test"

var testpw string = "test"

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, []byte("ALIVE"))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func auth(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		str_message := string(message)

		ret_code := "DEAD"

		cred_list := strings.Split(str_message, ":")

		uid := cred_list[0]

		upw := cred_list[1]

		if uid == testid && upw == testpw {

			ret_code = "ALIVE"

		}

		err = c.WriteMessage(mt, []byte(ret_code))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/enter", echo)
	http.HandleFunc("/enter_cred", auth)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type user_data struct {
	DATA []user_record
}

type user_record struct {
	UUID    string
	USER_ID string
	USER_PW string
	ACTIVE  int
}

func dbQuery(query string, args []string) (interface{}, error) {

	db, err := sql.Open("mysql", "seantywork:youdonthavetoknow@tcp(127.0.0.1:3306)/chfrank")

	if err != nil {
		return -1, err
	}

	for i := 0; i < len(args); i++ {

		query = strings.Replace(query, "?", args[i], 1)

	}

	defer db.Close()

	results, err := db.Query(query)

	if err != nil {

		return 1, err

	}

	return results, err

}

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

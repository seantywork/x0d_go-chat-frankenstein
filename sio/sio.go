package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"

	"github.com/gorilla/sessions"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	//	"crypto/sha256"
	// "encoding/hex"
)

type user_data struct {
	DATA []user_record
}

type user_record struct {
	UUID    string
	USID    string
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

	store := sessions.NewFilesystemStore("./sio_session", []byte("sio-test-key"))

	session, _ := store.Get(r, "sio-test-name")

	session.Values["what"] = "thehell"

	_ = session.Save(r, w)

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

		//sha_256 := sha256.New()

		//sha_256.Write([]byte(upw))

		//pw_checksum_byte := sha_256.Sum(nil)

		//pw_checksum_hex := hex.EncodeToString(pw_checksum_byte)

		q := "SELECT uuid, usid, user_pw FROM chfrank_user WHERE ACTIVE = 1 AND user_id = '?'"

		a := []string{uid}

		res, err := dbQuery(q, a)

		if err != nil {

			err = c.WriteMessage(mt, []byte(ret_code))
			if err != nil {
				log.Println("write:", err)
				continue
			}

		}

		result_rows := res.(*sql.Rows)

		var ud user_data

		for result_rows.Next() {
			var ur user_record

			err = result_rows.Scan(&ur.UUID, &ur.USID, &ur.USER_PW)
			if err != nil {
				panic(err.Error())
			}

			ud.DATA = append(ud.DATA, ur)

		}

		if len(ud.DATA) != 1 {

			err = c.WriteMessage(mt, []byte(ret_code))
			if err != nil {
				log.Println("write:", err)
				continue
			}

		}

		//if ud.DATA[0].USER_PW == pw_checksum_hex {

		//	ret_code = "ALIVE"

		//}

		if ud.DATA[0].USER_PW == upw {

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

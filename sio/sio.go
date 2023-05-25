package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"github.com/gorilla/sessions"

	"database/sql"

	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	_ "github.com/go-sql-driver/mysql"
)

type UserData struct {
	DATA []UserRecord
}

type UserRecord struct {
	UUID    string
	USID    string
	USER_ID string
	USER_PW string
	ACTIVE  int
}

var ADDR = flag.String("addr", "localhost:8889", "http service address")

var UPGRADER = websocket.Upgrader{} // use default options

var UID_CONNECTION = make(map[string]*websocket.Conn)

var SRV_CONNECTION = make(map[string]*websocket.Conn)

var L = log.New(os.Stdout, "", 0)

func eventLogger(msg string) {
	L.SetPrefix(time.Now().UTC().Format("2006-01-02 15:04:05.000") + " [INFO] ")
	L.Print(msg)
}

var DB *sql.DB

func dbQuery(query string, args []any) (*sql.Rows, error) {

	var empty_rows *sql.Rows

	results, err := DB.Query(query, args[0:]...)

	if err != nil {

		return empty_rows, err

	}

	return results, err

}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func usidChecker(session *sessions.Session) (string, int) {

	var res *sql.Rows
	var err error

	session_val := session.Values["usid"]

	str_session_val, okay := session_val.(string)

	if okay == false {

		return "", 0

	}

	q := "SELECT uuid, usid, user_pw FROM chfrank_user WHERE ACTIVE = 1 AND usid = ?"

	a := []any{str_session_val}

	res, err = dbQuery(q, a)

	if err != nil {

		return "", -2

	}

	var ud UserData

	defer res.Close()

	for res.Next() {
		var ur UserRecord

		err = res.Scan(&ur.UUID, &ur.USID, &ur.USER_PW)
		if err != nil {
			panic(err.Error())
		}

		ud.DATA = append(ud.DATA, ur)

	}

	if len(ud.DATA) != 1 {

		return "", -1

	}

	return ud.DATA[0].UUID, 1

}

func Auth(w http.ResponseWriter, r *http.Request, c *websocket.Conn, message []byte) int {

	store := sessions.NewFilesystemStore("./sio_session", []byte("sio-test-key"))

	session, _ := store.Get(r, "sio-test-name")

	uuid_key, usid_code := usidChecker(session)

	log.Println("auth session: ", usid_code)

	if usid_code == 1 {

		if old_c, okay := UID_CONNECTION[uuid_key]; okay {

			_ = old_c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Old Connection Close"))

			UID_CONNECTION[uuid_key] = c

		} else {

			UID_CONNECTION[uuid_key] = c

		}

		return usid_code

	} else if usid_code == -2 {

		return usid_code
	}

	str_message := string(message)

	log.Println(str_message)

	cred_list := strings.Split(str_message, ":")

	if len(cred_list) != 2 {

		return -1
	}

	uid := cred_list[0]

	upw := cred_list[1]

	sha_256 := sha256.New()

	sha_256.Write([]byte(upw))

	pw_checksum_byte := sha_256.Sum(nil)

	pw_checksum_hex := hex.EncodeToString(pw_checksum_byte)

	q := "SELECT uuid, usid, user_pw FROM chfrank_user WHERE ACTIVE = 1 AND user_id = ?"

	a := []any{uid}

	res, err := dbQuery(q, a)

	if err != nil {

		return -2

	}

	var ud UserData

	for res.Next() {
		var ur UserRecord

		err = res.Scan(&ur.UUID, &ur.USID, &ur.USER_PW)
		if err != nil {
			panic(err.Error())
		}

		ud.DATA = append(ud.DATA, ur)

	}

	if len(ud.DATA) != 1 {

		return -1

	}

	if ud.DATA[0].USER_PW == pw_checksum_hex {

		new_session_val, _ := randomHex(16)

		q := "UPDATE chfrank_user SET usid = ? WHERE ACTIVE = 1 AND uuid = ?"

		a := []any{new_session_val, ud.DATA[0].UUID}

		_, err = dbQuery(q, a)

		if err != nil {

			return -2

		}

		uuid_key = ud.DATA[0].UUID

		session.Values["usid"] = new_session_val

		_ = session.Save(r, w)

		if old_c, okay := UID_CONNECTION[uuid_key]; okay {

			_ = old_c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Old Connection Close"))

			UID_CONNECTION[uuid_key] = c

		} else {

			UID_CONNECTION[uuid_key] = c

		}

		return 1

	}

	return -1

}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := UPGRADER.Upgrade(w, r, nil)
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

func access(w http.ResponseWriter, r *http.Request) {

	log.Println("Access")

	c, err := UPGRADER.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		auth_code := Auth(w, r, c, message)

		if auth_code != 1 {

			err = c.WriteMessage(mt, []byte("DEAD"))
			if err != nil {
				log.Println("write:", err)
				return
			}

			log.Println("Deny")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Connection Close"))
			if err != nil {
				log.Println("write close:", err)
				return
			}

			return

		} else {
			err = c.WriteMessage(mt, []byte("ALIVE"))
			if err != nil {
				log.Println("write:", err)
				return
			}
			break
		}

	}

	log.Println("Accept")

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		log.Println("Accept")
		err = c.WriteMessage(mt, []byte("Do whatever you desire"))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}

}

func srv_message(w http.ResponseWriter, r *http.Request) {

	log.Println("Server Access")
	log.Printf("UID_CONN : %d", len(UID_CONNECTION))
	log.Printf("SRV_CONN : %d", len(SRV_CONNECTION))

	UPGRADER.CheckOrigin = func(r *http.Request) bool { return true }

	c, err_up := UPGRADER.Upgrade(w, r, nil)
	if err_up != nil {
		log.Print("upgrade:", err_up)
		return
	}

	defer c.Close()

	var mt int
	var message []byte
	var err error

	for {
		mt, message, err = c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		str_message := string(message)

		if str_message != "imtheserver" {

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Connection Close"))
			if err != nil {
				log.Println("write close:", err)
				return
			}

			return

		} else {
			break
		}

	}

	log.Println("Validation Success")

	err = c.WriteMessage(mt, []byte("READY"))
	if err != nil {
		log.Println("write:", err)
		return
	}

	for {
		mt_read, message_read, err_read := c.ReadMessage()
		if err_read != nil {
			log.Println("read:", err_read)
			return
		}
		log.Printf("recv: %s", message_read)

		str_message := string(message_read)

		if user_c, okay := UID_CONNECTION[str_message]; okay {

			if _, okay := SRV_CONNECTION[str_message]; !okay {
				SRV_CONNECTION[str_message] = c
				log.Println("Connection added")
			} else if srv_c, okay := SRV_CONNECTION[str_message]; okay {

				if c != srv_c {
					err = srv_c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Connection Close"))
					if err != nil {
						log.Println("write close:", err)
						return
					}
					SRV_CONNECTION[str_message] = c
					log.Println("Connection overwritten")

				} else {

					log.Println("Connection persistent")

				}

			}

			err_read = user_c.WriteMessage(mt_read, []byte("Server is saying hello"))
			if err_read != nil {
				log.Println("write:", err_read)
				return
			}

			err_read = c.WriteMessage(mt_read, []byte("SUCCESS"))
			if err != nil {
				log.Println("write:", err)
				return
			}

		} else {

			err_read = c.WriteMessage(mt_read, []byte("NOSESSION"))
			if err != nil {
				log.Println("write:", err)
				return
			}

		}

	}

}

func main() {

	var err error

	DB, err = sql.Open("mysql", "seantywork:youdonthavetoknow@tcp(127.0.0.1:3306)/chfrank")

	if err != nil {

		eventLogger("DB connection failed")
		return

	}

	DB.SetConnMaxLifetime(time.Second * 10)
	DB.SetConnMaxIdleTime(time.Second * 5)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/enter", echo)
	http.HandleFunc("/access", access)
	http.HandleFunc("/serverside-message", srv_message)
	log.Fatal(http.ListenAndServe(*ADDR, nil))
}

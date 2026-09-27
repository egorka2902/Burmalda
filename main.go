package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
)

type Us struct {
	Login    string `json:"Login"`
	Password string `json:"Password"`
	SecKey   string `json:"SecKey"`
}

type ChatMessage struct {
	From    string `json:"From"`
	To      string `json:"To"`
	Content string `json:"Content"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var clients = make(map[string]*websocket.Conn)
var mu sync.RWMutex

var login string = ""

func AuthFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(200)
		return
	}
	u := Us{}
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &u)

	dburl := "postgres://zov:123@localhost:5432/burmalda?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), dburl)
	if err != nil {
		w.WriteHeader(502)
		return
	}
	defer conn.Close(context.Background())
	var ss string
	var log string
	row := conn.QueryRow(context.Background(), "SELECT seckey FROM base.users  WHERE login = $1", u.Login).Scan(&log)
	if row != nil {
		_, e := conn.Exec(context.Background(), "INSERT INTO base.users (login,password,seckey) VALUES ($1, $2,$3)", u.Login, u.Password, u.SecKey)
		if e != nil {
			fmt.Println(e)
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(u.SecKey))
		return
	}
	rww := conn.QueryRow(context.Background(), " SELECT seckey FROM base.users WHERE login = $1 AND password = $2", u.Login, u.Password).Scan(&ss)
	if rww != nil {
		w.WriteHeader(504)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(ss))
}

func AkeyFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(200)
		return
	}
	dburl := "postgres://zov:123@localhost:5432/burmalda?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), dburl)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	u := Us{}
	body, e := io.ReadAll(r.Body)
	if e != nil {
		w.WriteHeader(409)
		return
	}
	json.Unmarshal(body, &u)
	var s string
	row := conn.QueryRow(context.Background(), "SELECT login FROM base.users WHERE seckey = $1", u.SecKey).Scan(&s)
	if row != nil {
		w.WriteHeader(http.StatusFailedDependency)

		return
	}
	login = s
	w.WriteHeader(200)
	w.Write([]byte(s))
	return
}

func WebSocketChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(200)
		return
	}

	dburl := "postgres://zov:123@localhost:5432/burmalda?sslmode=disable"
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(504)
		return
	}
	defer conn.Close()
	var log string
	token := r.URL.Query().Get("key")
	dbcon, err := pgx.Connect(context.Background(), dburl)
	if err != nil {
		w.WriteHeader(501)
		return
	}
	err = dbcon.QueryRow(context.Background(), "SELECT login FROM base.users WHERE seckey = $1", token).Scan(&log)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	dbcon.Close(context.Background())
	fmt.Println(log)

	mu.Lock()
	clients[log] = conn
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(clients, log)
		mu.Unlock()
	}()
	for {
		var msg ChatMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			conn.WriteJSON(map[string]string{"error": "failed to ReadJSON"})
			break
		}
		fmt.Println(msg.Content)
		con, err := pgx.Connect(context.Background(), dburl)
		if err != nil {
			w.WriteHeader(503)
			return
		}
		msg.From = log
		_, err = con.Exec(context.Background(), "INSERT INTO base.messages (from_user, to_user, content) VALUES ($1,$2,$3)", msg.From, msg.To, msg.Content)
		if err != nil {
			w.WriteHeader(502)
			return
		}
		con.Close(context.Background())
		mu.RLock()
		rec, online := clients[msg.To]
		mu.RUnlock()
		if online {
			err := rec.WriteJSON(msg)
			if err != nil {
				w.WriteHeader(404)
				return
			}
		}
	}
}

func ChangeData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(200)
		return
	}
	dburl := "postgres://zov:123@localhost:5432/burmalda?sslmode=disable"
	conn, e := pgx.Connect(context.Background(), dburl)
	if e != nil {
		w.WriteHeader(409)
		return
	}
	u := Us{}
	body, er := io.ReadAll(r.Body)
	if er != nil {
		w.WriteHeader(502)
		return
	}
	json.Unmarshal(body, &u)
	_, errr := conn.Exec(context.Background(), "UPDATE base.users SET password = $1, seckey = $2 WHERE (login = $3 AND password = $1) OR (login = $3 AND seckey = $2)", u.Password, u.SecKey, u.Login)
	if errr != nil {
		fmt.Println(errr)
		w.WriteHeader(424)
		return
	}

}

func returnLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(200)
		return
	}
	dburl := "postgres://zov:123@localhost:5432/burmalda?sslmode=disable"
	conn, e := pgx.Connect(context.Background(), dburl)
	if e != nil {
		w.WriteHeader(502)
		return
	}
	defer conn.Close(context.Background())
	var login string
	body, er := io.ReadAll(r.Body)
	if er != nil {
		w.WriteHeader(424)
		return
	}
	u := Us{}
	json.Unmarshal(body, &u)
	row := conn.QueryRow(context.Background(), "SELECT login FROM base.users WHERE seckey = $1", u.SecKey).Scan(&login)
	if row != nil {
		w.WriteHeader(504)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(login))
	return
}

func searchUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(200)
		return
	}
	db := "postgres://zov:123@localhost:5432/burmalda?sslmode=disable"
	conn, e := pgx.Connect(context.Background(), db)
	if e != nil {
		w.WriteHeader(502)
		return
	}
	defer conn.Close(context.Background())
	u := Us{}
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &u)
	rows, err := conn.Query(context.Background(), "SELECT login FROM base.users WHERE login LIKE $1", ("%" + u.Login + "%"))
	if err != nil {
		w.WriteHeader(504)
		return
	}
	var logins []string
	for rows.Next() {
		var login string
		if err := rows.Scan(&login); err != nil {
			w.WriteHeader(404)
			return
		}
		logins = append(logins, login)
	}
	loginsJson, _ := json.Marshal(logins)
	w.WriteHeader(200)
	w.Write(loginsJson)
}

func Unique(slice []string) []string {
	seen := make(map[string]struct{}, len(slice))
	result := make([]string, 0, len(slice))

	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func returnChats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(200)
		return
	}
	db := "postgres://zov:123@localhost:5432/burmalda?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), db)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	//var u string
	var chas []string
	lg := r.URL.Query().Get("lg")
	// body, _ := io.ReadAll(r.Body)
	// json.Unmarshal(body, &u)
	rows, err := conn.Query(context.Background(), "SELECT to_user,from_user FROM base.messages WHERE (from_user = $1 OR to_user = $1)", lg)
	if err != nil {
		w.WriteHeader(502)
		return
	}
	for rows.Next() {
		var a string
		var b string
		if err := rows.Scan(&a, &b); err != nil {
			w.WriteHeader(424)
			return
		}
		chas = append(chas, (a + "/" + b))
	}
	chas = Unique(chas)
	json_chas, _ := json.Marshal(chas)
	w.WriteHeader(200)
	w.Write(json_chas)
}

type kpc struct {
	Content string `json:"Content"`
	User    string `json:"User"`
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(200)
		return
	}
	dburl := "postgres://zov:123@localhost:5432/burmalda?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), dburl)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	u := ChatMessage{}
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &u)
	rows, err := conn.Query(context.Background(), "SELECT from_user, content FROM base.messages WHERE (from_user = $1 AND to_user = $2) OR (from_user = $2 AND to_user = $1)", u.From, u.To)
	if err != nil {
		w.WriteHeader(502)
		return
	}
	var msgs []kpc
	for rows.Next() {
		var s string
		var userZ string
		if e := rows.Scan(&userZ, &s); e != nil {
			w.WriteHeader(424)
			return
		}
		msgs = append(msgs, kpc{Content: s, User: userZ})
	}
	w.WriteHeader(200)
	jsms, _ := json.Marshal(msgs)
	w.Write(jsms)
}

func main() {
	http.HandleFunc("/ws", WebSocketChat)
	http.HandleFunc("/searchUser", searchUser)
	http.HandleFunc("/authkey", AkeyFunc)
	http.HandleFunc("/auth", AuthFunc)
	http.HandleFunc("/changeData", ChangeData)
	http.HandleFunc("/returnLogin", returnLogin)
	http.HandleFunc("/returnChats", returnChats)
	http.HandleFunc("/getMessages", getMessages)
	http.ListenAndServe(":7777", nil)
}

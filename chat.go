package main

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	port := 3000
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	Usersockets := make(map[string]string)
	NumberOfusers := 0

	server.On("connection", func(so socketio.Socket) {
		so.Join("chat")

		so.On("add user", func(username string) {
			Usersockets[so.Id()] = username
			NumberOfusers = NumberOfusers + 1
			so.Emit("login", map[string]int{"numUsers": NumberOfusers})
			so.BroadcastTo("chat", "user joined", map[string]interface{}{"username": username, "numUsers": NumberOfusers})
		})

		so.On("typing", func(msg interface{}) {
			so.BroadcastTo("chat", "typing", map[string]string{"username": Usersockets[so.Id()]})
		})

		so.On("stop typing", func(msg string) {
			so.BroadcastTo("chat", "stop typing", map[string]string{"username": Usersockets[so.Id()]})
		})

		so.On("new message", func(msg string) {
			so.BroadcastTo("chat", "new message", map[string]string{"username": Usersockets[so.Id()], "message": msg})
		})

		so.On("disconnection", func() {
			if NumberOfusers > 0 {
				NumberOfusers = NumberOfusers - 1
			}
			so.BroadcastTo("chat", "user left", map[string]interface{}{"username": Usersockets[so.Id()], "numUsers": NumberOfusers})
			delete(Usersockets, so.Id())
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("Server listening at port", port)
	log.Fatal(http.ListenAndServe(strings.Join([]string{"", strconv.Itoa(port)}, ":"), nil))
}

package main

import (
	"fmt"
	"log"
    "os"

	"golang.org/x/net/websocket"
)

var origin = "http://localhost/"
var url = "ws://localhost:8080/echo"

func main() {
    args := os.Args[1]
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	message := []byte(args)
	_, err = ws.Write(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Send: %s\n", message)

	var msg = make([]byte, 512)
	_, err = ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive: %s\n", msg)
}

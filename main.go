package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"sync"
)

// EncodedConn can Publish any raw Go type using the registered Encoder
type person struct {
	Name     string
	Address  string
	Age      int
}

func main() {
	if len(os.Args) == 1 {
		nope()
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "pub":	publish()
	case "sub":	subscribe()
	default: nope()
	}
}

func publish() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer nc.Close()

	//if err := c.Publish("foo", []byte("Hello World")); err != nil {
	//	log.Fatal(err)
	//}

	me := &person{Name: "derek", Age: 22, Address: "140 New Montgomery Street, San Francisco, CA"}
	// Go type Publisher
	c.Publish("hello", me)
}

func subscribe() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer nc.Close()

	var wg sync.WaitGroup

	wg.Add(1)
	//nc.Subscribe("foo", func(m *nats.Msg) {
	//	wg.Done()
	//	fmt.Printf("Received a message: %s\n", string(m.Data))
	//})

	c.Subscribe("hello", func(p *person) {
		fmt.Printf("Received a person: %+v\n", p)
		wg.Done()
	})

	wg.Wait()
}

func nope() {
	fmt.Printf("nope")
}

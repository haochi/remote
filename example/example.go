package main

import (
	"flag"
	"fmt"
	"github.com/haochi/remote"
	"time"
)

const (
	printOne = iota
	printTwo
	echo
)

func main() {
	flag.Parse()
	var command = flag.Arg(0)
	var port = flag.Arg(1)

	var r = remote.New(port, command == "w")
	var size = 1000
	var messengers = make(chan []byte, size*3)

	r.Register(printOne, func(request []byte) []byte {
		return []byte(time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
	})

	r.Register(printTwo, func(request []byte) []byte {
		return []byte("2222")
	})

	r.Register(echo, func(request []byte) []byte {
		return request
	})

	r.Run(func() {
		for i := 0; i < size; i++ {
			go r.Go(printOne, []byte(fmt.Sprintf("hello %d", i)), messengers)
		}

		for i := 0; i < size; i++ {
			go r.Go(printTwo, []byte(fmt.Sprintf("hello %d", i)), messengers)
		}

		for i := 0; i < size; i++ {
			go r.Go(echo, []byte(fmt.Sprintf("hello %d", i)), messengers)
		}

		for i := 0; i < size*3; i++ {
			v := <-messengers
			fmt.Println(string(v))
		}
	})
}

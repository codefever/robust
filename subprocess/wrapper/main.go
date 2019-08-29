package main

import (
	"flag"
	"log"

	"github.com/codefever/robust/subprocess"
)

var command = flag.String("command", "", "Command")

func main() {
	flag.Parse()

	_, errc, _ := subprocess.RunCommand(*command)
	if err := <-errc; err != nil {
		log.Fatal(err)
	}
}

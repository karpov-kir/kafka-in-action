package main

import (
	"flag"
	"fmt"

	"github.com/karpov-kir/kafka-in-action/consumers"
	"github.com/karpov-kir/kafka-in-action/producers"
)

var (
	typeFlag string
)

func main() {
	flag.StringVar(&typeFlag, "type", "producer", "Type to run")
	flag.Parse()

	if typeFlag == "producer" {
		fmt.Println("Producing")
		producers.Produce()
	}

	if typeFlag == "consumer" {
		fmt.Println("Consuming")
		consumers.Consume()
	}
}

package main

import (
	"flag"
	"fmt"

	"github.com/karpov-kir/kafka-in-action/producers"
)

var (
	typeFlag string
)

func main() {
	flag.StringVar(&typeFlag, "type", "alert-producer", "Type to run")
	flag.Parse()

	if typeFlag == "alert-producer" {
		fmt.Println("Producing alert data")
		producers.ProduceAlertData("localhost:9091,localhost:9092,localhost:9093")
	}
}

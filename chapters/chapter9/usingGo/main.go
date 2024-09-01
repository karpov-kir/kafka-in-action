package main

import (
	"flag"
	"fmt"

	"github.com/karpov-kir/kafka-in-action/admin"
	consumers "github.com/karpov-kir/kafka-in-action/consumer"
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

	if typeFlag == "create-selfservice-topic" {
		fmt.Println("Creating selfservice topic")
		admin.CreateTopic()
	}

	if typeFlag == "alert-consumer" {
		fmt.Println("Consuming alert data")
		consumers.ConsumeAlertData()
	}
}

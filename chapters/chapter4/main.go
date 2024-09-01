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
	flag.StringVar(&typeFlag, "type", "alert-producer", "Type to run")
	flag.Parse()

	if typeFlag == "alert-producer" {
		fmt.Println("Producing alert data")
		producers.ProduceAlertData()
	}

	if typeFlag == "alert-trending-producer" {
		fmt.Println("Producing alert trending data")
		producers.ProduceAlertTrendingData()
	}

	if typeFlag == "audit-producer" {
		fmt.Println("Producing audit data")
		producers.ProduceAuditData()
	}

	if typeFlag == "alert-consumer" {
		fmt.Println("Consuming alert data")
		consumers.ConsumeAlertData()
	}
}

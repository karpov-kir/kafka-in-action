package main

import (
	"flag"
	"fmt"

	"github.com/karpov-kir/kafka-in-action/consumers"
	"github.com/karpov-kir/kafka-in-action/processors"
	"github.com/karpov-kir/kafka-in-action/producers"
)

var (
	typeFlag string
)

func main() {
	flag.StringVar(&typeFlag, "type", "processor", "Type to run")
	flag.Parse()

	if typeFlag == "transaction-processor" {
		fmt.Println("Processing transactions")
		processors.ProcessTransactions()
	}

	if typeFlag == "transaction-request-producer" {
		fmt.Println("Producing transaction requests")
		producers.ProduceTransactionRequests()
	}

	if typeFlag == "transaction-success-consumer" {
		fmt.Println("Consuming success transactions")
		consumers.ConsumeTransactionSuccess()
	}

	if typeFlag == "transaction-failed-consumer" {
		fmt.Println("Consuming failed transactions")
		consumers.ConsumeTransactionFailed()
	}

	if typeFlag == "latest-transactions-view" {
		processors.RunLatestTransactionsView()
	}
}

package main

import (
	"flag"
	"fmt"

	"github.com/karpov-kir/kafka-in-action/consumers"
)

var (
	typeFlag string
)

func main() {
	flag.StringVar(&typeFlag, "type", "alert-trending-consumer", "Type to run")
	flag.Parse()

	if typeFlag == "alert-trending-consumer" {
		fmt.Println("Consuming alert trending data")
		consumers.ConsumeAlertTrendingData()
	}

	if typeFlag == "views-consumer" {
		fmt.Println("Consuming views data")
		consumers.ConsumeViewsData()
	}

	if typeFlag == "audit-consumer" {
		fmt.Println("Consuming audit data")
		consumers.ConsumeAuditData()
	}

	if typeFlag == "promo-consumer" {
		fmt.Println("Consuming promo data")
		consumers.ConsumePromoData()
	}

	if typeFlag == "promo-for-magic-calculation-consumer" {
		fmt.Println("Consuming promo data for magic calculation")
		consumers.ConsumePromoDataForMagicCalculation()
	}
}

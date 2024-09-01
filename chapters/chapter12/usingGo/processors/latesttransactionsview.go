package processors

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/karpov-kir/kafka-in-action/models"
	"github.com/karpov-kir/kafka-in-action/utils"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/web/index"
	"github.com/lovoo/goka/web/monitor"
)

func RunLatestTransactionsView() {
	var brokers = []string{"localhost:9091", "localhost:9092", "localhost:9093"}
	serializer, deserializer := newAvroSerdes()
	transactionRequestCodec := newAvroCodec[models.Transaction](serializer, deserializer, TransactionRequestTopic)

	latestTransactionsView, err := goka.NewView(brokers, goka.Table(TransactionRequestTopic), transactionRequestCodec, goka.WithViewAutoReconnect())
	if err != nil {
		log.Fatalf("Cannot create view: %v", err)
	}
	// Context we'll use to run the view and the state change observer
	ctx, cancel := context.WithCancel(context.Background())

	// Channel used to wait for the view to finish
	done := make(chan struct{})
	go func() {
		defer close(done)
		err := latestTransactionsView.Run(ctx)
		if err != nil {
			log.Printf("View finished with error: %v", err)
		}
	}()

	// Get a state change observer and
	go func() {
		observer := latestTransactionsView.ObserveStateChanges()
		defer observer.Stop()
		for {
			select {
			case state, ok := <-observer.C():
				if !ok {
					return
				}
				log.Printf("View is in state: %v", goka.ViewState(state))
			case <-ctx.Done():
				return
			}
		}
	}()

	router := mux.NewRouter()

	router.HandleFunc("/latest-transactions", func(w http.ResponseWriter, r *http.Request) {
		viewIterator, err := latestTransactionsView.Iterator()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		latestTransactions := make(map[string]interface{})
		for viewIterator.Next() {
			transaction, err := viewIterator.Value()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			adjustedTransaction := make(map[string]interface{})
			v := reflect.ValueOf(transaction)
			t := v.Type()

			for i := 0; i < v.NumField(); i++ {
				key := t.Field(i).Name
				value := v.Field(i).Interface()

				// It's better to create a separate struct and adjust the `Amount` field there,
				// but since it's a sample code, we can do it here.
				if key == "Amount" {
					adjustedTransaction[key] = utils.AvroBytesToDecimal(value.(models.Bytes))
					continue
				}
				adjustedTransaction[key] = value
			}

			latestTransactions[viewIterator.Key()] = adjustedTransaction
		}

		latestTransactionsJson, err := json.Marshal(latestTransactions)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(latestTransactionsJson)
	})

	monitorServer := monitor.NewServer("/monitor", router)
	monitorServer.AttachView(latestTransactionsView)

	idxServer := index.NewServer("/", router)
	idxServer.AddComponent(monitorServer, "Monitor")

	server := &http.Server{Addr: "0.0.0.0:8877", Handler: router}

	go func() {
		fmt.Printf("Starting server on %s\n", server.Addr)
		err = server.ListenAndServe()
		if err != http.ErrServerClosed {
			panic(err)
		}
	}()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigchan:
	case <-done:
	}
	cancel()
	<-done

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
}

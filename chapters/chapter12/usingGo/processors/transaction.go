package processors

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"github.com/karpov-kir/kafka-in-action/models"
	"github.com/karpov-kir/kafka-in-action/utils"
	"github.com/lovoo/goka"
)

const TransactionRequestTopic string = "transaction-request"
const TransactionSuccessTopic string = "transaction-success"
const TransactionFailedTopic string = "transaction-failed"
const TransactionProcessorGroup string = "transaction-processor"

func ProcessTransactions() {
	printTransactionRequest := func(ctx goka.Context, transaction models.Transaction) {
		log.Printf("Transaction request key = %s, msg = %v", ctx.Key(), transaction)
	}

	transformTransaction := func(ctx goka.Context, transaction models.Transaction) {
		getFunds := func() models.Funds {
			if val := ctx.Value(); val != nil {
				funds, ok := val.(models.Funds)

				if !ok {
					panic(fmt.Sprintf("could not cast value of type %T to %T in getFunds", val, models.Funds{}))
				}

				return funds
			}

			return models.Funds{
				Account: transaction.Account,
				Balance: utils.DecimalToAvroBytes(0),
			}
		}

		storeFunds := func(funds models.Funds) {
			ctx.SetValue(funds)
		}

		transformer := newTransactionTransformer(getFunds, storeFunds)
		transactionResult := transformer.transformTransaction(transaction)

		if transactionResult.Success {
			fmt.Printf("Total funds for account %f\n", utils.AvroBytesToDecimal(transactionResult.Funds.Balance))
			ctx.Emit(goka.Stream(TransactionSuccessTopic), transaction.Account, transactionResult)
		} else {
			fmt.Printf("Transaction failed: %s\n", transactionResult.ErrorType.ErrorType.String())
			ctx.Emit(goka.Stream(TransactionFailedTopic), transaction.Account, transactionResult)
		}
	}

	handleTransactionRequest := func(ctx goka.Context, message interface{}) {
		transaction, ok := message.(models.Transaction)

		if !ok {
			panic(fmt.Sprintf("could not cast message of type %T to %T in handleTransactionRequest", message, models.Transaction{}))
		}

		printTransactionRequest(ctx, transaction)
		transformTransaction(ctx, transaction)
	}

	serializer, deserializer := newAvroSerdes()
	transactionRequestCodec := newAvroCodec[models.Transaction](serializer, deserializer, TransactionRequestTopic)
	transactionSuccessCodec := newAvroCodec[models.TransactionResult](serializer, deserializer, TransactionSuccessTopic)
	transactionFailedCodec := newAvroCodec[models.TransactionResult](serializer, deserializer, TransactionFailedTopic)

	// Define a new processor group. The group defines all inputs, outputs, and
	// serialization formats. There is no a direct alternative to Kafka Streams in Go, so we use Goka
	// as the closest replacement.
	groupGraph := goka.DefineGroup(
		goka.Group(TransactionProcessorGroup),
		goka.Input(goka.Stream(TransactionRequestTopic), transactionRequestCodec, handleTransactionRequest),
		goka.Persist(newAvroCodec[models.Funds](serializer, deserializer, "transaction-funds")),
		goka.Output(goka.Stream(TransactionSuccessTopic), transactionSuccessCodec),
		goka.Output(goka.Stream(TransactionFailedTopic), transactionFailedCodec),
	)

	processor, err := goka.NewProcessor([]string{"localhost:9091", "localhost:9092", "localhost:9093"},
		groupGraph,
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
	)
	if err != nil {
		log.Fatalf("error creating processor: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err = processor.Run(ctx); err != nil {
			log.Printf("error running processor: %v", err)
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
}

func newAvroSerdes() (*avro.SpecificSerializer, *avro.SpecificDeserializer) {
	client, err := schemaregistry.NewClient(schemaregistry.NewConfig("http://localhost:8081"))

	if err != nil {
		fmt.Printf("Failed to create schema registry client: %s\n", err)
		os.Exit(1)
	}

	serializer, err := avro.NewSpecificSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())

	if err != nil {
		fmt.Printf("Failed to create serializer: %s\n", err)
		os.Exit(1)
	}

	deserializer, err := avro.NewSpecificDeserializer(client, serde.ValueSerde, avro.NewDeserializerConfig())

	if err != nil {
		fmt.Printf("Failed to create deserializer: %s\n", err)
		os.Exit(1)
	}

	return serializer, deserializer
}

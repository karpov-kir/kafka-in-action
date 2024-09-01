# Chapter 3 (using Avro schema)

To create the `kinaction_alert_avro` topic from the console execute on a broker:

```bash
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_alert_avro --partitions 3 --replication-factor 3
```

---

To install a tool and generate Go structs from JSON Avro schemas, execute on the host:

```bash
go install github.com/hamba/avro/v2/cmd/avrogen@v2.25.0
avrogen -pkg models -o ./models/alert.go -tags json:snake,yaml:upper-camel ./schemas/kinactionAlert.avsc
```

`avrogen` only generates Go structs from Avro schemas, so it can be used only with `github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro.GenericDeserializer` and `github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro.GenericSerializer`. Also, `avrogen` is quite primitive (e.g. it does not generate enums).

To generate methods (and some other stuff line enums) as well to use with the `github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro.SpecificDeserializer` and `github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro.SpecificSerializer` another tool is needed e.g. https://github.com/actgardner/gogen-avro. This will generate methods (`Serialize` and `Schema`) to implement the `github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro.SpecificAvroMessage` interface, the methods are used by `github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro.SpecificDeserializer` and `github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro.SpecificSerializer`. Refer to the [chapter 11](../chapter11/README.md) for an example.

---

To produce from the Go application and consume in the Go application execute on the host:

```bash
go run ./main.go --type consumer
go run ./main.go --type producer
```

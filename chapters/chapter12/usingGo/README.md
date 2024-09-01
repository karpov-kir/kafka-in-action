# Chapter 12 (using Go)

To create the topics from the console execute on a broker:

```bash
# This needs to be compacted because it's used as a table topic (refer to https://github.com/lovoo/goka/wiki/Tips). 
kafka-topics --create --bootstrap-server broker3:29093 --topic transaction-request --partitions 3 --replication-factor 3 --config cleanup.policy=compact
kafka-topics --create --bootstrap-server broker3:29093 --topic transaction-success --partitions 3 --replication-factor 3
kafka-topics --create --bootstrap-server broker3:29093 --topic transaction-failed --partitions 3 --replication-factor 3
```


To install a tool and generate Go structs and methods from JSON Avro schemas, execute on the host:

```bash
go install github.com/actgardner/gogen-avro/v10/cmd/...@v10.2.1

gogen-avro --package=models ./models ./schemas/*.avsc
```

To run the processor and the related view from the Go application execute on the host:

```bash
go run ./main.go --type transaction-processor
go run ./main.go --type latest-transactions-view
```

There might be a problem with Goka's cache (the producer throws an error in this case), to flush the cache execute on the host:

```bash
rm -rf /tmp/goka
```

To run test consumers and producers from the Go application execute on the host:

```bash
go run ./main.go --type transaction-success-consumer
go run ./main.go --type transaction-failed-consumer
go run ./main.go --type transaction-request-producer
```

After that you can access the view at http://localhost:8877/latest-transactions and see some data.

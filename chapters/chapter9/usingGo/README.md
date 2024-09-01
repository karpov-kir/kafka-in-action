# Chapter 9 (using Go)

To create the `kinaction_selfserviceTopic` topic via the admin client in Go execute on the host:

```bash
go run ./main.go --type create-selfservice-topic
```

---

To test the interceptors create the topic `kinaction_alert` by executing on a broker:

```bash
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_alert --partitions 3 --replication-factor 3
```

Then execute on the host:

```bash
go run ./main.go --type alert-consumer
go run ./main.go --type alert-producer
```

# Chapter 4

To create the topics from the console execute on a broker:

```bash
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_alert --partitions 3 --replication-factor 3
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_alerttrend --partitions 3 --replication-factor 3
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_audit --partitions 3 --replication-factor 3
```

---

To produce alert data from the Go application and consume in the Go application execute on the host:

```bash
go run ./main.go --type alert-consumer
go run ./main.go --type alert-producer
```

---

To produce alert trending data from the Go application and consume in the console execute on a broker:

```bash
kafka-console-consumer --bootstrap-server broker3:29093 --topic kinaction_alerttrend --from-beginning
```

Then on the host:

```bash
go run ./main.go --type alert-trending-producer
```

---

To produce audit data from the Go application and consume in the console execute on a broker:

```bash
kafka-console-consumer --bootstrap-server broker3:29093 --topic kinaction_audit --from-beginning
```

Then on the host:

```bash
go run ./main.go --type audit-producer
```

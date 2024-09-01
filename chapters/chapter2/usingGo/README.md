# Chapter 2 (using Go)

To create the `kinaction_helloworld` topic from the console execute on a broker:

```bash
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_helloworld --partitions 3 --replication-factor 3
```

---

To produce from the Go application and consume in the Go application:

```bash
go run ./main.go --type consumer
go run ./main.go --type producer
```

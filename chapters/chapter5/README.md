# Chapter 5

To create the topics from the console execute on a broker:

```bash
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_alerttrend --partitions 3 --replication-factor 3
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_views --partitions 3 --replication-factor 3
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_audit --partitions 3 --replication-factor 3
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_promo --partitions 3 --replication-factor 3
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_webconsumer --partitions 3 --replication-factor 3
```

---

To produce alert trending data from the console and consume in the Go application execute on the host:

```bash
go run ./main.go --type alert-trending-consumer
```

Then on a broker:

```bash
echo 'Stage 1:WARNING' | kafka-console-producer --property "parse.key=true" --property "key.separator=:" --broker-list broker3:29093 --topic kinaction_alerttrend
```

---

To produce views data from the console and consume in the Go application execute on the host:

```bash
go run ./main.go --type views-consumer
```

Then on a broker:

```bash
echo 'Page view 1' | kafka-console-producer --broker-list broker3:29093 --topic kinaction_views
```

---

To produce audit data from the console and consume in the Go application execute on the host:

```bash
go run ./main.go --type audit-consumer
```

Then on a broker:

```bash
echo 'audit event' | kafka-console-producer --broker-list broker3:29093 --topic kinaction_audit
```

---

To produce promo data from the console and consume in the Go application execute on the host:

```bash
go run ./main.go --type promo-consumer
```

Then on a broker:

```bash
echo 'promo event' | kafka-console-producer --broker-list broker3:29093 --topic kinaction_promo
```

---

To produce promo data for magic calculation from the console and consume in the Go application execute on the host:

```bash
go run ./main.go --type promo-for-magic-calculation-consumer
```

Then on a broker:

```bash
echo '100' | kafka-console-producer --broker-list broker3:29093 --topic kinaction_promo
```

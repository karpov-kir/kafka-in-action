# Chapter 6

To connect to the ZooKeeper shell and get the list of topics or the controller execute or on a broker:

```bash
zookeeper-shell zookeeper:2181
ls /brokers/topics
get /controller
```

---

Grafana and Prometheus are configured via Docker Compose. Refer to the root [README.md](../../README.md) file for more information.

---

To confirm the health of a topic create the topic by executing on a broker:

```bash
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_replica_test --partitions 3 --replication-factor 3
```

Then you check the status of the topic by executing on a broker:

```bash
kafka-topics --describe --bootstrap-server broker3:29093 --topic kinaction_replica_test
```

Then you stop one of the brokers to simulate an outage by executing on the host:

```bash
docker-compose stop broker2
```

Then produce a message and check the under-replicated partitions by executing on a broker:

```bash
echo 'Hello world!' | kafka-console-producer --broker-list broker3:29093 --topic kinaction_replica_test
kafka-topics --describe --bootstrap-server broker3:29093 --under-replicated-partitions
```

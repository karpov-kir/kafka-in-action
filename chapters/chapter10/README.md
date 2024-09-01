# Chapter 10

To create the `kinaction_test_ssl` topic from the console execute on a broker:

```bash
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_test_ssl --partitions 2 --replication-factor 2
```

---

SSL is already configured in the Docker compose setup. To produce from the console and consume in the console over SSL execute on a broker:

```bash
kafka-console-consumer --bootstrap-server broker3:29993 --topic kinaction_test_ssl --consumer.config /kafka-in-action/chapter10/kinaction-ssl.properties
kafka-console-producer --bootstrap-server broker3:29993 --topic kinaction_test_ssl --producer.config /kafka-in-action/chapter10/kinaction-ssl.properties
```


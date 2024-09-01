# Chapter 2 (using console)

To create the `kinaction_helloworld` topic from the console execute on a broker:

```bash
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_helloworld --partitions 3 --replication-factor 3
```

---

To describe the `kinaction_helloworld` topic from the console execute on a broker:

```bash
kafka-topics --describe --bootstrap-server broker3:29093 --topic kinaction_helloworld
```

---

To produce from the console and consume in the console execute on a broker:

```bash
kafka-console-producer --bootstrap-server broker3:29093 --topic kinaction_helloworld
kafka-console-consumer --bootstrap-server broker3:29093 --topic kinaction_helloworld --from-beginning
```

# Chapter 9 (using console)

To create the `kinaction_selfserviceTopic` topic from the console execute on a broker:

```bash
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_selfserviceTopic --partitions 2 --replication-factor 2
```

---

`kcat` is installed on the brokers (refer to the [Broker.Dockerfile](../../../docker/Broker.Dockerfile)).

To produce from the console and consume in the console via `kcat` execute on a broker:

```bash
# There might be problems with detecting EOF by `kcat` even if using CTRL+D.
# To workaround it you can use `cat |` in front of the `kcat` command and then use CTRL+D
# when you are done with the input. Otherwise, `kcat` might hang (this is probably due to Docker console implementation).
cat | kcat -C -b broker3:29093 -t kinaction_selfserviceTopic
echo "Test message" | kcat -P -b broker3:29093 -t kinaction_selfserviceTopic
```

---

Confluent REST proxy is available in the Docker compose setup.

To get the list of topics execute on the host:

```bash
curl -X GET  -H "Accept: application/vnd.kafka.v2+json" localhost:8086/topics
```

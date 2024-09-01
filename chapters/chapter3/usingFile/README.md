# Chapter 3 (using file)

To produce from the file and consume in the console execute on a broker:

```bash
connect-standalone /kafka-in-action/chapter3/usingFile/connect-standalone.properties /kafka-in-action/chapter3/usingFile/alert-source.properties
kafka-console-consumer --bootstrap-server broker3:29093 --topic kinaction_alert_connect --from-beginning
```

Then on the host:

```bash
echo "This is a test alert" >> ./alert.txt
```

---

To produce from the file and sink to another file execute on a broker:

```bash
connect-standalone /kafka-in-action/chapter3/connect-standalone.properties /kafka-in-action/chapter3/alert-source.properties /kafka-in-action/chapter3/alert-sink.properties
```

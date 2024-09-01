# Chapter 3 (using DB)

JDBC connector is installed on the brokers (refer to the [Broker.Dockerfile](../../../docker/Broker.Dockerfile)).

To produce from the DB and consume in the console execute on the host (ensure you have SQLite installed):

```bash
./initInvoicesDb.sh
```

Then on a broker:

```bash
connect-standalone /kafka-in-action/chapter3/usingDb/connect-standalone.properties /kafka-in-action/chapter3/usingDb/sqlite-source.properties
kafka-console-consumer --bootstrap-server broker3:29093 --topic kinaction-sqlite-jdbc-invoices --from-beginning
```

Then on the host:

```bash
./insertToInvoicesTable.sh 2
./insertToInvoicesTable.sh 3
...
```

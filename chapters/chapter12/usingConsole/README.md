# Chapter 12 (using console)

To use commands in this chapter ensure you have the processor, transaction success consumer, and transaction request producer of [../usingGo](../usingGo) running. They will be producing data and ensure that the schemas are registered. 

To create the `account` topic from the console execute on a broker:

```bash
kafka-topics --create --bootstrap-server broker3:29093 --topic account --partitions 3 --replication-factor 3
```

To seed the `account` topic with some data from the console execute on a broker:

```bash
echo '123#{"number":1,"firstName":"John Doe","lastName":"Doe","numberAddress":"10","streetAddress":"Main Street","cityAddress":"Springfield","countryAddress":"USA","creationDate":1612137600000,"updateDate":1612137600000}\n456#{"number":2,"firstName":"Jane Doe","lastName":"Doe","numberAddress":"20","streetAddress":"Main Street","cityAddress":"Springfield","countryAddress":"USA","creationDate":1612137600000,"updateDate":1612137600000}\n789#{"number":3,"firstName":"John Smith","lastName":"Smith","numberAddress":"30","streetAddress":"Main Street","cityAddress":"Springfield","countryAddress":"USA","creationDate":1612137600000,"updateDate":1612137600000}' | kafka-avro-console-producer --topic account --bootstrap-server broker3:29093 --property schema.registry.url=http://schema-registry:8081 --property value.schema="$(< /kafka-in-action/chapter12/usingGo/schemas/account.avsc)" --property "parse.key=true" --property "key.separator=#" --property key.serializer=org.apache.kafka.common.serialization.StringSerializer

kafka-avro-console-consumer --topic account --bootstrap-server broker3:29093 --property schema.registry.url=http://schema-registry:8081 --from-beginning --property print.key=true --property key.deserializer=org.apache.kafka.common.serialization.StringDeserializer --property key.separator="#"
```

To create tables and streams in ksqlDB execute on the ksqldb-cli:

```bash
ksql http://ksqldb-server:8088

SET 'auto.offset.reset' = 'earliest';

CREATE STREAM TRANSACTION_SUCCESS (
  numkey string KEY,
  transaction STRUCT<guid STRING, account STRING, amount DECIMAL(9, 2), type STRING, currency STRING, country STRING>,
  funds STRUCT<account STRING, balance DECIMAL(9, 2)>
) WITH (KAFKA_TOPIC='transaction-success', VALUE_FORMAT='avro', PARTITIONS=3, REPLICAS=1);

CREATE SOURCE TABLE ACCOUNT (
  numkey string PRIMARY KEY,
  number INT,
  firstName STRING,
  lastName STRING,
  numberAddress STRING,
  streetAddress STRING,
  cityAddress STRING,
  countryAddress STRING,
  creationDate BIGINT,
  updateDate BIGINT
) WITH (KAFKA_TOPIC = 'account', VALUE_FORMAT='avro', PARTITIONS=3, REPLICAS=1);

CREATE STREAM TRANSACTION_STATEMENT AS
  SELECT *
  FROM TRANSACTION_SUCCESS
  LEFT JOIN ACCOUNT ON TRANSACTION_SUCCESS.numkey = ACCOUNT.numkey
  EMIT CHANGES;

SELECT * FROM TRANSACTION_SUCCESS EMIT CHANGES LIMIT 10;
SELECT * FROM ACCOUNT EMIT CHANGES LIMIT 3;
SELECT * FROM TRANSACTION_STATEMENT EMIT CHANGES;
```

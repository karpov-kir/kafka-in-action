# Chapter 11

Some Schema registry REST API:

```bash
curl -X GET http://localhost:8081/config
curl -X GET http://localhost:8081/subjects
curl -X GET http://localhost:8081/subjects/{subject}/versions
curl -X GET http://localhost:8081/subjects/{subject}/versions/{version}
curl -X DELETE http://localhost:8081/subjects/{subject}
curl -X POST -H "Content-Type: application/vnd.schemaregistry.v1+json"  --data '{schema JSON}' http://localhost:8081/compatibility/subjects/kinaction_schematest-value/versions/latest
```

---

To create the `kinaction_schematest` topic from the console execute on a broker:

```bash
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_schematest --partitions 3 --replication-factor 3
```

---

To install a tool and generate Go structs and methods from JSON Avro schemas, execute on the host:

```bash
go install github.com/actgardner/gogen-avro/v10/cmd/...@v10.2.1

gogen-avro --package=models ./models ./schemas/kinactionAlert.avsc
```

---

To produce from the Go application and consume in the Go application execute on the host:

```bash
go run ./main.go --type consumer
go run ./main.go --type producer
```

This will create the `kinactionAlert.avsc` schema in the Schema registry automatically. To verify execute on the host:

```bash
curl -X GET http://localhost:8081/subjects
curl -X GET http://localhost:8081/subjects/kinaction_schematest-value/versions
curl -X GET http://localhost:8081/subjects/kinaction_schematest-value/versions/1
```

To generate the next version and update the schema in Schema registry execute on the host:

```bash
gogen-avro --package=models ./models ./schemas/kinactionAlert.v2.avsc
go run ./main.go --type consumer
go run ./main.go --type producer
```

To verify execute on the host:

```bash
curl -X GET http://localhost:8081/subjects/kinaction_schematest-value/versions
curl -X GET http://localhost:8081/subjects/kinaction_schematest-value/versions/2
```

---

To check a schema candidate compatibility against a topic in the Schema registry execute on the host:

```bash
curl -X POST -H "Content-Type: application/vnd.schemaregistry.v1+json" --data '{ "schema": "{ \"type\": \"record\", \"name\": \"Alert\", \"fields\": [{ \"name\": \"notafield\", \"type\": \"long\" } ]}" }' http://localhost:8081/compatibility/subjects/kinaction_schematest-value/versions/latest
```

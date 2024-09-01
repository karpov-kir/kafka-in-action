# Chapter 7 (using Go)

To run the partitioner integration tests execute the following command on the host:

```bash
go clean -testcache && go test  ./...
```

The first execution might take a while because it will download the `confluentinc/confluent-local` image.

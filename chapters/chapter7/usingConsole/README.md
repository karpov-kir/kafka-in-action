# Chapter 7 (using console)

To delete a topic create the topic fist by executing the following command on a broker:

```bash
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_topicandpart --partitions 2 --replication-factor 2
```

Then delete the topic by executing the following command on a broker:

```bash
kafka-topics --delete --bootstrap-server broker3:29093 --topic kinaction_topicandpart
```

---

To create and test a compacted topic execute the following commands on the first broker:

```bash
# Create a topic and assign the leader partition to the first broker and the replica partition to the second broker.
kafka-topics --create --bootstrap-server broker3:29093 --topic kinaction_compact --config cleanup.policy=compact --replica-assignment 0:1

# Write 3 messages.
echo "Stage 1:WARNING\nStage 1:WARNING\nStage 1:WARNING" | kafka-console-producer --property "parse.key=true" --property "key.separator=:" --broker-list broker3:29093 --topic kinaction_compact
# Check the existing logs and the topic content. You should see 3 messages.
kcat -C -e -o beginning -b broker3:29093 -t kinaction_compact -K:
# Note: these files exist only on the first broker.
kafka-dump-log --print-data-log --files $(ls -1 /var/lib/kafka/data/kinaction_compact-0/*.log | tr '\n' ',')

# Change the topic settings in a way such that compaction should kick in after 10 seconds.
kafka-configs --alter --bootstrap-server broker3:29093 --entity-type topics --entity-name kinaction_compact --add-config max.compaction.lag.ms=10000,min.cleanable.dirty.ratio=0,segment.ms=10000,delete.retention.ms=10000
# Wait for the last segment to outdate.
sleep 11

# Check the existing logs and the topic content again. You should still see the same 3 messages. To have compaction running you need to have at least 2 segment files (one finished and one running).
kcat -C -e -o beginning -b broker3:29093 -t kinaction_compact -K:
kafka-dump-log --print-data-log --files $(ls -1 /var/lib/kafka/data/kinaction_compact-0/*.log | tr '\n' ',')

# Write 3 messages again.
echo "Stage 1:WARNING\nStage 1:WARNING\nStage 1:WARNING" | kafka-console-producer --property "parse.key=true" --property "key.separator=:" --broker-list broker3:29093 --topic kinaction_compact
sleep 11

# Check the existing logs and the topic content again. You should see only 4 messages instead of 6.
kcat -C -e -o beginning -b broker3:29093 -t kinaction_compact -K:
kafka-dump-log --print-data-log --files $(ls -1 /var/lib/kafka/data/kinaction_compact-0/*.log | tr '\n' ',')
```

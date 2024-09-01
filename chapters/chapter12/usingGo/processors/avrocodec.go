package processors

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
)

type avroCodec[T any] struct {
	serializer   *avro.SpecificSerializer
	deserializer *avro.SpecificDeserializer
	topic        string
}

func (ac *avroCodec[T]) Encode(message interface{}) ([]byte, error) {
	model, ok := message.(T)

	if !ok {
		return nil, fmt.Errorf("could not cast message to %T in avroCodec.Encode", model)
	}

	return ac.serializer.Serialize(ac.topic, &model)
}

func (ac *avroCodec[T]) Decode(message []byte) (interface{}, error) {
	var model T
	err := ac.deserializer.DeserializeInto(ac.topic, message, &model)
	return model, err
}

func newAvroCodec[T any](serializer *avro.SpecificSerializer, deserializer *avro.SpecificDeserializer, topic string) *avroCodec[T] {
	return &avroCodec[T]{
		serializer:   serializer,
		deserializer: deserializer,
		topic:        topic,
	}
}

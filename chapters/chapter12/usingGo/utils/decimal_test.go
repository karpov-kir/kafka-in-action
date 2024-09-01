package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestDecimalToAvroBytes(t *testing.T) {
	expected := int16(1050)
	result := DecimalToAvroBytes(10.5000)
	resultInt := bytesToInt16(result)

	if expected != resultInt {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestAvroBytesToDecimal(t *testing.T) {
	bytes := int16ToBytes(1050)
	result := AvroBytesToDecimal(bytes)

	expected := 10.5

	if result != expected {
		t.Errorf("Expected %f, got %f", expected, result)
	}
}

func bytesToInt16(data []byte) int16 {
	var resultInt int16
	buffer := bytes.NewReader(data)
	err := binary.Read(buffer, binary.BigEndian, &resultInt)

	if err != nil {
		panic(fmt.Errorf("Error converting bytes to decimal: %v", err))
	}

	return resultInt
}

func int16ToBytes(data int16) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, data)

	if err != nil {
		panic(fmt.Errorf("Error converting decimal to bytes: %v", err))
	}

	return buffer.Bytes()
}

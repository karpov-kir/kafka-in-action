package utils

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/big"
)

func DecimalToAvroBytes(value float64) []byte {
	scale := 2
	// Convert the decimal to an unscaled integer
	scaleFactor := new(big.Float).SetFloat64(math.Pow(10, float64(scale)))
	unscaledValue := new(big.Float).Mul(new(big.Float).SetFloat64(value), scaleFactor)

	unscaledInt, _ := unscaledValue.Int(nil)

	byteBuffer := new(bytes.Buffer)
	binary.Write(byteBuffer, binary.BigEndian, unscaledInt.Bytes())

	return byteBuffer.Bytes()
}

func AvroBytesToDecimal(data []byte) float64 {
	scale := 2
	unscaledInt := new(big.Int).SetBytes(data)

	unscaledFloat := new(big.Float).SetInt(unscaledInt)

	scaleFactor := new(big.Float).SetFloat64(math.Pow(10, float64(scale)))
	decimalValue := new(big.Float).Quo(unscaledFloat, scaleFactor)

	result, _ := decimalValue.Float64()

	return result
}

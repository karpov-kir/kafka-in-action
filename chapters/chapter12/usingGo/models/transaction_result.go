// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCES:
 *     account.avsc
 *     funds.avsc
 *     transaction.avsc
 *     transaction_result.avsc
 */
package models

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/v10/compiler"
	"github.com/actgardner/gogen-avro/v10/vm"
	"github.com/actgardner/gogen-avro/v10/vm/types"
)

var _ = fmt.Printf

type TransactionResult struct {
	Transaction Transaction `json:"transaction"`

	Funds Funds `json:"funds"`
	// success
	Success bool `json:"success"`

	ErrorType *UnionNullErrorType `json:"errorType"`
}

const TransactionResultAvroCRC64Fingerprint = "\xdejP}\x99Ÿ\xed"

func NewTransactionResult() TransactionResult {
	r := TransactionResult{}
	r.Transaction = NewTransaction()

	r.Funds = NewFunds()

	return r
}

func DeserializeTransactionResult(r io.Reader) (TransactionResult, error) {
	t := NewTransactionResult()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeTransactionResultFromSchema(r io.Reader, schema string) (TransactionResult, error) {
	t := NewTransactionResult()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeTransactionResult(r TransactionResult, w io.Writer) error {
	var err error
	err = writeTransaction(r.Transaction, w)
	if err != nil {
		return err
	}
	err = writeFunds(r.Funds, w)
	if err != nil {
		return err
	}
	err = vm.WriteBool(r.Success, w)
	if err != nil {
		return err
	}
	err = writeUnionNullErrorType(r.ErrorType, w)
	if err != nil {
		return err
	}
	return err
}

func (r TransactionResult) Serialize(w io.Writer) error {
	return writeTransactionResult(r, w)
}

func (r TransactionResult) Schema() string {
	return "{\"fields\":[{\"name\":\"transaction\",\"type\":{\"fields\":[{\"doc\":\"The unique transaction guid\",\"name\":\"guid\",\"type\":\"string\"},{\"avro.java.string\":\"java.lang.String\",\"doc\":\"Account name\",\"name\":\"account\",\"type\":\"string\"},{\"name\":\"amount\",\"type\":{\"logicalType\":\"decimal\",\"precision\":9,\"scale\":2,\"type\":\"bytes\"}},{\"name\":\"type\",\"type\":{\"name\":\"TransactionType\",\"symbols\":[\"DEPOSIT\",\"WITHDRAW\"],\"type\":\"enum\"}},{\"doc\":\"Transaction currency\",\"name\":\"currency\",\"type\":\"string\"},{\"doc\":\"Transaction country\",\"name\":\"country\",\"type\":\"string\"}],\"name\":\"Transaction\",\"namespace\":\"org.kafkainaction\",\"type\":\"record\"}},{\"name\":\"funds\",\"type\":{\"fields\":[{\"doc\":\"Account name\",\"name\":\"account\",\"type\":\"string\"},{\"name\":\"balance\",\"type\":{\"logicalType\":\"decimal\",\"precision\":9,\"scale\":2,\"type\":\"bytes\"}}],\"name\":\"Funds\",\"namespace\":\"org.kafkainaction\",\"type\":\"record\"}},{\"doc\":\"success\",\"name\":\"success\",\"type\":\"boolean\"},{\"name\":\"errorType\",\"type\":[\"null\",{\"name\":\"ErrorType\",\"symbols\":[\"INSUFFICIENT_FUNDS\"],\"type\":\"enum\"}]}],\"includeSchemas\":[{\"name\":\"Transaction\"},{\"name\":\"Funds\"}],\"name\":\"org.kafkainaction.TransactionResult\",\"type\":\"record\"}"
}

func (r TransactionResult) SchemaName() string {
	return "org.kafkainaction.TransactionResult"
}

func (_ TransactionResult) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ TransactionResult) SetInt(v int32)       { panic("Unsupported operation") }
func (_ TransactionResult) SetLong(v int64)      { panic("Unsupported operation") }
func (_ TransactionResult) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ TransactionResult) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ TransactionResult) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ TransactionResult) SetString(v string)   { panic("Unsupported operation") }
func (_ TransactionResult) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *TransactionResult) Get(i int) types.Field {
	switch i {
	case 0:
		r.Transaction = NewTransaction()

		w := types.Record{Target: &r.Transaction}

		return w

	case 1:
		r.Funds = NewFunds()

		w := types.Record{Target: &r.Funds}

		return w

	case 2:
		w := types.Boolean{Target: &r.Success}

		return w

	case 3:
		r.ErrorType = NewUnionNullErrorType()

		return r.ErrorType
	}
	panic("Unknown field index")
}

func (r *TransactionResult) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *TransactionResult) NullField(i int) {
	switch i {
	case 3:
		r.ErrorType = nil
		return
	}
	panic("Not a nullable field index")
}

func (_ TransactionResult) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ TransactionResult) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ TransactionResult) HintSize(int)                     { panic("Unsupported operation") }
func (_ TransactionResult) Finalize()                        {}

func (_ TransactionResult) AvroCRC64Fingerprint() []byte {
	return []byte(TransactionResultAvroCRC64Fingerprint)
}

func (r TransactionResult) MarshalJSON() ([]byte, error) {
	var err error
	output := make(map[string]json.RawMessage)
	output["transaction"], err = json.Marshal(r.Transaction)
	if err != nil {
		return nil, err
	}
	output["funds"], err = json.Marshal(r.Funds)
	if err != nil {
		return nil, err
	}
	output["success"], err = json.Marshal(r.Success)
	if err != nil {
		return nil, err
	}
	output["errorType"], err = json.Marshal(r.ErrorType)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (r *TransactionResult) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	var val json.RawMessage
	val = func() json.RawMessage {
		if v, ok := fields["transaction"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Transaction); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for transaction")
	}
	val = func() json.RawMessage {
		if v, ok := fields["funds"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Funds); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for funds")
	}
	val = func() json.RawMessage {
		if v, ok := fields["success"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Success); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for success")
	}
	val = func() json.RawMessage {
		if v, ok := fields["errorType"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.ErrorType); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for errorType")
	}
	return nil
}

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

type Funds struct {
	// Account name
	Account string `json:"account"`

	Balance Bytes `json:"balance"`
}

const FundsAvroCRC64Fingerprint = "W\xc0\xa2X\xf2\r\x97\xbc"

func NewFunds() Funds {
	r := Funds{}
	return r
}

func DeserializeFunds(r io.Reader) (Funds, error) {
	t := NewFunds()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeFundsFromSchema(r io.Reader, schema string) (Funds, error) {
	t := NewFunds()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeFunds(r Funds, w io.Writer) error {
	var err error
	err = vm.WriteString(r.Account, w)
	if err != nil {
		return err
	}
	err = vm.WriteBytes(r.Balance, w)
	if err != nil {
		return err
	}
	return err
}

func (r Funds) Serialize(w io.Writer) error {
	return writeFunds(r, w)
}

func (r Funds) Schema() string {
	return "{\"fields\":[{\"doc\":\"Account name\",\"name\":\"account\",\"type\":\"string\"},{\"name\":\"balance\",\"type\":{\"logicalType\":\"decimal\",\"precision\":9,\"scale\":2,\"type\":\"bytes\"}}],\"name\":\"org.kafkainaction.Funds\",\"type\":\"record\"}"
}

func (r Funds) SchemaName() string {
	return "org.kafkainaction.Funds"
}

func (_ Funds) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ Funds) SetInt(v int32)       { panic("Unsupported operation") }
func (_ Funds) SetLong(v int64)      { panic("Unsupported operation") }
func (_ Funds) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ Funds) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ Funds) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ Funds) SetString(v string)   { panic("Unsupported operation") }
func (_ Funds) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *Funds) Get(i int) types.Field {
	switch i {
	case 0:
		w := types.String{Target: &r.Account}

		return w

	case 1:
		w := BytesWrapper{Target: &r.Balance}

		return w

	}
	panic("Unknown field index")
}

func (r *Funds) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *Funds) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ Funds) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ Funds) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ Funds) HintSize(int)                     { panic("Unsupported operation") }
func (_ Funds) Finalize()                        {}

func (_ Funds) AvroCRC64Fingerprint() []byte {
	return []byte(FundsAvroCRC64Fingerprint)
}

func (r Funds) MarshalJSON() ([]byte, error) {
	var err error
	output := make(map[string]json.RawMessage)
	output["account"], err = json.Marshal(r.Account)
	if err != nil {
		return nil, err
	}
	output["balance"], err = json.Marshal(r.Balance)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (r *Funds) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	var val json.RawMessage
	val = func() json.RawMessage {
		if v, ok := fields["account"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Account); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for account")
	}
	val = func() json.RawMessage {
		if v, ok := fields["balance"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Balance); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for balance")
	}
	return nil
}

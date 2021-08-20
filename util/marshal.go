package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

// Marshal is the SHIP serialization
func Marshal(v interface{}) ([]byte, error) {
	e := make([]map[string]interface{}, 0)

	for _, f := range structs.Fields(v) {
		if !f.IsExported() {
			continue
		}

		jsonTag := f.Tag("json")
		if f.IsZero() && strings.HasSuffix(jsonTag, ",omitempty") {
			continue
		}

		key := f.Name()
		if jsonTag != "" {
			key = strings.TrimSuffix(jsonTag, ",omitempty")
		}

		m := map[string]interface{}{key: f.Value()}
		e = append(e, m)
	}

	return json.Marshal(e)
}

// Unmarshal is the SHIP de-serialization
func Unmarshal(data []byte, v interface{}) error {
	var ar []map[string]json.RawMessage

	// convert input to json array
	if data[0] != byte('[') {
		data = append([]byte{'['}, append(data, ']')...)
	}
	if err := json.Unmarshal(data, &ar); err != nil {
		return err
	}

	// convert array elements to struct members
	for _, ae := range ar {
		if len(ae) > 1 {
			return fmt.Errorf("unmarshal: invalid map %v", ae)
		}

		// extract 1-element map
		var key string
		var val json.RawMessage
		for k, v := range ae {
			key = k
			val = v
		}

		// fmt.Println("json:", string(val))

		// find field
		var field *structs.Field
		for _, f := range structs.Fields(v) {
			name := f.Name()
			if jsonTag := f.Tag("json"); jsonTag != "" {
				name = strings.TrimSuffix(jsonTag, ",omitempty")
			}

			if name == key {
				field = f
				break
			}
		}

		if field == nil {
			return fmt.Errorf("unmarshal: field not found: %s", key)
		}

		// convert value into pointer to value as interface
		iface := reflect.New(reflect.TypeOf(field.Value())).Interface()

		// use pointer-interface to unmarshal into target type
		if err := json.Unmarshal(val, iface); err != nil {
			return err
		}

		// de-reference
		iface = reflect.ValueOf(iface).Elem().Interface()

		// fmt.Printf("set: %s=%+v (%T)\n", field.Name(), iface, iface)
		if err := field.Set(iface); err != nil {
			// fmt.Printf("err: %v\n", err)
			return err
		}
	}

	return nil
}

package feature

import "reflect"

// populatedFields finds the first non-nil field name for given struct
func populatedFields(cmd interface{}) string {
	res := "Unknown"

	cmdFields := reflect.TypeOf(cmd)
	cmdValues := reflect.ValueOf(cmd)
	for i := 0; i < cmdFields.NumField(); i++ {
		if !cmdValues.Field(i).IsNil() {
			res = cmdFields.Field(i).Name
			break
		}
	}

	return res
}

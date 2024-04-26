
package calculateTax

import (
	"encoding/json"
	// "fmt"
	"reflect"
	"errors"
)

func validateKey(jsonBytes []byte) error {
	var data map[string]interface{}

	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return err
	}

	// Custom validation
	if err := validateMap(data); err != nil {
		return err
	}


	return nil
}

func validateMap(data map[string]interface{}) error {
	

	// Convert the map to a struct
	var incomeData IncomeData
	for key, _:= range data {
		if !hasField(&incomeData, key) {
			return errors.New("Invalid JSON data " + key)
		}
	}
	return nil
}

func hasField(s interface{}, fieldName string) bool {
	t := reflect.TypeOf(s)
	for i := 0; i < t.Elem().NumField(); i++ {
		if t.Elem().Field(i).Tag.Get("json") == fieldName {
			return true
		}
	}
	return false
}

func validateKeyReqAdmin(jsonBytes []byte) error {
	var data map[string]interface{}

	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return err
	}

	// Custom validation
	if err := validateMapReqAdmin(data); err != nil {
		return err
	}


	return nil
}

func validateMapReqAdmin(data map[string]interface{}) error {
	

	// Convert the map to a struct
	var incomeData Request_amount_new 
	for key, _:= range data {
		if !hasField(&incomeData, key) {
			return errors.New("Invalid JSON data " + key)
		}
	}
	return nil
}
package helper

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"time"
)

func Init(serviceName string) {}

func StructToMap(obj interface{}, size int) (map[string]interface{}, error) {
	result := make(map[string]interface{}, size)

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func MapToStruct(m map[string]string, obj interface{}) error {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("obj must be a non-nil pointer to a struct")
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return errors.New("obj must be a pointer to a struct")
	}

	typ := val.Type()

	// Create a map[string]interface{} for JSON unmarshaling
	convertedMap := make(map[string]interface{})

	for key, value := range m {
		var matchedField reflect.StructField
		found := false

		// Match map key to struct field using JSON tags
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == key {
				matchedField = field
				found = true
				break
			}
		}

		if !found {
			// Key does not match any struct field; keep it as string
			convertedMap[key] = value
			continue
		}

		// Perform type conversion based on the field type
		fieldType := matchedField.Type.Kind()
		switch fieldType {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// Convert to int
			if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
				convertedMap[key] = intValue
			} else {
				convertedMap[key] = value // fallback to string if conversion fails
			}
		case reflect.Float32, reflect.Float64:
			// Convert to float
			if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
				convertedMap[key] = floatValue
			} else {
				convertedMap[key] = value // fallback to string if conversion fails
			}
		case reflect.Bool:
			// Convert to bool
			if boolValue, err := strconv.ParseBool(value); err == nil {
				convertedMap[key] = boolValue
			} else {
				convertedMap[key] = value // fallback to string if conversion fails
			}
		default:
			// Default to string
			convertedMap[key] = value
		}
	}

	// Marshal and unmarshal to populate the struct
	jsonBytes, err := json.Marshal(convertedMap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonBytes, obj)
	if err != nil {
		return err
	}

	return nil
}

func TimeNow() (*time.Time, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, err
	}

	now, err := time.Parse("2006-01-02 15:04:05", time.Now().In(loc).Format("2006-01-02 15:04:05"))
	if err != nil {
		return nil, err
	}

	return &now, nil
}

func StringToTime(dateStr string) (time.Time, error) {
	date_str := "2006-01-02 15:04:05"
	date_object, err := time.Parse(date_str, dateStr)
	return date_object, err
}

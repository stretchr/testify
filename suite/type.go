package suite

import (
	"reflect"
)

// GetEmbeddedValue accesses an embedded struct from a given struct value.
func GetEmbeddedValue(value interface{}, embeddedType reflect.Type) reflect.Value {
	// Get the reflect.Value of the input value
	val := reflect.ValueOf(value)

	// Ensure it's a pointer to the struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Ensure the value is a struct
	if val.Kind() != reflect.Struct {
		return reflect.Value{}
	}

	// Iterate through the fields of the struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)

		// Check if the field is anonymous (embedded)
		if field.Anonymous {
			// Check if this field is of the type we're looking for
			if val.Field(i).Type() == embeddedType {
				return val.Field(i)
			}

			// If it's not the exact match, recurse into the embedded struct
			if val.Field(i).Kind() == reflect.Struct {
				// Recurse into embedded struct
				embeddedVal := GetEmbeddedValue(val.Field(i).Interface(), embeddedType)
				if embeddedVal.IsValid() {
					return embeddedVal
				}
			}
		}
	}

	return reflect.Value{}
}

/*
// Function to check if a type embeds another type and return the embedded struct's value
func GetEmbeddedValue(value reflect.Value, embeddedType reflect.Type) reflect.Value {
	valType := value.Type()

	if valType.Kind() == reflect.Pointer {
		value = value.Elem()
		valType = value.Type()
	}

	// Ensure the value is a struct
	if valType.Kind() != reflect.Struct {
		fmt.Printf("is not struct (is %v): %v\n", valType.Kind(), value)
		return reflect.Value{}
	}

	// Iterate through the fields of the struct
	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)
		fmt.Printf("checking if field %#q embeds type %v\n", field.Name, embeddedType)

		// If the field is an anonymous (embedded) struct, check if it matches the embedded type
		if field.Anonymous {
			fieldValue := reflect.ValueOf(value).Field(i)
			if field.Type == embeddedType {
				fmt.Printf("field %#q embeds type %v (kind: %v)\n", field.Name, embeddedType, fieldValue.Type())
				return fieldValue.Elem()
			}
			// Recursively check if the embedded field contains the embedded type
			// We need to recurse into the field's type if it's a struct
			if fieldValue.Kind() == reflect.Struct {
				if embeddedValue := GetEmbeddedValue(fieldValue, embeddedType); embeddedValue.IsValid() {
					return embeddedValue
				}
			}
		}
	}

	// If not found, return an invalid reflect.Value
	return reflect.Value{}
}
*/

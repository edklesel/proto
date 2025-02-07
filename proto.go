package proto

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
)

func templateFields(v reflect.Value, nonce int) {

	// Get the type of the input struct
	objType := v.Type()

	// Cycle through the fields of the struct
	for i := 0; i < objType.NumField(); i++ {

		// Get the name of the field
		field := objType.Field(i)
		val := v.FieldByName(field.Name)
		valType := val.Type()

		var valKind reflect.Kind   // The kind of the value
		var valToSet reflect.Value // The valut to use for reflect.Set
		var fieldType reflect.Type // The type of the field

		// If the value is a pointer, pull out the value that the
		// pointer points to.
		if valType.Kind() == reflect.Pointer {
			valKind = valType.Elem().Kind()
			val.Set(reflect.New(valType.Elem()))
			valToSet = val.Elem()
			fieldType = field.Type.Elem()
			// Otherwise if it's just a value, use this instead
		} else {
			valKind = valType.Kind()
			valToSet = val
			fieldType = field.Type
		}

		// Set the different values for different types of the field
		var stringVal, intVal, boolVal reflect.Value

		var fieldTemplate string = field.Tag.Get("proto")

		stringVal = reflect.ValueOf(fmt.Sprintf("%s_%d", field.Name, nonce))
		intVal = reflect.ValueOf(nonce)
		boolVal = reflect.ValueOf(true)

		// This is the error message, in case we need it later
		var typeError string = fmt.Sprintf("Cannot provide \"%s\" as value for field %s of kind %s", fieldTemplate, field.Name, field.Type.Kind())

		// Set the value based on the type.
		// In each type check block, we also check if the template
		// tag is set, and if it is, use this value instead.

		switch valKind {
		// For structs, recurse
		case reflect.Struct:
			templateFields(valToSet, nonce)

			// For strings, we format the name of the field and the nonce
		case reflect.String:

			switch fieldTemplate == "" {
			case true:
				valToSet.Set(stringVal)
			case false:
				valToSet.Set(reflect.ValueOf(fieldTemplate))
			}

		// For ints, we just use the nonce
		case reflect.Int:
			switch fieldTemplate == "" {
			case true:
				valToSet.Set(intVal)
			case false:
				intVal, err := strconv.Atoi(fieldTemplate)
				if err != nil {
					panic(typeError)
				}
				valToSet.Set(reflect.ValueOf(intVal))

			}

		// For bools, use true as the default
		case reflect.Bool:
			switch fieldTemplate == "" {
			case true:
				valToSet.Set(boolVal)
			case false:
				boolVal, err := strconv.ParseBool(fieldTemplate)
				if err != nil {
					panic(err)
				}
				valToSet.Set(reflect.ValueOf(boolVal))
			}

		// For slices, we create a slice and append a single
		// element of the relevant type to it.
		// If the template tag is set on the field, split the
		// value by command and append each element to an
		// empty slice.
		case reflect.Slice:

			valToSet.Set(reflect.MakeSlice(fieldType, 0, 0))

			// Determine the kind of the slice elements
			sliceType := reflect.TypeOf(valToSet.Interface()).Elem()
			sliceKind := sliceType.Kind()

			switch sliceKind {
			case reflect.String:

				switch fieldTemplate == "" {
				case true:
					valToSet.Set(reflect.Append(valToSet, stringVal))
				case false:
					stringVals := strings.Split(fieldTemplate, ",")
					for i := 0; i < len(stringVals); i++ {
						stringVal = reflect.ValueOf(stringVals[i])
						valToSet.Set(reflect.Append(valToSet, stringVal))
					}
				}

			case reflect.Int:

				switch fieldTemplate == "" {
				case true:
					valToSet.Set(reflect.Append(valToSet, intVal))
				case false:
					intVals := strings.Split(fieldTemplate, ",")

					for i := 0; i < len(intVals); i++ {
						intVal, err := strconv.Atoi(intVals[i])
						if err != nil {
							panic(typeError)
						}
						valToSet.Set(reflect.Append(valToSet, reflect.ValueOf(intVal)))
					}
				}

				// For slices of structs, we template the struct and append
				// to the slice.
			case reflect.Struct:
				newStruct := reflect.New(sliceType).Elem()
				_template(newStruct, nonce)
				valToSet.Set(reflect.Append(valToSet, newStruct))
			}
		}

	}

}

func Template(v interface{}) int {

	var nonce int = rand.Intn(1e5)

	objValue := reflect.ValueOf(v).Elem()
	_template(objValue, nonce)

	return nonce
}

func _template(v reflect.Value, nonce int) {
	templateFields(v, nonce)
}

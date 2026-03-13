package proto

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"strconv"
	"strings"

	version "github.com/hashicorp/go-version"
)

func modifyFields(v reflect.Value, ver *version.Version) {

	objType := v.Type()

	for i := 0; i < objType.NumField(); i++ {

		field := objType.Field(i)
		val := v.FieldByName(field.Name)
		valType := val.Type()

		var valKind reflect.Kind  
		var valToSet reflect.Value
		var fieldType reflect.Type

		switch valType.Kind() {
		case reflect.Pointer:
			valKind = valType.Elem().Kind()
			valToSet = val.Elem()
			fieldType = field.Type.Elem()
		default:
			valKind = valType.Kind()
			valToSet = val
			fieldType = field.Type
		}

		var protoTag string = field.Tag.Get("proto")
		var modTag string = field.Tag.Get("proto.modify")
		var constraintTag string = field.Tag.Get("proto.constraint")

		// This is the error message, in case we need it later
		var typeError string = fmt.Sprintf("Cannot provide \"%s\" as value for field %s of kind %s", modTag, field.Name, field.Type.Kind())

		// If the proto tag is specified but not the modify tag,
		// then skip modifying it.
		if protoTag != "" && modTag == "" {
			continue
		}

		// If the version is provided and the constraint tag is also provided,
		// compare the version to the version in the proto tag
		if ver != nil && constraintTag != "" {
			constraints, err := version.NewConstraint(constraintTag)
			if err != nil {
				panic(fmt.Sprintf("Unable to parse version constraint for field %s: %s", field.Name, err))
			}

			// If the version doesn't match the constraint, then we skip this field
			if !constraints.Check(ver) {
				continue
			}
		}

		currentVal := fmt.Sprintf("%v", valToSet.Interface())

		switch valKind {
		case reflect.Struct:
			modifyFields(valToSet, ver)

		case reflect.String:
			var newVal string
			switch modTag == "" {
			case true:
				newVal = fmt.Sprintf("%v_Updated", currentVal)
			case false:
				newVal = modTag
			}
			valToSet.Set(reflect.ValueOf(newVal))
		
		case reflect.Int:
			var newVal int
			switch modTag == "" {
			case true:
				newVal, _ = strconv.Atoi(currentVal)
				newVal += 1
			case false:
				var err error
				newVal, err = strconv.Atoi(modTag)
				if err != nil {
					panic(typeError)
				}
			}
			valToSet.Set(reflect.ValueOf(newVal))
		
		case reflect.Bool:
			var newVal bool
			switch modTag == "" {
			case true:
				currentVal := valToSet.Interface()
				if currentVal == true {
					newVal = false
				} else {
					newVal = true
				}
			case false:
				var err error
				newVal, err = strconv.ParseBool(modTag)
				if err != nil {
					panic(err)
				}
			}
			valToSet.Set(reflect.ValueOf(newVal))
		
		case reflect.Slice:

			sliceType := reflect.TypeOf(valToSet.Interface()).Elem()
			sliceKind := sliceType.Kind()

			switch sliceKind {
			case reflect.String:
				var newVal string
				switch modTag == "" {
				case true:
					switch valToSet.Len() {
					case 0: newVal = fmt.Sprintf("%s_%d_Updated", field.Name, rand.IntN(1e5))
					default: newVal = fmt.Sprintf("%s_Updated", valToSet.Index(0).Interface())
					}
					valToSet.Set(reflect.Append(valToSet, reflect.ValueOf(newVal)))
				case false:
					valToSet.Set(reflect.MakeSlice(fieldType, 0, 0))
					stringVals := strings.Split(modTag, ",")
					for i := 0; i < len(stringVals); i++ {
						stringVal := reflect.ValueOf(stringVals[i])
						valToSet.Set(reflect.Append(valToSet, stringVal))
					}
				}
			case reflect.Int:

				switch modTag == "" {
				case true:
					var newVal int
					switch valToSet.Len() {
					case 0: newVal = rand.IntN(1e5)
					default:
						newVal, _ = strconv.Atoi(fmt.Sprintf("%v",valToSet.Index(0).Interface()))
						newVal += 1
					}
					valToSet.Set(reflect.Append(valToSet, reflect.ValueOf(newVal)))
				case false:
					valToSet.Set(reflect.MakeSlice(fieldType, 0, 0))
					intVals := strings.Split(modTag, ",")
					for i := 0; i < len(intVals); i++ {
						intVal, err := strconv.Atoi(intVals[i])
						if err != nil {
							panic(typeError)
						}
						valToSet.Set(reflect.Append(valToSet, reflect.ValueOf(intVal)))
					}
				}

				// The modify tag isn't supported for slices of structs
			case reflect.Struct:
				newVal := reflect.New(sliceType).Elem()
				_prototype(newVal, rand.IntN(1e3), nil)
				valToSet.Set(reflect.Append(valToSet, newVal))
			}
		}
	}
}

func _modify(v reflect.Value, ver *version.Version) {
	modifyFields(v, ver)
}

func Modify(v interface{}) {

	objValue := reflect.ValueOf(v).Elem()
	_modify(objValue, nil)

}

func ModifyWithVersion(v interface{}, ver *version.Version) {

	objValue := reflect.ValueOf(v).Elem()
	_modify(objValue, ver)

}
package proto

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"strconv"
)

func modifyFields(v reflect.Value) {

	objType := v.Type()

	for i := 0; i < objType.NumField(); i++ {

		field := objType.Field(i)
		val := v.FieldByName(field.Name)
		valType := val.Type()

		var valKind reflect.Kind  
		var valToSet reflect.Value

		switch valType.Kind() {
		case reflect.Pointer:
			valKind = valType.Elem().Kind()
			valToSet = val.Elem()
		default:
			valKind = valType.Kind()
			valToSet = val
		}

		currentVal := fmt.Sprintf("%v", valToSet.Interface())

		// fmt.Println(currentVal)

		switch valKind {
		case reflect.Struct:
			modifyFields(valToSet)

		case reflect.String:
			var newVal string
			newVal = fmt.Sprintf("%v_Updated", currentVal)
			valToSet.Set(reflect.ValueOf(newVal))
		
		case reflect.Int:
			var newVal int
			newVal, _ = strconv.Atoi(currentVal)
			newVal += 1
			valToSet.Set(reflect.ValueOf(newVal))
		
		case reflect.Bool:
			currentVal := valToSet.Interface()
			var newVal bool
			if currentVal == true {
				newVal = false
			} else {
				newVal = true
			}
			valToSet.Set(reflect.ValueOf(newVal))
		
		case reflect.Slice:

			sliceType := reflect.TypeOf(valToSet.Interface()).Elem()
			sliceKind := sliceType.Kind()
			// var newVal any
			switch sliceKind {
			case reflect.String:
				var newVal string
				switch valToSet.Len() {
				case 0: newVal = fmt.Sprintf("%s_%d_Updated", field.Name, rand.IntN(1e5))
				default: newVal = fmt.Sprintf("%s_Updated", valToSet.Index(0).Interface())
				}
				valToSet.Set(reflect.Append(valToSet, reflect.ValueOf(newVal)))
			case reflect.Int:
				var newVal int
				switch valToSet.Len() {
				case 0: newVal = rand.IntN(1e5)
				default:
					newVal, _ = strconv.Atoi(fmt.Sprintf("%v",valToSet.Index(0).Interface()))
					newVal += 1
				}
				valToSet.Set(reflect.Append(valToSet, reflect.ValueOf(newVal)))
			case reflect.Struct:
				newVal := reflect.New(sliceType).Elem()
				_template(newVal, rand.IntN(1e3))
				valToSet.Set(reflect.Append(valToSet, newVal))
			}
		}
	}
}

func _modify(v reflect.Value) {
	modifyFields(v)
}

func Modify(v interface{}) {

	objValue := reflect.ValueOf(v).Elem()
	_modify(objValue)

}
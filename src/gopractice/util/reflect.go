package util

import (
	"fmt"
	"reflect"
)

func SetStructByJson(obj interface{}, mapData map[string]interface{}) error {

	for key, value := range mapData {
		if err := setFiled(obj, key, value); err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}

//这里主要是判断obj的字段类型和value值匹配
func setFiled(obj interface{}, key string, value interface{}) error {
	structData := reflect.ValueOf(obj).Elem()
	fieldValue := structData.FieldByName(key)

	if !fieldValue.IsValid() {
		return fmt.Errorf("utils.setField() No such field: %s in obj ", key)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value ", key)
	}

	fieldType := fieldValue.Type()
	val       := reflect.ValueOf(value)

	valTypeStr   := val.Type().String()
	fieldTypeStr := fieldType.String()
	if valTypeStr == "float64" && fieldTypeStr == "int" {
		val = val.Convert(fieldType)
	} else if fieldType != val.Type() {
		return fmt.Errorf("Provided value type " + valTypeStr + " didn't match obj field type " + fieldTypeStr)
	}
	fieldValue.Set(val)
	return nil
}

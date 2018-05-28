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
	struckData := reflect.ValueOf(obj).Elem()
	filedValue := struckData.FieldByName(key) //这里获取obj的属性
	if !filedValue.IsValid() {
		return fmt.Errorf("util.setFiled() no such filed: %s in obj ", key)
	}

	if !filedValue.CanSet() {
		return fmt.Errorf("cannot set %s filed value ", key)
	}

	filedType := filedValue.Type()
	val := reflect.ValueOf(value)
	valTypeStr := val.Type().String()
	filedTypeStr := filedType.String()

	if valTypeStr == "float64" && filedTypeStr == "int" {
		val.Convert(filedType)
	} else if filedType != val.Type() {
		return fmt.Errorf("provide filed type %v didn't match obj filed type %v", valTypeStr, filedType)
	}

	filedValue.Set(val)
	return nil
}

package main

import (
	"reflect"
	"errors"
	"fmt"
)

type People struct {
	Name string
	Age int
	Sex int
}

var people People

func main() {
	mapData := make(map[string]interface{})
	mapData["Name"] = "张三"
	mapData["Age"] = 18
	mapData["Sex"] = 1
	setPeopleByReflect(&people,mapData)
	fmt.Println(people)
}

func setPeopleByReflect(object interface{}, mapData map[string]interface{}) error {
	for key,value := range mapData {
		if err := setFiled(object,key,value);err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func setFiled(object interface{},key string,value interface{}) error {
	structData := reflect.ValueOf(object).Elem() // 反射获取对应的struct
	filedValue := structData.FieldByName(key)   //获取key对应struct的字段

	if !filedValue.IsValid()  {     //如果没有找到则返回错误
		return errors.New("err1")
	}

	if !filedValue.CanSet() {       //如果不是可设置的则返回错误
		return errors.New("err2")
	}

	filedType := filedValue.Type()   //字段的类型
	val := reflect.ValueOf(value)    //value的值
	valType := val.Type()            //value的类型
	filedTypeStr := filedType.String()  //将类型转成字符串  为了好比较
	valTypeStr := valType.String()    //将类型转成字符串
	if valTypeStr == "float64" || filedTypeStr == "int" {  //判断struct的字段的属性是int  和 value的属性是float64
		val = val.Convert(filedType)   //将value的类型转成struct对应属性的类型
	}else if valType != filedType {
		return errors.New("err3")   //如果不是对应的类型则抛出error
	}
	filedValue.Set(val)                  //给struct的属性设置值
	return nil
}

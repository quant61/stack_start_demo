package main

import (
	"fmt"
	"reflect"
)

func printNonZeroFields(regs interface{}) {
	ref := reflect.ValueOf(regs)
	switch ref.Kind() {
	case reflect.Struct:
		for i := 0; i < ref.NumField(); i++ {
			if !ref.Field(i).IsZero() {
				t := ref.Type().Field(i)
				f := ref.Field(i)
				fmt.Println(i, t.Name, f)
			}
		}
	default:
		fmt.Println("not supported")
	}
}

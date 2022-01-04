package main

import (
	"fmt"
	"io"
	"io/ioutil"
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


func dumpStack(mem io.ReaderAt, pos int64) {
	all := make([]byte, 65536)
	n, _ := mem.ReadAt(all, pos)
	fmt.Println(n, uintptr(pos)+uintptr(n))
	ioutil.WriteFile("stack", all[:n], 0644)

	all = make([]byte, 65536)
	n, _ = mem.ReadAt(all, pos&^4095)
	fmt.Println(n, uintptr(pos)+uintptr(n))
	ioutil.WriteFile("stack_aligned", all[:n], 0644)
}


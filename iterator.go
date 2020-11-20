package wasm

import (
	"fmt"
	"syscall/js"
)

// Iterator permite iterar sobre un js.Value cuando se sabe de
// antemano que este equivale a un array/slice
type Iterator js.Value

// El método String obtiene un slice de strings de un js.Value
func (e Iterator) String() []string {
	return extractStringAndAdd(js.Value(e))
}

// El método String obtiene un slice de int64 de un js.Value
func (e Iterator) Int() []int64 {
	return extractIntAndAdd(js.Value(e))
}

// El método String obtiene un slice de float64 de un js.Value
func (e Iterator) Float() []float64 {
	return extractFloatAndAdd(js.Value(e))
}

func extractStringAndAdd(values js.Value) []string {
	var attribute []string
	if !values.IsNull() {
		for i := 0; i < values.Length(); i++ {
			value := values.Get(fmt.Sprintf("%v", i)).String()
			attribute = append(attribute, value)
		}
		return attribute
	}

	return attribute
}

func extractIntAndAdd(values js.Value) []int64 {
	var attribute []int64
	if !values.IsNull() {
		for i := 0; i < values.Length(); i++ {
			value := int64(values.Get(fmt.Sprintf("%v", i)).Int())
			attribute = append(attribute, value)
		}
		return attribute
	}

	return attribute
}

func extractFloatAndAdd(values js.Value) []float64 {
	var attribute []float64
	if !values.IsNull() {
		for i := 0; i < values.Length(); i++ {
			value := values.Get(fmt.Sprintf("%v", i)).Float()
			attribute = append(attribute, value)
		}
		return attribute
	}

	return attribute
}

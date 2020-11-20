package wasm

import (
	"encoding/json"
	"errors"
	"fmt"
	"syscall/js"
)

// SessionStorage llama a las funciones nativas del objeto javascript sessionStorage
//
// Example. La siguiente función
//
//  SessionStorage("setItem", "myItem", "itemValue")
//
// es equivalente a la función js
//
//  sessionStorage.setItem('myItem', 'itemValue')
//
func SessionStorage(action string, args ...interface{}) js.Value {
	return js.Global().Get("sessionStorage").Call(action, args...)

}

// LocalStorage llama a las funciones nativas del objeto javascript sessionStorage
//
// Example. La siguiente función
//
//  LocalStorage("setItem", "myItem", "itemValue")
//
// es equivalente a la función js
//
//  localStorage.setItem('myItem', 'itemValue')
//
func LocalStorage(action string, args ...interface{}) js.Value {
	return js.Global().Get("localStorage").Call(action, args...)

}

// ObjectToStorage es capaz de guardar distintos valores en el storage del navegador.
// Está comprobada su eficacia con structs.
//
// El parámetro storageType selecciona el tipo de storage y solo recibe los valores
// "localStorage" o "sessionStorage". El parámetro nameItem nombra al item en el
// storage elegido y puede ser cualquier string. El parámetro object es el valor
// que será introducido en el storage como un string.
//
// WARNING: valores distintos a "localStorage" o a "sessionStorage" invocarán un panic()
//
func ObjectToStorage(storageType, nameItem string, object interface{}) {
	storage := js.Global().Get(storageType)
	data, _ := json.Marshal(object)

	storage.Call("setItem", nameItem, string(data))

}

// @Experimental
//
// ImportToStorage construir un objeto a partir de un item del storage.
func ImportToStorage(storageType, nameItem string, v interface{}) error {
	item := js.Global().Get(storageType).Call("getItem", nameItem)
	fmt.Println(item.String())
	if !item.Truthy() {
		return errors.New("invalid storage operation")
	}

	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)

}

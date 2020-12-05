// +build js,wasm

// This package contains some functions to be simple WebAssembly applications developed with go language.
// Allows to access and manipulate the LocalStorage, IndexedDB and DOM.
//
package wasm

import "syscall/js"

// Lanza una alerta del navegador
func Alert(msj string) {
	js.Global().Call("alert", msj)
}

// Console imprime mensaje en la consola. El parámetro action define el tipo de mensaje:
// -- "log"		ejecutará un console.log()
// -- "error" 	ejecutará un console.error()
// -- "warn" 	ejecutará un console.warn()
// WARNING: un valor distinto a los anteriores como argumento action invocará un panic()
func Console(action string, msj interface{}) {
	js.Global().Get("console").Call(action, msj)
}

// Recarga la página
func Reload() {
	js.Global().Get("location").Call("reload")
}

// Devuelve un js.Value que corresponde al objeto HTML con el id señalado.
func GetElementById(id string) js.Value {
	return js.Global().Get("document").Call("getElementById", id)
}

// Devuelve el valor de un input en un formulario HTML.
func InputValue(id string) string {
	return js.Global().Get("document").Call("getElementById", id).Get("value").String()
}

// Selecciona un DOMObject y le inserta un contenido.
//
// Example. La función
//
//  wasm.InnerHTML("#my_div", "<p>Hello Word</p>")
//
// Es equivalente al siguiente script js:
//
//  document.querySelector('#my_div').innerHTML = '<p>Hello Word</p>'
//
func InnerHTML(selector, content string) {
	js.Global().Get("document").Call("querySelector", selector).Set("innerHTML", content)
}

// Selecciona un DOMObject y le añade código HTML.
//
// Example. La función
//
//  wasm.AddHTML("#my_div", "<p>Hello Word</p>")
//
// Es equivalente al siguiente script js:
//
//  document.querySelector('#my_div').innerHTML += '<p>Hello Word</p>'
//
func AddHTML(selector, content string) {
	c := js.Global().Get("document").Call("querySelector", selector).Get("innerHTML").String()
	c += content
	InnerHTML(selector, c)
}

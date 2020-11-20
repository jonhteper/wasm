//+build js,wasm

package wasm

import (
	"fmt"
	"syscall/js"
)

// IndexeDB es un objeto que permite realizar consultas en la API de navegador web del mismo nombre.
// Para más información revisar la documentación sobre indexeDB en:
//  https://developer.mozilla.org/es/docs/Web/API/IndexedDB_API/Usando_IndexedDB.
type IndexeDB struct {
	Name string
}

// IndexConsult realiza la consulta a una indexeDB y permite manipular el objeto encontrado. Esta función
//	sigue el mismo esquema que consulta realizada en javascript vanilla, por lo tanto es necesario
// 	ingresar el nombre del objectStore y el índice para realizar la consulta. Así mismo, para poder manipular
//	los datos obtenidos, es necesario pasar una función a la medida.
func (iDB IndexeDB) IndexConsult(objectStore, index string, onsuccessFunction func(this js.Value)) {
	db := js.Global().Get("window").Get("indexedDB").Call("open", iDB.Name, 1)

	db.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		DB := db.Get("result")
		data := DB.Call("transaction", iDB.Name, "readonly")
		objSt := data.Call("objectStore", objectStore)
		request := objSt.Call("get", index)

		request.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			onsuccessFunction(request)

			return nil
		}))

		return nil
	}))

}

// CustomConsult realiza la consulta a una indexeDB y permite manipular el objeto encontrado. La búsqueda
// se realiza en el índice que coincida con el valor del parámetro index y el valor del parámetro indexValue.
// El parámetro onsuccessFunction permite manipular el objeto encontrado.
func (iDB IndexeDB) CustomConsult(objectStore, index, indexValue string, onsuccessFunction func(this js.Value)) {
	db := js.Global().Get("window").Get("indexedDB").Call("open", iDB.Name, 1)

	db.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		DB := db.Get("result")
		data := DB.Call("transaction", iDB.Name, "readonly")
		objSt := data.Call("objectStore", objectStore)
		indexDB := objSt.Call("index", index)
		request := indexDB.Call("get", indexValue)

		request.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			onsuccessFunction(request)

			return nil
		}))

		return nil
	}))
}

// GetItems permite manipular los objetos dentro de una indexeDB. Sigue el mismo esquema que una consulta
// con js, por lo que es necesario indicar el objectStore. El parámetro onsuccessFunction recibe como argumento
// el cursor resultante (el resultado del evento onsuccess), sobre el mismo se puede acceder a cada uno de
// los objetos guardados en la indexeDB.
func (iDB IndexeDB) GetItems(objectStore string, onsuccessFunction func(e js.Value)) {
	db := js.Global().Get("window").Get("indexedDB").Call("open", iDB.Name, 1)
	db.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		DB := db.Get("result")
		objSt := DB.Call("transaction", iDB.Name).Call("objectStore", objectStore)
		objSt.Call("openCursor").Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			cursor := args[0].Get("target").Get("result") //<-- args[0] Es el evento
			onsuccessFunction(cursor)
			return nil
		}))

		return nil
	}))
}

// NewElement es un método que permite agregar o actualizar un elemento de una indexeDB
func (iDB IndexeDB) NewElement(objectStore string, element map[string]interface{}, log Log) {
	db := js.Global().Get("window").Get("indexedDB").Call("open", iDB.Name, 1)

	db.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		msj := fmt.Sprintf("impossible open %v db", iDB.Name)
		_ = log.Add(msj)
		Console("warn", msj)

		return nil
	}))

	db.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		DB := db.Get("result")
		trans := DB.Call("transaction", iDB.Name, "readwrite")
		objSt := trans.Call("objectStore", objectStore)

		objSt.Call("add", element)

		trans.Set("oncomplete", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			fmt.Println("Element type ", iDB.Name, " created")
			return nil
		}))

		trans.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			updateElement(iDB.Name, objectStore, DB, element, log)
			return nil
		}))

		return nil
	}))
}

// Intenta crear una indexDB con un objectStore, en caso de existir una con el mismo nombre será actualizada.
// Los parámetros name y objectStore nombrarán la base de datos y al objectStore creados. Por defecto al
// objectStore se le asigna un keyPath autoincrementable con valor "id". El parámetro onupgradeneededFunc es
// una función que recibirá como argumento el objectStore resultante de la consulta. Con base en él es posible
// crear los índices de la base de datos. El parámetro log se encarga de obtener los errores que puedan
// resultar.
func CreateDB(name, objectStore string, onupgradeneededFunc func(objSt js.Value), log Log) {
	db := js.Global().Get("window").Get("indexedDB").Call("open", name, 1)

	db.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		msj := fmt.Sprintf("impossible build %v db", name)
		_ = log.Add(msj)
		Console("warn", msj)
		return nil
	}))

	db.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println(name, " db available")
		return nil
	}))

	db.Set("onupgradeneeded", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		DB := args[0].Get("target").Get("result")
		params := map[string]interface{}{"keyPath": "id", "autoIncrement": true}
		objSt := DB.Call("createObjectStore", objectStore, params)

		onupgradeneededFunc(objSt)

		fmt.Println(name, " db created")
		return nil
	}))
}

func updateElement(nameDB, objectStore string, result js.Value, element map[string]interface{}, log Log) {
	trans := result.Call("transaction", nameDB, "readwrite")
	objSt := trans.Call("objectStore", objectStore)

	objSt.Call("put", element)

	trans.Set("oncomplete", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Element type ", nameDB, " updated")
		return nil
	}))

	trans.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		msj := fmt.Sprintf("failed transaccion in %v DB", nameDB)
		_ = log.Add(msj)
		Console("warn", msj)

		return nil
	}))
}

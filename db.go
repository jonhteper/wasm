//+build js,wasm

package wasm

import (
	"fmt"
	"sync"
	"syscall/js"
)

// IndexedDB es un objeto que permite realizar consultas en la API de navegador web del mismo nombre.
// Para más información revisar la documentación sobre indexedDB en:
//  https://developer.mozilla.org/es/docs/Web/API/IndexedDB_API/Usando_IndexedDB.
type IndexedDB struct {
	Name        string
	ObjectStore string
}

// Count devuelve el número de elementos disponibles en la indexedDB
func (iDB IndexedDB) Count(log Log) int {
	var wg sync.WaitGroup
	c := make(chan int, 1)
	wg.Add(1)
	simpleConsult(iDB, "readonly", func(o js.Value) {
		t := o.Call("count")
		t.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			n := t.Get("result").Int()
			c <- n
			wg.Done()
			return nil
		}))

	}, log)
	wg.Wait()
	return <-c

}

// IndexConsult realiza la consulta a una indexedDB y permite manipular el objeto encontrado. Esta función
//	sigue el mismo esquema que consulta realizada en javascript vanilla, por lo tanto es necesario
// 	ingresar el nombre del objectStore y el índice para realizar la consulta. Así mismo, para poder manipular
//	los datos obtenidos, es necesario pasar una función a la medida.
func (iDB IndexedDB) IndexConsult(indexValue string, onsuccessFunction func(this js.Value), log Log) {
	simpleConsult(iDB, "readonly", func(o js.Value) {
		request := o.Call("get", indexValue)

		request.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			onsuccessFunction(request)
			return nil
		}))

		request.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			_ = log.Add("item " + indexValue + " not exist")
			return nil
		}))
	}, log)
}

// CustomConsult realiza la consulta a una indexedDB y permite manipular el objeto encontrado. La búsqueda
// se realiza en el índice que coincida con el valor del parámetro index y el valor del parámetro indexValue.
// El parámetro onsuccessFunction permite manipular el objeto encontrado.
func (iDB IndexedDB) CustomConsult(index string, indexValue interface{}, onsuccessFunction func(this js.Value), log Log) {
	simpleConsult(iDB, "readonly", func(o js.Value) {
		indexDB := o.Call("index", index)
		request := indexDB.Call("get", indexValue)

		request.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			onsuccessFunction(request)

			return nil
		}))

		request.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			_ = log.Add(fmt.Sprintf("item %v not exist", indexValue))
			return nil
		}))
	}, log)

}

// GetItems permite manipular los objetos dentro de una indexeDB. Sigue el mismo esquema que una consulta
// con js, por lo que es necesario indicar el objectStore. El parámetro onsuccessFunction recibe como argumento
// el cursor resultante (el resultado del evento onsuccess), sobre el mismo se puede acceder a cada uno de
// los objetos guardados en la indexedDB. Tomar en cuenta que la transaction tiene permisos 'readwrite'.
func (iDB IndexedDB) GetItems(onsuccessFunction func(e js.Value), log Log) {
	db := js.Global().Get("window").Get("indexedDB").Call("open", iDB.Name, 1)

	db.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		msj := fmt.Sprintf("impossible open %v db", iDB.Name)
		_ = log.Add(msj)
		Console("warn", msj)
		return nil
	}))

	db.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		DB := db.Get("result")
		objSt := DB.Call("transaction", iDB.ObjectStore, "readwrite").Call("objectStore", iDB.ObjectStore)
		objSt.Call("openCursor").Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			cursor := args[0].Get("target").Get("result") //<-- args[0] Es el evento
			onsuccessFunction(cursor)
			return nil
		}))

		return nil
	}))
}

// NewElement es un método que permite agregar o actualizar un elemento de una indexedDB.
func (iDB IndexedDB) NewElement(element map[string]interface{}, log Log) {
	db := js.Global().Get("window").Get("indexedDB").Call("open", iDB.Name, 1)

	db.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		msj := fmt.Sprintf("impossible open %v db", iDB.Name)
		_ = log.Add(msj)
		Console("warn", msj)

		return nil
	}))

	db.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		DB := db.Get("result")
		trans := DB.Call("transaction", iDB.ObjectStore, "readwrite")
		objSt := trans.Call("objectStore", iDB.ObjectStore)

		objSt.Call("add", element)

		trans.Set("oncomplete", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			fmt.Println("Element type ", iDB.Name, " created")
			return nil
		}))

		trans.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			updateElement(iDB.Name, iDB.ObjectStore, DB, element, log)
			return nil
		}))

		return nil
	}))
}

// DeleteElement permite eliminar items de la base de datos. La eliminación solo puede realizarse con base en
// el índice principal.
//
// WARNING: Actualmente el método delete siempre retorna success. Por lo mismo, el mensaje no debe ser tomado
// por cierto.
func (iDB IndexedDB) DeleteElement(indexValue string, log Log) {
	simpleConsult(iDB, "readwrite", func(o js.Value) {
		request := o.Call("delete", indexValue)

		request.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			fmt.Println("Element ", indexValue, " has been deleted")
			return nil
		}))

	}, log)
}

// Clear borra todos los elementos de la indexedDB.
func (iDB IndexedDB) Clear(log Log) {
	simpleConsult(iDB, "readwrite", func(o js.Value) {
		o.Call("clear")
	}, log)
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

func simpleConsult(iDB IndexedDB, permission string, onsuccessFunc func(objSt js.Value), log Log) {
	db := js.Global().Get("window").Get("indexedDB").Call("open", iDB.Name, 1)
	db.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		msj := fmt.Sprintf("impossible open %v db", iDB.Name)
		_ = log.Add(msj)
		Console("warn", msj)

		return nil
	}))

	db.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		DB := db.Get("result")
		trans := DB.Call("transaction", iDB.ObjectStore, permission)
		objSt := trans.Call("objectStore", iDB.ObjectStore)

		onsuccessFunc(objSt)

		return nil
	}))
}

func updateElement(nameDB, objectStore string, result js.Value, element map[string]interface{}, log Log) {
	trans := result.Call("transaction", objectStore, "readwrite")
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

//+build js,wasm

package main

import (
	"fmt"
	"github.com/jonhteper/wasm"
	"syscall/js"
)

func main() {
	fmt.Println("Hello World from WASM")
	log := wasm.NewLog("tryLog", "authLog", 10, func(data string) error {
		fmt.Println(data)
		return nil
	})

	// Exported functions

	js.Global().Set("testLog", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("starting Log test...")
		go func() {
			_ = log.Add("try message in log")
		}()
		wasm.InnerHTML("#test_box_log", "Log pass test 😀")

		return nil
	}))

	js.Global().Set("testDB", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("starting DB tests...")
		go func() {
			wasm.CreateDB("try", "try", func(objSt js.Value) {
				objSt.Call("createIndex", "name", "name", map[string]interface{}{"unique": false})
			}, log)

			db := wasm.IndexedDB{Name: "try", ObjectStore: "try"}
			db.NewElement(map[string]interface{}{"id": "myData1", "name": "MyName"}, log)
			db.NewElement(map[string]interface{}{"id": "myData2", "name": "MyName2"}, log)
		}()
		wasm.InnerHTML("#test_box_db", "DB pass test 😀")

		return nil
	}))

	js.Global().Set("testDB2", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("starting DB tests 2...")
		go func() {
			db := wasm.IndexedDB{Name: "try", ObjectStore: "try"}
			n := db.Count(log)
			fmt.Println(n)
			wasm.InnerHTML("#test_box_db2", "DB pass test 😀")
		}()
		return nil
	}))
	<-make(chan interface{})
}

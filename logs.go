package wasm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Log es un objeto que permite crear archivos de registro para las aplicaciones web, utiliza el
// localStorage. También es utilizado en este paquete para recoger errores en ciertas funciones.
// Cada Log tiene dos items del localStorage, uno que se guarda los registros y otro que guarda
// la autorización de enviar estos datos; los nombres de cada item son guardados en los atributos
// nameItem y authItem respectivamente. Se asume que los usuarios de la app serán notificados y
// consultados sobre la recolección de datos.
//
// Debido a que el localStorage tiene un límite de almacenamiento, se exige que cada Log tenga un
// número máximo de registros. Este número es guardado en el atributo maxLength.
//
// WARNING: Se recuerda a los desarrolladores consultar las leyes vigentes en los países de sus
// usuarios sobre el almacenar y enviar datos sensibles o que pudieran vulnerar la privacidad de los
// mismos usuarios. Así mismo, se les recuerda que en algunos países podría ser ilegal no informar
// de manera transparente sobre el uso de la telemetría.
type Log struct {
	nameItem  string                  // nombre del item del localStorage que contiene los registros
	authItem  string                  // nombre del item del localStorage que contiene la autorización del usuario
	maxLength int                     // número de registros máximos, nunca debe ser alcanzado
	sendFunc  func(data string) error // función con la que se enviarán los datos recabados
}

// Obtiene el nombre del item del localStorage donde se alojan los registros.
func (l *Log) NameItem() string {
	return l.nameItem
}

// Constructor del objeto Log.
func NewLog(nameItem, auth string, maxLength int, sendFunc func(data string) error) Log {
	return Log{nameItem: nameItem, maxLength: maxLength, sendFunc: sendFunc, authItem: auth}
}

// Agregar un registro al log si no se ha alcanzado el tamaño máximo. Cada registro sigue el esquema:
//  <mensaje>-<fecha hora>
func (l Log) Add(msj string) error {
	logs := LocalStorage("getItem", l.nameItem)

	m := setTime(msj)
	if !logs.Truthy() {
		LocalStorage("setItem", l.nameItem, m)
	} else {
		logs := strings.Split(logs.String(), ",")
		logs = append(logs, m)
		var r string
		for _, v := range logs {
			if v == "" {
				continue
			}

			r += fmt.Sprintf("%v,", v)
		}

		LocalStorage("setItem", l.nameItem, r)

		err := checkSize(len(logs), l)
		if err != nil {
			return l.SendToServer()
		}
	}

	return nil
}

// Obtiene todos los registros y los borra; si hay autorización de parte del usuario se envían con la
// el parámetro sendFunc.
func (l Log) SendToServer() error {
	data := LocalStorage("getItem", l.nameItem).String()
	LocalStorage("setItem", l.nameItem, "")

	auth, err := strconv.ParseBool(LocalStorage("getItem", l.authItem).String())
	if err != nil || !auth {
		return errors.New("forbidden send data")
	}

	err = l.sendFunc(data)
	if err != nil {
		return err
	}

	return nil
}

func setTime(msj string) string {
	t := time.Now()
	return fmt.Sprintf("%v-%v", msj, t.String()[:19])
}

func checkSize(n int, l Log) error {
	if n < l.maxLength {
		return nil
	}

	return fmt.Errorf("log length larger than allowed, log length: %v; log not add", n)
}

package Procesos

import (
	"fmt"
	"io/ioutil"
)

func AbrirArchivoM(ruta string) {
	fmt.Println(ruta)
	b, err := ioutil.ReadFile(ruta) // just pass the file name
	if err != nil {
		fmt.Println("error --")
		fmt.Print(err)
		return
	}

	fmt.Println(b) // print the content as 'bytes'
	str := string(b)
	fmt.Println(str)
}

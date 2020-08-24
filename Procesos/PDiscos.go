package Procesos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
)

type discoDuro struct {
	SizeD int64
}

func CrearDisco(sizeD int64, pathD string, nameD string) {

	if _, err := os.Stat(pathD); os.IsNotExist(err) {
		os.MkdirAll(pathD, os.ModePerm)
	}

	discoNuevo := discoDuro{SizeD: sizeD}

	archivo, err := os.Create(pathD + nameD)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Cannot create the file")
		return
	}
	var cero int64 = 0

	direccion := &cero
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, direccion)
	escribirBytes(archivo, binario.Bytes())
	//Nos posicionamos en el byte 1023 (primera posicion es 0)
	archivo.Seek(discoNuevo.SizeD, 0)

	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, cero)
	escribirBytes(archivo, binario2.Bytes())
	archivo.Seek(0, 0)
}

func EliminarDisco(pathD string) {
	var pru string
	nameD := strings.Split(pathD, "/")
	fmt.Printf("Desea Eliminar %s este disco? Y/N: ", nameD[len(nameD)-1])
	fmt.Scanf("%s", &pru)
	if strings.ToLower(pru) == "y" {
		if _, err := os.Stat(pathD); !os.IsNotExist(err) {
			estadoDelete := os.Remove(pathD)
			if estadoDelete != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Se elimino archivo")
		} else {
			fmt.Println("El archivo no existe")
		}
	} else {
		fmt.Println("No quiso eliminar archivo")
	}
}

func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

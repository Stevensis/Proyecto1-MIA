package Procesos

import "os"

func CrearDisco(sizeD int64, pathD string, nameD string) {
	if _, err := os.Stat(pathD); os.IsNotExist(err) {
		os.MkdirAll(pathD, os.ModePerm)
	}

}

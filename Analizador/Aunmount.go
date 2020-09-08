package Analizador

import (
	"fmt"
	"strconv"
	"strings"
)

func UnmountAnalizador(contenido []string) {
	banderaIsParticion := false
	indexParticion := 0
	for _, elemento := range contenido {
		divi := strings.Split(elemento, "->")
		if string((divi[0])[0])+string((divi[0])[1]) != "id" {
			fmt.Println("Error, se necesita el parametro idn para el comando Unmount")
			return
		}
		entero := strings.Split(divi[0], "id")
		indexDisco, _ := strconv.ParseInt(entero[1], 10, 64)
		indexDisco--
		for i := 0; i < len(Mount_Discos[indexDisco].Lst_PaticionM); i++ {

			if Mount_Discos[indexDisco].Lst_PaticionM[i].Id == divi[1] {
				banderaIsParticion = true
				indexParticion = i
				break
			}

		}
		if banderaIsParticion {
			Mount_Discos[indexDisco].Lst_PaticionM[indexParticion].StatePM = '0'
			fmt.Println("se desmonto, aqui la lista de particiones montadas")
			MostrarMount()
		} else {
			fmt.Println("La particion no esta montada " + divi[1] + " en el disco " + Mount_Discos[indexDisco].NameMD)
		}
	}
}

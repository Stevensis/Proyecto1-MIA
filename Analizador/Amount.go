package Analizador

import (
	"Archivos/Proyecto/Procesos"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func AnalizarMount(contenido []string) {
	pathD := ""
	namePM := ""
	bPath := false
	bName := false
	var particionT Procesos.Particion
	existeP := false
	for i := 0; i < len(contenido); i++ {
		divi := strings.Split(contenido[i], "->")
		switch strings.ToLower(divi[0]) {
		case "path":
			pathD = strings.ReplaceAll(divi[1], "\"", "")
			if _, err := os.Stat(pathD); !os.IsNotExist(err) {
				//fmt.Printf("Si existe este archivo %s", pathP)
			} else {
				ErrorT(pathD, " Este archivo no existe ")
				return
			}
			bPath = true
		case "name":
			namePM = strings.ReplaceAll(divi[1], "\"", "")
			existeP, particionT = Procesos.VerificarParticionMount(pathD, namePM)
			if !existeP {
				fmt.Println("La particion aun no esta creada")
				return
			}
			bName = true
		default:
			ErrorT(divi[0], " No es un parametro de path ")
		}
	}
	if bPath && bName {
		mountPrincipal(namePM, pathD, particionT)
	} else {
		fmt.Println("Falta un path o un path")
		return
	}
}

func mountPrincipal(namaPM string, pathDM string, particionM Procesos.Particion) {
	creadoMD, index := EstaMontadoDisco(pathDM)

	if !creadoMD {
		index = MontarDisco(pathDM)
	}
	MontarParticion(index, particionM)
}

func MontarDisco(pathDM string) int {
	tempMbr := Procesos.OptenerMbr(pathDM)
	for i := 0; i < len(Mount_Discos); i++ {
		if Mount_Discos[i].StateMD == '0' {
			nameD := strings.Split(pathDM, "/")
			Mount_Discos[i].NameMD = nameD[len(nameD)-1]
			Mount_Discos[i].PathMD = pathDM
			Mount_Discos[i].SizeMD = int(tempMbr.SizeD)
			Mount_Discos[i].StateMD = '1'
			return i
		}
	}
	return 0
}

func MontarParticion(index int, particionM Procesos.Particion) {
	for i := 0; i < len(Mount_Discos[index].Lst_PaticionM); i++ {
		if Mount_Discos[index].Lst_PaticionM[i].StatePM == '0' {
			fmt.Println("Monta Particion")
			Mount_Discos[index].Lst_PaticionM[i].PariticionM = particionM
			Mount_Discos[index].Lst_PaticionM[i].StatePM = '1'
			Mount_Discos[index].Lst_PaticionM[i].Id = obtenerId(i, Mount_Discos[index].IdentificadorMD)
			return
		}
	}
}

func obtenerId(index int, id string) string {
	numero := strconv.Itoa(index + 1)
	//fmt.Println(" -*- " + numero)
	return ("vd" + id + numero)
}

func EstaMontadoDisco(pathD string) (bool, int) {
	for i := 0; i < len(Mount_Discos); i++ {
		if Mount_Discos[i].PathMD == pathD {
			return true, i
		}
	}
	return false, 0
}

func MostrarMount() {
	fmt.Println("Lista de particiones montadas")
	for _, elemento := range Mount_Discos {
		if elemento.StateMD == '1' {
			for _, elemt := range elemento.Lst_PaticionM {
				if elemt.StatePM == '1' {
					fmt.Printf("id->%s -path->%s -name->%s\n", elemt.Id, elemento.PathMD, elemt.PariticionM.Name)
				}

			}
		}
	}
}

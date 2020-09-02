package Analizador

import (
	"Archivos/Proyecto/Procesos"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var contador int = 0
var tamno = 0

func OpenArchivo(rutaA string) {

	contador = 0
	analizador(rutaA)
}

func analizador(contenido string) {
	var auxiliar string = ""
	var caso int = 0
	tamno = len(contenido)
	for contador < len(contenido) {
		var c byte = contenido[contador]
		switch caso {
		case 0:
			if unicode.IsLetter(rune(c)) {
				auxiliar += string(c)
				caso = 1
			} else if unicode.IsSpace(rune(c)) {
				caso = 0
			} else {
				ErrorT(string(c), "letra-")
				contador = tamno
			}
		case 1:
			if unicode.IsLetter(rune(c)) {
				auxiliar += string(c)
				caso = 1
			} else if unicode.IsSpace(rune(c)) {
				contador++
				reconocerPalabra(auxiliar, contenido)
				auxiliar = ""
				caso = 0
			} else {
				ErrorT(string(c), "letra--")
				contador = tamno
			}
		default:
			fmt.Printf("este caso no existe \n")
		}
		contador++
	}
}

//Metodo que verifica las palabras reservadas
func reconocerPalabra(palabraR string, contenido string) {
	var te1, pathB string
	switch strings.ToLower(palabraR) {
	case "exec":
		if contenido[contador] == '-' {
			contador++
			if strings.ToLower(extraerString(contenido)) == "path" {
				if contenido[contador] == '-' && contenido[contador+1] == '>' {
					contador += 2
					abrirArchivoM(extraerPath(contenido))
					contador = 50000
				} else {
					ErrorT(string(contenido[contador-1]), "->")
					contador = tamno
				}
			} else {
				ErrorT("n", "se esperaba un path")
			}
		} else {
			ErrorT(string(contenido[contador]), "se esperaba un -")
		}
	case "pause":
		var pru string
		fmt.Print("Pause-:")
		fmt.Scanf("%s", &pru)
		contador--
	case "mkdisk":
		contador++
		if strings.ToLower(extraerString(contenido)) == "size" {
			contador += 2
			var sizeA int64 = extraerInt(contenido)
			if sizeA > 0 {
				contador++
				if contenido[contador] == '-' {
					contador++
					if strings.ToLower(extraerString(contenido)) == "path" {
						if contenido[contador] == '-' && contenido[contador+1] == '>' {
							contador += 2
							pathB = extraerPath(contenido)
							contador += 2
							t1 := strings.ToLower(extraerString(contenido))
							if t1 == "name" {
								if contenido[contador] == '-' && contenido[contador+1] == '>' {
									contador += 2
									idA := extraerId(contenido)
									if strings.Contains(idA, ".dsk") {
										contador++
										if contador >= tamno || contenido[contador] == '\n' {
											sizeA = sizeA * 1024 * 1024
											Procesos.CrearDisco(sizeA, pathB, idA) //Se manda a crear el disco
											return
										}
										contador++
										t2 := strings.ToLower(extraerString(contenido))
										if t2 == "unit" {
											contador += 2
											if contenido[contador] == 'k' {
												sizeA = sizeA * 1024
												Procesos.CrearDisco(sizeA, pathB, idA) //Se manda a crear el disco

											} else if contenido[contador] == 'm' {
												sizeA = sizeA * 1024 * 1024
												Procesos.CrearDisco(sizeA, pathB, idA) //Se manda a crear el disco
											} else {
												ErrorT(string(contenido[contador]), "una letra k o m")
											}
											contador++
										} else {
											ErrorT(t2, "unit")
										}
									} else {
										ErrorT(idA, "no contiene la extencion .dsk ")
									}
								} else {
									ErrorT(string(contenido[contador-1]), "->")
								}
							} else {
								ErrorT(t1, "name")
							}

						} else {
							ErrorT(string(contenido[contador-1]), "->")

						}
					} else {
						ErrorT("n2", "se esperaba un path")
					}
				} else {
					ErrorT(string(contenido[contador]), "se esperaba un --")
				}
			} else {
				ErrorT(string(contenido[contador-1]), "se esperaba un numero positivo mayor que 0")
			}
		} else {
			ErrorT(string(contenido[contador]), "size")
		}
	case "rmdisk":
		contador++
		te1 = strings.ToLower(extraerString(contenido))
		if te1 == "path" {
			contador += 2
			pathB = extraerPath(contenido)
			Procesos.EliminarDisco(pathB)
		} else {
			ErrorT(te1, "path")
		}

	case "fdisk":
		contador++
		linea := extrarLInea(contenido)
		startfdisk(strings.Split(linea, " -"))
	default:
		fmt.Printf("este caso no existe- %s\n", palabraR)
	}
}

//Metodo para trabajar fdisk
func startfdisk(contenido []string) {
	var bUnit, bType, bFit, bDelete, bAdd bool
	var typeP byte
	sizeP := 0
	pathP := ""
	var fitP byte
	deleteP := ""
	nameP := ""
	addS := 0
	//fmt.Printf("%q\n", contenido)
	for i := 0; i < len(contenido); i++ {
		divi := strings.Split(contenido[i], "->")
		switch strings.ToLower(divi[0]) {
		case "size":
			sizeP = extraerInt2(divi[1])
		case "unit":
			if strings.ToLower(divi[1]) == "k" {
				sizeP = sizeP * 1024
			} else if strings.ToLower(divi[1]) == "m" {
				sizeP = sizeP * 1024 * 1024
			} else {
				ErrorT(divi[1], "No es un parametro para unit")
				i = len(contenido)
			}
			bUnit = true
		case "path":
			pathP = strings.ReplaceAll(divi[1], "\"", "")
			if _, err := os.Stat(pathP); !os.IsNotExist(err) {
				//fmt.Printf("Si existe este archivo %s", pathP)
			} else {
				ErrorT(pathP, " Este archivo no existe ")
				i = len(contenido)
			}
		case "type":
			typeP = strings.ToLower(divi[1])[0]
			if typeP == 'p' || typeP == 'e' || typeP == 'l' {
				bType = true
			} else {
				ErrorT(string(typeP), " P,E o L ")
				i = len(contenido)
			}
		case "fit":
			fitPT := strings.ToLower(divi[1])
			if fitPT == "bf" || fitPT == "ff" || fitPT == "wf" {
				switch fitPT {
				case "bf":
					fitP = 'b'
				case "ff":
					fitP = 'f'
				case "wf":
					fitP = 'w'
				}
				bFit = true
			} else {
				ErrorT(fitPT, "No es un parametro de fit")
				i = len(contenido)
			}
		case "delete":
			deleteP = strings.ToLower(divi[1])
			if deleteP == "full" || deleteP == "fast" {
				bDelete = true
			} else {
				ErrorT(deleteP, "no es un parametro de delete")
				i = len(contenido)
			}
		case "name":
			nameP = divi[1]
		case "add":
			addS = extraerInt2(divi[1])
		default:
			ErrorT(divi[0], "No es parametro de fdisk")
			i = len(contenido)
		}
	}

	if !bUnit { //se verifica que la bandera unit
		sizeP = sizeP * 1024
	}

	if !bType {
		typeP = 'p'
	}

	if !bFit {
		fitP = 'w'
	}

	if !bDelete && !bAdd {
		Procesos.CrearParticion(sizeP, pathP, typeP, fitP, nameP)
		return
	}

	if bDelete {
		fmt.Println("\n Eliminar particion" + nameP)
		return
	}

	if bAdd {
		fmt.Println("\n Modificar tama;o particion " + nameP + " " + string(addS))
		return
	}

}

func extraerInt2(contenido string) int {
	i, err := strconv.Atoi(contenido)
	if err == nil {
		return int(i)
	}
	ErrorT(contenido, "No es un numero")
	return 0
}

func extrarLInea(contenido string) string {
	auxiliar := ""
	for contenido[contador] != '\n' && contador < tamno {
		auxiliar = auxiliar + string(contenido[contador])
		contador++
	}
	return auxiliar
}

func extraerId(contenido string) string {
	auxiliar := ""
	for unicode.IsLetter(rune(contenido[contador])) || unicode.IsDigit(rune(contenido[contador])) || contenido[contador] == '_' || contenido[contador] == '.' {
		auxiliar += string(contenido[contador])
		contador++
		if contador >= tamno {
			return auxiliar
		}
	}
	return auxiliar
}

func extraerInt(contenido string) int64 {
	auxiliar := ""
	if contenido[contador] == '-' {
		return 0
	}
	for unicode.IsDigit(rune(contenido[contador])) {
		auxiliar += string(contenido[contador])
		contador++
	}
	i, err := strconv.Atoi(auxiliar)
	if err == nil {
		return int64(i)
	}
	return 0
}

func extraerPath(contenido string) string {
	auxiliar := ""
	var limitador byte
	if contenido[contador] == '"' {
		limitador = '"'
		contador++
		for contenido[contador] != limitador {
			auxiliar += string(contenido[contador])
			contador++
			if contador >= tamno {
				return auxiliar
			}
		}
		contador++
	} else {
		for !unicode.IsSpace(rune(contenido[contador])) {
			auxiliar += string(contenido[contador])
			contador++
			if contador >= tamno {
				return auxiliar
			}
		}
	}

	return auxiliar
}

func extraerString(contenido string) string {
	auxiliar := ""
	for unicode.IsLetter(rune(contenido[contador])) {
		auxiliar += string(contenido[contador])
		contador++
	}
	return auxiliar
}

func ErrorT(caracter string, esperaba string) {
	contador = tamno
	fmt.Printf("Se encontro un: %s y se esperaba: %s\n", caracter, esperaba)
}

func abrirArchivoM(ruta string) {

	b, err := ioutil.ReadFile(ruta) // just pass the file name
	if err != nil {
		fmt.Println("error --")
		fmt.Print(err)
		return
	}
	// convertimos los bits a string
	str := string(b)
	str = strings.ReplaceAll(str, "\\*\n", "") //quitamos del archivo los " \* \n " que pueden venir
	OpenArchivo(str)
}

package Analizador

import (
	"Archivos/Proyecto/Procesos"
	"fmt"
	"strings"
	"unicode"
)

var contador = 0
var tamno = 0

func OpenArchivo(rutaA string) {
	analizador(rutaA)
	contador = 0
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
	switch strings.ToLower(palabraR) {
	case "exec":
		if contenido[contador] == '-' {
			contador++
			if strings.ToLower(extraerString(contenido)) == "path" {
				if contenido[contador] == '-' && contenido[contador+1] == '>' {
					contador += 2
					Procesos.AbrirArchivoM(extraerPath(contenido))
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
	default:
		fmt.Printf("este caso no existe- %s\n", palabraR)
	}
}

func extraerPath(contenido string) string {
	auxiliar := ""
	for contenido[contador] != '\n' {
		auxiliar += string(contenido[contador])
		contador++
		if contador >= tamno {
			return auxiliar
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
	fmt.Printf("Se encontro un: %s y se esperaba: %s\n", caracter, esperaba)
}

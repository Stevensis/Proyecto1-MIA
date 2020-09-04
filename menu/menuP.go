package menu

import (
	"Archivos/Proyecto/Analizador"
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func InicioP() {
	Ca := true
	for Ca {
		fmt.Println("   ----*- Coloque ruta del archivo a seleccionar")
		leer := bufio.NewReader(os.Stdin)
		texto, errL := leer.ReadString('\n')
		if errL != nil {
			fmt.Println(errL)
		} else if texto == "p\n" {
			Ca = false
		} else {
			texto = strings.TrimRight(texto, "\n")
			Analizador.OpenArchivo(texto)
			//fmt.Println("Lectura exitosa")
		}
	}

}

func quitarsaltodelinea(txt string) string {
	if runtime.GOOS == "windows" {
		txt = strings.TrimRight(txt, "\r\n")
	} else {
		txt = strings.TrimRight(txt, "\n")
	}
	return txt
}

package Procesos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	"unsafe"
)

type mbr struct {
	SizeD       int64
	Signature   int8
	Tiempo      [25]byte
	Particiones [4]Particion
}

type Particion struct {
	Status byte
	Type   byte
	Fit    byte
	Start  int64
	Size   int64
	Name   [16]byte
}

type Disk_mount struct {
	PathMD          string
	StateMD         byte
	NameMD          string
	indexMD         int
	SizeMD          int
	IdentificadorMD string
	Lst_PaticionM   [15]Particion_mount
}

type Particion_mount struct {
	Id          string
	PariticionM Particion
	StatePM     byte
}

func CrearDisco(sizeD int64, pathD string, nameD string) {
	if _, err := os.Stat(pathD); os.IsNotExist(err) {
		os.MkdirAll(pathD, os.ModePerm)
	}

	discoNuevo := mbr{SizeD: sizeD}
	discoNuevo.Signature = int8(rand.Intn(254))
	date := time.Now()
	dates := date.Format("01-02-2006 15:04:00")
	copy(discoNuevo.Tiempo[:], dates)
	archivo, err := os.Create(pathD + nameD)
	for i := 3; i > -1; i-- {
		discoNuevo.Particiones[i].Status = 'f'
	}
	if err != nil {
		log.Fatal(err)
		fmt.Println("No se pudo crear el archivo")
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
	binary.Write(&binario2, binary.BigEndian, direccion)
	escribirBytes(archivo, binario2.Bytes())

	//Escribimos el mbr
	archivo.Seek(0, 0)

	estructura := &discoNuevo
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, estructura)
	escribirBytes(archivo, binario3.Bytes())
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

func CrearParticion(sizeP int, pathP string, typeP byte, fitP byte, nameP string) {
	tempM := readFileBinary(pathP)
	mbr_size := unsafe.Sizeof(tempM)
	//fmt.Println(tempM)
	e, p, l, espacioO := verificarParticion(tempM.Particiones)
	if e >= 1 && typeP == 'e' {
		fmt.Println("La particion " + nameP + " es una extendida y ya existe una en este disco")
		return
	}

	if e+p >= 4 && l == 0 { //verificamos que no existan ya 4 particiones creadas
		fmt.Println("Ya existen 4 particiones en este disco")
		return
	}
	existencia, _ := particionExiste(tempM.Particiones, nameP)
	if existencia { //verificamos que no exista una particion con este nombre
		fmt.Print("Ya existe una particion con este nombre " + nameP + " en el disco " + pathP)
		return
	}

	if int(mbr_size)+espacioO+sizeP >= int(tempM.SizeD) {
		fmt.Println("El espacio de la particion " + nameP + " exede el tamaño del disco " + pathP)
		return
	}

	for i := 0; i < len(tempM.Particiones); i++ {
		if tempM.Particiones[i].Status == 'f' {
			switch i {
			case 0:
				if sizeP >= int(tempM.Particiones[1].Size) && tempM.Particiones[1].Size != 0 {
					fmt.Println("El tamaño de la parcion 0 es pequeño para alojar esta particon " + nameP + " del disco " + pathP)
					return
				}
				tempM.Particiones[0].Start = int64(mbr_size)
			default:
				llegaP := tempM.Particiones[i-1].Start + tempM.Particiones[i-1].Size + int64(sizeP)
				iniciaP := tempM.Particiones[i-1].Start + tempM.Particiones[i-1].Size
				if i != 3 {
					if llegaP >= tempM.Particiones[i+1].Start && tempM.Particiones[i+1].Start != 0 {
						fmt.Println("No se puede escribir la particion " + nameP + " en la particion No. " + string(i))
						return
					}
				}

				if llegaP >= tempM.SizeD {
					fmt.Println("esta particion " + nameP + " exede el tamaño del disco " + pathP)
					return
				}

				tempM.Particiones[i].Start = iniciaP
			}
			tempM.Particiones[i].Fit = fitP
			copy(tempM.Particiones[i].Name[:], nameP)
			tempM.Particiones[i].Size = int64(sizeP)
			tempM.Particiones[i].Status = 't'
			tempM.Particiones[i].Type = typeP
			actualizarMbr(pathP, tempM)
			return
		}
	}
}

func EliminarParticion(pathD string, nameP string) {
	tempM := readFileBinary(pathD)
	existencia, index := particionExiste(tempM.Particiones, nameP)
	if !existencia { //verificamos que no exista una particion con este nombre
		fmt.Print("No existe una particion con este nombre " + nameP + " en el disco " + pathD)
		return
	}

	var pru string
	fmt.Printf("Desea Eliminar %s esta particion? Y/N: ", nameP)
	fmt.Scanf("%s", &pru)
	if strings.ToLower(pru) == "y" {
		tempM.Particiones[index].Fit = '0'
		copy(tempM.Particiones[index].Name[:], "")
		tempM.Particiones[index].Size = 0
		tempM.Particiones[index].Start = 0
		tempM.Particiones[index].Status = 'f'
		tempM.Particiones[index].Type = '0'
		actualizarMbr(pathD, tempM)
		fmt.Println("Se elimino la particion")
	} else {
		fmt.Println("No quiso eliminar la particion")
	}

}

func VerificarParticionMount(pathD string, nameP string) (bool, Particion) {
	tempM := readFileBinary(pathD)
	existencia, index := particionExiste(tempM.Particiones, nameP)
	return existencia, tempM.Particiones[index]
}

func OptenerMbr(pathD string) mbr {
	return readFileBinary(pathD)
}

func PruebaContenido(pathD string) {
	tempM := readFileBinary(pathD)
	nameD := strings.Split(pathD, "/") //extrae el nombre del disco
	fmt.Printf("Nombre: %s\n", nameD[len(nameD)-1])
	fmt.Printf("mbr_tamanio: %d\n", tempM.SizeD)
	fmt.Printf("mbr_fecha: %s\n", tempM.Tiempo)
	fmt.Printf("mbt_signature %d\n", tempM.Signature)
	for i := 0; i < len(tempM.Particiones); i++ {
		fmt.Printf("Nombre_Particion_%d: %s\n", i, tempM.Particiones[i].Name)
		fmt.Printf("Particion%d_estado: %c\n", i, tempM.Particiones[i].Status)
		fmt.Printf("Particion%d_tipo: %c\n", i, tempM.Particiones[i].Type)
		fmt.Printf("Particion%d_fit %c\n", i, tempM.Particiones[i].Fit)
		fmt.Printf("Particion%d_start: %d\n", i, tempM.Particiones[i].Start)
		fmt.Printf("Particion%d_size %d\n", i, tempM.Particiones[i].Size)
	}
}

func particionExiste(tempP [4]Particion, nameP string) (bool, int) {
	var nameT [16]byte
	copy(nameT[:], nameP)
	for i := 0; i < len(tempP); i++ {
		if tempP[i].Status != 'f' && tempP[i].Name == nameT {
			return true, i
		}
	}
	return false, 0
}

func verificarParticion(tempP [4]Particion) (int, int, int, int) {
	primaria := 0
	extendida := 0
	libre := 0
	espacioO := 0
	for i := 0; i < len(tempP); i++ {
		if tempP[i].Status == 'f' {
			libre++
		} else if tempP[i].Type == 'p' {
			primaria++
			espacioO = espacioO + int(tempP[i].Size)
		} else if tempP[i].Type == 'e' {
			extendida++
			espacioO = espacioO + int(tempP[i].Size)
		}
	}
	return extendida, primaria, libre, espacioO
}

func readFileBinary(pathAB string) mbr {
	//Abrimos/creamos un archivo.
	file, err := os.Open(pathAB)
	defer file.Close()
	if err != nil { //validar que no sea nulo.
		log.Fatal(err)
	}

	//Declaramos variable de tipo mbr
	m := mbr{}
	//Obtenemos el tamanio del mbr
	var size int = int(unsafe.Sizeof(m))

	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, size)
	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	return m
}

func leerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number) //array de bytes

	_, err := file.Read(bytes) // Leido -> bytes
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func actualizarMbr(pathDisco string, tempMbr mbr) {
	file, err := os.Create(pathDisco)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		fmt.Println("No existe la direccion " + pathDisco)
	}
	ss := &tempMbr

	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, ss)
	escribirBytes(file, binario.Bytes())

	var cero int8 = 0
	s := &cero
	file.Seek(tempMbr.SizeD-1, 0)
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s)
	escribirBytes(file, binario2.Bytes())
}

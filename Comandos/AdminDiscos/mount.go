package Admindiscos

import (
	"MIA_2S_P1_201513656/Herramientas"
	"MIA_2S_P1_201513656/Structs"
	"fmt"
	"os"
	"strings"
)

func Mount(entrada []string)  {
	var name string
	var pathE string
	Valido := true
	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro,"")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			fmt.Println("ERROR MKDIS, valor desconocido de parametros ",valores[1])
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			break
		}

		//******************* PATH *************
		if strings.ToLower(valores[0]) == "path" {
			pathE = strings.ReplaceAll(valores[1],"\"","")
			fmt.Println(pathE)
			_, err := os.Stat(pathE)
			if os.IsNotExist(err) {
				fmt.Println("FDISK Error: El disco no existe")
				Valido = false
				break // Terminar el bucle porque encontramos un nombre único
			}
		//********************  NAME *****************
		} else if strings.ToLower(valores[0]) == "name" {
			// Eliminar comillas
			name = strings.ReplaceAll(valores[1], "\"", "")
			// Eliminar espacios en blanco al final
			name = strings.TrimSpace(name)
		
		//******************* ERROR EN LOS PARAMETROS *************
		} else {
			fmt.Println("MKDISK Error: Parametro desconocido: ", valores[0])
			break //por si en el camino reconoce algo invalido de una vez se sale
		}

		
	}

	if Valido{
		if pathE != ""{
			if name != ""{
				// Abrir y cargar el disco
				disco, err := Herramientas.OpenFile(pathE)
				if err != nil {
					fmt.Println("MOUNT Error: No se pudo leer el disco")
					return
				}

				//Se crea un mbr para cargar el mbr del disco
				var mbr Structs.MBR
				//Guardo el mbr leido
				if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
					return
				}

				

				// cerrar el archivo del disco
				defer disco.Close()
			}else{
				fmt.Println("ERROR: FALTA NAME PATH EN MOUNT")				
			}
		}else{
			fmt.Println("ERROR: FALTA PARAMETRO PATH EN MOUNT")
		}
	}

	
}
package Admindiscos

import (
	"MIA_2S_P1_201513656/Herramientas"
	"MIA_2S_P1_201513656/Structs"
	"fmt"
	"os"
	"strings"
)

func Mount(entrada []string) string {
	var respuesta string
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
					respuesta += "ERROR NO SE PUEDE LEER EL DISCO " + err.Error()+ "\n"
					return  respuesta
				}

				//Se crea un mbr para cargar el mbr del disco
				var mbr Structs.MBR
				//Guardo el mbr leido
				if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
					respuesta += "ERROR Read " + err.Error()+ "\n"
					return  respuesta
				}				

				montar := true //usar si se van a montar logicas
				//reportar := false
				for i := 0; i < 4; i++ {
					nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
					if nombre == name{
						montar = false
						if string(mbr.Partitions[i].Status[:]) != "A" {
							if string(mbr.Partitions[i].Type[:]) != "E" {
								//COMIENZO PARA MONTAR
								Structs.PrintMBR(mbr)
								IdMbr := Structs.GetIdMBR(mbr)
								fmt.Println(IdMbr)

								mount := make([]Structs.Mount, 0)//Creamos un slice para contener los datos de los discos montados
								InitMount := false
								cont :=1
								
								//Si el slice no esta vacio, buscamos si ya existe el disco
								if len(mount) > 0{
									for _,montado := range mount{
										fmt.Println("Id: ",montado.Id, ", cont: ", montado.Cont, ", Letra: ", montado.Letter)
										if montado.Id ==IdMbr{
											InitMount = true										
											break 
										}
									}

									if InitMount{
										fmt.Println("Existe")
									}else{
										nuevaLetra := mount[len(mount)-1].Letter
										ultima := nuevaLetra[0]
										ultima++									
									
										mount = append(mount, Structs.Mount{Id: int32(IdMbr),Letter: [1]byte{ultima},Cont: int32(cont)})								
									}
								//Si el slice esta vacio sera el primer dato en agregar
								}else{
									mount = append(mount, Structs.Mount{Id: int32(IdMbr),Letter: [1]byte{'A'},Cont: int32(cont)})
									nuevaLetra := mount[len(mount)-1].Letter
									ultima := nuevaLetra[0]
									ultima++									
									
									mount = append(mount, Structs.Mount{Id: int32(IdMbr),Letter: [1]byte{ultima},Cont: int32(cont)})
									respuesta+="Se monto la particion, sin agregar a Disco"
								}

								for _,montado := range mount{
									fmt.Println("Id: ",montado.Id, ", cont: ", montado.Cont, ", Letra: ", string(montado.Letter[:]))}
											
							}else{
								fmt.Println("MOUNT Error. No se puede montar una particion extendida")
								respuesta += "MOUNT Error. No se puede montar una particion extendida"
								Structs.PrintMBR(mbr)
							}
						}
					}
				}

				if montar {
					fmt.Println("MOUNT Error. No se pudo montar la particion ", name)
					fmt.Println("MOUNT Error. No se encontro la particion")
					respuesta += "MOUNT Error. NO SE ENCONTRO LA PARTICION " + name
					respuesta += "\nNO SE PUDO MONTAR LA PARICION \n"
				}
				

				// cerrar el archivo del disco
				defer disco.Close()
			}else{
				fmt.Println("ERROR: FALTA NAME  EN MOUNT")	
				respuesta += "ERROR: FALTA NAME  EN MOUNT"			
			}
		}else{
			fmt.Println("ERROR: FALTA PARAMETRO PATH EN MOUNT")
			respuesta += "ERROR: FALTA PATH EN MOUNT"	
		}
	}

	return respuesta
	
}
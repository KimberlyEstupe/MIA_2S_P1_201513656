package Admindiscos

import (
	"MIA_2S_P1_201513656/Herramientas"
	"MIA_2S_P1_201513656/Structs"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var Pmontaje []Structs.Mount//GUarda en Ram las particones montadas
func Mount(entrada []string) (string){
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
				reportar := false
				for i := 0; i < 4; i++ {
					nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
					if nombre == name{
						montar = false
						if string(mbr.Partitions[i].Status[:]) != "A" {
							if string(mbr.Partitions[i].Type[:]) != "E" {
								//COMIENZO PARA MONTAR
								
								var id string
								var ultima byte // ultima letra a agregar
								ultima = 65 // A
								contador :=0 //Cantidad de particiones montadas

								if len(Pmontaje) > 0{
									var nuevaLetra [1]byte
									
									
									for _,montado := range Pmontaje{
										if montado.MPath ==pathE{
											nuevaLetra = montado.Letter	
											contador = int(montado.Cont)	
											contador++							
											break 
										}
									}

									if contador!=0{									
										ultima = nuevaLetra[0]
									}else{
										nuevaLetra := Pmontaje[len(Pmontaje)-1].Letter
										ultima = nuevaLetra[0]
										ultima++	
										contador++																
									}
								//Si el slice esta vacio sera el primer dato en agregar
								}else{
									contador++																		
								}

								id = "56"+strconv.Itoa(contador)+string(ultima) //Id de particion
								//ingresar al struck de particiones montadas
								Pmontaje = append(Pmontaje, Structs.Mount{MPath: pathE ,Letter: [1]byte{ultima},Cont: int32(contador), Id: id})

								//modificar la particion que se va a montar
								copy(mbr.Partitions[i].Status[:], "A")
								copy(mbr.Partitions[i].Id[:], id)

								//sobreescribir el mbr para guardar los cambios
								if err := Herramientas.WriteObject(disco, mbr, 0); err != nil { //Sobre escribir el mbr
									respuesta += "Error "
								}
								reportar = true

								respuesta+="Particion con nombre "+ name+ " montada correctamente. ID: "+id
								fmt.Println("Particion con nombre ", name, " montada correctamente. ID: ",id)

											
							}else{
								fmt.Println("MOUNT Error. No se puede montar una particion extendida")
								respuesta += "MOUNT Error. No se puede montar una particion extendida"
								Structs.PrintMBR(mbr)
							}
						}
					}
				}

				// cerrar el archivo del disco
				defer disco.Close()

				if montar {
					fmt.Println("MOUNT Error. No se pudo montar la particion ", name)
					fmt.Println("MOUNT Error. No se encontro la particion")
					respuesta += "MOUNT Error. NO SE ENCONTRO LA PARTICION " + name
					respuesta += "\nNO SE PUDO MONTAR LA PARICION \n"
				}

				if reportar {
					fmt.Println("\nLISTA DE PARTICIONES MONTADAS\n ")
					for i := 0; i < 4; i++ {
						estado := string(mbr.Partitions[i].Status[:])
						if estado == "A" {
							fmt.Printf("Partition %d: name: %s, status: %s, id: %s, tipo: %s, start: %d, size: %d, fit: %s, correlativo: %d\n", i, string(mbr.Partitions[i].Name[:]), string(mbr.Partitions[i].Status[:]), string(mbr.Partitions[i].Id[:]), string(mbr.Partitions[i].Type[:]), mbr.Partitions[i].Start, mbr.Partitions[i].Size, string(mbr.Partitions[i].Fit[:]), mbr.Partitions[i].Correlative)
						}
					}
				}
			}else{
				fmt.Println("ERROR: FALTA NAME  EN MOUNT")	
				respuesta += "ERROR: FALTA NAME  EN MOUNT"			
			}
		}else{
			fmt.Println("ERROR: FALTA PARAMETRO PATH EN MOUNT")
			respuesta += "ERROR: FALTA PATH EN MOUNT"	
		}
	}

	for _,montado := range Pmontaje{
		fmt.Println("Path: ",montado.MPath, "Letra: ",string(montado.Letter[:])," Contador: ",montado.Cont, " Id: ", montado.Id)
	}

	return respuesta
	
}
package Usuarios

import (
	"MIA_2S_P1_201513656/Herramientas"
	"MIA_2S_P1_201513656/Structs"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func Mkgrp(entrada []string) string{
	var respuesta string
	var name string
	UsuarioA := Structs.UsuarioActual
	
	if !UsuarioA.Status {
		respuesta += "ERROR MKGRP: NO HAY SECION INICIADA"+ "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR"+ "\n"
		return respuesta
	}

	for _,parametro :=range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			fmt.Println("ERROR MKGRP, valor desconocido de parametros ",valores[1])
			respuesta += "ERROR MKGRP, valor desconocido de parametros " + valores[1]+ "\n"
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			return respuesta
		}
		tmp = strings.TrimRight(valores[1],"")
		valores[1] = tmp

		//********************  NAME *****************
		if strings.ToLower(valores[0]) == "name" {
			tmp = strings.ReplaceAll(valores[1],"\"","")
			name = (tmp)
			//validar maximo 10 caracteres
			if len(name) > 10 {
				fmt.Println("MKGRP ERROR: name debe tener maximo 10 caracteres")
				return "ERROR MKGRP: name debe tener maximo 10 caracteres"
			}
		//******************* ERROR EN LOS PARAMETROS *************
		} else {
			fmt.Println("LOGIN ERROR: Parametro desconocido: ", valores[0])
			//por si en el camino reconoce algo invalido de una vez se sale
			return "LOGIN ERROR: Parametro desconocido: "+valores[0] + "\n"
		}
	}

	if UsuarioA.Nombre == "root"{
		file, err := Herramientas.OpenFile(UsuarioA.PathD)
		if err != nil {
			return "ERROR REP SB OPEN FILE "+err.Error()+ "\n"
		}

		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR REP SB READ FILE "+err.Error()+ "\n"
		}

		// Close bin file
		defer file.Close()

		//Encontrar la particion correcta
		AddNewUser := false
		part := -1
		for i := 0; i < 4; i++ {		
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == UsuarioA.IdPart {
				part = i
				AddNewUser = true
				break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
			}
		}

		if AddNewUser{
			var superBloque Structs.Superblock
			errREAD := Herramientas.ReadObject(file, &superBloque, int64(mbr.Partitions[part].Start))
			if errREAD != nil {
				fmt.Println("REP Error. Particion sin formato")
				return "REP Error. Particion sin formato"+ "\n"
			}

			var inodo Structs.Inode		
			//Le agrego una structura de inodo para ver el user.txt que esta en el primer inodo del sb
			Herramientas.ReadObject(file, &inodo, int64(superBloque.S_inode_start + int32(binary.Size(Structs.Inode{}))))
			
			//leer los datos del user.txt
			var contenido string
			var fileBlock Structs.Fileblock
			var idFb int32 //id/numero de ultimo fileblock para trabajar sobre ese
			for _, item := range inodo.I_block {
				if item != -1 {
					Herramientas.ReadObject(file, &fileBlock, int64(superBloque.S_block_start+(item*int32(binary.Size(Structs.Fileblock{})))))
						contenido += string(fileBlock.B_content[:])
						idFb = item
				}
			}

			lineaID := strings.Split(contenido, "\n")

			//Verificar si el grupo ya existe
			for _, registro := range lineaID[:len(lineaID)-1] {
				datos := strings.Split(registro, ",")
				if len(datos) == 3 {
					if datos[2] == name {
						fmt.Println("MKGRP ERROR: El grupo ya existe")
						return "MKGRP ERROR: El grupo ya existe"
					}
				}
			}


			//Buscar el ultimo ID activo desde el ultimo hasta el primero (ignorando los eliminado (0))
			//desde -2 porque siempre se crea un salto de linea al final generando una linea vacia al final del arreglo
			id := -1        //para guardar el nuevo ID
			var errId error //para la conversion a numero del ID
			for i := len(lineaID) - 2; i >= 0; i--{
				registro := strings.Split(lineaID[i], ",")
				//valido que sea un grupo
				if registro[1] == "G"{
					//valido que el id sea distinto a 0 (eliminado)
					if registro[0] != "0"{
						//convierto el id en numero para sumarle 1 y crear el nuevo id
						id, errId = strconv.Atoi(registro[0])
						if errId != nil {
							fmt.Println("MKGRP ERROR: No se pudo obtener un nuevo id para el nuevo grupo")
							return "MKGRP ERROR: No se pudo obtener un nuevo id para el nuevo grupo"
						}
						id++
						break
					}
				}
			}
			

			//valido que se haya encontrado un nuevo id
			if id != -1 {				
				contenidoActual := string(fileBlock.B_content[:])
				posicionNulo := strings.IndexByte(contenidoActual, 0)			
				data := fmt.Sprintf("%d,G,%s\n", id, name)
				//Aseguro que haya al menos un byte libre
				if posicionNulo != -1 {
					libre := 64 - (posicionNulo + len(data))
					if libre > 0 {
						copy(fileBlock.B_content[posicionNulo:], []byte(data))
						//Escribir el fileblock con espacio libre
						Herramientas.WriteObject(file, fileBlock, int64(superBloque.S_block_start+(idFb*int32(binary.Size(Structs.Fileblock{})))))
					}else{
						//Si es 0 (quedó exacta), entra aqui y crea un bloque vacío que podrá usarse para el proximo registro
						data1 := data[:len(data)+libre]
						//Ingreso lo que cabe en el bloque actual
						copy(fileBlock.B_content[posicionNulo:], []byte(data1))
						Herramientas.WriteObject(file, fileBlock, int64(superBloque.S_block_start+(idFb*int32(binary.Size(Structs.Fileblock{})))))

						//Creo otro fileblock para el resto de la informacion
						guardoInfo := true

						for i, item := range inodo.I_block{
							//i es el indice en el arreglo inodo.Iblock
							// DIferencia i/item:  inodo.I_block[i] = item
							if item == -1 {
								guardoInfo = false
								//agrego el apuntador del bloque al inodo
								inodo.I_block[i] = superBloque.S_first_blo
								//actualizo el superbloque
								superBloque.S_free_blocks_count -= 1
								superBloque.S_first_blo += 1
								data2 := data[len(data)+libre:]
								//crear nuevo fileblock
								var newFileBlock Structs.Fileblock
								copy(newFileBlock.B_content[:], []byte(data2))

								//escribir las estructuras para guardar los cambios
								// Escribir el superbloque
								Herramientas.WriteObject(file, superBloque, int64(mbr.Partitions[part].Start))

								//escribir el bitmap de bloques (se uso un bloque). inodo.I_block[i] contiene el numero de bloque que se uso
								Herramientas.WriteObject(file, byte(1), int64(superBloque.S_bm_block_start+inodo.I_block[i]))

								//escribir inodes (es el inodo 1, porque es donde esta users.txt)
								Herramientas.WriteObject(file, inodo, int64(superBloque.S_inode_start+int32(binary.Size(Structs.Inode{}))))

								//Escribir bloques
								Herramientas.WriteObject(file, newFileBlock, int64(superBloque.S_block_start+(inodo.I_block[i]*int32(binary.Size(Structs.Fileblock{})))))
								break
							}
						}

						if guardoInfo {
							fmt.Println("MKGRP ERROR: Espacio insuficiente para nuevo registro")
							return "MKGRP ERROR: Espacio insuficiente para nuevo registro. "
						}
					}

					
					fmt.Println("Se ha agregado el grupo '"+name+"' exitosamente. ")
					respuesta = "Se ha agregado el grupo '"+name+"' exitosamente."
					for k:=0; k<len(lineaID)-1; k++{
						fmt.Println(lineaID[k])
					}
					return respuesta
				}
			}
		//FIn Add new Usuario
		}else{	
			fmt.Println("ERROR INESPERADO CON LA PARCION EN MKGRP")
			respuesta += "ERROR INESPERADO CON LA PARCION EN MKGRP"
		}

	}else{
		fmt.Println("ERROR FALTA DE PERMISOS, NO ES EL USUARIO ROOT")
		respuesta += "ERROR MKGRO: ESTE USUARIO NO CUENTA CON LOS PERMISOS PARA REALIZAR ESTA ACCION"
	}

	return respuesta	
}
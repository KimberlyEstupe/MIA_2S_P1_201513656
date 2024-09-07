package Usuarios

import (
	"MIA_2S_P1_201513656/Herramientas"
	"MIA_2S_P1_201513656/Structs"
	"encoding/binary"
	"fmt"
	"strings"
)

func Rmgrp(entrada []string) string {
	var respuesta string
	var name string
	UsuarioA := Structs.UsuarioActual
	if !UsuarioA.Status {
		respuesta += "ERROR RMGRP: NO HAY SECION INICIADA"+ "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR"+ "\n"
		return respuesta
	}

	for _,parametro :=range entrada[1:]{
		tmp := strings.TrimRight(parametro,"")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			fmt.Println("ERROR RMGRP, valor desconocido de parametros ",valores[1])
			respuesta += "ERROR RMGRP, valor desconocido de parametros " + valores[1]+ "\n"
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			return respuesta
		}

		//********************  NAME *****************
		if strings.ToLower(valores[0]) == "name" {
			name = (valores[1])
			//validar maximo 10 caracteres
			if len(name) > 10 {
				fmt.Println("RMGRP ERROR: name debe tener maximo 10 caracteres")
				return "ERROR RMGRP: name debe tener maximo 10 caracteres"
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
		delete := false
		part := -1 //particion a utilizar y modificar
		for i := 0; i < 4; i++ {		
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == UsuarioA.IdPart {
				part = i
				delete = true
				break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
			}
		}

		if delete{
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
			for _, item := range inodo.I_block {
				if item != -1 {
					Herramientas.ReadObject(file, &fileBlock, int64(superBloque.S_block_start+(item*int32(binary.Size(Structs.Fileblock{})))))
					contenido += string(fileBlock.B_content[:])
				}
			}

			lineaID := strings.Split(contenido, "\n")

			for _, registro := range lineaID[:len(lineaID)-1]{
				datos := strings.Split(registro, ",")
				if len(datos) == 3 {
					if datos[2] == name {
						fmt.Println(datos)
					}
				}
			}
		}


	}else{
		fmt.Println("ERROR FALTA DE PERMISOS, NO ES EL USUARIO ROOT")
		respuesta += "ERROR MKGRO: ESTE USUARIO NO CUENTA CON LOS PERMISOS PARA REALIZAR ESTA ACCION"
	}
	
	return respuesta
}